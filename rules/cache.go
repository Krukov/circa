package rules

import (
	"circa/message"
	"circa/storages"
	"time"
)

type CacheRule struct {
	TTL time.Duration
}

func (r *CacheRule) String() string {
	return "cache"
}

func (r *CacheRule) Process(request *message.Request, key string, storage storages.Storage, call message.Requester) (resp *message.Response, hit bool, err error) {
	resp, err = storage.Get(key)
	if err == nil {
		resp.CachedKey = key
		hit = true
		return
	} else {
		err = nil
		resp, err = call(request)
		if err == nil {
			storage.Set(key, resp, r.TTL)
		}
	}
	return
}
