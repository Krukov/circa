package rules

import (
	"circa/message"
	"circa/storages"
	"time"
)

type FailRule struct {
	TTL time.Duration
}

func (r *FailRule) Process(request *message.Request, key string, storage storages.Storage, call message.Requester) (resp *message.Response, hit bool, err error) {
	var errStorage error
	resp, err = call(request)
	if err != nil {
		request.Logger.Debug().Msg("error on call api, try ")

		resp, errStorage = storage.Get(key)
		if errStorage == nil {
			resp.SetHeader("X-Circa-Cache-Key", key)
			resp.SetHeader("X-Circa-Cache-Storage", storage.String())
			request.Logger = request.Logger.With().Str("cache_key", key).Logger()
			hit = true
			err = nil
			return
		}
	}
	if resp != nil && !request.Skip {
		storage.Set(key, resp, r.TTL)
	}
	return
}
