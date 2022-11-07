package fastreq

import (
	"bufio"
	"encoding/json"
	"github.com/tidwall/gjson"
	"github.com/valyala/fasthttp"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
)

// Response ...
type Response struct {
	*fasthttp.Response
}

func NewResponse() *Response {
	return &Response{
		fasthttp.AcquireResponse(),
	}
}

func NewResponseFromFastHTTP(resp *fasthttp.Response) *Response {
	return &Response{
		resp,
	}
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

// ================================= Get Body ===================================

func (r *Response) Json(v interface{}) error {
	body, err := r.BodyUncompressed()
	if err != nil {
		return err
	}

	return json.Unmarshal(body, v)
}

func (r *Response) JsonGet(path string) gjson.Result {
	return gjson.GetBytes(r.Body(), path)
}

func (r *Response) JsonGetMany(path ...string) []gjson.Result {
	return gjson.GetManyBytes(r.Body(), path...)
}

// ================================= Get Body End ===============================

func (r *Response) Copy() *Response {
	resp := fasthttp.AcquireResponse()
	r.CopyTo(resp)

	return NewResponseFromFastHTTP(resp)
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
	fasthttp.ReleaseResponse(r.Response)
}
