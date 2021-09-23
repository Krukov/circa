package server

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"

	"circa/handler"
	"circa/message"

	"github.com/valyala/fasthttp"
)

func Run(cancel context.CancelFunc, h *handler.Runner, port string) *fasthttp.Server {
	requestHandler := func(ctx *fasthttp.RequestCtx) {
		start := time.Now()
		request := requestFromHttpRequest(ctx)

		logger := log.With().
			Str("method", request.Method).
			Str("path", request.Path).
			Logger()
		request.Logger = logger
		request.Logger.Info().Msg("-> Request")
		response, err := h.Handle(request)
		if err != nil {
			if err == handler.NotFound {
				responseNotFound(ctx)
				request.Logger.Warn().Msg("<- Route not found")
			} else {
				responseError(ctx, err)
				request.Logger.Error().Err(err).Msg("<- Response error")
			}
			requestsLatency.WithLabelValues(request.Method, request.Route, "error").Observe(time.Since(start).Seconds())
			return
		}
		request.Logger.Info().Msg("<- Response")
		requestsLatency.WithLabelValues(request.Method, request.Route, strconv.Itoa(response.Status)).Observe(time.Since(start).Seconds())
		// response.SetHeader("X-Circa-Proxy-Spend", strconv.Itoa(int(time.Since(start).Milliseconds())))
		responseFor(ctx, response)
	}

	// Start HTTP server.
	log.Info().Str("port", port).Msg("Start server")
	srv := fasthttp.Server{Handler: requestHandler, Name: "circa", ReadTimeout: time.Second * 10}
	go func() {
		defer cancel()
		if err := srv.ListenAndServe(fmt.Sprintf(":%s", port)); err != nil {
			log.Fatal().Err(err).Msg("Can't start proxy")
		}
		log.Info().Msg("Shutdown server")
	}()
	return &srv
}

func requestFromHttpRequest(ctx *fasthttp.RequestCtx) *message.Request {
	headers := map[string]string{}
	ctx.Request.Header.VisitAll(func(key, value []byte) {
		headers[string(key)] = string(value)
	})
	ctx.QueryArgs()
	return &message.Request{
		Method:  string(ctx.Method()),
		Path:    string(ctx.Path()),
		Headers: headers,
	}
}
