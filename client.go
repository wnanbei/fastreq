package fastreq

import (
	"encoding/json"
	"sync"

	"github.com/gofiber/fiber/v2/utils"
)

var defaultClient Client

type Client struct {
	mutex                    sync.RWMutex
	UserAgent                string
	NoDefaultUserAgentHeader bool
	JSONEncoder              utils.JSONMarshal
	JSONDecoder              utils.JSONUnmarshal
}

func Get(url string) *Agent { return defaultClient.Get(url) }

func (c *Client) Get(url string) *Agent {
	return c.createAgent(MethodGet, url)
}

func Head(url string) *Agent { return defaultClient.Head(url) }

func (c *Client) Head(url string) *Agent {
	return c.createAgent(MethodHead, url)
}

func Post(url string) *Agent { return defaultClient.Post(url) }

func (c *Client) Post(url string) *Agent {
	return c.createAgent(MethodPost, url)
}

func Put(url string) *Agent { return defaultClient.Put(url) }

func (c *Client) Put(url string) *Agent {
	return c.createAgent(MethodPut, url)
}

func Patch(url string) *Agent { return defaultClient.Patch(url) }

func (c *Client) Patch(url string) *Agent {
	return c.createAgent(MethodPatch, url)
}

func Delete(url string) *Agent { return defaultClient.Delete(url) }

func (c *Client) Delete(url string) *Agent {
	return c.createAgent(MethodDelete, url)
}

func (c *Client) createAgent(method, url string) *Agent {
	a := AcquireAgent()
	a.req.Header.SetMethod(method)
	a.req.SetRequestURI(url)

	c.mutex.RLock()
	a.Name = c.UserAgent
	a.NoDefaultUserAgentHeader = c.NoDefaultUserAgentHeader
	a.jsonDecoder = c.JSONDecoder
	a.jsonEncoder = c.JSONEncoder
	if a.jsonDecoder == nil {
		a.jsonDecoder = json.Unmarshal
	}
	c.mutex.RUnlock()

	if err := a.Parse(); err != nil {
		a.errs = append(a.errs, err)
	}

	return a
}
