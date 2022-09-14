package fastreq

import (
	"encoding/json"
	"io"
	"net"

	"github.com/valyala/fasthttp"
)

// Response ...
type Response struct {
	resp *fasthttp.Response
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

func ReleaseResponse(r *Response) {
	fasthttp.ReleaseResponse(r.resp)
}

func (r *Response) StatusCode() int {
	return r.resp.StatusCode()
}

func (r *Response) RemoteAddr() net.Addr {
	return r.resp.RemoteAddr()
}

// ================================= Get Body ===================================

func (r *Response) Body() []byte {
	return r.resp.Body()
}

func (r *Response) BodyString() string {
	return r.resp.String()
}

func (r *Response) BodyGunzip() ([]byte, error) {
	return r.resp.BodyGunzip()
}

func (r *Response) BodyUncompressed() ([]byte, error) {
	return r.resp.BodyUncompressed()
}

func (r *Response) BodyWriteTo(w io.Writer) error {
	return r.resp.BodyWriteTo(w)
}

func (r *Response) Json(v interface{}) error {
	body, err := r.resp.BodyUncompressed()
	if err != nil {
		return err
	}

	return json.Unmarshal(body, v)
}

// ================================= Get Body End ===============================

func (r *Response) Copy() *Response {
	resp := fasthttp.AcquireResponse()
	r.resp.CopyTo(resp)

	return NewResponseFromFastHTTP(resp)
}
