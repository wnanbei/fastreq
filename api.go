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

func Post(url string, opts ...ReqOption) (*Response, error) {
	return defaultClient.Post(url, opts...)
}

func Put(url string, opts ...ReqOption) (*Response, error) {
	return defaultClient.Put(url, opts...)
}

func Patch(url string, opts ...ReqOption) (*Response, error) {
	return defaultClient.Patch(url, opts...)
}

func Delete(url string, opts ...ReqOption) (*Response, error) {
	return defaultClient.Delete(url, opts...)
}

func Options(url string, opts ...ReqOption) (*Response, error) {
	return defaultClient.Options(url, opts...)
}

func Connect(url string, opts ...ReqOption) (*Response, error) {
	return defaultClient.Connect(url, opts...)
}

func Do(req *Request, opts ...ReqOption) (*Response, error) {
	return defaultClient.Do(req, opts...)
}

func DownloadFile(req *Request, path, filename string) error {
	return defaultClient.DownloadFile(req, path, filename)
}

func SetHTTPProxy(proxy string) {
	defaultClient.SetHTTPProxy(proxy)
}

func SetSocks5Proxy(proxy string) {
	defaultClient.SetSocks5Proxy(proxy)
}

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
// The default unmarshal is from the github.com/json-iterator/go package.
func SetJsonUnmarshal(unmarshalFunc func([]byte, any) error) {
	jsonUnmarshal = unmarshalFunc
}
