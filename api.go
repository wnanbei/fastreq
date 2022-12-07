package fastreq

import jsoniter "github.com/json-iterator/go"

var defaultClient = NewClient()

func SetDefaultClient(client *Client) {
	defaultClient = client
}

func Get(url string, opts ...ReqOption) (*Response, error) {
	return defaultClient.Get(url, opts...)
}

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

func SetOauth1(o *Oauth1) {
	defaultClient.SetOauth1(o)
}

func SetDebugLevel(lvl DebugLevel) {
	defaultClient.SetDebugLevel(lvl)
}

type Releaser interface {
	Release()
}

func Release(releasers ...Releaser) {
	for _, r := range releasers {
		r.Release()
	}
}

var jsonMarshal = jsoniter.ConfigCompatibleWithStandardLibrary.Marshal
var jsonUnmarshal = jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal

// SetJsonMarshal can set json marshal function.
// Default Marshal is github.com/json-iterator/go
func SetJsonMarshal(f func(any) ([]byte, error)) {
	jsonMarshal = f
}

// SetJsonUnmarshal can set json unmarshal function.
// Default Unarshal is github.com/json-iterator/go
func SetJsonUnmarshal(f func([]byte, any) error) {
	jsonUnmarshal = f
}
