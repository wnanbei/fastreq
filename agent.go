package fastreq

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2/utils"
	"github.com/valyala/fasthttp"
)

type Agent struct {
	Name                     string
	NoDefaultUserAgentHeader bool
	*fasthttp.HostClient
	req               *Request
	resp              *Response
	dest              []byte
	args              *Args
	timeout           time.Duration
	errs              []error
	formFiles         []*FormFile
	debugWriter       io.Writer
	mw                multipartWriter
	jsonEncoder       utils.JSONMarshal
	jsonDecoder       utils.JSONUnmarshal
	maxRedirectsCount int
	boundary          string
	reuse             bool
	parsed            bool
}

// Parse initializes URI and HostClient.
func (a *Agent) Parse() error {
	if a.parsed {
		return nil
	}
	a.parsed = true

	uri := a.req.URI()

	isTLS := false
	scheme := uri.Scheme()
	if bytes.Equal(scheme, strHTTPS) {
		isTLS = true
	} else if !bytes.Equal(scheme, strHTTP) {
		return fmt.Errorf("unsupported protocol %q. http and https are supported", scheme)
	}

	name := a.Name
	if name == "" && !a.NoDefaultUserAgentHeader {
		name = defaultUserAgent
	}

	a.HostClient = &fasthttp.HostClient{
		Addr:                     addMissingPort(string(uri.Host()), isTLS),
		Name:                     name,
		NoDefaultUserAgentHeader: a.NoDefaultUserAgentHeader,
		IsTLS:                    isTLS,
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
