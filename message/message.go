package message

import (
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

type Requester func(*Request) (*Response, error)

type Response struct {
	Status  int
	Body    []byte
	headers map[string]string
	hmutex  sync.RWMutex
}

func NewResponse(status int, body []byte, headers map[string]string) *Response {
	if _, ok := headers["Content-Type"]; !ok {
		headers["Content-Type"] = "application/json; charset=utf8"
	}
	return &Response{Status: status, Body: body, headers: headers, hmutex: sync.RWMutex{}}
}

func (r *Response) SetHeader(name string, value string) {
	r.hmutex.Lock()
	defer r.hmutex.Unlock()
	r.headers[name] = value
}

func (r *Response) GetHeaders(all bool) map[string]string {
	res := map[string]string{}
	r.hmutex.RLock()
	defer r.hmutex.RUnlock()
	for k, v := range r.headers {
		if all || !strings.HasPrefix(k, "X-Circa") {
			res[k] = v
		}
	}
	return res
}

func (r *Response) GetHeader(name string) string {
	r.hmutex.RLock()
	defer r.hmutex.RUnlock()
	return r.headers[name]
}

type Request struct {
	Method   string
	Path     string
	QueryStr string
	Host     string
	FullPath string

	Body []byte

	Route  string
	Params map[string]string

	Query   map[string][]string
	Headers map[string]string

	Timeout time.Duration
	Skip    bool

	Logger zerolog.Logger
}

func (r *Request) GetHeader(name string) string {
	// fast path
	if v, ok := r.Headers[name]; ok {
		return v
	}
	lName := strings.ToLower(name)
	for header, value := range r.Headers {
		if lName == strings.ToLower(header) {
			return value
		}
	}
	return ""
}
