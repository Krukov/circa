package storages

import (
	"context"
	"github.com/go-redis/redis/v8"
	"strconv"
	"time"

	"circa/message"
)

type Redis struct {
	client *redis.Client
	timeout time.Duration
}

func (s *Redis) String() string {
	return "redis"
}

func (s *Redis) Set(key string, value *message.Response, ttl time.Duration) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()
	values := map[string]string{}
	for header, hValue := range value.Headers {
		// todo: do not store all headers ( only
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
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()
	deleted, err := s.client.Del(ctx, "key").Result()
	if err != nil {
		return false, err
	}
	return deleted > 0, nil
}

func (s *Redis) Get(key string) (*message.Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()
	keys, err := s.client.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	if len(key) == 0 {
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
	return &message.Response{Status: statusInt, Body: []byte(body), Headers: keys}, nil
}

