package message

import (
	"sync"
	"time"

	"github.com/rs/zerolog"
)

type Requester func(*Request) (*Response, error)

type Response struct {
	Status    int
	Body      []byte
	headers   map[string]string
	hmutex    sync.RWMutex
	CachedKey string
}

func NewResponse(status int, body []byte, headers map[string]string) *Response {
	if _, ok := headers["Content-Type"]; !ok {
		headers["Content-Type"] = "application/json; charset=utf8"
	}
	return &Response{Status: status, Body: body, headers: headers, hmutex: sync.RWMutex{}}
}

func (r *Response) AddHandlers(headers map[string]string) {
	r.hmutex.Lock()
	defer r.hmutex.Unlock()
	for h, value := range headers {
		r.headers[h] = value
	}
}

func (r *Response) SetHeader(name string, value string) {
	r.hmutex.Lock()
	defer r.hmutex.Unlock()
	r.headers[name] = value
}

func (r *Response) GetHeaders() map[string]string {
	res := make(map[string]string)
	r.hmutex.RLock()
	defer r.hmutex.RUnlock()
	for k, v := range r.headers {
		res[k] = v
	}
	return res
}

func (r *Response) GetHeader(name string) string {
	r.hmutex.RLock()
	defer r.hmutex.RUnlock()
	return r.headers[name]
}

type Request struct {
	Method  string
	Path    string
	Route   string
	Headers map[string]string
	Query   map[string][]string

	Host string
	Body []byte

	Params map[string]string

	Timeout time.Duration
	Skip    bool

	Logger zerolog.Logger
}
