package fastreq

import (
	"mime/multipart"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttputil"
)

func Test_Request_Json_Body(t *testing.T) {
	ln := fasthttputil.NewInmemoryListener()
	s := &fasthttp.Server{
		Handler: func(ctx *fasthttp.RequestCtx) {
			require.Equal(t, "{\"Param1\":100,\"Param2\":\"hello world\"}", string(ctx.Request.Body()))
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

	body := map[string]interface{}{
		"Param1": 100,
		"Param2": "hello world",
	}

	resp, err := client.Post("http://make.fasthttp.great", NewJsonBody(body))
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

	resp, err := client.Post(
		"http://make.fasthttp.great",
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

	mf := NewMultipartForm(
		boundary,
		"foo", "bar",
	)
	resp, err := client.Post("http://make.fasthttp.great", mf)
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
			defer func(f multipart.File) {
				err := f.Close()
				if err != nil {
					t.FailNow()
				}
			}(f)
			buf := make([]byte, mf.File["txt"][0].Size)
			_, err = f.Read(buf)
			require.NoError(t, err)
			require.Equal(t, "fastreq", string(buf))

			f2, err := mf.File["file2"][0].Open()
			require.NoError(t, err)
			defer func(f2 multipart.File) {
				err := f2.Close()
				if err != nil {
					t.FailNow()
				}
			}(f2)
			buf = make([]byte, mf.File["file2"][0].Size)
			_, err = f2.Read(buf)
			require.NoError(t, err)
			require.Equal(t, "<p>fastreq</p>", string(buf))
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

	req := NewRequest(GET, "http://make.fasthttp.great")
	err := req.SetBoundary(boundary)
	require.NoError(t, err)
	err = req.AddMFField("foo", "bar")
	require.NoError(t, err)
	err = req.AddMFFile("txt", ".github/testdata/file1.txt")
	require.NoError(t, err)
	err = req.AddMFFile("", ".github/testdata/index.html")
	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)
	require.Equal(t, fasthttp.StatusOK, resp.StatusCode())
}

func Test_RequestOption_AutoRelease(t *testing.T) {
	ln := fasthttputil.NewInmemoryListener()
	s := &fasthttp.Server{
		Handler: func(ctx *fasthttp.RequestCtx) {
			require.Equal(t, fasthttp.MethodGet, string(ctx.Request.Header.Method()))
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

	params := NewQueryParams("hello", "world")
	form := NewPostForm("hello", "world")
	body := NewBody([]byte("hello world"))
	jsonBody := NewJsonBody(map[string]string{"hello": "world"})
	mf := NewMultipartForm("fastreq", "hello", "world")
	timeout := NewTimeout(time.Second)

	req := NewRequest(GET, "http://make.fasthttp.great")
	resp, err := client.Do(req, params, form, body, jsonBody, mf, timeout)
	require.NoError(t, err)
	require.Equal(t, fasthttp.StatusOK, resp.StatusCode())

	require.Nil(t, params.Args)
	require.Nil(t, form.Args)
	require.Empty(t, body.body)
	require.Empty(t, jsonBody.body)
	require.Empty(t, mf.Boundary)
	require.Nil(t, mf.Args)
	require.Zero(t, timeout.timeout)
}

func Test_Request_Headers(t *testing.T) {
	ln := fasthttputil.NewInmemoryListener()
	s := &fasthttp.Server{
		Handler: func(ctx *fasthttp.RequestCtx) {
			require.Equal(t, "fastreq", string(ctx.Request.Header.Peek("Hello")))
			require.Equal(t, "111", string(ctx.Request.Header.ContentType()))
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

	headers := NewHeaders("Content-Type", "111")
	headers.Add("Hello", "fastreq")
	resp, err := client.Post("http://make.fasthttp.great", headers)
	require.NoError(t, err)
	require.Equal(t, fasthttp.StatusOK, resp.StatusCode())
}

func Test_Request_Cookies(t *testing.T) {
	ln := fasthttputil.NewInmemoryListener()
	s := &fasthttp.Server{
		Handler: func(ctx *fasthttp.RequestCtx) {
			require.Equal(t, "fastreq", string(ctx.Request.Header.Cookie("name")))
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

	cookies := NewCookies("name", "fastreq")
	resp, err := client.Post("http://make.fasthttp.great", cookies)
	require.NoError(t, err)
	require.Equal(t, fasthttp.StatusOK, resp.StatusCode())
}
