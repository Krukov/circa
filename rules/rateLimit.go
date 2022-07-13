package rules

import (
	"circa/message"
	"circa/storages"
	"time"
)

type RateLimitRule struct {
	TTL   time.Duration
	Limit int
}

func (r *RateLimitRule) Process(request *message.Request, key string, storage storages.Storage, call message.Requester) (*message.Response, bool, error) {
	count, _ := storage.Incr(key)
	if count == 1 {
		storage.Expire(key, r.TTL)
	}
	if count >= r.Limit {
		headers := map[string]string{
			"X-Circa-Rate-Key": key,
		}
		return message.NewResponse(429, []byte(`{"message":"rate limit"}`), headers), true, nil
	}
	return simpleCall(request, call)
}
