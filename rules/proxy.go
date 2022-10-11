package rules

import (
	"circa/message"
	"circa/storages"
	"strings"
)

type ProxyRule struct {
	Target string
	Method string
	Path   string
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
		fp := strings.Replace(request.FullPath, request.Path, r.Path, 1)
		request.Path = r.Path
		request.FullPath = fp
	}
	return simpleCall(request, call)
}
