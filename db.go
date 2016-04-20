package eventsourcedb

const (
	ErrCloseTimeout = Error("database close timeout")
	ErrDBClosed     = Error("database already closed")
)

type DB interface {
	Insert(events ...Event) error
	Fetch(firstID, lastID uint64) ([]Event, error)
	Close() error
}
