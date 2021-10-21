package storages

import (
	"errors"
	"fmt"
	"net/url"
	"time"

	"circa/message"
)

type Options struct {
	Name string
	DSN  string
}

type Storage interface {
	String() string

	Get(key string) (*message.Response, error)
	Set(key string, value *message.Response, ttl time.Duration) (bool, error)
	Incr(key string) (int, error)
	Del(key string) (bool, error)

	// SetRow(key string, value string)
	// GetRow(key string, value string)

	Expire(key string, ttl time.Duration) error
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
