package storages

import (
	"circa/message"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
)

type Options struct {
	Name string
	DSN  string
}

type Storage interface {
	String() string

	Set(key string, value *message.Response, ttl time.Duration) (bool, error)
	Del(key string) (bool, error)
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
		return &Memory{map[string]*message.Response{}, 100}, nil
	case "redis":
		return createRedisStorage(sURL)
	default:
		return nil, fmt.Errorf("unknown storage type: %s", sURL.Scheme)
	}
}

func createRedisStorage(sURL *url.URL) (*Redis, error) {
	p, _ := sURL.User.Password()
	DB, err := strconv.Atoi(sURL.Path[1:])
	if err != nil {
		return nil, err
	}
	host := sURL.Host
	if !strings.Contains(host, ":") {
		host += ":6379"
	}
	var poolSize int
	if _, ok := sURL.Query()["pool_size"]; ok {
		poolSize, err = strconv.Atoi(sURL.Query()["pool_size"][0])
		if err != nil {
			return nil, err
		}
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     host,
		Password: p,  // no password set
		DB:       DB, // use default DB
		PoolSize: poolSize,
	})

	return &Redis{client: rdb, timeout: time.Millisecond * 30}, nil
}
