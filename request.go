package fastreq

import (
	"encoding/xml"
	"github.com/valyala/fasthttp"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
)

// Request ...
type Request struct {
	*fasthttp.Request
	mw           *multipart.Writer
	formFilesNum int
}

func NewRequest(method HTTPMethod, url string) *Request {
	req := &Request{Request: fasthttp.AcquireRequest()}
	req.SetMethod(method)
	req.SetRequestURI(url)
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

	content, err := os.ReadFile(filepath.Clean(filePath))
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
