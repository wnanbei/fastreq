package fastreq

import (
	"net"
	"testing"
	"time"

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
	ln := fasthttputil.NewInmemoryListener()
	s := &fasthttp.Server{
		Handler: func(ctx *fasthttp.RequestCtx) {
			require.Equal(t, "{\"Param1\":100,\"Param2\":\"hello world\"}", string(ctx.Request.Body()))
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
	ln := fasthttputil.NewInmemoryListener()
	s := &fasthttp.Server{
		Handler: func(ctx *fasthttp.RequestCtx) {
			require.Equal(t, "hello world", string(ctx.Request.Body()))
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

func Test_Request_Multipart_Form(t *testing.T) {
	boundary := "fastreq"

	ln := fasthttputil.NewInmemoryListener()
	s := &fasthttp.Server{
		Handler: func(ctx *fasthttp.RequestCtx) {
			require.Equal(t, "multipart/form-data; boundary="+boundary, string(ctx.Request.Header.ContentType()))
			mf, err := ctx.Request.MultipartForm()
			require.NoError(t, err)
			require.Equal(t, "bar", mf.Value["foo"][0])
		},
	}
	go s.Serve(ln) //nolint:errcheck

	client := NewClient()
	client.Dial = func(addr string) (net.Conn, error) {
		return ln.Dial()
	}

	mf := NewMultipartForm(
		boundary,
		"foo", "bar",
	)
	resp, err := client.Post("http://fastreq.com", mf)
	require.NoError(t, err)
	require.Equal(t, fasthttp.StatusOK, resp.StatusCode())
}

func Test_Request_Multipart_Form_Files(t *testing.T) {
	boundary := "fastreq"

	ln := fasthttputil.NewInmemoryListener()
	s := &fasthttp.Server{
		Handler: func(ctx *fasthttp.RequestCtx) {
			require.Equal(t, "multipart/form-data; boundary="+boundary, string(ctx.Request.Header.ContentType()))
			mf, err := ctx.Request.MultipartForm()
			require.NoError(t, err)
			require.Equal(t, "bar", mf.Value["foo"][0])
			require.Equal(t, "file1.txt", mf.File["txt"][0].Filename)

			f, err := mf.File["txt"][0].Open()
			require.NoError(t, err)
			defer f.Close()
			buf := make([]byte, mf.File["txt"][0].Size)
			_, err = f.Read(buf)
			require.NoError(t, err)
			require.Equal(t, "fastreq", string(buf))

			f2, err := mf.File["file2"][0].Open()
			require.NoError(t, err)
			defer f2.Close()
			buf = make([]byte, mf.File["file2"][0].Size)
			_, err = f2.Read(buf)
			require.NoError(t, err)
			require.Equal(t, "<p>fastreq</p>", string(buf))
		},
	}
	go s.Serve(ln) //nolint:errcheck

	client := NewClient()
	client.Dial = func(addr string) (net.Conn, error) {
		return ln.Dial()
	}

	req := NewRequest(GET, "http://fastreq.com")
	req.SetBoundary(boundary)
	req.AddMFField("foo", "bar")
	req.AddMFFile("txt", ".github/testdata/file1.txt")
	req.AddMFFile("", ".github/testdata/index.html")

	resp, err := client.Do(req)
	require.NoError(t, err)
	require.Equal(t, fasthttp.StatusOK, resp.StatusCode())
}

func Test_RequestOption_AutoRelease(t *testing.T) {
	params := NewQueryParams("hello", "world")
	form := NewPostForm("hello", "world")
	body := NewBody([]byte("hello world"))
	jsonBody := NewJsonBody(map[string]string{"hello": "world"})
	mf := NewMultipartForm("fastreq", "hello", "world")
	timeout := NewTimeout(time.Second)

	NewRequest(GET, "", params, form, body, jsonBody, mf, timeout)

	require.Nil(t, params.Args)
	require.Nil(t, form.Args)
	require.Empty(t, body.body)
	require.Empty(t, jsonBody.body)
	require.Empty(t, mf.Boundary)
	require.Nil(t, mf.Args)
	require.Zero(t, timeout.timeout)
}
