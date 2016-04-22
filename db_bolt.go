package eventsourcedb

import (
	"bytes"
	"encoding/binary"
	"sync"
	"time"

	"golang.org/x/net/context"

	"github.com/boltdb/bolt"
)

var (
	// make sure *BoltDB implements DB interface
	_           DB = (*BoltDB)(nil)
	eventsName     = []byte("events")
	streamsName    = []byte("streams")
	typesName      = []byte("types")
)

type BoltDB struct {
	db   *bolt.DB
	wg   *sync.WaitGroup
	ctx  context.Context
	done context.CancelFunc
}

func (b *BoltDB) Insert(events ...Event) (uint64, error) {
	select {
	case <-b.ctx.Done():
		return 0, ErrDBClosed
	default:
		b.wg.Add(1)
	}

	var lastID uint64

	err := b.db.Update(func(tx *bolt.Tx) error {
		defer b.wg.Done()

		buEvents := bucket(tx, eventsName)
		buStreams := bucket(tx, streamsName)
		buTypes := bucket(tx, typesName)

		for _, e := range events {
			e.ID, _ = buEvents.NextSequence()
			lastID = e.ID
			evt_bytes, err := e.MarshalJSON()
			if err != nil {
				return err
			}

			key := make([]byte, 8)
			binary.PutUvarint(key, e.ID)
			buStream, _ := buStreams.CreateBucketIfNotExists([]byte(e.Stream))
			buType, _ := buTypes.CreateBucketIfNotExists([]byte(e.Type))

			err = buEvents.Put(key, evt_bytes)
			if err != nil {
				return err
			}

			err = buStream.Put(key, nil)
			if err != nil {
				return err
			}

			err = buType.Put(key, nil)
			if err != nil {
				return err
			}
		}
		return nil
	})

	return lastID, err
}

func (b *BoltDB) Fetch(firstID, lastID uint64) (Cursor, error) {
	select {
	case <-b.ctx.Done():
		return nil, ErrDBClosed
	default:
		b.wg.Add(1)
		defer b.wg.Done()
	}

	min := make([]byte, 8)
	max := make([]byte, 8)
	binary.PutUvarint(min, firstID)
	binary.PutUvarint(max, lastID)

	cur := &boltCursor{
		next: min,
		max:  max,
		db:   b.db,
		ctx:  b.ctx,
	}

	return cur, nil
}

func (b *BoltDB) Close() error {
	select {
	case <-b.ctx.Done():
		return nil
	default:
	}

	b.done()
	b.wg.Wait()

	return b.db.Close()
}

func (b *BoltDB) createBuckets() {
	b.db.Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists(eventsName)
		tx.CreateBucketIfNotExists(streamsName)
		tx.CreateBucketIfNotExists(typesName)
		return nil
	})
}

func OpenBolt(fname string) (*BoltDB, error) {
	db, err := bolt.Open(fname, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, err
	}
	ctx, cancelFunc := context.WithCancel(context.Background())

	b := &BoltDB{
		db:   db,
		ctx:  ctx,
		done: cancelFunc,
		wg:   new(sync.WaitGroup),
	}
	b.createBuckets()
	return b, nil
}

func bucket(tx *bolt.Tx, name []byte) *bolt.Bucket {
	bu := tx.Bucket(name)
	bu.FillPercent = 0.99
	return bu
}

type boltCursor struct {
	next []byte
	max  []byte
	db   *bolt.DB
	ctx  context.Context
}

func (c *boltCursor) Next() (Event, error) {
	var (
		e Event
	)

	select {
	case <-c.ctx.Done():
		return e, ErrDBClosed
	default:
	}

	err1 := c.db.View(func(tx *bolt.Tx) error {
		cur := tx.Bucket(eventsName).Cursor()
		k, v := cur.Seek(c.next)

		if k == nil || bytes.Compare(k, c.max) > 0 {
			return EOS
		}

		if err := e.UnmarshalJSON(v); err != nil {
			return err
		}

		next, _ := cur.Next()
		c.next = make([]byte, len(next))
		copy(c.next, next)

		return nil
	})

	return e, err1
}
