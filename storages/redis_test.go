package storages

import (
	"circa/message"
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
)

func TestRedisStoreSetAndGet(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	s := Redis{client: rdb}
	ctx := context.Background()
	defer rdb.Del(ctx, "key").Err()

	resp := &message.Response{Status: 201, Body: []byte(`data`), Headers: map[string]string{"Request-Id": "test"}}

	_, err := s.Get("key")
	if err == nil {
		t.Error("Value already set")
		return
	}

	setted, err := s.Set("key", resp, time.Second)
	if err != nil {
		t.Errorf("Error on setting value %v", err)
		return
	}
	if !setted {
		t.Error("Value rewrote")
		return
	}

	fromRedis, err := s.Get("key")
	if err != nil {
		t.Errorf("Error at geting resp %v", err)
		return
	}
	if !reflect.DeepEqual(fromRedis, resp) {
		t.Errorf("got = %v, want %v", fromRedis, resp)
	}
}
