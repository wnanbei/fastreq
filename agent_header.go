package fastreq

func (a *Agent) Set(k, v string) *Agent {
	a.req.Header.Set(k, v)

	return a
}

func (a *Agent) SetBytesK(k []byte, v string) *Agent {
	a.req.Header.SetBytesK(k, v)

	return a
}

func (a *Agent) SetBytesV(k string, v []byte) *Agent {
	a.req.Header.SetBytesV(k, v)

	return a
}

func (a *Agent) SetBytesKV(k []byte, v []byte) *Agent {
	a.req.Header.SetBytesKV(k, v)

	return a
}

func (a *Agent) Add(k, v string) *Agent {
	a.req.Header.Add(k, v)

	return a
}

func (a *Agent) AddBytesK(k []byte, v string) *Agent {
	a.req.Header.AddBytesK(k, v)

	return a
}

func (a *Agent) AddBytesV(k string, v []byte) *Agent {
	a.req.Header.AddBytesV(k, v)

	return a
}

func (a *Agent) AddBytesKV(k []byte, v []byte) *Agent {
	a.req.Header.AddBytesKV(k, v)

	return a
}

func (a *Agent) ConnectionClose() *Agent {
	a.req.Header.SetConnectionClose()

	return a
}

func (a *Agent) UserAgent(userAgent string) *Agent {
	a.req.Header.SetUserAgent(userAgent)

	return a
}

func (a *Agent) UserAgentBytes(userAgent []byte) *Agent {
	a.req.Header.SetUserAgentBytes(userAgent)

	return a
}

func (a *Agent) Cookie(key, value string) *Agent {
	a.req.Header.SetCookie(key, value)

	return a
}

func (a *Agent) CookieBytesK(key []byte, value string) *Agent {
	a.req.Header.SetCookieBytesK(key, value)

	return a
}

func (a *Agent) CookieBytesKV(key, value []byte) *Agent {
	a.req.Header.SetCookieBytesKV(key, value)

	return a
}

func (a *Agent) Cookies(kv ...string) *Agent {
	for i := 1; i < len(kv); i += 2 {
		a.req.Header.SetCookie(kv[i-1], kv[i])
	}

	return a
}

func (a *Agent) CookiesBytesKV(kv ...[]byte) *Agent {
	for i := 1; i < len(kv); i += 2 {
		a.req.Header.SetCookieBytesKV(kv[i-1], kv[i])
	}

	return a
}

func (a *Agent) Referer(referer string) *Agent {
	a.req.Header.SetReferer(referer)

	return a
}

func (a *Agent) RefererBytes(referer []byte) *Agent {
	a.req.Header.SetRefererBytes(referer)

	return a
}

func (a *Agent) ContentType(contentType string) *Agent {
	a.req.Header.SetContentType(contentType)

	return a
}

func (a *Agent) ContentTypeBytes(contentType []byte) *Agent {
	a.req.Header.SetContentTypeBytes(contentType)

	return a
}
