package fastreq

import (
	"net"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttputil"
)

func Benchmark_Header_Set(b *testing.B) {
	key := "user-agent"
	value := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.0.0 Safari/537.36 Edg/105.0.1343.27"

	req := fasthttp.AcquireRequest()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req.Header.Set(key, value)
	}
	fasthttp.ReleaseRequest(req)
}

func Benchmark_Header_Set_Bytes(b *testing.B) {
	key := []byte("user-agent")
	value := []byte("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.0.0 Safari/537.36 Edg/105.0.1343.27")

	req := fasthttp.AcquireRequest()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req.Header.SetBytesKV(key, value)
	}
	fasthttp.ReleaseRequest(req)
}

func Benchmark_Header_Set_Convert(b *testing.B) {
	key := []byte("user-agent")
	value := []byte("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.0.0 Safari/537.36 Edg/105.0.1343.27")

	req := fasthttp.AcquireRequest()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req.Header.Set(string(key), string(value))
	}
	fasthttp.ReleaseRequest(req)
}

type testStruct struct {
	Param1 int
	Param2 string
}

func Test_Request_Json_Body(t *testing.T) {
	t.Parallel()

	ln := fasthttputil.NewInmemoryListener()
	s := &fasthttp.Server{
		Handler: func(ctx *fasthttp.RequestCtx) {
			require.Equal(t, "{\"Param1\":100,\"Param2\":\"hello world\"}", string(ctx.Request.Body()))
			ctx.Response.SetStatusCode(fasthttp.StatusOK)
		},
	}
	go s.Serve(ln) //nolint:errcheck

	client := NewClient()
	client.Dial = func(addr string) (net.Conn, error) {
		return ln.Dial()
	}

	body := &testStruct{
		Param1: 100,
		Param2: "hello world",
	}

	resp, err := client.Post("http://fastreq.com", NewJsonBody(body))
	require.NoError(t, err)
	require.Equal(t, fasthttp.StatusOK, resp.StatusCode())
}

func Test_Request_Body(t *testing.T) {
	t.Parallel()

	ln := fasthttputil.NewInmemoryListener()
	s := &fasthttp.Server{
		Handler: func(ctx *fasthttp.RequestCtx) {
			require.Equal(t, "hello world", string(ctx.Request.Body()))
			ctx.Response.SetStatusCode(fasthttp.StatusOK)
		},
	}
	go s.Serve(ln) //nolint:errcheck

	client := NewClient()
	client.Dial = func(addr string) (net.Conn, error) {
		return ln.Dial()
	}

	resp, err := client.Post(
		"http://fastreq.com",
		NewBody([]byte("hello world")),
	)
	require.NoError(t, err)
	require.Equal(t, fasthttp.StatusOK, resp.StatusCode())
}
