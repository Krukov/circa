package rules

import (
	"circa/message"
	"circa/storages"
	"time"
)

type CacheRule struct {
	TTL            time.Duration
	Duration       time.Duration
	ResponseStatus int
}

func (r *CacheRule) Process(request *message.Request, key string, storage storages.Storage, call message.Requester) (resp *message.Response, hit bool, err error) {
	resp, err = storage.Get(key)
	if err == nil {
		hit = true
		resp.SetHeader("X-Circa-Cache-Key", key)
		resp.SetHeader("X-Circa-Cache-Storage", storage.String())
		request.Logger = request.Logger.With().Str("cache_key", key).Logger()
		return
	} else {
		err = nil
		start := time.Now()
		resp, err = call(request)
		elapsed := time.Since(start)
		if err == nil {
			if r.Duration != 0 && r.Duration < elapsed {
				request.Logger.Info().Msgf("skip cache in duration condition")
				return
			}
			if r.ResponseStatus != 0 && resp.Status != r.ResponseStatus {
				request.Logger.Info().Msgf("skip cache in status condition")
				return
			}
			storage.Set(key, resp, r.TTL)
		}
	}
	return
}
