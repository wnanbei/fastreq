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

func (c *ReClient) Get(url string) *ReqResp {
	return c.createAgent(MethodGet, url)
}

func (c *ReClient) Head(url string) *ReqResp {
	return c.createAgent(MethodHead, url)
}

func (c *ReClient) Post(url string) *ReqResp {
	return c.createAgent(MethodPost, url)
}

func (c *ReClient) Put(url string) *ReqResp {
	return c.createAgent(MethodPut, url)
}

func (c *ReClient) Patch(url string) *ReqResp {
	return c.createAgent(MethodPatch, url)
}

func (c *ReClient) Delete(url string) *ReqResp {
	return c.createAgent(MethodDelete, url)
}

func (c *ReClient) createAgent(method, url string) *ReqResp {
	a := AcquireAgent()
	a.req.Header.SetMethod(method)
	a.req.SetRequestURI(url)

	if err := a.Parse(); err != nil {
		a.errs = append(a.errs, err)
	}

	return &ReqResp{}
}
