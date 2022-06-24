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

func (r *CircuitBreaker) Process(request *message.Request, key string, storage storages.Storage, call message.Requester) (resp *message.Response, hit bool, err error) {
	// Check if circuit_breaker is open
	// 2 ** 16  - seze for one intager
	// 1440  intager for a day
	// 1440 * 16 = 23040
	// If close
	//      resp, err = call(request)
	//      if err:
	//           count error
	//           if errors more then ErrorRate
	//                open circuit_breaker
	//      else:
	//            return resp, false, err
	//
	//
	resp, err = call(request)
	hit = false
	return
}
