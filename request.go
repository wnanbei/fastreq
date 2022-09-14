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

// ReRequest
type ReRequest struct {
	req *fasthttp.Request
	mw  *multipart.Writer
}

func NewRequest() *ReRequest {
	return &ReRequest{
		req: fasthttp.AcquireRequest(),
	}
}

func NewRequestFromFastHTTP(req *fasthttp.Request) *ReRequest {
	return &ReRequest{
		req: req,
	}
}

func ReleaseRequest(r *ReRequest) {
	fasthttp.ReleaseRequest(r.req)
}

// ================================= Set uri =====================================

func (r *ReRequest) SetRequestURI(url string) {
	r.req.SetRequestURI(url)
}

func (r *ReRequest) SetHost(host string) {
	r.req.URI().SetHost(host)
}

func (r *ReRequest) SetQueryString(queryString string) {
	r.req.URI().SetQueryString(queryString)
}

func (r *ReRequest) SetBasicAuth(username, password string) {
	r.req.URI().SetUsername(username)
	r.req.URI().SetPassword(password)
}

// ================================= Set uri End===================================

// ================================= Set Header ===================================

func (r *ReRequest) SetMethod(method HTTPMethod) {
	r.req.Header.SetMethod(string(method))
}

func (r *ReRequest) SetUserAgent(userAgent string) {
	r.req.Header.SetUserAgent(userAgent)
}

func (r *ReRequest) SetReferer(referer string) {
	r.req.Header.SetReferer(referer)
}

func (r *ReRequest) SetContentType(contentType string) {
	r.req.Header.SetContentType(contentType)
}

func (r *ReRequest) SetHeaders(kv ...string) {
	for i := 1; i < len(kv); i += 2 {
		r.req.Header.Set(kv[i-1], kv[i])
	}
}

func (r *ReRequest) SetHeader(k, v string) {
	r.req.Header.Set(k, v)
}

func (r *ReRequest) AddHeader(k, v string) {
	r.req.Header.Add(k, v)
}

func (r *ReRequest) SetCookies(kv ...string) {
	for i := 1; i < len(kv); i += 2 {
		r.req.Header.SetCookie(kv[i-1], kv[i])
	}
}

func (r *ReRequest) SetCookie(key, value string) {
	r.req.Header.SetCookie(key, value)
}

// ================================= Set Header End ===================================

// ================================= Set Body =========================================

func (r *ReRequest) SetBody(body []byte) {
	r.req.SetBody(body)
}

func (r *ReRequest) SetBodyStream(bodyStream io.Reader, bodySize int) {
	r.req.SetBodyStream(bodyStream, bodySize)
}

func (r *ReRequest) SetBodyJSON(v interface{}) error {
	r.req.Header.SetContentType(MIMEApplicationJSON)

	body, err := json.Marshal(v)
	if err != nil {
		return err
	}

	r.req.SetBody(body)
	return nil
}

func (r *ReRequest) SetBodyXML(v interface{}) error {
	r.req.Header.SetContentType(MIMEApplicationXML)

	body, err := xml.Marshal(v)
	if err != nil {
		return err
	}

	r.req.SetBody(body)
	return nil
}

func (r *ReRequest) SetBodyForm(args *fasthttp.Args) {
	r.req.Header.SetContentType(MIMEApplicationForm)

	if args != nil {
		r.req.SetBody(args.QueryString())
	}

	fasthttp.ReleaseArgs(args)
}

// ================================= Set Body End ===================================

func (r *ReRequest) SetBodyBoundary(boundary string) {
	if r.mw == nil {
		r.mw = multipart.NewWriter(r.req.BodyWriter())
	}

	r.mw.SetBoundary(boundary)
}

func (r *ReRequest) AddBodyField(field, value string) error {
	if r.mw == nil {
		r.mw = multipart.NewWriter(r.req.BodyWriter())
	}

	if err := r.mw.WriteField(field, value); err != nil {
		return err
	}

	return nil
}

func (r *ReRequest) AddBodyFile(fieldName, filePath string) error {
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

func (r *ReRequest) Copy() *ReRequest {
	req := fasthttp.AcquireRequest()
	r.req.CopyTo(req)

	return NewRequestFromFastHTTP(req)
}
