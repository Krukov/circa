package rules

import (
	"circa/message"
	"circa/storages"
)

type RuleProcessor interface {
	Process(request *message.Request, key string, storage storages.Storage, call message.Requester) (*message.Response, bool, error)
}

type Rule struct {
	Name      string  // like kind
	Key       string  // template
	Route     string 
	Methods   map[string]bool
	
	Storage   storages.Storage
	Processor RuleProcessor
}


