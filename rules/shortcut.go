package rules

import "circa/message"

func simpleCall(request *message.Request, call message.Requester) (*message.Response, bool, error) {
	resp, err := call(request)
	return resp, false, err
}
