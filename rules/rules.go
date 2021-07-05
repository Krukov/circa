package rules

import (
	"circa/message"
	"circa/storages"
)

type Rule interface {
	String() string
	Process(request *message.Request, key string, storage storages.Storage, call message.Requester) (*message.Response, error)
}

