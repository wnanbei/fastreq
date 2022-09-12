package fastreq

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"mime/multipart"

	"github.com/valyala/fasthttp"
)

// ReRequest
type ReRequest struct {
	req *fasthttp.Request
}

func NewRequest() *ReRequest {
	return &ReRequest{
		fasthttp.AcquireRequest(),
	}
}

func NewRequestFromFastHTTP(req *fasthttp.Request) *ReRequest {
	return &ReRequest{
		req,
	}
}

func ReleaseRequest(r *ReRequest) {
	fasthttp.ReleaseRequest(r.req)
}

// ================================= Set uri ===================================

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

// ================================================================================

// ================================= Set Header ===================================

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

func (r *ReRequest) SetJSONBody(v interface{}) error {
	r.req.Header.SetContentType(MIMEApplicationJSON)

	body, err := json.Marshal(v)
	if err != nil {
		return err
	}

	r.req.SetBody(body)
	return nil
}

func (r *ReRequest) SetXMLBody(v interface{}) error {
	r.req.Header.SetContentType(MIMEApplicationXML)

	body, err := xml.Marshal(v)
	if err != nil {
		return err
	}

	r.req.SetBody(body)
	return nil
}

func (r *ReRequest) SetFormBody(args *fasthttp.Args) {
	r.req.Header.SetContentType(MIMEApplicationForm)

	if args != nil {
		r.req.SetBody(args.QueryString())
	}

	fasthttp.ReleaseArgs(args)
}

func (r *ReRequest) SetMultipartForm(f *multipart.Form, boundary string) error {
	return fasthttp.WriteMultipartForm(r.req.BodyWriter(), f, boundary)
}

// ================================= Set Body End ===================================
