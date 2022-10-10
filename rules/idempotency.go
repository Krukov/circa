package rules

import (
	"circa/message"
	"circa/storages"
	"time"
)

type IdempotencyRule struct {
	TTL time.Duration
}

func (r *IdempotencyRule) Process(request *message.Request, key string, storage storages.Storage, call message.Requester) (*message.Response, bool, error) {
	if set, err := storage.Set(key, &message.Response{}, r.TTL); !set && err == nil {
		request.Skip = true
		return message.NewResponse(409, []byte(`{"message":"dublicated"}`), map[string]string{}), true, nil
	}
	return simpleCall(request, call)
}
