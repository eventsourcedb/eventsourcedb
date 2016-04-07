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

func TestHub_Pub(t *testing.T) {
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
