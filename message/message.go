package message

import (
	"github.com/rs/zerolog"
	"time"
)

type Requester func(*Request) (*Response, error)

type Response struct {
	Status    int
	Body      []byte
	Headers   map[string]string
	CachedKey string
}

type Request struct {
	Method  string
	Path    string
	Headers map[string]string

	Host string
	Body []byte

	Params map[string]string
	//ProxyHeaders map[string]string

	Timeout time.Duration

	Logger zerolog.Logger
}
