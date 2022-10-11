package rules

import (
	"circa/message"
	"circa/storages"
	"time"
)

type RuleProcessor interface {
	Process(request *message.Request, key string, storage storages.Storage, call message.Requester) (*message.Response, bool, error)
}

type Condition struct {
	Status int

	SkipHeaderValue   string
	ShouldHeaderValue string
	Header            string

	Duration time.Duration
}

type Rule struct {
	Name        string // like kind
	Key         string // template
	Route       string
	Methods     map[string]bool
	StorageName string //

	Condition Condition

	Storage   storages.Storage
	Processor RuleProcessor
}

func (r *Rule) Process(request *message.Request, key string, storage storages.Storage, call message.Requester) (*message.Response, bool, error) {
	if r.skipCondition(request) {
		request.Logger.Debug().Msgf("skip by condition %s", r.Name)
		resp, err := call(request)
		return resp, false, err
	}
	return r.Processor.Process(request, key, storage, call)
}

func (r *Rule) skipCondition(request *message.Request) bool {
	if r.Condition.ShouldHeaderValue != "" && request.Headers[r.Condition.Header] != r.Condition.ShouldHeaderValue {
		return true
	}
	if r.Condition.SkipHeaderValue != "" && request.Headers[r.Condition.Header] == r.Condition.SkipHeaderValue {
		return true
	}
	return false
}
