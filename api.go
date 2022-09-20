package fastreq

var defaultClient = NewClient()

func Get(url string, params *Args) (*Response, error) {
	return defaultClient.Get(url, params)
}

func Head(url string, params *Args) (*Response, error) {
	return defaultClient.Head(url, params)
}

func Post(url string, body *Args) (*Response, error) {
	return defaultClient.Post(url, body)
}

func Put(url string, body *Args) (*Response, error) {
	return defaultClient.Put(url, body)
}

func Patch(url string, params *Args) (*Response, error) {
	return defaultClient.Patch(url, params)
}

func Delete(url string, params *Args) (*Response, error) {
	return defaultClient.Delete(url, params)
}

func Do(req *Request) (*Response, error) {
	return defaultClient.Do(req)
}

type Releaser interface {
	Release()
}

func Release(releasers ...Releaser) {
	for _, r := range releasers {
		r.Release()
	}
}
