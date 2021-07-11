package server

import (
	"circa/message"
	"strconv"
	"time"

	"github.com/valyala/fasthttp"
)

var httpClient = fasthttp.Client{ReadTimeout: time.Second * 5}

func MakeRequest(request *message.Request) (*message.Response, error) {
	start := time.Now()
	logger := request.Logger.With().
		Str("host", request.Host).
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
		request_.Header.Set(header, request.Headers[header])
	}

	if err := httpClient.DoTimeout(request_, response_, request.Timeout); err != nil {
		return nil, err
	}

	data := response_.Body()
	logger.Info().Str("status", strconv.Itoa(response_.StatusCode())).Msg("<<- Response from target")
	headers := map[string]string{}
	response_.Header.VisitAll(func(key, value []byte) {
		headers[string(key)] = string(value)
	})
	headers["X-Circa-Requester-Spend"] = strconv.Itoa(int(time.Since(start).Milliseconds()))
	return &message.Response{Body: data, Status: response_.StatusCode(), Headers: headers}, nil
}
