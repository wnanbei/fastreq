package fastreq

import "github.com/valyala/fasthttp"

type Cookies struct {
	cookies []*fasthttp.Cookie
}

func NewCookies(kv ...string) *Cookies {
	cookies := &Cookies{}

	for i := 1; i < len(kv); i += 2 {
		cookie := fasthttp.AcquireCookie()
		cookie.SetKey(kv[i-1])
		cookie.SetValue(kv[i])
		cookies.cookies = append(cookies.cookies, cookie)
	}

	return cookies
}

func (c *Cookies) BindRequest(req *Request) error {
	for i := range c.cookies {
		req.Header.SetCookieBytesKV(c.cookies[i].Key(), c.cookies[i].Value())
	}
	return nil
}

func (c *Cookies) Release() {
	for i := range c.cookies {
		fasthttp.ReleaseCookie(c.cookies[i])
	}
	c.cookies = c.cookies[:0]
}
