package fastreq

import (
	"context"
	"sync"

	"github.com/valyala/fasthttp"
)

// Ctx represents the Context which hold the HTTP request and response.
type Ctx struct {
	Request         *Request
	Response        *Response
	ctx             context.Context
	client          *Client
	indexMiddleware int
}

// Next ..
func (c *Ctx) Next() (err error) {
	// Increment handler index
	c.indexMiddleware++
	// Did we execute all route handlers?
	if c.indexMiddleware < len(c.client.middlewares) {
		// Continue route stack
		return c.client.middlewares[c.indexMiddleware](c)
	} else {
		// Continue handler stack
		return do(c)
	}
}

func (c *Ctx) Release() {
	c.Request.Release()
	c.Response.Release()
	c.client = nil
	c.indexMiddleware = 0

	ctxPool.Put(c)
}

func (c *Ctx) fastClient() *fasthttp.Client {
	return c.client.Client
}

func (c *Ctx) fastRequest() *fasthttp.Request {
	return c.Request.Request
}

func (c *Ctx) fastResponse() *fasthttp.Response {
	return c.Response.Response
}

var ctxPool sync.Pool

func AcquireCtx() *Ctx {
	v := ctxPool.Get()
	if v == nil {
		return &Ctx{}
	}
	return v.(*Ctx)
}
