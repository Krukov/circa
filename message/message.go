package message

import (
	"time"

	"github.com/rs/zerolog"
)

type Requester func(*Request) (*Response, error)

type Response struct {
	Status    int
	Body      []byte
	Headers   map[string]string
	CachedKey string
}

func NewResponse(status int, body []byte, headers map[string]string) *Response {
	headers["Content-Type"] = "application/json; charset=utf8"
	return &Response{Status: status, Body: body, Headers: headers}
}

type Request struct {
	Method  string
	Path    string
	Route   string
	Headers map[string]string

	Host string
	Body []byte

	Params map[string]string
	//ProxyHeaders map[string]string

	Timeout time.Duration
	Skip    bool

	Logger zerolog.Logger
}
