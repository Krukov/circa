package rules

import (
	"circa/message"
	"circa/storages"
	"time"
)

type RateLimitRule struct {
	TTL time.Duration
	Limit int
}

func (r *RateLimitRule) String() string {
	return "rate-limit"
}

func (r *RateLimitRule) Process(request *message.Request, key string, storage storages.Storage, call message.Requester) (*message.Response, bool, error) {
	count, _ := storage.Incr(key)
	if count >= r.Limit {
		return &message.Response{Status: 429, Body: []byte(`rate limited`)}, true, nil
	}
	return simpleCall(request, call)
}
