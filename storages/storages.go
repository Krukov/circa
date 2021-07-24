package storages

import (
	"circa/message"
	"errors"
	"fmt"
	"net/url"
	"time"
)

type Options struct {
	Name string
	DSN  string
}

type Storage interface {
	String() string

	Set(key string, value *message.Response, ttl time.Duration) (bool, error)
	Del(key string) (bool, error)
	Incr(key string) (int, error)
	Get(key string) (*message.Response, error)
}

var NotFound = errors.New("key not found")

func StorageFormDSN(DSN string) (Storage, error) {
	sURL, err := url.Parse(DSN)
	if err != nil {
		return nil, fmt.Errorf("can't parse storage DSN: %v", err)
	}
	switch sURL.Scheme {
	case "mem":
		return NewMemStorageFromURL(sURL)
	case "redis":
		return NewRedisStorageFormURL(sURL)
	default:
		return nil, fmt.Errorf("unknown storage type: %s", sURL.Scheme)
	}
}
