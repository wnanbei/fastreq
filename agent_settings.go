package fastreq

import (
	"crypto/tls"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2/utils"
	"github.com/valyala/fasthttp"
)

func (a *Agent) Debug(w ...io.Writer) *Agent {
	a.debugWriter = os.Stdout
	if len(w) > 0 {
		a.debugWriter = w[0]
	}

	return a
}

// Timeout sets request timeout duration.
func (a *Agent) Timeout(timeout time.Duration) *Agent {
	a.timeout = timeout

	return a
}

// Reuse enables the Agent instance to be used again after one request.
//
// If agent is reusable, then it should be released manually when it is no
// longer used.
func (a *Agent) Reuse() *Agent {
	a.reuse = true

	return a
}

// InsecureSkipVerify controls whether the Agent verifies the server
// certificate chain and host name.
func (a *Agent) InsecureSkipVerify() *Agent {
	if a.Client.TLSConfig == nil {
		/* #nosec G402 */
		a.Client.TLSConfig = &tls.Config{InsecureSkipVerify: true} // #nosec G402
	} else {
		/* #nosec G402 */
		a.Client.TLSConfig.InsecureSkipVerify = true
	}

	return a
}

// TLSConfig sets tls config.
func (a *Agent) TLSConfig(config *tls.Config) *Agent {
	a.Client.TLSConfig = config

	return a
}

// MaxRedirectsCount sets max redirect count for GET and HEAD.
func (a *Agent) MaxRedirectsCount(count int) *Agent {
	a.maxRedirectsCount = count

	return a
}

// JSONEncoder sets custom json encoder.
func (a *Agent) JSONEncoder(jsonEncoder utils.JSONMarshal) *Agent {
	a.jsonEncoder = jsonEncoder

	return a
}

// JSONDecoder sets custom json decoder.
func (a *Agent) JSONDecoder(jsonDecoder utils.JSONUnmarshal) *Agent {
	a.jsonDecoder = jsonDecoder

	return a
}

// Request returns Agent request instance.
func (a *Agent) Request() *Request {
	return a.req
}

// SetResponse sets custom response for the Agent instance.
//
// It is recommended obtaining custom response via AcquireResponse and release it
// manually in performance-critical code.
func (a *Agent) SetResponse(customResp *Response) *Agent {
	a.resp = customResp

	return a
}

// Dest sets custom dest.
//
// The contents of dest will be replaced by the response body, if the dest
// is too small a new slice will be allocated.
func (a *Agent) Dest(dest []byte) *Agent {
	a.dest = dest

	return a
}

// RetryIf controls whether a retry should be attempted after an error.
//
// By default, will use isIdempotent function from fasthttp
func (a *Agent) RetryIf(retryIf fasthttp.RetryIfFunc) *Agent {
	a.Client.RetryIf = retryIf
	return a
}

/************************** End Agent Setting **************************/

// Bytes returns the status code, bytes body and errors of url.
//
// it's not safe to use Agent after calling [Agent.Bytes]
func (a *Agent) Bytes() (code int, body []byte, errs []error) {
	defer a.release()
	return a.bytes()
}

func (a *Agent) bytes() (code int, body []byte, errs []error) {
	if errs = append(errs, a.errs...); len(errs) > 0 {
		return
	}

	var (
		req     = a.req
		resp    *Response
		nilResp bool
	)

	if a.resp == nil {
		resp = AcquireResponse()
		nilResp = true
	} else {
		resp = a.resp
	}

	defer func() {
		if a.debugWriter != nil {
			printDebugInfo(req, resp, a.debugWriter)
		}

		if len(errs) == 0 {
			code = resp.StatusCode()
		}

		body = append(a.dest, resp.Body()...)

		if nilResp {
			ReleaseResponse(resp)
		}
	}()

	if a.timeout > 0 {
		if err := a.Client.DoTimeout(req, resp, a.timeout); err != nil {
			errs = append(errs, err)
			return
		}
	} else if a.maxRedirectsCount > 0 && (string(req.Header.Method()) == MethodGet || string(req.Header.Method()) == MethodHead) {
		if err := a.Client.DoRedirects(req, resp, a.maxRedirectsCount); err != nil {
			errs = append(errs, err)
			return
		}
	} else if err := a.Client.Do(req, resp); err != nil {
		errs = append(errs, err)
	}

	return
}

func printDebugInfo(req *Request, resp *Response, w io.Writer) {
	msg := fmt.Sprintf("Connected to %s(%s)\r\n\r\n", req.URI().Host(), resp.RemoteAddr())
	_, _ = w.Write(utils.UnsafeBytes(msg))
	_, _ = req.WriteTo(w)
	_, _ = resp.WriteTo(w)
}

