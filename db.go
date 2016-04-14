package eventsourcedb

type DB interface {
	Insert(events ...Event) error
	Fetch(firstID, lastID uint64) ([]Event, error)
	Close() error
}
