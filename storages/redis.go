package storages

import (
	"context"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"

	"circa/message"
)

type Redis struct {
	client *redis.Client
}

func NewRedisStorageFormURL(sURL *url.URL) (*Redis, error) {
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

	timeout := "30ms"
	if _, ok := sURL.Query()["timeout"]; ok {
		timeout = sURL.Query()["timeout"][0]
	}
	readTimeout, err := time.ParseDuration(timeout)
	if err != nil {
		return nil, err
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:        host,
		Password:    p,  // no password set
		DB:          DB, // use default DB
		PoolSize:    poolSize,
		ReadTimeout: readTimeout,
	})

	return &Redis{client: rdb}, err
}

func (s *Redis) String() string {
	return "redis"
}

func (s *Redis) Set(key string, value *message.Response, ttl time.Duration) (bool, error) {
	ctx := context.Background()
	start := time.Now()
	defer func() {
		operationHistogram.WithLabelValues(s.String(), "set").Observe(time.Since(start).Seconds())
	}()
	values := map[string]string{}
	for header, hValue := range value.GetHeaders(false) {
		values[header] = hValue
	}
	values["body"] = string(value.Body)
	values["status"] = strconv.Itoa(value.Status)
	addedCount, err := s.client.HSet(ctx, key, values).Result()
	if err != nil {
		return false, err
	}
	added := int(addedCount) == len(values)
	if added {
		s.client.Expire(ctx, key, ttl)
	}
	return added, nil
}

func (s *Redis) Del(key string) (bool, error) {
	ctx := context.Background()
	start := time.Now()
	defer func() {
		operationHistogram.WithLabelValues(s.String(), "del").Observe(time.Since(start).Seconds())
	}()
	deleted, err := s.client.Del(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return deleted > 0, nil
}

func (s *Redis) Get(key string) (*message.Response, error) {
	start := time.Now()
	ctx := context.Background()
	defer func() {
		operationHistogram.WithLabelValues(s.String(), "get").Observe(time.Since(start).Seconds())
	}()
	keys, err := s.client.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	if len(keys) == 0 {
		return nil, NotFound
	}
	body := keys["body"]
	delete(keys, "body")
	status := keys["status"]

	delete(keys, "status")
	statusInt, err := strconv.Atoi(status)
	if err != nil {
		s.client.Del(ctx, key)
		return nil, NotFound
	}
	return message.NewResponse(statusInt, []byte(body), keys), nil
}

func (s *Redis) Incr(key string) (int, error) {
	start := time.Now()
	defer func() {
		operationHistogram.WithLabelValues(s.String(), "incr").Observe(time.Since(start).Seconds())
	}()
	count, err := s.client.Incr(context.Background(), key).Result()
	return int(count), err
}

func (s *Redis) Expire(key string, ttl time.Duration) error {
	start := time.Now()
	defer func() {
		operationHistogram.WithLabelValues(s.String(), "expire").Observe(time.Since(start).Seconds())
	}()
	_, err := s.client.Expire(context.Background(), key, ttl).Result()
	return err
}
