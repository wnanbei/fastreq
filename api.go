package fastreq

import jsoniter "github.com/json-iterator/go"

var defaultClient = NewClient()

func Get(url string, params *Args) (*Ctx, error) {
	return defaultClient.Get(url, params)
}

func Head(url string, params *Args) (*Ctx, error) {
	return defaultClient.Head(url, params)
}

func Post(url string, body *Args) (*Ctx, error) {
	return defaultClient.Post(url, body)
}

func Put(url string, body *Args) (*Ctx, error) {
	return defaultClient.Put(url, body)
}

func Patch(url string, params *Args) (*Ctx, error) {
	return defaultClient.Patch(url, params)
}

func Delete(url string, params *Args) (*Ctx, error) {
	return defaultClient.Delete(url, params)
}

func Do(req *Request) (*Ctx, error) {
	return defaultClient.Do(req)
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
