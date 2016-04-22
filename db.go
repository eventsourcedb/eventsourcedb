package eventsourcedb

const (
	ErrCloseTimeout = Error("database close timeout")
	ErrDBClosed     = Error("database already closed")
)

type DB interface {
	Insert(events ...Event) (uint64, error)
	Fetch(firstID, lastID uint64) (Cursor, error)
	Close() error
}

type Cursor interface {
	Next() (Event, error)
}
