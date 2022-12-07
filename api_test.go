package fastreq

import (
	"net"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttputil"
)

func TestGet(t *testing.T) {
	t.Parallel()

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

func TestPost(t *testing.T) {
	t.Parallel()

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
