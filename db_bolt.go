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

func (b *BoltDB) Insert(events ...Event) error {
	select {
	case <-b.ctx.Done():
		return ErrDBClosed
	default:
		b.wg.Add(1)
	}

	return b.db.Update(func(tx *bolt.Tx) error {
		defer b.wg.Done()

		buEvents := bucket(tx, eventsName)
		buStreams := bucket(tx, streamsName)
		buTypes := bucket(tx, typesName)

		for _, e := range events {
			e.ID, _ = buEvents.NextSequence()
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
}

func (b *BoltDB) Fetch(firstID, lastID uint64) ([]Event, error) {
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
	result := make([]Event, 0)

	b.db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket(eventsName).Cursor()

		for k, v := c.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, v = c.Next() {
			var e Event
			if err := e.UnmarshalJSON(v); err != nil {
				return err
			}
			result = append(result, e)
		}

		return nil
	})

	return result, nil
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
