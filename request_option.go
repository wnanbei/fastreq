package fastreq

import (
	"github.com/valyala/fasthttp"
	"time"
)

type ReqOption interface {
	Releaser
	BindRequest(req *Request) error
	AutoRelease(bool)
	isAutoRelease() bool
}

type QueryParams struct {
	*fasthttp.Args
	notAutoRelease bool
}

func NewQueryParams(kv ...string) *QueryParams {
	q := &QueryParams{Args: fasthttp.AcquireArgs()}

	for i := 1; i < len(kv); i += 2 {
		q.Add(kv[i-1], kv[i])
	}

	return q
}

func (q *QueryParams) BindRequest(req *Request) error {
	req.Request.URI().SetQueryStringBytes(q.Args.QueryString())
	return nil
}

func (q *QueryParams) Release() {
	fasthttp.ReleaseArgs(q.Args)
	q.Args = nil
	q.notAutoRelease = false
}

func (q *QueryParams) AutoRelease(auto bool) {
	q.notAutoRelease = !auto
}

func (q *QueryParams) isAutoRelease() bool {
	return !q.notAutoRelease
}

type PostForm struct {
	*fasthttp.Args
	notAutoRelease bool
}

func NewPostForm(kv ...string) *PostForm {
	f := &PostForm{Args: fasthttp.AcquireArgs()}

	for i := 1; i < len(kv); i += 2 {
		f.Add(kv[i-1], kv[i])
	}

	return f
}

func (f *PostForm) BindRequest(req *Request) error {
	req.SetPostForm(f)
	return nil
}

func (f *PostForm) Release() {
	fasthttp.ReleaseArgs(f.Args)
	f.Args = nil
	f.notAutoRelease = false
}

func (f *PostForm) AutoRelease(auto bool) {
	f.notAutoRelease = !auto
}

func (f *PostForm) isAutoRelease() bool {
	return !f.notAutoRelease
}

type Body struct {
	body           []byte
	notAutoRelease bool
}

func NewBody(b []byte) *Body {
	body := Body{body: b}
	return &body
}

func (b *Body) BindRequest(req *Request) error {
	req.SetBody(b.body)
	return nil
}

func (b *Body) Release() {
	b.body = nil
	b.notAutoRelease = false
}

func (b *Body) AutoRelease(auto bool) {
	b.notAutoRelease = !auto
}

func (b *Body) isAutoRelease() bool {
	return !b.notAutoRelease
}

type JsonBody struct {
	body           interface{}
	notAutoRelease bool
}

func NewJsonBody(b interface{}) *JsonBody {
	return &JsonBody{body: b}
}

func (b *JsonBody) BindRequest(req *Request) error {
	return req.SetJSON(b.body)
}

func (b *JsonBody) Release() {
	b.body = nil
	b.notAutoRelease = false
}

func (b *JsonBody) AutoRelease(auto bool) {
	b.notAutoRelease = !auto
}

func (b *JsonBody) isAutoRelease() bool {
	return !b.notAutoRelease
}

type Timeout struct {
	timeout        *time.Duration
	notAutoRelease bool
}

func NewTimeout(t time.Duration) *Timeout {
	return &Timeout{timeout: &t}
}

func (t *Timeout) BindRequest(req *Request) error {
	req.SetTimeout(*t.timeout)
	return nil
}

func (t *Timeout) Release() {
	t.timeout = nil
	t.notAutoRelease = false
}

func (t *Timeout) AutoRelease(auto bool) {
	t.notAutoRelease = !auto
}

func (t *Timeout) isAutoRelease() bool {
	return !t.notAutoRelease
}

type MultipartForm struct {
	*fasthttp.Args
	Boundary       string
	notAutoRelease bool
}

func NewMultipartForm(boundary string, kv ...string) *MultipartForm {
	f := &MultipartForm{
		Boundary: boundary,
		Args:     fasthttp.AcquireArgs(),
	}

	for i := 1; i < len(kv); i += 2 {
		f.Add(kv[i-1], kv[i])
	}

	return f
}

func (mf *MultipartForm) BindRequest(req *Request) error {
	if err := req.SetBoundary(mf.Boundary); err != nil {
		return err
	}

	var err error
	if mf.Args != nil {
		mf.Args.VisitAll(func(key, value []byte) {
			if addErr := req.AddMFField(unsafeB2S(key), unsafeB2S(value)); err != nil {
				err = addErr
				return
			}
		})
	}
	if err != nil {
		return err
	}

	if err := req.mw.Close(); err != nil {
		return err
	}
	return nil
}

func (mf *MultipartForm) Release() {
	fasthttp.ReleaseArgs(mf.Args)
	mf.Args = nil
	mf.notAutoRelease = false
	mf.Boundary = ""
}

func (mf *MultipartForm) AutoRelease(auto bool) {
	mf.notAutoRelease = !auto
}

func (mf *MultipartForm) isAutoRelease() bool {
	return !mf.notAutoRelease
}

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
