package eventsourcedb

type stubDB struct {
	_Insert      func(...Event) error
	_InsertCalls int
	_Fetch       func(uint64, uint64) ([]Event, error)
	_FetchCalls  int
	_Close       func() error
	_CloseCalls  int
}

func (db *stubDB) Insert(evts ...Event) error {
	db._InsertCalls++
	return db._Insert(evts...)
}

func (db *stubDB) Fetch(firstID, lastID uint64) ([]Event, error) {
	db._FetchCalls++
	return db._Fetch(firstID, lastID)
}
func (db *stubDB) Close() error {
	db._CloseCalls++
	return db._Close()
}

func NewSimplestubDB() *stubDB {
	return &stubDB{
		_Insert: func(...Event) error {
			return nil
		},
		_Fetch: func(uint64, uint64) ([]Event, error) {
			return nil, nil
		},
		_Close: func() error {
			return nil
		},
	}
}

func InMemStubDB() *stubDB {
	evtMap := make(map[uint64]Event)
	db := &stubDB{}

	db._Insert = func(evts ...Event) error {
		for _, e := range evts {
			evtMap[e.ID] = e
		}
		return nil
	}

	db._Fetch = func(low uint64, up uint64) ([]Event, error) {
		evts := make([]Event, 0)
		for id, e := range evtMap {
			if id >= low || id <= up {
				evts = append(evts, e)
			}
		}
		return evts, nil
	}

	db._Close = func() error {
		for id, _ := range evtMap {
			delete(evtMap, id)
		}
		return nil
	}

	return db
}
