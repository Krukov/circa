package rules

import (
	"circa/message"
	"circa/storages"
	"time"
)

type HitRule struct {
	TTL             time.Duration
	Hits            int
	UpdateAfterHits int
}

func (r *HitRule) String() string {
	return "hit"
}

func (r *HitRule) Process(request *message.Request, key string, storage storages.Storage, call message.Requester) (resp *message.Response, hit bool, err error) {
	hits, err := storage.Incr(key + ":hits")
	if err != nil {
		// Error at storage connection or so on : just proxy request
		return simpleCall(request, call)
	}
	if hits == 1 {
		// No cache record - call and save result with ttl
		return callAndSet(request, key, storage, call, r.TTL)
	}

	if hits > r.Hits {
		_, _ = storage.Del(key + ":hits")
		return callAndSet(request, key, storage, call, r.TTL)
	} else {
		if hits > r.UpdateAfterHits {
			go func() {
				_, _ = storage.Del(key + ":hits")
				callAndSet(request, key, storage, call, r.TTL)
			}()
		}
		resp, err = storage.Get(key)
		if err == nil {
			return resp, true, err
		}
	}
	return simpleCall(request, call)
}

func callAndSet(request *message.Request, key string, storage storages.Storage, call message.Requester, ttl time.Duration) (*message.Response, bool, error) {
	resp, _, err := simpleCall(request, call)
	if err != nil {
		return nil, false, err
	}
	_, setErr := storage.Set(key, resp, ttl)
	if setErr != nil {
		request.Logger.Warn().Msgf("error on set value %v", setErr)
	}
	return resp, false, err
}
