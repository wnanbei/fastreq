package fastreq

import (
	"io"
	"time"

	"github.com/gofiber/fiber/v2/utils"
	"github.com/valyala/fasthttp"
)

// ReClient ...
type ReClient struct {
	*fasthttp.Client
	UserAgent        string
	MaxRedirectCount int
	Timeout          time.Duration
	jsonEncoder      utils.JSONMarshal
	jsonDecoder      utils.JSONUnmarshal
	debugWriter      io.Writer
	mw               multipartWriter
}

// NewClient ...
func NewClient() *ReClient {
	return &ReClient{
		Client: &fasthttp.Client{},
	}
}

// NewClientFromFastHTTP ...
func NewClientFromFastHTTP(client *fasthttp.Client) *ReClient {
	return &ReClient{
		Client: client,
	}
}

func (c *ReClient) Get(url string) (*ReResponse, error) {
	req := NewRequest()
	req.SetMethod(GET)
	req.SetRequestURI(url)

	return c.do(req)
}

func (c *ReClient) Head(url string) (*ReResponse, error) {
	req := NewRequest()
	req.SetMethod(HEAD)
	req.SetRequestURI(url)

	return c.do(req)
}

func (c *ReClient) Post(url string) (*ReResponse, error) {
	req := NewRequest()
	req.SetMethod(POST)
	req.SetRequestURI(url)

	return c.do(req)
}

func (c *ReClient) Put(url string) (*ReResponse, error) {
	req := NewRequest()
	req.SetMethod(PUT)
	req.SetRequestURI(url)

	return c.do(req)
}

func (c *ReClient) Patch(url string) (*ReResponse, error) {
	req := NewRequest()
	req.SetMethod(PATCH)
	req.SetRequestURI(url)

	return c.do(req)
}

func (c *ReClient) Delete(url string) (*ReResponse, error) {
	req := NewRequest()
	req.SetMethod(DELETE)
	req.SetRequestURI(url)

	return c.do(req)
}

func (c *ReClient) Do(req *ReRequest) (*ReResponse, error) {
	return c.do(req)
}

func (c *ReClient) do(req *ReRequest) (*ReResponse, error) {
	resp := AcquireResponse()

	if err := c.Client.DoTimeout(req.req, resp, c.Timeout); err != nil {
		fasthttp.ReleaseResponse(resp)
		return nil, err
	}

	return &ReResponse{resp: resp}, nil
}
