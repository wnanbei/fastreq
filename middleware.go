package fastreq

type Middleware func(ctx *Ctx) error

func MiddlewareOauth1(o *Oauth1) Middleware {
	return func(ctx *Ctx) error {
		auth := o.GenHeader(ctx.Request)
		ctx.Request.Header.SetBytesV("Authorization", auth)
		return ctx.Next()
	}
}
