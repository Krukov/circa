package rules

import (
	"circa/message"
	"circa/storages"
	"strconv"
	"time"
)

var earlyExpiredHeader = "_early_expired"

type EarlyCacheRule struct {
	TTL      time.Duration
	EarlyTTL time.Duration
}

func (r *EarlyCacheRule) Process(request *message.Request, key string, storage storages.Storage, call message.Requester) (resp *message.Response, hit bool, err error) {
	resp, err = storage.Get(key)
	if err == nil {
		early, errInt := strconv.ParseInt(resp.GetHeader(earlyExpiredHeader), 10, 64)
		if errInt == nil && early < time.Now().UTC().Unix() {
			incr, errIncr := storage.Incr(key + ":lock")
			if errIncr == nil && incr == 1 {
				go r.callAndSet(request, key, storage, call)
			}
		}
		resp.SetHeader("X-Circa-Cache-Key", key)
		resp.SetHeader("X-Circa-Cache-Storage", storage.String())
		request.Logger = request.Logger.With().Str("cache_key", key).Logger()
		hit = true
		return
	} else {
		resp, err = r.callAndSet(request, key, storage, call)
	}
	return
}

func (r *EarlyCacheRule) callAndSet(request *message.Request, key string, storage storages.Storage, call message.Requester) (*message.Response, error) {
	resp, err := call(request)
	if err == nil {
		early := strconv.FormatInt(time.Now().UTC().Add(r.EarlyTTL).Unix(), 10)
		resp.SetHeader(earlyExpiredHeader, early)
		storage.Set(key, resp, r.TTL)
		storage.Del(key + ":lock")
	}
	return resp, err
}
