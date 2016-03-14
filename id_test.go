package eventsourcedb

import "testing"

func BenchmarkNewID(b *testing.B) {
	var oldID ID
	for i := 0; i < b.N; i++ {
		id := New(uint64(i % 100))
		if id == oldID {
			b.Error(id)
		}
		oldID = id
	}
}
