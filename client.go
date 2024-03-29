package fastreq

import (
	"crypto/tls"
	"fmt"
	"time"

	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpproxy"
)

var defaultClientConfig = ClientConfig{
	Timeout:           time.Second * 30,
	DebugLevel:        DebugClose,
	MaxRedirectsCount: 10,
	DefaultUserAgent:  defaultUserAgent,
}

// ClientConfig Client Config
type ClientConfig struct {
	// Timeout global timeout
	Timeout time.Duration

	// DebugLevel ...
	DebugLevel DebugLevel

	// MaxRedirectsCount ...
	MaxRedirectsCount int

	// DefaultUserAgent ...
	DefaultUserAgent string
}

type Client struct {
	*fasthttp.Client
	defaultUserAgent  []byte
	maxRedirectsCount int
	timeout           time.Duration
	debugLevel        DebugLevel
	auth              Oauth1
	middlewares       []Middleware
}

// NewClient creates a new instance of a Client.
// If no configuration is provided, the default configuration is used.
func NewClient(config ...*ClientConfig) *Client {
	// Use the default client config if no config is provided
	var realConfig *ClientConfig
	if len(config) > 0 {
		realConfig = config[0]
	} else {
		realConfig = &defaultClientConfig
	}

	// Create a new fasthttp client with the provided config
	client := &Client{
		Client: &fasthttp.Client{
			Name:                          "",
			NoDefaultUserAgentHeader:      false,
			Dial:                          nil,
			DialDualStack:                 false,
			TLSConfig:                     nil,
			MaxConnsPerHost:               0,
			MaxIdleConnDuration:           0,
			MaxConnDuration:               0,
			MaxIdemponentCallAttempts:     0,
			ReadBufferSize:                0,
			WriteBufferSize:               0,
			ReadTimeout:                   0,
			WriteTimeout:                  0,
			MaxResponseBodySize:           0,
			DisableHeaderNamesNormalizing: false,
			DisablePathNormalizing:        false,
			MaxConnWaitTimeout:            0,
			RetryIf:                       nil,
			ConnPoolStrategy:              0,
			ConfigureClient:               nil,
		},
		middlewares:       []Middleware{},
		timeout:           realConfig.Timeout,
		debugLevel:        realConfig.DebugLevel,
		maxRedirectsCount: realConfig.MaxRedirectsCount,
		defaultUserAgent:  unsafeS2B(realConfig.DefaultUserAgent),
	}

	return client
}

// Get performs an HTTP GET request to the specified URL with optional request options.
func (c *Client) Get(url string, opts ...ReqOption) (*Response, error) {
	req := NewRequest(GET, url)
	return c.Do(req, opts...)
}

// Head sends a HEAD request to the specified URL and returns the response.
func (c *Client) Head(url string, opts ...ReqOption) (*Response, error) {
	req := NewRequest(HEAD, url)
	return c.Do(req, opts...)
}

// Post sends an HTTP POST request to the specified URL with the provided request options.
func (c *Client) Post(url string, opts ...ReqOption) (*Response, error) {
	req := NewRequest(POST, url)
	return c.Do(req, opts...)
}

// Put sends an HTTP PUT request to the specified URL with the provided request options.
func (c *Client) Put(url string, opts ...ReqOption) (*Response, error) {
	req := NewRequest(PUT, url)
	return c.Do(req, opts...)
}

// Patch sends an HTTP PATCH request to the specified URL with the provided request options.
func (c *Client) Patch(url string, opts ...ReqOption) (*Response, error) {
	req := NewRequest(PATCH, url)
	return c.Do(req, opts...)
}

// Delete sends an HTTP DELETE request to the specified URL with the provided request options.
func (c *Client) Delete(url string, opts ...ReqOption) (*Response, error) {
	req := NewRequest(DELETE, url)
	return c.Do(req, opts...)
}

// Options sends an HTTP OPTIONS request to the specified URL with the provided request options.
func (c *Client) Options(url string, opts ...ReqOption) (*Response, error) {
	req := NewRequest(OPTIONS, url)
	return c.Do(req, opts...)
}

// Connect sends an HTTP CONNECT request to the specified URL with the provided request options.
func (c *Client) Connect(url string, opts ...ReqOption) (*Response, error) {
	req := NewRequest(CONNECT, url)
	return c.Do(req, opts...)
}

