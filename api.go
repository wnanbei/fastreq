package fastreq

var defaultClient ReClient

func Get(url string) (*Response, error)    { return defaultClient.Get(url) }
func Head(url string) (*Response, error)   { return defaultClient.Head(url) }
func Post(url string) (*Response, error)   { return defaultClient.Post(url) }
func Put(url string) (*Response, error)    { return defaultClient.Put(url) }
func Patch(url string) (*Response, error)  { return defaultClient.Patch(url) }
func Delete(url string) (*Response, error) { return defaultClient.Delete(url) }
