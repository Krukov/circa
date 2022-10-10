package rules

import (
	"circa/message"
	"circa/storages"
)

type StaticRule struct {
	Response string
}

func (r *StaticRule) String() string {
	return "static"
}

func (r *StaticRule) Process(request *message.Request, key string, storage storages.Storage, call message.Requester) (*message.Response, bool, error) {
	return message.NewResponse(200, []byte(r.Response), map[string]string{}), true, nil
}
