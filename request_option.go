package fastreq

import (
	"time"

	"github.com/valyala/fasthttp"
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

// NewQueryParams creates a new QueryParams object with key-value pairs passed as arguments.
func NewQueryParams(kv ...string) *QueryParams {
	q := &QueryParams{Args: fasthttp.AcquireArgs()}

	for i := 1; i < len(kv); i += 2 {
		q.Add(kv[i-1], kv[i])
	}

	return q
}

// BindRequest binds the query parameters to a request by setting the request's URI query string.
func (q *QueryParams) BindRequest(req *Request) error {
	req.Request.URI().SetQueryStringBytes(q.Args.QueryString())
	return nil
}

// Release frees the resources held by query
func (q *QueryParams) Release() {
	fasthttp.ReleaseArgs(q.Args)
	q.Args = nil
	q.notAutoRelease = false
}

// AutoRelease sets whether query parameters should be automatically released when the
// associated object is destroyed.
func (q *QueryParams) AutoRelease(auto bool) {
	q.notAutoRelease = !auto
}

// isAutoRelease returns true if the QueryParams instance is set to auto-release.
func (q *QueryParams) isAutoRelease() bool {
	return !q.notAutoRelease
}

type PostForm struct {
	*fasthttp.Args
	notAutoRelease bool
}

// NewPostForm creates a new PostForm object with key-value pairs passed as arguments.
func NewPostForm(kv ...string) *PostForm {
	f := &PostForm{Args: fasthttp.AcquireArgs()}

	for i := 1; i < len(kv); i += 2 {
		f.Add(kv[i-1], kv[i])
	}

	return f
}

// BindRequest binds the PostForm to a Request object.
func (f *PostForm) BindRequest(req *Request) error {
	req.SetPostForm(f)
	return nil
}

// Release frees the resources held by PostForm
func (f *PostForm) Release() {
	fasthttp.ReleaseArgs(f.Args)
	f.Args = nil
	f.notAutoRelease = false
}

// AutoRelease sets whether PostForm should be automatically released when the
// associated object is destroyed.
func (f *PostForm) AutoRelease(auto bool) {
	f.notAutoRelease = !auto
}

// isAutoRelease returns true if the PostForm instance is set to auto-release.
func (f *PostForm) isAutoRelease() bool {
	return !f.notAutoRelease
}

type Body struct {
	body           []byte
	notAutoRelease bool
}

// NewBody creates a new Body object
func NewBody(b []byte) *Body {
	body := Body{body: b}
	return &body
}

// BindRequest binds the Body to a Request object
func (b *Body) BindRequest(req *Request) error {
	req.SetBody(b.body)
	return nil
}

// Release frees the resources held by Body
func (b *Body) Release() {
	b.body = nil
	b.notAutoRelease = false
}

// AutoRelease sets whether Body should be automatically released when the
func (b *Body) AutoRelease(auto bool) {
	b.notAutoRelease = !auto
}

// isAutoRelease returns true if the Body instance is set to auto-release.
func (b *Body) isAutoRelease() bool {
	return !b.notAutoRelease
}

type JsonBody struct {
	body           interface{}
	notAutoRelease bool
}

// NewJsonBody creates a new JsonBody object
func NewJsonBody(b interface{}) *JsonBody {
	return &JsonBody{body: b}
}

// BindRequest binds the JsonBody to a Request object
func (b *JsonBody) BindRequest(req *Request) error {
	return req.SetJSON(b.body)
}

// Release frees the resources held by JsonBody
func (b *JsonBody) Release() {
	b.body = nil
	b.notAutoRelease = false
}

// AutoRelease sets whether JsonBody should be automatically released when the
// associated object is destroyed.
func (b *JsonBody) AutoRelease(auto bool) {
	b.notAutoRelease = !auto
}

// isAutoRelease returns true if the JsonBody instance is set to auto-release.
func (b *JsonBody) isAutoRelease() bool {
	return !b.notAutoRelease
}

type Timeout struct {
	timeout        *time.Duration
	notAutoRelease bool
}

// NewTimeout creates a new Timeout object
func NewTimeout(t time.Duration) *Timeout {
	return &Timeout{timeout: &t}
}

// BindRequest binds the Timeout to a Request object
func (t *Timeout) BindRequest(req *Request) error {
	req.SetTimeout(*t.timeout)
	return nil
}

// Release frees the resources held by Timeout
func (t *Timeout) Release() {
	t.timeout = nil
	t.notAutoRelease = false
}

// AutoRelease sets whether Timeout should be automatically released when the
// associated object is destroyed.
func (t *Timeout) AutoRelease(auto bool) {
	t.notAutoRelease = !auto
}

// isAutoRelease returns true if the Timeout instance is set to auto-release.
func (t *Timeout) isAutoRelease() bool {
	return !t.notAutoRelease
}

type MultipartForm struct {
	*fasthttp.Args
	Boundary       string
	notAutoRelease bool
}

// NewMultipartForm creates a new MultipartForm object
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

// BindRequest binds the MultipartForm to a Request object
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

// Release frees the resources held by MultipartForm
func (mf *MultipartForm) Release() {
	fasthttp.ReleaseArgs(mf.Args)
	mf.Args = nil
	mf.notAutoRelease = false
	mf.Boundary = ""
}

// AutoRelease sets whether MultipartForm should be automatically released when the
// associated object is destroyed.
func (mf *MultipartForm) AutoRelease(auto bool) {
	mf.notAutoRelease = !auto
}

// isAutoRelease returns true if the MultipartForm instance is set to auto-release.
func (mf *MultipartForm) isAutoRelease() bool {
	return !mf.notAutoRelease
}

type Cookies struct {
	cookies        []*fasthttp.Cookie
	notAutoRelease bool
}

// NewCookies creates a new Cookies object
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

// BindRequest binds the Cookies to a Request object
func (c *Cookies) BindRequest(req *Request) error {
	for i := range c.cookies {
		req.Header.SetCookieBytesKV(c.cookies[i].Key(), c.cookies[i].Value())
	}
	return nil
}

// Release frees the resources held by Cookies
func (c *Cookies) Release() {
	for i := range c.cookies {
		fasthttp.ReleaseCookie(c.cookies[i])
	}
	c.cookies = c.cookies[:0]
}

// AutoRelease sets whether Cookies should be automatically released when the
func (c *Cookies) AutoRelease(auto bool) {
	c.notAutoRelease = !auto
}

// isAutoRelease returns true if the Cookies instance is set to auto-release.
func (c *Cookies) isAutoRelease() bool {
	return !c.notAutoRelease
}
