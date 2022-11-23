package fastreq

import (
	"encoding/xml"
	"io/ioutil"
	"mime/multipart"
	"path/filepath"

	"github.com/valyala/fasthttp"
)

// Request ...
type Request struct {
	*fasthttp.Request
	mw *multipart.Writer
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

func (r *Request) SetQueryParams(args *Args) {
	r.URI().SetQueryStringBytes(args.QueryString())
}

func (r *Request) SetBasicAuth(username, password string) {
	r.URI().SetUsername(username)
	r.URI().SetPassword(password)
}

// ================================= Set uri End===================================

// ================================= Set Header ===================================

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

// ================================= Set Header End ===================================

// ================================= Set Body =========================================

func (r *Request) SetBodyJSON(v interface{}) error {
	r.Header.SetContentType(MIMEApplicationJSON)

	body, err := jsonMarshal(v)
	if err != nil {
		return err
	}

	r.SetBody(body)
	return nil
}

func (r *Request) SetBodyXML(v interface{}) error {
	r.Header.SetContentType(MIMEApplicationXML)

	body, err := xml.Marshal(v)
	if err != nil {
		return err
	}

	r.SetBody(body)
	return nil
}

func (r *Request) SetBodyForm(args *Args) {
	r.Header.SetContentType(MIMEApplicationForm)

	if args != nil {
		r.SetBody(args.QueryString())
	}

	Release(args)
}

// ================================= Set Body End ===================================

// ============================== Set Multipart Form ================================

func (r *Request) SetBodyBoundary(boundary string) {
	if r.mw == nil {
		r.mw = multipart.NewWriter(r.BodyWriter())
	}

	r.mw.SetBoundary(boundary)
}

func (r *Request) AddBodyField(field, value string) error {
	if r.mw == nil {
		r.mw = multipart.NewWriter(r.BodyWriter())
	}

	if err := r.mw.WriteField(field, value); err != nil {
		return err
	}

	return nil
}

func (r *Request) AddBodyFile(fieldName, filePath string) error {
	if r.mw == nil {
		r.mw = multipart.NewWriter(r.BodyWriter())
	}

	if fieldName == "" {
		// fieldname = "file" + strconv.Itoa(len(a.formFiles)+1) // TODO
	}

	content, err := ioutil.ReadFile(filepath.Clean(filePath))
	if err != nil {
		return err
	}

	w, err := r.mw.CreateFormFile(fieldName, filepath.Base(filePath))
	if err != nil {
		return err
	}
	if _, err = w.Write(content); err != nil {
		return err
	}

	return nil
}

// ============================== Set Multipart Form End ============================

func (r *Request) Copy() *Request {
	req := fasthttp.AcquireRequest()
	r.CopyTo(req)

	return NewRequestFromFastHTTP(req)
}

func (r *Request) Release() {
	fasthttp.ReleaseRequest(r.Request)
	r.mw = nil
}
