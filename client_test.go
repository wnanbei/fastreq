package fastreq

import (
	"net"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttputil"
)

func TestClientGet(t *testing.T) {
	ln := fasthttputil.NewInmemoryListener()
	s := &fasthttp.Server{
		Handler: func(ctx *fasthttp.RequestCtx) {
			require.Equal(t, fasthttp.MethodGet, string(ctx.Request.Header.Method()))
			require.Equal(t, "hello=world&params=2", string(ctx.Request.URI().QueryString()))
		},
	}
	go func() {
		err := s.Serve(ln)
		if err != nil {
			return
		}
	}()

	client := NewClient()
	client.Dial = func(addr string) (net.Conn, error) {
		return ln.Dial()
	}

	params := NewQueryParams()
	params.Add("hello", "world")
	params.Add("params", "2")

	resp, err := client.Get("http://make.fasthttp.great", params)
	require.NoError(t, err)
	require.Equal(t, fasthttp.StatusOK, resp.StatusCode())
}

func TestClientPost(t *testing.T) {
	ln := fasthttputil.NewInmemoryListener()
	s := &fasthttp.Server{
		Handler: func(ctx *fasthttp.RequestCtx) {
			require.Equal(t, fasthttp.MethodPost, string(ctx.Request.Header.Method()))
			require.Equal(t, "hello=world&params=2", string(ctx.Request.Body()))
		},
	}
	go func() {
		err := s.Serve(ln)
		if err != nil {
			return
		}
	}()

	client := NewClient()
	client.Dial = func(addr string) (net.Conn, error) {
		return ln.Dial()
	}

	body := NewPostForm()
	body.Add("hello", "world")
	body.Add("params", "2")

	resp, err := client.Post("http://make.fasthttp.great", body)
	require.NoError(t, err)
	require.Equal(t, fasthttp.StatusOK, resp.StatusCode())
}

func TestClientDefaultUserAgent(t *testing.T) {
	ln := fasthttputil.NewInmemoryListener()
	s := &fasthttp.Server{
		Handler: func(ctx *fasthttp.RequestCtx) {
			_, err := ctx.Write(ctx.Request.Header.UserAgent())
			require.NoError(t, err)
		},
	}
	go func() {
		err := s.Serve(ln)
		if err != nil {
			return
		}
	}()

	client := NewClient()
	client.Dial = func(addr string) (net.Conn, error) {
		return ln.Dial()
	}

	resp, err := client.Get("http://make.fasthttp.great")
	require.NoError(t, err)
	require.Equal(t, fasthttp.StatusOK, resp.StatusCode())
	require.Equal(t, defaultUserAgent, resp.BodyString())

	resp, err = client.Get("http://make.fasthttp.great", NewHeaders("user-agent", "111"))
	require.NoError(t, err)
	require.Equal(t, fasthttp.StatusOK, resp.StatusCode())
	require.Equal(t, "111", resp.BodyString())

	client.SetDefaultUserAgent("hello world")
	resp, err = client.Get("http://make.fasthttp.great")
	require.NoError(t, err)
	require.Equal(t, fasthttp.StatusOK, resp.StatusCode())
	require.Equal(t, "hello world", resp.BodyString())

	resp, err = client.Get("http://make.fasthttp.great", NewHeaders("user-agent", "222"))
	require.NoError(t, err)
	require.Equal(t, fasthttp.StatusOK, resp.StatusCode())
	require.Equal(t, "222", resp.BodyString())
}
