package rules

import (
	"circa/message"
	"circa/storages"
)


type ProxyRule struct {}

func (r *ProxyRule) String() string {
	return "proxy"
}

func (r *ProxyRule) Process(request *message.Request, key string, storage storages.Storage, call message.Requester) (*message.Response, error) {
	return call(request)
}


