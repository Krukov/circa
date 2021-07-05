package rules

import (
	"circa/message"
	"circa/storages"
	"time"
)

type FailRule struct {
	TTL      time.Duration
}

func (r *FailRule) String() string {
	return "fail"
}

func (r *FailRule) Process(request *message.Request, key string, storage storages.Storage, call message.Requester) (resp *message.Response, err error) {
	resp, err = call(request)
	if err != nil || resp.Status >= 500 {
		request.Logger.Debug().Msg("error on call api, try to find in cache")

		resp, err = storage.Get(key)
		if err == nil {
			resp.CachedKey = key
			return
		} else {
			if err != storages.NotFound {
				request.Logger.Warn().Msgf("error on get value %v", err)
			}
		}
		_, setErr := storage.Set(key, resp, r.TTL)
		if setErr != nil {
			request.Logger.Warn().Msgf("error on set value %v", err)
		}
	}
	return resp, err
}


