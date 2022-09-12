package fastreq

import (
	"encoding/json"
	"io"
	"net"

	"github.com/valyala/fasthttp"
)

// ReResponse ...
type ReResponse struct {
	resp *fasthttp.Response
}

func NewResponse() *ReResponse {
	return &ReResponse{
		fasthttp.AcquireResponse(),
	}
}

func NewResponseFromFastHTTP(resp *fasthttp.Response) *ReResponse {
	return &ReResponse{
		resp,
	}
}

func ReleaseReResponse(r *ReResponse) {
	fasthttp.ReleaseResponse(r.resp)
}

func (r *ReResponse) StatusCode() int {
	return r.resp.StatusCode()
}

func (r *ReResponse) RemoteAddr() net.Addr {
	return r.resp.RemoteAddr()
}

// ================================= Get Body ===================================

func (r *ReResponse) Body() []byte {
	return r.resp.Body()
}

func (r *ReResponse) BodyString() string {
	return r.resp.String()
}

func (r *ReResponse) BodyGunzip() ([]byte, error) {
	return r.resp.BodyGunzip()
}

func (r *ReResponse) BodyUncompressed() ([]byte, error) {
	return r.resp.BodyUncompressed()
}

func (r *ReResponse) BodyWriteTo(w io.Writer) error {
	return r.resp.BodyWriteTo(w)
}

func (r *ReResponse) Json(v interface{}) error {
	body, err := r.resp.BodyUncompressed()
	if err != nil {
		return err
	}

	return json.Unmarshal(body, v)
}

// ================================= Get Body End ===============================

func (r *ReResponse) Copy() *ReResponse {
	resp := fasthttp.AcquireResponse()
	r.resp.CopyTo(resp)

	return NewResponseFromFastHTTP(resp)
}
