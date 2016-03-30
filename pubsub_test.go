package eventsourcedb

import "testing"

var (
	stream1 = "topic1"
	evt1    = Event{
		ID:     42,
		Stream: stream1,
		Type:   "evt1",
		Body:   []byte("body1"),
	}
)

func TestPub(t *testing.T) {
	store := &noopBackend{}

	f := func(h *hub) {
		h.store = store
	}

	h := newhub(f)
	h.Pub([]byte(stream1), evt1)

	if store.PersistCalled != 1 {
		t.Error("didn't called Persist")
	}
}

type noopBackend struct {
	PersistCalled int
}

func (b *noopBackend) Persist(stream []byte, e ...Event) error {
	b.PersistCalled += 1
	return nil
}
