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

func (r *FailRule) Process(request *message.Request, key string, storage storages.Storage, call message.Requester) (resp *message.Response, hit bool, err error) {
	var errStorage error
	resp, err = call(request)
	if err != nil {
		request.Logger.Debug().Msg("error on call api, try to find in cache")

		resp, errStorage = storage.Get(key)
		if errStorage == nil {
			resp.CachedKey = key
			hit = true
			err = nil
			return
		} else {
			if errStorage != storages.NotFound {
				request.Logger.Warn().Msgf("error on get value %v", err)
			}
		}
	}
	if resp != nil && resp.CachedKey == "" {
		_, setErr := storage.Set(key, resp, r.TTL)
		if setErr != nil {
			request.Logger.Warn().Msgf("error on set value %v", setErr)
		}
	}
	return
}


