package fastreq

type Middleware func(ctx *Ctx) error

// MiddlewareOauth1 generates a middleware function that adds an OAuth1
// authorization header to incoming requests. The middleware uses the given
// Oauth1 object o to generate the header.
func MiddlewareOauth1(o *Oauth1) Middleware {
	return func(ctx *Ctx) error {
		auth := o.GenHeader(ctx.Request)
		ctx.Request.Header.SetBytesV("Authorization", auth)
		return ctx.Next()
	}
}
