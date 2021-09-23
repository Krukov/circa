package rules

import (
	"circa/message"
	"circa/storages"
	"time"
)

type CircuitBreaker struct {
	ErrorRate int
	MinCalls  int
	TTL       time.Duration
	HalfOpen  time.Duration
}

func (r *CircuitBreaker) String() string {
	return "circuit_breaker"
}

func (r *CircuitBreaker) Process(request *message.Request, key string, storage storages.Storage, call message.Requester) (resp *message.Response, hit bool, err error) {
	// Check if circuit_breaker is open
	// If close
	//      err = call (request)
	//       if err
	//
	resp, err = call(request)
	hit = false
	return
}
