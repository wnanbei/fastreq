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

// Response ...
type Response struct {
	*fasthttp.Response
	Request *fasthttp.Request
	Dom     *goquery.Document
}

func NewResponse() *Response {
	return &Response{
		Response: fasthttp.AcquireResponse(),
	}
}

func (r *Response) BodyString() string {
	return unsafeB2S(r.Body())
}

func (r *Response) Json(v interface{}) error {
	body, err := r.BodyUncompressed()
	if err != nil {
		return err
	}

	return jsonUnmarshal(body, v)
}

func (r *Response) JsonGet(path string) gjson.Result {
	return gjson.GetBytes(r.Body(), path)
}

func (r *Response) JsonGetMany(path ...string) []gjson.Result {
	return gjson.GetManyBytes(r.Body(), path...)
}

func (r *Response) JsonGetPartOf(path string, v interface{}) error {
	part := gjson.GetBytes(r.Body(), path)
	if part.Raw == "" {
		return nil
	}
	return jsonUnmarshal(unsafeS2B(part.Raw), v)
}

func (r *Response) BuildDom() error {
	if r.Dom == nil {
		doc, err := goquery.NewDocumentFromReader(bytes.NewReader(r.Body()))
		if err != nil {
			return err
		}
		r.Dom = doc
	}
	return nil
}

func (r *Response) Copy() *Response {
	resp := fasthttp.AcquireResponse()
	r.CopyTo(resp)

	return &Response{Response: resp}
}

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

func (r *Response) Release() {
	if r.Request != nil {
		fasthttp.ReleaseRequest(r.Request)
	}
	fasthttp.ReleaseResponse(r.Response)
}
