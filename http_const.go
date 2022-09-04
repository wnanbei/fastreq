package fastreq

const (
	MethodGet     = "GET"     // RFC 7231, 4.3.1
	MethodHead    = "HEAD"    // RFC 7231, 4.3.2
	MethodPost    = "POST"    // RFC 7231, 4.3.3
	MethodPut     = "PUT"     // RFC 7231, 4.3.4
	MethodPatch   = "PATCH"   // RFC 5789
	MethodDelete  = "DELETE"  // RFC 7231, 4.3.5
	MethodConnect = "CONNECT" // RFC 7231, 4.3.6
	MethodOptions = "OPTIONS" // RFC 7231, 4.3.7
	MethodTrace   = "TRACE"   // RFC 7231, 4.3.8
	methodUse     = "USE"
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
