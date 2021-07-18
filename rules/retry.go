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

func (r *RetryRule) String() string {
	return "retry"
}

func (r *RetryRule) Process(request *message.Request, key string, storage storages.Storage, call message.Requester) (*message.Response, bool, error) {
	resp, err := call(request)
	retry := 0
	hit := false
	for retry < r.Count && err != nil {
		request.Logger.Info().Msgf("== Retrying request %v, err %v", retry, err)
		retry += 1

		hit = true
		resp, err = call(request)
	}
	return resp, hit, err
}
