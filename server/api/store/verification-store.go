package store

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type VerificationStore struct {
	Map             map[userId]VerificationCode
	TimeoutDuration time.Duration
	M               sync.Mutex
}

type VerificationCode struct {
	CreatedAt time.Time
	Code      string
}

type userId string

func CreateStore(validityDuration time.Duration, reapDuration time.Duration) *VerificationStore {
	store := VerificationStore{
		TimeoutDuration: validityDuration,
	}
	go store.reapLoop(reapDuration)
	return &store
}

func (s *VerificationStore) Set(id userId) {
	s.M.Lock()
	defer s.M.Unlock()
	s.Map[id] = VerificationCode{
		CreatedAt: time.Now().UTC(),
		Code:      GenerateOTC(),
	}
}

func (s *VerificationStore) Get(id userId) (string, error) {
	s.M.Lock()
	defer s.M.Unlock()

	obj, exists := s.Map[id]
	if !exists {
		return "", fmt.Errorf("ID %s doesn't exist in store", id)
	}

	expiryTime := obj.CreatedAt.Add(s.TimeoutDuration)
	now := time.Now().UTC()

	if expiryTime.Before(now) {
		return "", fmt.Errorf("otc for %s has timed out", id)
	}
	return obj.Code, nil
}

func (s *VerificationStore) Delete(id userId) {
	s.M.Lock()
	defer s.M.Unlock()

	delete(s.Map, id)
}

func (s *VerificationStore) reap() {
	s.M.Lock()
	defer s.M.Unlock()
	var count int

	now := time.Now().UTC()
	for id, val := range s.Map {
		if val.CreatedAt.Before(now) {
			delete(s.Map, id)
			count++
		}
	}
	if count > 0 {
		log.Printf("Reaped %d Entries from VerificationStore", count)
	}
}

func (s *VerificationStore) reapLoop(interval time.Duration) {
	tick := time.NewTicker(interval)
	for range tick.C {
		s.reap()
	}
}
