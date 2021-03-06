package eventsourcedb

import "sync"

const (
	ErrNoEventsData = Error("no events data")
	EOS             = Error("end of stream")
)

var (
	subBufSize = 10
)

type HubOpt func(*Hub)

func NewHub(opts ...HubOpt) *Hub {
	h := &Hub{
		subs:       make(map[Sub]struct{}),
		subBufSize: subBufSize,
	}

	for _, o := range opts {
		o(h)
	}
	return h
}

type Hub struct {
	rw         sync.RWMutex
	db         DB
	subBufSize int
	subs       map[Sub]struct{}
}

func (h *Hub) Pub(evts ...Event) error {
	if len(evts) == 0 {
		return ErrNoEventsData
	}

	lastID, err := h.db.Insert(evts...)
	if err != nil {
		return err
	}

	h.rw.RLock()
	defer h.rw.RUnlock()

	for s, _ := range h.subs {
		select {
		case s.Events <- lastID:
		default:
		}
	}

	return nil
}

func (h *Hub) Sub() *Sub {
	sub := Sub{
		Events: make(chan uint64, h.subBufSize),
		db:     h.db,
	}
	h.rw.Lock()
	defer h.rw.Unlock()

	h.subs[sub] = struct{}{}
	return &sub
}

func (h *Hub) Cancel(sub *Sub) {
	h.rw.Lock()
	defer h.rw.Unlock()

	delete(h.subs, *sub)
	close(sub.Events)
}

type Sub struct {
	Events chan uint64
	db     DB
	lastID uint64
	mu     sync.Mutex
}

func (s *Sub) Next() (Cursor, error) {
	curID, ok := <-s.Events
	if !ok {
		return nil, EOS
	}

	s.mu.Lock()
	lastID := s.lastID
	s.lastID = curID
	s.mu.Unlock()

	return s.db.Fetch(lastID+1, curID)
}
