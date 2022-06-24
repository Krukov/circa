package rules

import (
	"circa/message"
	"circa/storages"
	"strconv"
	"time"
)

type HitRule struct {
	TTL             time.Duration
	Hits            int
	UpdateAfterHits int
}

func (r *HitRule) Process(request *message.Request, key string, storage storages.Storage, call message.Requester) (resp *message.Response, hit bool, err error) {
	hits, err := storage.Incr(key + ":hits")
	if err != nil {
		// Error at storage connection or so on : just proxy request
		return simpleCall(request, call)
	}
	hits -= 1

	if hits > r.Hits {
		return updateCache(request, key, storage, call, r.TTL)
	} else {
		if hits == r.UpdateAfterHits && hits != 0 {
			go func() {
				updateCache(request, key, storage, call, r.TTL)
			}()
		}
		resp, err = storage.Get(key)
		if err == nil {
			resp.SetHeader("X-Circa-Cache-Key", key)
			resp.SetHeader("X-Circa-Cache-Storage", storage.String())
			resp.SetHeader("X-Circa-Hits-To-Update", strconv.Itoa(r.Hits-hits))
			request.Logger = request.Logger.With().Str("cache_key", key).Logger()
			return resp, true, err
		}
		return updateCache(request, key, storage, call, r.TTL)
	}
}

func updateCache(request *message.Request, key string, storage storages.Storage, call message.Requester, ttl time.Duration) (*message.Response, bool, error) {
	_, _ = storage.Del(key + ":hits")
	return callAndSet(request, key, storage, call, ttl)
}

func callAndSet(request *message.Request, key string, storage storages.Storage, call message.Requester, ttl time.Duration) (*message.Response, bool, error) {
	resp, _, err := simpleCall(request, call)
	if err != nil {
		return nil, false, err
	}
	storage.Set(key, resp, ttl)
	storage.Incr(key + ":hits")
	storage.Expire(key+":hits", ttl)
	return resp, false, err
}
