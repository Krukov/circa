package rules

import (
	"circa/message"
	"circa/storages"
)

type InvalidateRule struct {
	Methods map[string]bool
}

func (r *InvalidateRule) String() string {
	return "invalidate"
}

func (r *InvalidateRule) Process(request *message.Request, key string, storage storages.Storage, call message.Requester) (resp *message.Response, hit bool, err error) {
	resp, hit, err = simpleCall(request, call)
	if err == nil && resp.Status < 300 && r.Methods[request.Method] {
		storage.Del(key)
	}
	return
}
