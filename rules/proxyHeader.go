package rules

import (
	"circa/message"
	"circa/storages"
)

type ProxyHeaderRule struct {
	HeaderName string
}

func (r *ProxyHeaderRule) Process(request *message.Request, key string, storage storages.Storage, call message.Requester) (*message.Response, bool, error) {
	value, ok := request.Headers[r.HeaderName]
	resp, hit, err := simpleCall(request, call)
	if err != nil {
		return nil, false, err
	}
	if ok {
		resp.SetHeader(r.HeaderName, value)
	}
	return resp, hit, err
}
