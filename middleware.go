package fastreq

import "fmt"

type RequestMiddleware func(req *Request)

type ResponseMiddleware func(resp *Response)

func MiddlewareOauth1(o *Oauth1) RequestMiddleware {
	return func(req *Request) {
		auth := o.GenHeader(req)
		fmt.Println(string(auth))
		req.req.Header.SetBytesV("Authorization", auth)
	}
}

func MiddlewareLog() RequestMiddleware {
	return func(req *Request) {
		fmt.Println(req.req.Header.RawHeaders())
	}
}
