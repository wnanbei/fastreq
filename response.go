package fastreq

import (
	"bufio"
	"bytes"
	"net/url"
	"os"
	"path/filepath"
	"regexp"

	"github.com/PuerkitoBio/goquery"
	"github.com/tidwall/gjson"
	"github.com/valyala/fasthttp"
)

// Response represents an HTTP response.
type Response struct {
	*fasthttp.Response
	Request *fasthttp.Request
	dom     *goquery.Document
}

// NewResponse initializes and returns a new Response object.
func NewResponse() *Response {
	return &Response{
		Response: fasthttp.AcquireResponse(),
	}
}

// BodyString returns the response body as a string.
func (r *Response) BodyString() string {
	return unsafeB2S(r.Body())
}

// Json decodes the body of the HTTP response and stores the result in the
// variable pointed to by v. The response body is assumed to be in JSON format.
func (r *Response) Json(v interface{}) error {
	body, err := r.BodyUncompressed()
	if err != nil {
		return err
	}

	return jsonUnmarshal(body, v)
}

// JsonGet retrieves a JSON value from the response body at the given JSON pointer path
// and returns it as a gjson.Result object.
func (r *Response) JsonGet(path string) gjson.Result {
	return gjson.GetBytes(r.Body(), path)
}

// JsonGetMany retrieves a JSON array from the response body at the given JSON pointer path
// and returns it as a slice of gjson.Result objects.
func (r *Response) JsonGetMany(path ...string) []gjson.Result {
	return gjson.GetManyBytes(r.Body(), path...)
}

// JsonGetPartOf returns a part of the JSON response body at the given path.
// If the path is not found or the part is empty, it returns nil.
// The function takes in a string path and an interface{} v, which can be any type that can be unmarshaled into JSON.
// It returns an error if there is a problem with unmarshaling the JSON.
func (r *Response) JsonGetPartOf(path string, v interface{}) error {
	part := gjson.GetBytes(r.Body(), path)
	if part.Raw == "" {
		return nil
	}
	return jsonUnmarshal(unsafeS2B(part.Raw), v)
}

// Dom returns a parsed HTML document from the response body.
// If the response body is not a valid HTML document, an error is returned.
// github.com/PuerkitoBio/goquery is used to parse the HTML document.
func (r *Response) Dom() (*goquery.Document, error) {
	if r.dom == nil {
		doc, err := goquery.NewDocumentFromReader(bytes.NewReader(r.Body()))
		if err != nil {
			return nil, err
		}
		r.dom = doc
	}

	return r.dom, nil
}

// Copy creates a new Response instance that is a copy of the current one.
func (r *Response) Copy() *Response {
	resp := fasthttp.AcquireResponse()
	r.CopyTo(resp)

	return &Response{Response: resp}
}

// FileName extracts the filename from the Content-Disposition header in the
// response. If there is no such header or it does not contain a filename, an
// empty string is returned.
func (r *Response) FileName() string {
	disposition := r.Header.Peek("Content-Disposition")
	if len(disposition) == 0 {
		return ""
	}

	matches := regexp.MustCompile(`filename[^;=\n]*=(['"]*.*?['"]*)$`).FindSubmatch(disposition)
	if len(matches) == 0 {
		return ""
	}
	n := regexp.MustCompile(`['"]`).ReplaceAll(matches[1], []byte{})

	un, err := url.QueryUnescape(unsafeB2S(n))
	if err != nil {
		return unsafeB2S(n)
	}

	return un
}

// Saves the response body to a file at the given path with the given filename.
func (r *Response) SaveToFile(path, filename string) error {
	if filename == "" {
		filename = r.FileName()
	}

	file, err := os.OpenFile(filepath.Join(path, filename), os.O_WRONLY|os.O_CREATE, 0664)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	if err := r.BodyWriteTo(w); err != nil {
		return err
	}

	if err := w.Flush(); err != nil {
		return err
	}

	return nil
}

// Release releases the resources associated with the Response by releasing both
// the Request and Response objects.
func (r *Response) Release() {
	if r.Request != nil {
		fasthttp.ReleaseRequest(r.Request)
	}
	fasthttp.ReleaseResponse(r.Response)
}
