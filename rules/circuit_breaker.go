package rules

import (
	"circa/message"
	"circa/storages"
	"time"
)

type CircuitBreaker struct {
	ErrorRate int
	MinCalls  int
	TTL       time.Duration
	OpenTTL   time.Duration
}

func (r *CircuitBreaker) Process(request *message.Request, key string, storage storages.Storage, call message.Requester) (resp *message.Response, hit bool, err error) {
	if ex, _ := storage.Exists(key + ":open"); ex {
		return message.NewResponse(500, []byte(`{"message":"CircuitBreakerOpen"}`), map[string]string{}), true, nil
	}
	total, _ := storage.Incr(key + ":total")
	if total == 1 {
		storage.Expire(key+":total", r.TTL)
		storage.SetRaw(key+":fails", "0", r.TTL)
	}
	resp, err = call(request)
	if err != nil {
		fails, _ := storage.Incr(key + ":fails")
		if total > r.MinCalls && r.ErrorRate <= (fails*100/total) {
			storage.SetRaw(key+":open", "1", r.OpenTTL)
		}
	}
	hit = false
	return
}
