package fastreq

import "fmt"

type Middleware func(ctx *Ctx) error

func MiddlewareOauth1(o *Oauth1) Middleware {
	return func(ctx *Ctx) error {
		auth := o.GenHeader(ctx.Request)
		ctx.Request.Header.SetBytesV("Authorization", auth)
		return ctx.Next()
	}
}

func MiddlewareLog() Middleware {
	return func(ctx *Ctx) error {
		fmt.Printf("%s\n", ctx.Request.URI().FullURI())

		if err := ctx.Next(); err != nil {
			return err
		}

		fmt.Printf("%d\n", ctx.Response.StatusCode())

		return nil
	}
}
