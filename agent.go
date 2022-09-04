package fastreq

import (
	"io"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2/utils"
	"github.com/valyala/fasthttp"
)

type Agent struct {
	*fasthttp.Client

	Name                     string
	NoDefaultUserAgentHeader bool
	req                      *Request
	resp                     *Response
	dest                     []byte
	args                     *Args
	timeout                  time.Duration
	errs                     []error
	formFiles                []*FormFile
	debugWriter              io.Writer
	mw                       multipartWriter
	jsonEncoder              utils.JSONMarshal
	jsonDecoder              utils.JSONUnmarshal
	maxRedirectsCount        int
	boundary                 string
	reuse                    bool
	parsed                   bool
}

func (a *Agent) Parse() error {
	if a.parsed {
		return nil
	}
	a.parsed = true

	name := a.Name
	if name == "" && !a.NoDefaultUserAgentHeader {
		name = defaultUserAgent
	}

	a.Client = &fasthttp.Client{
		Name:                     name,
		NoDefaultUserAgentHeader: a.NoDefaultUserAgentHeader,
	}

	return nil
}

func addMissingPort(addr string, isTLS bool) string {
	n := strings.Index(addr, ":")
	if n >= 0 {
		return addr
	}
	port := 80
	if isTLS {
		port = 443
	}
	return net.JoinHostPort(addr, strconv.Itoa(port))
}
