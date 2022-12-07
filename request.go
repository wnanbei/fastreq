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

	Release(form)
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
	BindRequest(req *Request) error
}

type QueryParams struct {
	*fasthttp.Args
}

func NewQueryParams(kv ...string) *QueryParams {
	q := &QueryParams{fasthttp.AcquireArgs()}

	for i := 1; i < len(kv); i += 2 {
		q.Add(kv[i-1], kv[i])
	}

	return q
}

func (o *QueryParams) BindRequest(req *Request) error {
	req.Request.URI().SetQueryStringBytes(o.Args.QueryString())
	fasthttp.ReleaseArgs(o.Args)
	return nil
}

func (a *QueryParams) Release() {
	fasthttp.ReleaseArgs(a.Args)
}

type PostForm struct {
	*fasthttp.Args
}

func NewPostForm(kv ...string) *PostForm {
	f := &PostForm{fasthttp.AcquireArgs()}

	for i := 1; i < len(kv); i += 2 {
		f.Add(kv[i-1], kv[i])
	}

	return f
}

func (a *PostForm) BindRequest(req *Request) error {
	req.SetPostForm(a)
	fasthttp.ReleaseArgs(a.Args)
	return nil
}

func (a *PostForm) Release() {
	fasthttp.ReleaseArgs(a.Args)
}

type Body []byte

func NewBody(b []byte) *Body {
	body := Body(b)
	return &body
}

func (a *Body) BindRequest(req *Request) error {
	req.SetBody(*a)
	return nil
}

type JsonBody struct {
	b interface{}
}

func NewJsonBody(b interface{}) *JsonBody {
	return &JsonBody{b: b}
}

func (a *JsonBody) BindRequest(req *Request) error {
	return req.SetJSON(a.b)
}

type Timeout time.Duration

func (t *Timeout) BindRequest(req *Request) {
	req.SetTimeout(time.Duration(*t))
}

type MultipartForm struct {
	*fasthttp.Args
	Boundary string
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

func (a *MultipartForm) BindRequest(req *Request) error {
	if err := req.SetBoundary(a.Boundary); err != nil {
		return err
	}

	var err error
	if a.Args != nil {
		a.Args.VisitAll(func(key, value []byte) {
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

func (a *MultipartForm) Release() {
	fasthttp.ReleaseArgs(a.Args)
}
