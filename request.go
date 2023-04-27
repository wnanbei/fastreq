package fastreq

import (
	"encoding/xml"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"

	"github.com/valyala/fasthttp"
)

// Request represents an HTTP request.
type Request struct {
	*fasthttp.Request
	mw           *multipart.Writer
	formFilesNum int
}

// NewRequest creates a new HTTP request with the given method and URL.
func NewRequest(method HTTPMethod, url string) *Request {
	req := &Request{Request: fasthttp.AcquireRequest()}
	req.SetMethod(method)
	req.SetRequestURI(url)
	return req
}

// NewRequestFromFastHTTP returns a new Request object created from the given
// fasthttp.Request object.
func NewRequestFromFastHTTP(req *fasthttp.Request) *Request {
	return &Request{
		Request: req,
	}
}

// SetHost sets the host of the request
func (r *Request) SetHost(host string) {
	r.URI().SetHost(host)
}

// SetQueryString sets the query string for the current request.
func (r *Request) SetQueryString(queryString string) {
	r.URI().SetQueryString(queryString)
}

// SetQueryParams sets the query string of the request's URI using the
// provided QueryParams object. The query string is constructed from the
// QueryString method of the QueryParams object.
func (r *Request) SetQueryParams(params *QueryParams) {
	r.URI().SetQueryStringBytes(params.QueryString())
}

// SetBasicAuth sets the username and password of the request
func (r *Request) SetBasicAuth(username, password string) {
	r.URI().SetUsername(username)
	r.URI().SetPassword(password)
}

// SetMethod sets the HTTP request method for the given Request object.
func (r *Request) SetMethod(method HTTPMethod) {
	r.Header.SetMethod(string(method))
}

// SetUserAgent sets the User-Agent header field in the request header to the given
// userAgent string.
func (r *Request) SetUserAgent(userAgent string) {
	r.Header.SetUserAgent(userAgent)
}

// SetReferer sets the 'Referer' header of the HTTP request to the provided value.
func (r *Request) SetReferer(referer string) {
	r.Header.SetReferer(referer)
}

// SetContentType sets the Content-Type header of the request to the given content type.
func (r *Request) SetContentType(contentType string) {
	r.Header.SetContentType(contentType)
}

// SetHeaders sets the headers of a request object with key-value pairs.
func (r *Request) SetHeaders(kv ...string) {
	for i := 1; i < len(kv); i += 2 {
		r.Header.Set(kv[i-1], kv[i])
	}
}

// SetHeader sets the header with the given key-value pair in the HTTP request.
func (r *Request) SetHeader(k, v string) {
	r.Header.Set(k, v)
}

// AddHeader adds a header to the request.
func (r *Request) AddHeader(k, v string) {
	r.Header.Add(k, v)
}

// SetCookies sets the cookies of the request
func (r *Request) SetCookies(kv ...string) {
	for i := 1; i < len(kv); i += 2 {
		r.Header.SetCookie(kv[i-1], kv[i])
	}
}

// SetCookie adds a cookie to the request
func (r *Request) SetCookie(key, value string) {
	r.Header.SetCookie(key, value)
}

// SetJSON sets the JSON body of the request
func (r *Request) SetJSON(v interface{}) error {
	r.Header.SetContentType(MIMEApplicationJSON)

	body, err := jsonMarshal(v)
	if err != nil {
		return err
	}

	r.SetBody(body)
	return nil
}

// SetXML sets the XML body of the request
func (r *Request) SetXML(v interface{}) error {
	r.Header.SetContentType(MIMEApplicationXML)

	body, err := xml.Marshal(v)
	if err != nil {
		return err
	}

	r.SetBody(body)
	return nil
}

// SetPostForm sets the PostForm body of the request
func (r *Request) SetPostForm(form *PostForm) {
	r.Header.SetContentType(MIMEApplicationForm)

	if form != nil {
		r.SetBody(form.QueryString())
	}
}

// SetBoundary sets the boundary of the request
func (r *Request) SetBoundary(boundary string) error {
	if r.mw == nil {
		r.mw = multipart.NewWriter(r.BodyWriter())
	}

	return r.mw.SetBoundary(boundary)
}

// AddMFField writes a form field to a multipart request. If the request's
// multipart writer is not initialized, it initializes it before writing the
// field.
func (r *Request) AddMFField(field, value string) error {
	if r.mw == nil {
		r.mw = multipart.NewWriter(r.BodyWriter())
	}

	if err := r.mw.WriteField(field, value); err != nil {
		return err
	}

	return nil
}

// AddMFFile adds a multipart/form-data file to the request body using the given
// field name and file path.
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

// Copy returns a new instance of the Request struct with the same values as r.
// The returned value should be properly released to the pool via Release() when no
// longer needed.
func (r *Request) Copy() *Request {
	req := fasthttp.AcquireRequest()
	r.CopyTo(req)

	return NewRequestFromFastHTTP(req)
}

// Release frees the resources associated with the Request object.
func (r *Request) Release() {
	fasthttp.ReleaseRequest(r.Request)
	r.mw = nil
}
