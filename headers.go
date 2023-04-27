package fastreq

import (
	"sync"

	"github.com/valyala/fasthttp"
)

var headersPool sync.Pool

// Headers is a wrapper for fasthttp.RequestHeader
type Headers struct {
	headers        fasthttp.RequestHeader
	notAutoRelease bool
}

// NewHeaders creates a new Headers
func NewHeaders(kv ...string) *Headers {
	var h *Headers
	v := headersPool.Get()
	if v != nil {
		h = v.(*Headers)
	} else {
		h = &Headers{}
	}

	for i := 1; i < len(kv); i += 2 {
		h.headers.Add(kv[i-1], kv[i])
	}

	return h
}

// BindRequest binds the headers to the given request
func (h *Headers) BindRequest(req *Request) error {
	h.headers.VisitAll(func(key, value []byte) {
		req.Request.Header.AddBytesKV(key, value)
	})
	return nil
}

// Release frees the resources held by header
func (h *Headers) Release() {
	h.headers.Reset()
	h.notAutoRelease = false
	headersPool.Put(h)
}

// AutoRelease sets whether headers should be automatically released when the
// associated object is destroyed.
func (h *Headers) AutoRelease(auto bool) {
	h.notAutoRelease = !auto
}

// isAutoRelease returns true if the Headers instance is set to auto-release.
// In other words, it returns true if the `notAutoRelease` field is false.
func (h *Headers) isAutoRelease() bool {
	return !h.notAutoRelease
}

// Add adds the given 'key: value' header.
//
// Multiple headers with the same key may be added with this function.
// Use Set for setting a single header for the given key.
func (h *Headers) Add(key, value string) {
	h.headers.Add(key, value)
}

// AddBytesK adds the given 'key: value' header.
//
// Multiple headers with the same key may be added with this function.
// Use SetBytesK for setting a single header for the given key.
func (h *Headers) AddBytesK(key []byte, value string) {
	h.headers.AddBytesK(key, value)
}

// AddBytesV adds the given 'key: value' header.
//
// Multiple headers with the same key may be added with this function.
// Use SetBytesV for setting a single header for the given key.
func (h *Headers) AddBytesV(key string, value []byte) {
	h.headers.AddBytesV(key, value)
}

// AddBytesKV adds the given 'key: value' header.
//
// Multiple headers with the same key may be added with this function.
// Use SetBytesKV for setting a single header for the given key.
//
// the Content-Type, Content-Length, Connection, Cookie,
// Transfer-Encoding, Host and User-Agent headers can only be set once
// and will overwrite the previous value.
func (h *Headers) AddBytesKV(key, value []byte) {
	h.headers.AddBytesKV(key, value)
}

// Set sets the given 'key: value' header.
//
// Use Add for setting multiple header values under the same key.
func (h *Headers) Set(key, value string) {
	h.headers.Set(key, value)
}

// SetBytesK sets the given 'key: value' header.
//
// Use AddBytesK for setting multiple header values under the same key.
func (h *Headers) SetBytesK(key []byte, value string) {
	h.headers.SetBytesK(key, value)
}

// SetBytesV sets the given 'key: value' header.
//
// Use AddBytesV for setting multiple header values under the same key.
func (h *Headers) SetBytesV(key string, value []byte) {
	h.headers.SetBytesV(key, value)
}

// SetBytesKV sets the given 'key: value' header.
//
// Use AddBytesKV for setting multiple header values under the same key.
func (h *Headers) SetBytesKV(key, value []byte) {
	h.headers.SetBytesKV(key, value)
}

// Del deletes header with the given key.
func (h *Headers) Del(key string) {
	h.headers.Del(key)
}

// DelBytes deletes header with the given key.
func (h *Headers) DelBytes(key []byte) {
	h.headers.DelBytes(key)
}

// Peek returns header value for the given key.
//
// The returned value is valid until the request is released,
// either though ReleaseRequest or your request handler returning.
// Do not store references to returned value. Make copies instead.
func (h *Headers) Peek(key string) []byte {
	return h.headers.Peek(key)
}

// PeekBytes returns header value for the given key.
//
// The returned value is valid until the request is released,
// either though ReleaseRequest or your request handler returning.
// Do not store references to returned value. Make copies instead.
func (h *Headers) PeekBytes(key []byte) []byte {
	return h.headers.PeekBytes(key)
}

// PeekAll returns all header value for the given key.
//
// The returned value is valid until the request is released,
// either though ReleaseRequest or your request handler returning.
// Any future calls to the Peek* will modify the returned value.
// Do not store references to returned value. Make copies instead.
func (h *Headers) PeekAll(key string) [][]byte {
	return h.headers.PeekAll(key)
}

// VisitAll calls f for each header.
//
// f must not retain references to key and/or value after returning.
// Copy key and/or value contents before returning if you need retaining them.
//
// To get the headers in order they were received use VisitAllInOrder.
func (h *Headers) VisitAll(f func(key, value []byte)) {
	h.headers.VisitAll(f)
}
