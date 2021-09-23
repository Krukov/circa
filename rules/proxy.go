package rules

import (
	"circa/message"
	"circa/storages"
)

type ProxyRule struct {
	Target string
	Method string
	Path   string
}

func (r *ProxyRule) String() string {
	return "proxy"
}

func (r *ProxyRule) Process(request *message.Request, key string, storage storages.Storage, call message.Requester) (*message.Response, bool, error) {
	if r.Target != "" {
		request.Host = r.Target
	}
	if r.Method != "" {
		request.Method = r.Method
		if r.Method == "GET" {
			request.Body = nil
		}
	}
	if r.Path != "" {
		request.Path = r.Path
	}
	return simpleCall(request, call)
}
