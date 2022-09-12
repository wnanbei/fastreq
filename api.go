package fastreq

var defaultClient ReClient

func Get(url string) (*ReResponse, error)    { return defaultClient.Get(url) }
func Head(url string) (*ReResponse, error)   { return defaultClient.Head(url) }
func Post(url string) (*ReResponse, error)   { return defaultClient.Post(url) }
func Put(url string) (*ReResponse, error)    { return defaultClient.Put(url) }
func Patch(url string) (*ReResponse, error)  { return defaultClient.Patch(url) }
func Delete(url string) (*ReResponse, error) { return defaultClient.Delete(url) }
