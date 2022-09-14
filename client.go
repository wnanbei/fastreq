package fastreq

import (
	"crypto/tls"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/valyala/fasthttp"
)

// Client ...
type Client struct {
	client            *fasthttp.Client
	userAgent         string
	maxRedirectsCount int
	timeout           time.Duration
	debugWriter       []io.Writer
}

// NewClient ...
func NewClient() *Client {
	return &Client{
		client:      &fasthttp.Client{},
		debugWriter: []io.Writer{os.Stdout},
	}
}

// NewClientFromFastHTTP ...
func NewClientFromFastHTTP(client *fasthttp.Client) *Client {
	return &Client{
		client: client,
	}
}

func (c *Client) Get(url string) (*Response, error) {
	req := NewRequest()
	req.SetMethod(GET)
	req.SetRequestURI(url)

	return c.do(req)
}

func (c *Client) Head(url string) (*Response, error) {
	req := NewRequest()
	req.SetMethod(HEAD)
	req.SetRequestURI(url)

	return c.do(req)
}

func (c *Client) Post(url string) (*Response, error) {
	req := NewRequest()
	req.SetMethod(POST)
	req.SetRequestURI(url)

	return c.do(req)
}

func (c *Client) Put(url string) (*Response, error) {
	req := NewRequest()
	req.SetMethod(PUT)
	req.SetRequestURI(url)

	return c.do(req)
}

func (c *Client) Patch(url string) (*Response, error) {
	req := NewRequest()
	req.SetMethod(PATCH)
	req.SetRequestURI(url)

	return c.do(req)
}

func (c *Client) Delete(url string) (*Response, error) {
	req := NewRequest()
	req.SetMethod(DELETE)
	req.SetRequestURI(url)

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

	if err := c.client.DoTimeout(req.req, resp, c.timeout); err != nil {
		fasthttp.ReleaseResponse(resp)
		return nil, err
	}

	return &Response{resp: resp}, nil
}

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

// ================================= Client Setting End =================================

func writeDebugInfo(req *fasthttp.Request, resp *fasthttp.Response, w io.Writer) {
	msg := fmt.Sprintf("Connected to %s(%s)\r\n\r\n", req.URI().Host(), resp.RemoteAddr())
	_, _ = w.Write(unsafeS2B(msg))
	_, _ = req.WriteTo(w)
	_, _ = resp.WriteTo(w)
}
