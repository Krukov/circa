package rules

import (
	"circa/message"
	"circa/storages"
	"fmt"
)

type Rule interface {
	fmt.Stringer
	Process(request *message.Request, key string, storage storages.Storage, call message.Requester) (*message.Response, bool, error)
}
