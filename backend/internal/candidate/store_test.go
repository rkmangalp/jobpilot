package candidate

import "testing"

func TestUpdateIncrementsVersion(t *testing.T) {
	s := NewMemoryStore()
	before := s.Get()
	after := s.Update(before)
	if after.Version != before.Version+1 {
		t.Fatalf("got %d, want %d", after.Version, before.Version+1)
	}
}