// String returns the status code, string body and errors of url.
//
// it's not safe to use Agent after calling [Agent.String]
func (a *Agent) String() (int, string, []error) {
	defer a.release()
	code, body, errs := a.bytes()

	return code, utils.UnsafeString(body), errs
}

// Struct returns the status code, bytes body and errors of url.
// And bytes body will be unmarshalled to given v.
//
// it's not safe to use Agent after calling [Agent.Struct]
func (a *Agent) Struct(v interface{}) (code int, body []byte, errs []error) {
	defer a.release()
	if code, body, errs = a.bytes(); len(errs) > 0 {
		return
	}

	if err := a.jsonDecoder(body, v); err != nil {
		errs = append(errs, err)
	}

	return
}

func (a *Agent) release() {
	if !a.reuse {
		ReleaseAgent(a)
	} else {
		a.errs = a.errs[:0]
	}
}

func (a *Agent) reset() {
	a.Client = nil
	a.req.Reset()
	a.resp = nil
	a.dest = nil
	a.timeout = 0
	a.args = nil
	a.errs = a.errs[:0]
	a.debugWriter = nil
	a.mw = nil
	a.reuse = false
	a.parsed = false
	a.maxRedirectsCount = 0
	a.boundary = ""
	a.Name = ""
	a.NoDefaultUserAgentHeader = false
	for i, ff := range a.formFiles {
		if ff.autoRelease {
			ReleaseFormFile(ff)
		}
		a.formFiles[i] = nil
	}
	a.formFiles = a.formFiles[:0]
}

var (
	agentPool = sync.Pool{
		New: func() interface{} {
			return &Agent{req: &Request{}}
		},
	}
	responsePool sync.Pool
	argsPool     sync.Pool
	formFilePool sync.Pool
)

// AcquireAgent returns an empty Agent instance from Agent pool.
//
// The returned Agent instance may be passed to ReleaseAgent when it is
// no longer needed. This allows Agent recycling, reduces GC pressure
// and usually improves performance.
func AcquireAgent() *Agent {
	return agentPool.Get().(*Agent)
}

// ReleaseAgent returns a acquired via AcquireAgent to Agent pool.
//
// It is forbidden accessing req and/or its' members after returning
// it to Agent pool.
func ReleaseAgent(a *Agent) {
	a.reset()
	agentPool.Put(a)
}

// AcquireResponse returns an empty Response instance from response pool.
//
// The returned Response instance may be passed to ReleaseResponse when it is
// no longer needed. This allows Response recycling, reduces GC pressure
// and usually improves performance.
// Copy from fasthttp
func AcquireResponse() *Response {
	v := responsePool.Get()
	if v == nil {
		return &Response{}
	}
	return v.(*Response)
}

// ReleaseResponse return resp acquired via AcquireResponse to response pool.
//
// It is forbidden accessing resp and/or its' members after returning
// it to response pool.
// Copy from fasthttp
func ReleaseResponse(resp *Response) {
	resp.Reset()
	responsePool.Put(resp)
}

// AcquireArgs returns an empty Args object from the pool.
//
// The returned Args may be returned to the pool with ReleaseArgs
// when no longer needed. This allows reducing GC load.
// Copy from fasthttp
func AcquireArgs() *Args {
	v := argsPool.Get()
	if v == nil {
		return &Args{}
	}
	return v.(*Args)
}

// ReleaseArgs returns the object acquired via AcquireArgs to the pool.
//
// String not access the released Args object, otherwise data races may occur.
// Copy from fasthttp
func ReleaseArgs(a *Args) {
	a.Reset()
	argsPool.Put(a)
}

// AcquireFormFile returns an empty FormFile object from the pool.
//
// The returned FormFile may be returned to the pool with ReleaseFormFile
// when no longer needed. This allows reducing GC load.
func AcquireFormFile() *FormFile {
	v := formFilePool.Get()
	if v == nil {
		return &FormFile{}
	}
	return v.(*FormFile)
}

// ReleaseFormFile returns the object acquired via AcquireFormFile to the pool.
//
// String not access the released FormFile object, otherwise data races may occur.
func ReleaseFormFile(ff *FormFile) {
	ff.Fieldname = ""
	ff.Name = ""
	ff.Content = ff.Content[:0]
	ff.autoRelease = false

	formFilePool.Put(ff)
}

var (
	strHTTP          = []byte("http")
	strHTTPS         = []byte("https")
	defaultUserAgent = "fiber"
)

type multipartWriter interface {
	Boundary() string
	SetBoundary(boundary string) error
	CreateFormFile(fieldname, filename string) (io.Writer, error)
	WriteField(fieldname, value string) error
	Close() error
}
