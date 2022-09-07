package fastreq

import (
	"testing"

	"github.com/valyala/fasthttp"
)

func BenchmarkHeaderSet(b *testing.B) {
	key := "user-agent"
	value := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.0.0 Safari/537.36 Edg/105.0.1343.27"

	req := fasthttp.AcquireRequest()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req.Header.Set(key, value)
	}
	fasthttp.ReleaseRequest(req)
}

func BenchmarkHeaderSetBytes(b *testing.B) {
	key := []byte("user-agent")
	value := []byte("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.0.0 Safari/537.36 Edg/105.0.1343.27")

	req := fasthttp.AcquireRequest()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req.Header.SetBytesKV(key, value)
	}
	fasthttp.ReleaseRequest(req)
}

func BenchmarkHeaderSetConvert(b *testing.B) {
	key := []byte("user-agent")
	value := []byte("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.0.0 Safari/537.36 Edg/105.0.1343.27")

	req := fasthttp.AcquireRequest()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req.Header.Set(string(key), string(value))
	}
	fasthttp.ReleaseRequest(req)
}
