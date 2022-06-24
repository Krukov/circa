package rules

import (
	"circa/message"
	"circa/storages"
)

type SkipRule struct{}

func (r *SkipRule) Process(request *message.Request, key string, storage storages.Storage, call message.Requester) (*message.Response, bool, error) {
	request.Skip = true
	return simpleCall(request, call)
}