// DownloadFile downloads a file from the specified URL.
func (c *Client) DownloadFile(req *Request, path, filename string) error {
	resp, err := c.Do(req)
	if err != nil {
		return err
	}

	if err := resp.SaveToFile(path, filename); err != nil {
		return err
	}

	resp.Release()

	return nil
}

// Do execute an HTTP request with optional request options
func (c *Client) Do(req *Request, opts ...ReqOption) (*Response, error) {
	// apply all request options
	for _, o := range opts {
		if err := o.BindRequest(req); err != nil {
			return nil, err
		}
		if o.isAutoRelease() {
			Release(o)
		}
	}

	// set default user agent if none is provided
	if len(req.Header.UserAgent()) == 0 {
		req.Header.SetUserAgentBytes(c.defaultUserAgent)
	}

	// create a context object with the request and client info
	ctx := NewCtx()
	ctx.Request = req
	ctx.client = c

	// apply the first middleware function if there are any
	if len(c.middlewares) > 0 {
		if err := c.middlewares[0](ctx); err != nil {
			return nil, err
		}
	} else {
		// execute the request
		if err := do(ctx); err != nil {
			return nil, err
		}
	}

	return ctx.Response, nil
}

// do executes the request
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
	ctx.Response = &Response{Response: resp, Request: ctx.Request.Request}

	debugAfterRequest(ctx, start)

	return nil
}

// SetHTTPProxy sets the HTTP proxy to use for requests
func (c *Client) SetHTTPProxy(proxy string) {
	c.Dial = fasthttpproxy.FasthttpHTTPDialer(proxy)
}

// SetSocks5Proxy sets the SOCKS5 proxy to use for requests
func (c *Client) SetSocks5Proxy(proxy string) {
	c.Dial = fasthttpproxy.FasthttpSocksDialer(proxy)
}

// SetEnvHTTPProxy using the env(HTTP_PROXY, HTTPS_PROXY and NO_PROXY) configured HTTP proxy for the default client.
func (c *Client) SetEnvHTTPProxy() {
	c.Dial = fasthttpproxy.FasthttpProxyHTTPDialer()
}

// SetTimeout sets the request timeout
func (c *Client) SetTimeout(timeout time.Duration) {
	c.timeout = timeout
}

// SetDefaultUserAgent sets the default user agent
func (c *Client) SetDefaultUserAgent(userAgent string) {
	c.defaultUserAgent = unsafeS2B(userAgent)
}

// SetDebugLevel sets the debug level
func (c *Client) SetDebugLevel(lvl DebugLevel) {
	c.debugLevel = lvl
}

// SetTLSConfig sets the TLS config
func (c *Client) SetTLSConfig(config *tls.Config) {
	c.TLSConfig = config
}

// SetMaxRedirectsCount sets the max redirects count
func (c *Client) SetMaxRedirectsCount(count int) {
	c.maxRedirectsCount = count
}

// SetRetryIf sets the RetryIf function for the HTTP client.
func (c *Client) SetRetryIf(retryIf fasthttp.RetryIfFunc) {
	c.RetryIf = retryIf
}

// SkipInsecureVerify sets whether the client should skip verification of the server's
// certificate chain and hostname. If isSkip is true, the client will skip verification.
// If c.TLSConfig is nil, a new TLS config with InsecureSkipVerify set to isSkip will be
// created and set as c.TLSConfig. Otherwise, InsecureSkipVerify will be set to isSkip in
// c.TLSConfig.

// @param isSkip: bool - whether to skip verification
func (c *Client) SkipInsecureVerify(isSkip bool) {
	if c.TLSConfig == nil {
		/* #nosec G402 */
		c.TLSConfig = &tls.Config{InsecureSkipVerify: isSkip} // #nosec G402
	} else {
		/* #nosec G402 */
		c.TLSConfig.InsecureSkipVerify = isSkip
	}
}

// SetOauth1 sets the Oauth1 middleware
func (c *Client) SetOauth1(o *Oauth1) {
	c.middlewares = append(c.middlewares, MiddlewareOauth1(o))
}

// AddMiddleware appends one or more Middleware functions to the Client's list of
// middlewares. These middlewares are called in the order they are provided when
// sending HTTP requests.
func (c *Client) AddMiddleware(middlewares ...Middleware) {
	c.middlewares = append(c.middlewares, middlewares...)
}

// debugBeforeRequest logs information about the incoming request when debugging is enabled.
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

// debugAfterRequest logs information about the incoming request when debugging is enabled.
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
