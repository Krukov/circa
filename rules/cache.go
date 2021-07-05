package rules

import (
	"circa/message"
	"circa/storages"
	"time"
)

type CacheRule struct {
	TTL  time.Duration
}

func (r *CacheRule) String() string {
	return "cache"
}


func (r *CacheRule) Process(request *message.Request, key string, storage storages.Storage, call message.Requester) (resp *message.Response, err error) {
	resp, err = storage.Get(key)
	if err == nil {
		resp.CachedKey = key
		return
	} else {
		if err != storages.NotFound {
			request.Logger.Warn().Msgf("error on get value %v", err)
		}
		resp, err = call(request)
		if err == nil {
			_, setErr := storage.Set(key, resp, r.TTL)
			if setErr != nil {
				request.Logger.Warn().Msgf("error on set value %v", err)
			}
		}
	}
	return resp, err
}


