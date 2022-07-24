package rules

import (
	"circa/message"
	"circa/storages"

	"github.com/google/uuid"
)

type RequestIDRule struct {
	SkipCheckReturn bool
	HeaderName      string
}

func (r *RequestIDRule) Process(request *message.Request, key string, storage storages.Storage, call message.Requester) (*message.Response, bool, error) {
	var requestID string
	if _, ok := request.Headers[r.HeaderName]; ok {
		requestID = request.Headers[r.HeaderName]
	} else {
		requestID = uuid.NewString()
	}
	request.Headers[r.HeaderName] = requestID
	request.Logger = request.Logger.With().Str("request_id", requestID).Logger()
	resp, hit, err := simpleCall(request, call)
	if err != nil {
		return nil, false, err
	}
	if !r.SkipCheckReturn && requestID != resp.GetHeader(r.HeaderName) {
		request.Logger.Warn().Msgf("Request id of response doesn't match request value %v != %v", resp.GetHeader(r.HeaderName), requestID)
	}
	resp.SetHeader(r.HeaderName, requestID)
	return resp, hit, err
}
