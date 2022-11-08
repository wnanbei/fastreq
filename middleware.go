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
		fmt.Printf(
			"REQUEST: %s %s %s\n%s",
			time.Now().Format(time.RFC3339),
			ctx.Request.Header.Method(),
			ctx.Request.URI().FullURI(),
			ctx.Request.String(),
		)

		if err := ctx.Next(); err != nil {
			return err
		}

		fmt.Printf(
			"RESPONSE: %s %s %s %d\n%s",
			time.Now().Format(time.RFC3339),
			ctx.Request.Header.Method(),
			ctx.Request.URI().FullURI(),
			ctx.Response.StatusCode(),
			ctx.Response.String(),
		)

		return nil
	}
}
