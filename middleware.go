package fastreq

import "fmt"

type Middleware func(ctx *Ctx) error

func MiddlewareOauth1(o *Oauth1) Middleware {
	return func(ctx *Ctx) error {
		auth := o.GenHeader(ctx.Request)
		ctx.Request.Header.SetBytesV("Authorization", auth)
		return nil
	}
}

func MiddlewareLog() Middleware {
	return func(ctx *Ctx) error {
		fmt.Printf("%s\n", req.req.URI().FullURI())
		fmt.Printf("%s\n", req.req.Header.RawHeaders())
	}
}
