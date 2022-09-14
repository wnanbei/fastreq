package fastreq

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"io/ioutil"
	"mime/multipart"
	"path/filepath"

	"github.com/valyala/fasthttp"
)

// Request
type Request struct {
	req *fasthttp.Request
	mw  *multipart.Writer
}

func NewRequest() *Request {
	return &Request{
		req: fasthttp.AcquireRequest(),
	}
}

func NewRequestFromFastHTTP(req *fasthttp.Request) *Request {
	return &Request{
		req: req,
	}
}

func ReleaseRequest(r *Request) {
	fasthttp.ReleaseRequest(r.req)
}

// ================================= Set uri =====================================

func (r *Request) SetRequestURI(url string) {
	r.req.SetRequestURI(url)
}

func (r *Request) SetHost(host string) {
	r.req.URI().SetHost(host)
}

func (r *Request) SetQueryString(queryString string) {
	r.req.URI().SetQueryString(queryString)
}

func (r *Request) SetBasicAuth(username, password string) {
	r.req.URI().SetUsername(username)
	r.req.URI().SetPassword(password)
}

// ================================= Set uri End===================================

// ================================= Set Header ===================================

func (r *Request) SetMethod(method HTTPMethod) {
	r.req.Header.SetMethod(string(method))
}

func (r *Request) SetUserAgent(userAgent string) {
	r.req.Header.SetUserAgent(userAgent)
}

func (r *Request) SetReferer(referer string) {
	r.req.Header.SetReferer(referer)
}

func (r *Request) SetContentType(contentType string) {
	r.req.Header.SetContentType(contentType)
}

func (r *Request) SetHeaders(kv ...string) {
	for i := 1; i < len(kv); i += 2 {
		r.req.Header.Set(kv[i-1], kv[i])
	}
}

func (r *Request) SetHeader(k, v string) {
	r.req.Header.Set(k, v)
}

func (r *Request) AddHeader(k, v string) {
	r.req.Header.Add(k, v)
}

func (r *Request) SetCookies(kv ...string) {
	for i := 1; i < len(kv); i += 2 {
		r.req.Header.SetCookie(kv[i-1], kv[i])
	}
}

func (r *Request) SetCookie(key, value string) {
	r.req.Header.SetCookie(key, value)
}

// ================================= Set Header End ===================================

// ================================= Set Body =========================================

func (r *Request) SetBody(body []byte) {
	r.req.SetBody(body)
}

func (r *Request) SetBodyStream(bodyStream io.Reader, bodySize int) {
	r.req.SetBodyStream(bodyStream, bodySize)
}

func (r *Request) SetBodyJSON(v interface{}) error {
	r.req.Header.SetContentType(MIMEApplicationJSON)

	body, err := json.Marshal(v)
	if err != nil {
		return err
	}

	r.req.SetBody(body)
	return nil
}

func (r *Request) SetBodyXML(v interface{}) error {
	r.req.Header.SetContentType(MIMEApplicationXML)

	body, err := xml.Marshal(v)
	if err != nil {
		return err
	}

	r.req.SetBody(body)
	return nil
}

func (r *Request) SetBodyForm(args *fasthttp.Args) {
	r.req.Header.SetContentType(MIMEApplicationForm)

	if args != nil {
		r.req.SetBody(args.QueryString())
	}

	fasthttp.ReleaseArgs(args)
}

// ================================= Set Body End ===================================

// ============================== Set Multipart Form ================================

func (r *Request) SetBodyBoundary(boundary string) {
	if r.mw == nil {
		r.mw = multipart.NewWriter(r.req.BodyWriter())
	}

	r.mw.SetBoundary(boundary)
}

func (r *Request) AddBodyField(field, value string) error {
	if r.mw == nil {
		r.mw = multipart.NewWriter(r.req.BodyWriter())
	}

	if err := r.mw.WriteField(field, value); err != nil {
		return err
	}

	return nil
}

func (r *Request) AddBodyFile(fieldName, filePath string) error {
	if r.mw == nil {
		r.mw = multipart.NewWriter(r.req.BodyWriter())
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
	r.req.CopyTo(req)

	return NewRequestFromFastHTTP(req)
}
