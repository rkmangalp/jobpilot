package candidate

import "sync"

type Store interface {
	Get() Profile
	Update(Profile) Profile
}
type MemoryStore struct {
	mu      sync.RWMutex
	profile Profile
}

func NewMemoryStore() *MemoryStore  { return &MemoryStore{profile: Seed()} }
func (s *MemoryStore) Get() Profile { s.mu.RLock(); defer s.mu.RUnlock(); return s.profile }
func (s *MemoryStore) Update(p Profile) Profile {
	s.mu.Lock()
	defer s.mu.Unlock()
	p.ID = s.profile.ID
	p.Version = s.profile.Version + 1
	s.profile = p
	return s.profile
}
