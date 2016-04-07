package eventsourcedb

type stubDB struct {
	_Insert      func(...Event) error
	_InsertCalls int
	_Fetch       func(uint64, uint64) ([]Event, error)
	_FetchCalls  int
}

func (db *stubDB) Insert(evts ...Event) error {
	db._InsertCalls++
	return db._Insert(evts...)
}

func (db *stubDB) Fetch(firstID, lastID uint64) ([]Event, error) {
	db._FetchCalls++
	return db._Fetch(firstID, lastID)
}

func NewSimplestubDB() *stubDB {
	return &stubDB{
		_Insert: func(...Event) error {
			return nil
		},
		_Fetch: func(uint64, uint64) ([]Event, error) {
			return nil, nil
		},
	}
}
