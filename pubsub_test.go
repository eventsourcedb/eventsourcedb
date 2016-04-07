package eventsourcedb

import "testing"

func TestNewHub(t *testing.T) {
	var optCalled bool
	NewHub(func(h *Hub) {
		optCalled = true
	})

	if !optCalled {
		t.Fatal("optional function wasn't called")
	}
}

func TestHub_Cancel(t *testing.T) {
	h := NewHub()
	s1 := h.Sub()
	go func() {
		h.Cancel(s1)
	}()
	if _, err := s1.Next(); err != EOS {
		t.Fatal(EOS, "is expected, got", err)
	}
}

func TestHub_Pub_errors(t *testing.T) {
	mockDB1 := NewSimplestubDB()
	mockDB1Opt := func(h *Hub) {
		h.db = mockDB1
	}

	test_tab := []struct {
		opts   []HubOpt
		events []Event
		err    error
	}{
		{
			[]HubOpt{},
			[]Event{},
			ErrNoEventsData,
		},
		{
			[]HubOpt{mockDB1Opt},
			[]Event{},
			ErrNoEventsData,
		},
		{
			[]HubOpt{mockDB1Opt},
			[]Event{{}},
			nil,
		},
	}

	for _, tc := range test_tab {
		h := NewHub(tc.opts...)

		if err := h.Pub(tc.events...); err != tc.err {
			t.Fatalf("%q is expected, got %q", tc.err, err)
		}
	}
}

func TestHub_Pub(t *testing.T) {
	mockDB1 := InMemStubDB()
	h := NewHub()
	h.db = mockDB1
	h.subBufSize = 2
	s1 := h.Sub()
	evt1 := Event{ID: 1234}

	if err1 := h.Pub(evt1); err1 != nil {
		t.Fatal(err1)
	}

	evts, err2 := s1.Next()
	if err2 != nil {
		t.Fatal(err2)
	}
	t.Log(evts)
}
