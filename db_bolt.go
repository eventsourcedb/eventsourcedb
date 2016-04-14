package eventsourcedb

// make sure *BoltDB implements DB interface
var _ DB = (*BoltDB)(nil)

type BoltDB struct{}

func (b *BoltDB) Insert(events ...Event) error {
	return nil
}

func (b *BoltDB) Fetch(firstID, lastID uint64) ([]Event, error) {
	return nil, nil
}

func (b *BoltDB) Close() error {
	return nil
}
