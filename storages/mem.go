package storages

import (
	"circa/message"
	"time"
)

type Memory struct {
	storage map[string]*message.Response
}

func (s *Memory) String() string {
	return "memory"
}

func (s *Memory) Set(key string, value *message.Response, ttl time.Duration) (bool, error) {
	_, ok := s.storage[key]
	s.storage[key] = value
	return ok, nil
}

func (s *Memory) Del(key string) (bool, error) {
	_, ok := s.storage[key]
	if !ok {
		return false, NotFound
	}
	delete(s.storage, key)
	return true, nil
}

func (s *Memory) Get(key string) (*message.Response, error) {
	value, ok := s.storage[key]
	if !ok {
		return nil, NotFound
	}
	return value, nil
}

