package fastreq

import "time"

// Version of fastreq
const Version = "0.2.1"

const defaultUserAgent = "fastreq@" + Version

// HTTPMethod http request method
type HTTPMethod string

const (
	GET     HTTPMethod = "GET"     // RFC 7231, 4.3.1
	HEAD    HTTPMethod = "HEAD"    // RFC 7231, 4.3.2
	POST    HTTPMethod = "POST"    // RFC 7231, 4.3.3
	PUT     HTTPMethod = "PUT"     // RFC 7231, 4.3.4
	PATCH   HTTPMethod = "PATCH"   // RFC 5789
	DELETE  HTTPMethod = "DELETE"  // RFC 7231, 4.3.5
	CONNECT HTTPMethod = "CONNECT" // RFC 7231, 4.3.6
	OPTIONS HTTPMethod = "OPTIONS" // RFC 7231, 4.3.7
)

const (
	MIMETextXML                          = "text/xml"
	MIMETextHTML                         = "text/html"
	MIMETextPlain                        = "text/plain"
	MIMEApplicationXML                   = "application/xml"
	MIMEApplicationJSON                  = "application/json"
	MIMEApplicationJavaScript            = "application/javascript"
	MIMEApplicationForm                  = "application/x-www-form-urlencoded"
	MIMEOctetStream                      = "application/octet-stream"
	MIMEMultipartForm                    = "multipart/form-data"
	MIMETextXMLCharsetUTF8               = "text/xml; charset=utf-8"
	MIMETextHTMLCharsetUTF8              = "text/html; charset=utf-8"
	MIMETextPlainCharsetUTF8             = "text/plain; charset=utf-8"
	MIMEApplicationXMLCharsetUTF8        = "application/xml; charset=utf-8"
	MIMEApplicationJSONCharsetUTF8       = "application/json; charset=utf-8"
	MIMEApplicationJavaScriptCharsetUTF8 = "application/javascript; charset=utf-8"
)

const defaultTimeout = time.Second * 30

// DebugLevel debug log level
type DebugLevel int

const (
	DebugClose DebugLevel = iota // close debug
	DebugSimple
	DebugDetail
)

// debugLimit limit length of debug logging output
const debugLimit = 10000
