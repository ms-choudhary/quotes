package store

import (
	"math/rand"
	"sync"
	"time"
)

type fakeStore struct {
	mu      sync.Mutex
	counter int
	db      map[int]Quote
}

func NewFakeStore() *fakeStore {
	return &fakeStore{counter: 0, db: map[int]Quote{}}
}

func (s *fakeStore) Ping() error {
	return nil
}

func (s *fakeStore) Create(q Quote) (Quote, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.counter += 1
	s.db[s.counter] = q
	q.ID = s.counter
	return q, nil
}

func (s *fakeStore) Get(id int) (Quote, error) {
	return s.db[id], nil
}

func (s *fakeStore) GetAll() ([]Quote, error) {
	res := []Quote{}
	for _, v := range s.db {
		res = append(res, v)
	}
	return res, nil
}

func (s *fakeStore) GetRandom() (Quote, error) {
	rand.Seed(time.Now().UnixNano())
	return s.Get(rand.Intn(s.counter-1) + 1)
}

func (s *fakeStore) Clean() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.counter = 0
	s.db = map[int]Quote{}
	return nil
}
