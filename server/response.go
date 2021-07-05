package server

import (
	"circa/message"
	"fmt"
	"github.com/valyala/fasthttp"
)

func responseNotFound(ctx *fasthttp.RequestCtx) {
	ctx.Response.SetStatusCode(404)
	ctx.SetContentType("application/json; charset=utf8")
	fmt.Fprintf(ctx, "{}\n\n")
}

func responseError(ctx *fasthttp.RequestCtx, err error) {
	ctx.Response.SetStatusCode(500)
	ctx.SetContentType("application/json; charset=utf8")
	fmt.Fprintf(ctx, "{\"error\": \"%v\"}\n\n", err)
}

func responseFor(ctx *fasthttp.RequestCtx, response *message.Response) {
	ctx.Response.SetStatusCode(response.Status)
	for header, value := range response.Headers {
		ctx.Response.Header.Set(header, value)
	}
	fmt.Fprintf(ctx, string(response.Body))
}

