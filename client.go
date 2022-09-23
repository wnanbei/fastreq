package fastreq

import (
	"crypto/tls"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpproxy"
)

// Client ...
type Client struct {
	client            *fasthttp.Client
	userAgent         string
	maxRedirectsCount int
	timeout           time.Duration
	debugWriter       []io.Writer
	auth              Oauth1
	reqMiddleware     []RequestMiddleware
	respMiddleware    []ResponseMiddleware
}

// NewClient ...
func NewClient() *Client {
	return &Client{
		client:        &fasthttp.Client{},
		debugWriter:   []io.Writer{os.Stdout},
		reqMiddleware: []RequestMiddleware{MiddlewareLog()},
	}
}

// NewClientFromFastHTTP ...
func NewClientFromFastHTTP(client *fasthttp.Client) *Client {
	return &Client{
		client:      client,
		debugWriter: []io.Writer{os.Stdout},
	}
}

// ================================= Client Send Request =================================

func (c *Client) Get(url string, params *Args) (*Response, error) {
	req := NewRequest(GET, url)

	if params != nil {
		req.SetQueryString(string(params.QueryString()))
	}

	return c.do(req)
}

func (c *Client) Head(url string, params *Args) (*Response, error) {
	req := NewRequest(HEAD, url)

	if params != nil {
		req.SetQueryString(string(params.QueryString()))
	}

	return c.do(req)
}

func (c *Client) Post(url string, body *Args) (*Response, error) {
	req := NewRequest(POST, url)

	if body != nil {
		req.SetBodyForm(body)
	}

	return c.do(req)
}

func (c *Client) Put(url string, body *Args) (*Response, error) {
	req := NewRequest(PUT, url)

	if body != nil {
		req.SetBodyForm(body)
	}

	return c.do(req)
}

func (c *Client) Patch(url string, params *Args) (*Response, error) {
	req := NewRequest(PATCH, url)

	if params != nil {
		req.SetQueryString(string(params.QueryString()))
	}

	return c.do(req)
}

func (c *Client) Delete(url string, params *Args) (*Response, error) {
	req := NewRequest(DELETE, url)

	if params != nil {
		req.SetQueryString(string(params.QueryString()))
	}

	return c.do(req)
}

func (c *Client) Do(req *Request) (*Response, error) {
	return c.do(req)
}

func (c *Client) do(req *Request) (*Response, error) {
	resp := fasthttp.AcquireResponse()

	if req.mw != nil {
		req.req.Header.SetMultipartFormBoundary(req.mw.Boundary())
		if err := req.mw.Close(); err != nil {
			return nil, err
		}
	}

	for _, rm := range c.reqMiddleware { // set request middleware
		rm(req)
	}

	if err := c.client.DoTimeout(req.req, resp, c.timeout); err != nil {
		fasthttp.ReleaseResponse(resp)
		return nil, err
	}

	response := &Response{resp: resp}

	for _, rm := range c.respMiddleware { // set response middleware
		rm(response)
	}

	return response, nil
}

// ================================= Client Send Request End ============================

// ================================= Client Send Proxy ===============================

func (c *Client) SetHTTPProxy(proxy string) {
	c.client.Dial = fasthttpproxy.FasthttpHTTPDialer(proxy)
}

func (c *Client) SetSocks5Proxy(proxy string) {
	c.client.Dial = fasthttpproxy.FasthttpSocksDialer(proxy)
}

func (c *Client) SetEnvHTTPProxy() {
	c.client.Dial = fasthttpproxy.FasthttpProxyHTTPDialer()
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
	c.client.TLSConfig = config
}

func (c *Client) SetMaxRedirectsCount(count int) {
	c.maxRedirectsCount = count
}

func (c *Client) SetRetryIf(retryIf fasthttp.RetryIfFunc) {
	c.client.RetryIf = retryIf
}

func (c *Client) SkipInsecureVerify(isSkip bool) {
	if c.client.TLSConfig == nil {
		/* #nosec G402 */
		c.client.TLSConfig = &tls.Config{InsecureSkipVerify: isSkip} // #nosec G402
	} else {
		/* #nosec G402 */
		c.client.TLSConfig.InsecureSkipVerify = isSkip
	}
}

func (c *Client) SetOauth1(o *Oauth1) {
	c.reqMiddleware = append(c.reqMiddleware, MiddlewareOauth1(o))
}

// ================================= Client Setting End =================================

func writeDebugInfo(req *fasthttp.Request, resp *fasthttp.Response, w io.Writer) {
	msg := fmt.Sprintf("Connected to %s(%s)\r\n\r\n", req.URI().Host(), resp.RemoteAddr())
	_, _ = w.Write(unsafeS2B(msg))
	_, _ = req.WriteTo(w)
	_, _ = resp.WriteTo(w)
}
