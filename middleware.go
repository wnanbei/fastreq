package fastreq

import (
	"fmt"
	"time"
)

type Middleware func(ctx *Ctx) error

func MiddlewareOauth1(o *Oauth1) Middleware {
	return func(ctx *Ctx) error {
		auth := o.GenHeader(ctx.Request)
		ctx.Request.Header.SetBytesV("Authorization", auth)
		return ctx.Next()
	}
}

func MiddlewareLogger() Middleware {
	return func(ctx *Ctx) error {
		start := time.Now()
		fmt.Printf(
			"REQUEST[%s]: %s %s\n%s",
			start.Format(time.RFC3339),
			ctx.Request.Header.Method(),
			ctx.Request.URI().FullURI(),
			ctx.Request.String(),
		)

		if err := ctx.Next(); err != nil {
			return err
		}

		end := time.Now()
		fmt.Printf(
			"RESPONSE[%s]: %d %dms %s %s\n%s",
			end.Format(time.RFC3339),
			ctx.Response.StatusCode(),
			end.Sub(start).Milliseconds(),
			ctx.Request.Header.Method(),
			ctx.Request.URI().FullURI(),
			ctx.Response.String(),
		)

		return nil
	}
}
