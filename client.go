package fastreq

import (
	"crypto/tls"
	"io"
	"os"
	"time"

	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpproxy"
	"github.com/wnanbei/fastreq/middleware/auth"
)

// Client ...
type Client struct {
	*fasthttp.Client
	userAgent         string
	maxRedirectsCount int
	timeout           time.Duration
	debugWriter       []io.Writer
	auth              auth.Oauth1
	middlewares       []Middleware
}

// NewClient ...
func NewClient() *Client {
	return &Client{
		Client:      &fasthttp.Client{},
		debugWriter: []io.Writer{os.Stdout},
		middlewares: []Middleware{},
	}
}

// NewClientFromFastHTTP ...
func NewClientFromFastHTTP(client *fasthttp.Client) *Client {
	return &Client{
		Client:      client,
		debugWriter: []io.Writer{os.Stdout},
	}
}

// ================================= Client Send Request =================================

func (c *Client) Get(url string, params *Args) (*Ctx, error) {
	req := NewRequest(GET, url)

	if params != nil {
		req.SetQueryParams(params)
	}

	return c.do(req)
}

func (c *Client) Head(url string, params *Args) (*Ctx, error) {
	req := NewRequest(HEAD, url)

	if params != nil {
		req.SetQueryParams(params)
	}

	return c.do(req)
}

func (c *Client) Post(url string, body *Args) (*Ctx, error) {
	req := NewRequest(POST, url)

	if body != nil {
		req.SetBodyForm(body)
	}

	return c.do(req)
}

func (c *Client) Put(url string, body *Args) (*Ctx, error) {
	req := NewRequest(PUT, url)

	if body != nil {
		req.SetBodyForm(body)
	}

	return c.do(req)
}

func (c *Client) Patch(url string, params *Args) (*Ctx, error) {
	req := NewRequest(PATCH, url)

	if params != nil {
		req.SetQueryParams(params)
	}

	return c.do(req)
}

func (c *Client) Delete(url string, params *Args) (*Ctx, error) {
	req := NewRequest(DELETE, url)

	if params != nil {
		req.SetQueryParams(params)
	}

	return c.do(req)
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

func (c *Client) Do(req *Request) (*Ctx, error) {
	return c.do(req)
}

func (c *Client) do(req *Request) (*Ctx, error) {
	ctx := AcquireCtx()
	ctx.Request = req
	ctx.client = c

	if err := c.middlewares[0](ctx); err != nil {
		return nil, err
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

	resp := fasthttp.AcquireResponse()
	if err := ctx.client.DoTimeout(ctx.fastRequest(), resp, ctx.client.timeout); err != nil {
		fasthttp.ReleaseResponse(resp)
		return err
	}

	ctx.Response = NewResponseFromFastHTTP(resp)
	return nil
}

// ================================= Client Send Request End ============================

// ================================= Client Send Proxy ===============================

func (c *Client) SetHTTPProxy(proxy string) {
	c.Dial = fasthttpproxy.FasthttpHTTPDialer(proxy)
}

func (c *Client) SetSocks5Proxy(proxy string) {
	c.Dial = fasthttpproxy.FasthttpSocksDialer(proxy)
}

func (c *Client) SetEnvHTTPProxy() {
	c.Dial = fasthttpproxy.FasthttpProxyHTTPDialer()
}

// ================================= Client Send Proxy End ===============================

// ================================= Client Settings ====================================

func (c *Client) SetTimeout(timeout time.Duration) {
	c.timeout = timeout
}

func (c *Client) SetUserAgent(userAgent string) {
	c.userAgent = userAgent
}

func (c *Client) SetDebugWriter(debugWriter ...io.Writer) {
	c.debugWriter = debugWriter
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

func (c *Client) SetOauth1(o *auth.Oauth1) {
	c.middlewares = append(c.middlewares, auth.MiddlewareOauth1(o))
}
