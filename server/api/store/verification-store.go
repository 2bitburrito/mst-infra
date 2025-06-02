package store

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
)

type VerificationStore struct {
	Map             map[uuid.UUID]verificationCode
	TimeoutDuration time.Duration
	M               sync.Mutex
}

type verificationCode struct {
	CreatedAt time.Time
	Code      string
}

func CreateVerificationStore(validityDuration time.Duration, reapDuration time.Duration) *VerificationStore {
	store := VerificationStore{
		Map:             make(map[uuid.UUID]verificationCode),
		TimeoutDuration: validityDuration,
	}
	go store.reapLoop(reapDuration)
	return &store
}

func (s *VerificationStore) New(id uuid.UUID) string {
	s.M.Lock()
	defer s.M.Unlock()
	otc := GenerateOTC()

	s.Map[id] = verificationCode{
		CreatedAt: time.Now().UTC(),
		Code:      otc,
	}
	return otc
}

func (s *VerificationStore) GetFromOTC(otc string) (uuid.UUID, string, error) {
	s.M.Lock()
	defer s.M.Unlock()
	log.Println("Checking VerificationStore using OTC")

	now := time.Now()

	for id, code := range s.Map {
		if code.Code == otc {
			expiryTime := code.CreatedAt.Add(s.TimeoutDuration)
			if expiryTime.Before(now) {
				return uuid.UUID{}, "", fmt.Errorf("otc: %s has timed out", otc)
			}
			return id, code.Code, nil
		}
	}

	return uuid.UUID{}, "", fmt.Errorf("otc %s doesn't exist in store", otc)
}

func (s *VerificationStore) Get(id uuid.UUID) (string, error) {
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

func (s *VerificationStore) Delete(id uuid.UUID) {
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
