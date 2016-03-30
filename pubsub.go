package eventsourcedb

func newhub(opts ...pubSubOpt) *hub {
	h := &hub{subs: make(chan sub)}
	for _, o := range opts {
		o(h)
	}
	return h
}

type pubSubOpt func(*hub)

type persister interface {
	Persist(stream []byte, events ...Event) error
}

type sub struct {
	events chan Event
	done   chan struct{}
}

type hub struct {
	store persister
	subs  chan sub
}

func (h *hub) Pub(stream []byte, e ...Event) error {
	err := h.store.Persist(stream, e...)
	return err
}
