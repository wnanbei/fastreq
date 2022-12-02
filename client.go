package fastreq

import (
	"crypto/tls"
	"fmt"
	"time"

	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpproxy"
)

// Client ...
type Client struct {
	*fasthttp.Client
	userAgent         string
	maxRedirectsCount int
	timeout           time.Duration
	debugLevel        DebugLevel
	auth              Oauth1
	middlewares       []Middleware
}

// NewClient ...
func NewClient() *Client {
	return &Client{
		Client:      &fasthttp.Client{},
		middlewares: []Middleware{},
	}
}

// NewClientFromFastHTTP ...
func NewClientFromFastHTTP(client *fasthttp.Client) *Client {
	return &Client{
		Client: client,
	}
}

func (c *Client) Get(url string, opts ...ReqOption) (*Ctx, error) {
	req := NewRequest(GET, url)
	return c.Do(req, opts...)
}

func (c *Client) Head(url string, opts ...ReqOption) (*Ctx, error) {
	req := NewRequest(HEAD, url)
	return c.Do(req, opts...)
}

func (c *Client) Post(url string, opts ...ReqOption) (*Ctx, error) {
	req := NewRequest(POST, url)
	return c.Do(req, opts...)
}

func (c *Client) Put(url string, opts ...ReqOption) (*Ctx, error) {
	req := NewRequest(PUT, url)
	return c.Do(req, opts...)
}

func (c *Client) Patch(url string, opts ...ReqOption) (*Ctx, error) {
	req := NewRequest(PATCH, url)
	return c.Do(req, opts...)
}

func (c *Client) Delete(url string, opts ...ReqOption) (*Ctx, error) {
	req := NewRequest(DELETE, url)
	return c.Do(req, opts...)
}

func (c *Client) Connect(url string, opts ...ReqOption) (*Ctx, error) {
	req := NewRequest(CONNECT, url)
	return c.Do(req, opts...)
}

func (c *Client) DownloadFile(req *Request, path, filename string) error {
	ctx, err := c.do(req)
	if err != nil {
		return err
	}

	if err := ctx.Response.SaveToFile(path, filename); err != nil {
		return err
	}

	ctx.Release()

	return nil
}

func (c *Client) Do(req *Request, opts ...ReqOption) (*Ctx, error) {
	for _, opt := range opts {
		opt.BindRequest(req)
	}

	return c.do(req)
}

func (c *Client) do(req *Request) (*Ctx, error) {
	ctx := AcquireCtx()
	ctx.Request = req
	ctx.client = c

	if len(c.middlewares) > 0 {
		if err := c.middlewares[0](ctx); err != nil {
			return nil, err
		}
	} else {
		if err := do(ctx); err != nil {
			return nil, err
		}
	}

	return ctx, nil
}

func do(ctx *Ctx) error {
	if ctx.Request.mw != nil {
		ctx.Request.Header.SetMultipartFormBoundary(ctx.Request.mw.Boundary())
		if err := ctx.Request.mw.Close(); err != nil {
			return err
		}
	}

	start := time.Now()
	debugBeforeRequest(ctx, start)

	resp := fasthttp.AcquireResponse()
	if err := ctx.client.DoTimeout(ctx.fastRequest(), resp, ctx.client.timeout); err != nil {
		fasthttp.ReleaseResponse(resp)
		return err
	}
	ctx.Response = NewResponseFromFastHTTP(resp)

	debugAfterRequest(ctx, start)

	return nil
}

func (c *Client) SetHTTPProxy(proxy string) {
	c.Dial = fasthttpproxy.FasthttpHTTPDialer(proxy)
}

func (c *Client) SetSocks5Proxy(proxy string) {
	c.Dial = fasthttpproxy.FasthttpSocksDialer(proxy)
}

func (c *Client) SetEnvHTTPProxy() {
	c.Dial = fasthttpproxy.FasthttpProxyHTTPDialer()
}

func (c *Client) SetTimeout(timeout time.Duration) {
	c.timeout = timeout
}

func (c *Client) SetUserAgent(userAgent string) {
	c.userAgent = userAgent
}

func (c *Client) SetDebugLevel(lvl DebugLevel) {
	c.debugLevel = lvl
}

func (c *Client) SetTLSConfig(config *tls.Config) {
	c.TLSConfig = config
}

func (c *Client) SetMaxRedirectsCount(count int) {
	c.maxRedirectsCount = count
}

func (c *Client) SetRetryIf(retryIf fasthttp.RetryIfFunc) {
	c.RetryIf = retryIf
}

func (c *Client) SkipInsecureVerify(isSkip bool) {
	if c.TLSConfig == nil {
		/* #nosec G402 */
		c.TLSConfig = &tls.Config{InsecureSkipVerify: isSkip} // #nosec G402
	} else {
		/* #nosec G402 */
		c.TLSConfig.InsecureSkipVerify = isSkip
	}
}

func (c *Client) SetOauth1(o *Oauth1) {
	c.middlewares = append(c.middlewares, MiddlewareOauth1(o))
}

func (c *Client) AddMiddleware(middlewares ...Middleware) {
	c.middlewares = append(c.middlewares, middlewares...)
}

func debugBeforeRequest(ctx *Ctx, start time.Time) {
	switch ctx.client.debugLevel {
	case DebugClose:
		return
	case DebugSimple:
		fmt.Printf(
			"REQUEST[%s]: %s %s\n",
			start.Format("2006-01-02 15:04:05"),
			ctx.Request.Header.Method(),
			ctx.Request.URI().FullURI(),
		)
	case DebugDetail:
		body := ctx.Request.String()
		if len(body) > debugLimit {
			body = body[:debugLimit] + "\n"
		}
		fmt.Printf(
			"REQUEST[%s]: %s %s\n%s",
			start.Format("2006-01-02 15:04:05"),
			ctx.Request.Header.Method(),
			ctx.Request.URI().FullURI(),
			body,
		)
	}
}

func debugAfterRequest(ctx *Ctx, start time.Time) {
	switch ctx.client.debugLevel {
	case DebugClose:
		return
	case DebugSimple:
		end := time.Now()
		fmt.Printf(
			"RESPONSE[%s]: %d %dms %s %s\n",
			end.Format("2006-01-02 15:04:05"),
			ctx.Response.StatusCode(),
			end.Sub(start).Milliseconds(),
			ctx.Request.Header.Method(),
			ctx.Request.URI().FullURI(),
		)
	case DebugDetail:
		end := time.Now()
		body := ctx.Response.String()
		if len(body) > debugLimit {
			body = body[:debugLimit] + "\n"
		}
		fmt.Printf(
			"RESPONSE[%s]: %d %dms %s %s\n%s",
			end.Format("2006-01-02 15:04:05"),
			ctx.Response.StatusCode(),
			end.Sub(start).Milliseconds(),
			ctx.Request.Header.Method(),
			ctx.Request.URI().FullURI(),
			body,
		)
	}
}
