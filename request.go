package fastreq

import (
	"encoding/xml"
	"io/ioutil"
	"mime/multipart"
	"path/filepath"
	"strconv"
	"time"

	"github.com/valyala/fasthttp"
)

// Request ...
type Request struct {
	*fasthttp.Request
	mw           *multipart.Writer
	formFilesNum int
}

func NewRequest(method HTTPMethod, url string, opts ...ReqOption) *Request {
	req := &Request{Request: fasthttp.AcquireRequest()}
	req.SetMethod(method)
	req.SetRequestURI(url)

	for _, opt := range opts {
		opt.BindRequest(req)
		if opt.isAutoRelease() {
			Release(opt)
		}
	}

	return req
}

func NewRequestFromFastHTTP(req *fasthttp.Request) *Request {
	return &Request{
		Request: req,
	}
}

// ================================= Set uri =====================================

func (r *Request) SetHost(host string) {
	r.URI().SetHost(host)
}

func (r *Request) SetQueryString(queryString string) {
	r.URI().SetQueryString(queryString)
}

func (r *Request) SetQueryParams(params *QueryParams) {
	r.URI().SetQueryStringBytes(params.QueryString())
}

func (r *Request) SetBasicAuth(username, password string) {
	r.URI().SetUsername(username)
	r.URI().SetPassword(password)
}

func (r *Request) SetMethod(method HTTPMethod) {
	r.Header.SetMethod(string(method))
}

func (r *Request) SetUserAgent(userAgent string) {
	r.Header.SetUserAgent(userAgent)
}

func (r *Request) SetReferer(referer string) {
	r.Header.SetReferer(referer)
}

func (r *Request) SetContentType(contentType string) {
	r.Header.SetContentType(contentType)
}

func (r *Request) SetHeaders(kv ...string) {
	for i := 1; i < len(kv); i += 2 {
		r.Header.Set(kv[i-1], kv[i])
	}
}

func (r *Request) SetHeader(k, v string) {
	r.Header.Set(k, v)
}

func (r *Request) AddHeader(k, v string) {
	r.Header.Add(k, v)
}

func (r *Request) SetCookies(kv ...string) {
	for i := 1; i < len(kv); i += 2 {
		r.Header.SetCookie(kv[i-1], kv[i])
	}
}

func (r *Request) SetCookie(key, value string) {
	r.Header.SetCookie(key, value)
}

func (r *Request) SetJSON(v interface{}) error {
	r.Header.SetContentType(MIMEApplicationJSON)

	body, err := jsonMarshal(v)
	if err != nil {
		return err
	}

	r.SetBody(body)
	return nil
}

func (r *Request) SetXML(v interface{}) error {
	r.Header.SetContentType(MIMEApplicationXML)

	body, err := xml.Marshal(v)
	if err != nil {
		return err
	}

	r.SetBody(body)
	return nil
}

func (r *Request) SetPostForm(form *PostForm) {
	r.Header.SetContentType(MIMEApplicationForm)

	if form != nil {
		r.SetBody(form.QueryString())
	}
}

func (r *Request) SetBoundary(boundary string) error {
	if r.mw == nil {
		r.mw = multipart.NewWriter(r.BodyWriter())
	}

	return r.mw.SetBoundary(boundary)
}

func (r *Request) AddMFField(field, value string) error {
	if r.mw == nil {
		r.mw = multipart.NewWriter(r.BodyWriter())
	}

	if err := r.mw.WriteField(field, value); err != nil {
		return err
	}

	return nil
}

func (r *Request) AddMFFile(fieldName, filePath string) error {
	if r.mw == nil {
		r.mw = multipart.NewWriter(r.BodyWriter())
	}

	content, err := ioutil.ReadFile(filepath.Clean(filePath))
	if err != nil {
		return err
	}

	if fieldName == "" { // default field name
		fieldName = "file" + strconv.Itoa(r.formFilesNum+1)
	}

	w, err := r.mw.CreateFormFile(fieldName, filepath.Base(filePath))
	if err != nil {
		return err
	}
	if _, err = w.Write(content); err != nil {
		return err
	}

	r.formFilesNum++
	return nil
}

func (r *Request) Copy() *Request {
	req := fasthttp.AcquireRequest()
	r.CopyTo(req)

	return NewRequestFromFastHTTP(req)
}

func (r *Request) Release() {
	fasthttp.ReleaseRequest(r.Request)
	r.mw = nil
}

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
	req.SetTimeout(time.Duration(*t.timeout))
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
