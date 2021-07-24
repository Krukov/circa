package rules

import (
	"circa/message"
	"circa/storages"
)

type ProxyRule struct {
	Target string
}

func (r *ProxyRule) String() string {
	return "proxy"
}

func (r *ProxyRule) Process(request *message.Request, key string, storage storages.Storage, call message.Requester) (*message.Response, bool, error) {
	if r.Target != "" {
		request.Host = r.Target
	}
	return simpleCall(request, call)
}
