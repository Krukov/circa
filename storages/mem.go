package storages

import (
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/bluele/gcache"

	"circa/message"
)

type Memory struct {
	gc       gcache.Cache
	incrLock sync.Mutex
}

func NewMemStorageFromURL(sURL *url.URL) (*Memory, error) {
	size := 10000
	var err error
	if _, ok := sURL.Query()["size"]; ok {
		size, err = strconv.Atoi(sURL.Query()["size"][0])
		if err != nil {
			return nil, err
		}
	}
	gc := gcache.New(size).LRU().Build()
	return &Memory{gc: gc, incrLock: sync.Mutex{}}, nil
}

func (s *Memory) String() string {
	return "memory"
}

func (s *Memory) Set(key string, value *message.Response, ttl time.Duration) (bool, error) {
	exists := s.gc.Has(key)
	err := s.gc.SetWithExpire(key, value, ttl)
	return !exists, err
}

func (s *Memory) Del(key string) (bool, error) {
	return s.gc.Remove(key), nil
}

func (s *Memory) Get(key string) (*message.Response, error) {
	resp, err := s.gc.Get(key)
	if err != nil {
		return nil, NotFound
	}
	return resp.(*message.Response), err
}

func (s *Memory) Incr(key string) (int, error) {
	countInt := 0
	s.incrLock.Lock()
	defer s.incrLock.Unlock()
	count, err := s.gc.Get(key)
	if err == nil {
		countInt = count.(int)
	}
	s.gc.Set(key, countInt+1)
	return countInt + 1, nil
}

func (s *Memory) Expire(key string, ttl time.Duration) error {
	value, err := s.gc.Get(key)
	if err != nil {
		return NotFound
	}
	s.gc.SetWithExpire(key, value, ttl)
	return nil
}
