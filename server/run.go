package server

import (
	"fmt"
	"github.com/rs/zerolog/log"

	"circa/message"
	"circa/handler"

	"github.com/valyala/fasthttp"
)

func Run(h *handler.Runner, port string)  {
	requestHandler := func(ctx *fasthttp.RequestCtx) {
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
				request.Logger.Warn().Msg("<- Response not found")
			} else {
				responseError(ctx, err)
				request.Logger.Error().Err(err).Msg("<- Response error")
			}
			return
		}
		request.Logger.Info().Msg("<- Response")
		responseFor(ctx, response)
	}

	// Start HTTP server.
	log.Info().Str("port", port).Msg("Start server")
	if err := fasthttp.ListenAndServe(fmt.Sprintf(":%s", port), requestHandler); err != nil {
		log.Fatal().Err(err).Msg("Can't start proxy")
	}
}

func requestFromHttpRequest(ctx *fasthttp.RequestCtx) *message.Request {
	headers := map[string]string{}
	ctx.Request.Header.VisitAll(func (key, value []byte) {
		headers[string(key)] = string(value)
	})
	return &message.Request{
		Method: string(ctx.Method()),
		Path: string(ctx.Path()),
		Headers: headers,
	}
}