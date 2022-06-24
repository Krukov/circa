package rules

import (
	"circa/message"
	"circa/storages"
	"time"
)

type RetryRule struct {
	Count   int
	Backoff time.Duration
}

func (r *RetryRule) Process(request *message.Request, key string, storage storages.Storage, call message.Requester) (*message.Response, bool, error) {
	resp, hit, err := simpleCall(request, call)
	retry := 0
	for retry < r.Count && err != nil {
		request.Logger.Info().Msgf("== Retrying request %v, err %v", retry, err)
		retry += 1
		hit = true
		time.Sleep(time.Duration(retry) * r.Backoff)
		resp, err = call(request)
	}
	return resp, hit, err
}
