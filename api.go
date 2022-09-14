package fastreq

var defaultClient = NewClient()

func Get(url string) (*Response, error) {
	return defaultClient.Get(url)
}

func Head(url string) (*Response, error) {
	return defaultClient.Head(url)
}

func Post(url string) (*Response, error) {
	return defaultClient.Post(url)
}

func Put(url string) (*Response, error) {
	return defaultClient.Put(url)
}

func Patch(url string) (*Response, error) {
	return defaultClient.Patch(url)
}

func Delete(url string) (*Response, error) {
	return defaultClient.Delete(url)
}

type Releaser interface {
	Release()
}

func Release(releasers ...Releaser) {
	for _, r := range releasers {
		r.Release()
	}
}
