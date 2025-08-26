package store

import (
	"log"
	"sync"
	"time"
)

type HandledEventsStore struct {
	Set map[string]time.Time
	M   sync.Mutex
}

func CreateHandledEventStore(reapDuration time.Duration) *HandledEventsStore {
	store := HandledEventsStore{
		Set: make(map[string]time.Time),
	}
	go store.reapLoop(reapDuration)
	return &store
}

func (s *HandledEventsStore) Add(eventID string) {
	s.M.Lock()
	defer s.M.Unlock()

	s.Set[eventID] = time.Now().UTC()
}

func (s *HandledEventsStore) Exists(eventID string) bool {
	s.M.Lock()
	defer s.M.Unlock()

	_, exists := s.Set[eventID]
	return exists
}

func (s *HandledEventsStore) reap() {
	s.M.Lock()
	defer s.M.Unlock()
	var count int

	now := time.Now().UTC()
	for eventID, val := range s.Set {
		if val.Before(now) {
			delete(s.Set, eventID)
			count++
		}
	}
	if count > 0 {
		log.Printf("Reaped %d Entries from HandledEventsStore", count)
	}
}

func (s *HandledEventsStore) reapLoop(interval time.Duration) {
	tick := time.NewTicker(interval)
	for range tick.C {
		s.reap()
	}
}
