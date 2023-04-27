package fastreq

import jsoniter "github.com/json-iterator/go"

var defaultClient = NewClient(&defaultClientConfig)

// SetDefaultClient sets the default client to be used for subsequent requests.
func SetDefaultClient(client *Client) {
	defaultClient = client
}

// Get performs an HTTP GET request to the specified URL with optional request options.
func Get(url string, opts ...ReqOption) (*Response, error) {
	return defaultClient.Get(url, opts...)
}

// Head sends a HEAD request to the specified URL and returns the response.
func Head(url string, opts ...ReqOption) (*Response, error) {
	return defaultClient.Head(url, opts...)
}

// Post sends an HTTP POST request to the specified URL with the provided request options.
func Post(url string, opts ...ReqOption) (*Response, error) {
	return defaultClient.Post(url, opts...)
}

// Put performs an HTTP PUT request to the given URL with the given request options.
func Put(url string, opts ...ReqOption) (*Response, error) {
	return defaultClient.Put(url, opts...)
}

// Patch performs an HTTP PATCH request to the given URL with the given request options.
func Patch(url string, opts ...ReqOption) (*Response, error) {
	return defaultClient.Patch(url, opts...)
}

// Delete performs an HTTP DELETE request to the given URL with the given request options.
func Delete(url string, opts ...ReqOption) (*Response, error) {
	return defaultClient.Delete(url, opts...)
}

// Options performs an HTTP OPTIONS request to the given URL with the given request options.
func Options(url string, opts ...ReqOption) (*Response, error) {
	return defaultClient.Options(url, opts...)
}

// Connect performs an HTTP CONNECT request to the given URL with the given request options.
func Connect(url string, opts ...ReqOption) (*Response, error) {
	return defaultClient.Connect(url, opts...)
}

// Do performs an HTTP request to the given URL with the given request options.
func Do(req *Request, opts ...ReqOption) (*Response, error) {
	return defaultClient.Do(req, opts...)
}

// DownloadFile downloads a file from the given path and filename using the default client.
func DownloadFile(req *Request, path, filename string) error {
	return defaultClient.DownloadFile(req, path, filename)
}

// SetHTTPProxy sets the HTTP proxy for the default client.
func SetHTTPProxy(proxy string) {
	defaultClient.SetHTTPProxy(proxy)
}

// SetSocks5Proxy sets the SOCKS5 proxy for the default client.
func SetSocks5Proxy(proxy string) {
	defaultClient.SetSocks5Proxy(proxy)
}

// SetEnvHTTPProxy using the env(HTTP_PROXY, HTTPS_PROXY and NO_PROXY) configured HTTP proxy for the default client.
func SetEnvHTTPProxy() {
	defaultClient.SetEnvHTTPProxy()
}

// SetOAuth1 sets the OAuth1 token for the default client.
func SetOauth1(o *Oauth1) {
	defaultClient.SetOauth1(o)
}

// SetDebugLevel sets the debug output level for the default client.
func SetDebugLevel(lvl DebugLevel) {
	defaultClient.SetDebugLevel(lvl)
}

var jsonMarshal = jsoniter.ConfigCompatibleWithStandardLibrary.Marshal
var jsonUnmarshal = jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal

// SetJsonMarshal sets the JSON marshal function.
// By default, github.com/json-iterator/go is used.
func SetJsonMarshal(marshalFunc func(any) ([]byte, error)) {
	jsonMarshal = marshalFunc
}

// SetJSONUnmarshal sets the JSON unmarshal function.
// By default, github.com/json-iterator/go is used.
func SetJsonUnmarshal(unmarshalFunc func([]byte, any) error) {
	jsonUnmarshal = unmarshalFunc
}
