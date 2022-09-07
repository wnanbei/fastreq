package fastreq

var defaultClient ReClient

func Get(url string) *ReqResp    { return defaultClient.Get(url) }
func Head(url string) *ReqResp   { return defaultClient.Head(url) }
func Post(url string) *ReqResp   { return defaultClient.Post(url) }
func Put(url string) *ReqResp    { return defaultClient.Put(url) }
func Patch(url string) *ReqResp  { return defaultClient.Patch(url) }
func Delete(url string) *ReqResp { return defaultClient.Delete(url) }
