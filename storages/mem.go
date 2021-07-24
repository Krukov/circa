package storages

import (
	"circa/message"
	"errors"
	"net/url"
	"sync"
	"time"
)

type Memory struct {
	storage    map[string]*message.Response
	intStorage map[string]int
	maxSize    int
	lock       *sync.Mutex
}

func NewMemStorageFromURL(sURL *url.URL) (*Memory, error) {
	return &Memory{map[string]*message.Response{}, map[string]int{}, 100, &sync.Mutex{}}, nil
}

func (s *Memory) String() string {
	return "memory"
}

func (s *Memory) Set(key string, value *message.Response, ttl time.Duration) (bool, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	_, ok := s.storage[key]
	if !ok && len(s.storage) > s.maxSize {
		return false, errors.New("overflow")
	}
	s.storage[key] = value
	return ok, nil
}

func (s *Memory) Del(key string) (bool, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	_, ok := s.storage[key]
	if !ok {
		_, ok = s.intStorage[key]
		if !ok {
			return false, NotFound
		}
		delete(s.intStorage, key)
	}
	delete(s.storage, key)

	return true, nil
}

func (s *Memory) Get(key string) (*message.Response, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	value, ok := s.storage[key]
	if !ok {
		return nil, NotFound
	}
	return value, nil
}

func (s *Memory) Incr(key string) (int, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.intStorage[key] = s.intStorage[key] + 1
	return s.intStorage[key], nil
}
