package server

import (
	"circa/message"
	"strconv"
	"time"

	"github.com/valyala/fasthttp"
)

var headersForbiddenToProxy = map[string]bool{"Connection": true, "Content-Length": true, "Keep-Alive": true}
var headersForbiddenToPass = map[string]bool{"Accept-Encoding": true, "Connection": true}

// var httpClient = fasthttp.Client{ReadTimeout: time.Second * 5, MaxConnsPerHost: 10} make it with closure

func MakeRequest(request *message.Request) (*message.Response, error) {
	start := time.Now()
	logger := request.Logger.With().
		Str("host", request.Host).
		Str("path", request.Path).
		Str("timeout", request.Timeout.String()).
		Logger()
	logger.Info().Msg("->> Forward request")

	request_ := fasthttp.AcquireRequest()
	response_ := fasthttp.AcquireResponse()

	defer func() {
		proxyLatency.WithLabelValues(request.Host, request.Method, request.Route, strconv.Itoa(response_.StatusCode())).Observe(time.Since(start).Seconds())
		fasthttp.ReleaseResponse(response_)
		fasthttp.ReleaseRequest(request_)
	}()

	request_.Header.SetMethod(request.Method)
	request_.SetRequestURI(request.Host + request.Path)
	for header := range request.Headers {
		if !headersForbiddenToPass[header] {
			request_.Header.Set(header, request.Headers[header])
		}
	}

	if err := fasthttp.DoTimeout(request_, response_, request.Timeout); err != nil {
		return nil, err
	}

	data := make([]byte, len(response_.Body()))
	copy(data, response_.Body())

	logger.Info().Str("status", strconv.Itoa(response_.StatusCode())).Msg("<<- Response from target")
	headers := map[string]string{}
	response_.Header.VisitAll(func(key, value []byte) {
		if !headersForbiddenToProxy[string(key)] {
			headers[string(key)] = string(value)
		}
	})
	// headers["X-Circa-Requester-Spend"] = strconv.Itoa(int(time.Since(start).Milliseconds()))
	return message.NewResponse(response_.StatusCode(), data, headers), nil
}
