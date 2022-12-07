package fastreq

import (
	"net"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttputil"
)

func Test_Method_Get(t *testing.T) {
	ln := fasthttputil.NewInmemoryListener()
	s := &fasthttp.Server{
		Handler: func(ctx *fasthttp.RequestCtx) {
			require.Equal(t, fasthttp.MethodGet, string(ctx.Request.Header.Method()))
			require.Equal(t, "hello=world&params=2", string(ctx.Request.URI().QueryString()))
		},
	}
	go s.Serve(ln) //nolint:errcheck

	client := NewClient()
	client.Dial = func(addr string) (net.Conn, error) {
		return ln.Dial()
	}
	SetDefaultClient(client)

	params := NewQueryParams()
	params.Add("hello", "world")
	params.Add("params", "2")

	resp, err := Get("http://fasreq.com", params)
	require.NoError(t, err)
	require.Equal(t, fasthttp.StatusOK, resp.StatusCode())
	Release(resp)
}

func Test_Method_Post(t *testing.T) {
	ln := fasthttputil.NewInmemoryListener()
	s := &fasthttp.Server{
		Handler: func(ctx *fasthttp.RequestCtx) {
			require.Equal(t, fasthttp.MethodPost, string(ctx.Request.Header.Method()))
			require.Equal(t, "hello=world&params=2", string(ctx.Request.Body()))
		},
	}
	go s.Serve(ln) //nolint:errcheck

	client := NewClient()
	client.Dial = func(addr string) (net.Conn, error) {
		return ln.Dial()
	}
	SetDefaultClient(client)

	formBody := NewPostForm()
	formBody.Add("hello", "world")
	formBody.Add("params", "2")

	resp, err := Post("http://httpbin.org/post", formBody)
	require.NoError(t, err)
	require.Equal(t, fasthttp.StatusOK, resp.StatusCode())
	Release(resp)
}

func Test_Method_Head(t *testing.T) {
	t.Parallel()

	ln := fasthttputil.NewInmemoryListener()
	s := &fasthttp.Server{
		Handler: func(ctx *fasthttp.RequestCtx) {
			require.Equal(t, fasthttp.MethodHead, string(ctx.Request.Header.Method()))
		},
	}
	go s.Serve(ln) //nolint:errcheck

	client := NewClient()
	client.Dial = func(addr string) (net.Conn, error) {
		return ln.Dial()
	}
	SetDefaultClient(client)

	resp, err := Head("http://httpbin.org/post")
	require.NoError(t, err)
	require.Equal(t, fasthttp.StatusOK, resp.StatusCode())
	Release(resp)
}

func Test_Method_Put(t *testing.T) {
	ln := fasthttputil.NewInmemoryListener()
	s := &fasthttp.Server{
		Handler: func(ctx *fasthttp.RequestCtx) {
			require.Equal(t, fasthttp.MethodPut, string(ctx.Request.Header.Method()))
		},
	}
	go s.Serve(ln) //nolint:errcheck

	client := NewClient()
	client.Dial = func(addr string) (net.Conn, error) {
		return ln.Dial()
	}
	SetDefaultClient(client)

	resp, err := Put("http://httpbin.org/post")
	require.NoError(t, err)
	require.Equal(t, fasthttp.StatusOK, resp.StatusCode())
	Release(resp)
}

func Test_Method_Patch(t *testing.T) {
	ln := fasthttputil.NewInmemoryListener()
	s := &fasthttp.Server{
		Handler: func(ctx *fasthttp.RequestCtx) {
			require.Equal(t, fasthttp.MethodPatch, string(ctx.Request.Header.Method()))
		},
	}
	go s.Serve(ln) //nolint:errcheck

	client := NewClient()
	client.Dial = func(addr string) (net.Conn, error) {
		return ln.Dial()
	}
	SetDefaultClient(client)

	resp, err := Patch("http://httpbin.org/post")
	require.NoError(t, err)
	require.Equal(t, fasthttp.StatusOK, resp.StatusCode())
	Release(resp)
}

func Test_Method_Delete(t *testing.T) {
	ln := fasthttputil.NewInmemoryListener()
	s := &fasthttp.Server{
		Handler: func(ctx *fasthttp.RequestCtx) {
			require.Equal(t, fasthttp.MethodDelete, string(ctx.Request.Header.Method()))
		},
	}
	go s.Serve(ln) //nolint:errcheck

	client := NewClient()
	client.Dial = func(addr string) (net.Conn, error) {
		return ln.Dial()
	}
	SetDefaultClient(client)

	resp, err := Delete("http://httpbin.org/post")
	require.NoError(t, err)
	require.Equal(t, fasthttp.StatusOK, resp.StatusCode())
	Release(resp)
}

func Test_Method_Connect(t *testing.T) {
	ln := fasthttputil.NewInmemoryListener()
	s := &fasthttp.Server{
		Handler: func(ctx *fasthttp.RequestCtx) {
			require.Equal(t, fasthttp.MethodConnect, string(ctx.Request.Header.Method()))
		},
	}
	go s.Serve(ln) //nolint:errcheck

	client := NewClient()
	client.Dial = func(addr string) (net.Conn, error) {
		return ln.Dial()
	}
	SetDefaultClient(client)

	resp, err := Connect("http://httpbin.org/post")
	require.NoError(t, err)
	require.Equal(t, fasthttp.StatusOK, resp.StatusCode())
	Release(resp)
}

func Test_Method_Options(t *testing.T) {
	ln := fasthttputil.NewInmemoryListener()
	s := &fasthttp.Server{
		Handler: func(ctx *fasthttp.RequestCtx) {
			require.Equal(t, fasthttp.MethodOptions, string(ctx.Request.Header.Method()))
		},
	}
	go s.Serve(ln) //nolint:errcheck

	client := NewClient()
	client.Dial = func(addr string) (net.Conn, error) {
		return ln.Dial()
	}
	SetDefaultClient(client)

	resp, err := Options("http://httpbin.org/post")
	require.NoError(t, err)
	require.Equal(t, fasthttp.StatusOK, resp.StatusCode())
	Release(resp)
}

func Test_Method_Do(t *testing.T) {
	ln := fasthttputil.NewInmemoryListener()
	s := &fasthttp.Server{
		Handler: func(ctx *fasthttp.RequestCtx) {
			require.Equal(t, fasthttp.MethodOptions, string(ctx.Request.Header.Method()))
		},
	}
	go s.Serve(ln) //nolint:errcheck

	client := NewClient()
	client.Dial = func(addr string) (net.Conn, error) {
		return ln.Dial()
	}
	SetDefaultClient(client)

	req := NewRequest(OPTIONS, "http://httpbin.org/post")
	resp, err := Do(req)
	require.NoError(t, err)
	require.Equal(t, fasthttp.StatusOK, resp.StatusCode())
	Release(resp)
}
