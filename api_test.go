package fastreq

import (
	"testing"
)

func TestGet(t *testing.T) {
	params := NewQueryParams()
	params.Add("hello", "world")
	params.Add("params", "2")
	resp, err := Get("http://httpbin.org/get", params)
	if err != nil {
		t.Fatal(err)
	}
	Release(resp)
}

func TestPost(t *testing.T) {
	formBody := NewQueryParams()
	formBody.Add("hello", "world")
	formBody.Add("params", "2")
	resp, err := Post("http://httpbin.org/post", formBody)
	if err != nil {
		t.Fatal(err)
	}
	Release(resp)
}

func TestSetHTTPProxy(t *testing.T) {
	SetHTTPProxy("localhost:8001")

	resp, err := Get("http://httpbin.org/ip")
	if err != nil {
		t.Fatal(err)
	}
	Release(resp)
}

func TestSetSocks5Proxy(t *testing.T) {
	SetSocks5Proxy("SOCKS5://localhost:1081")

	resp, err := Get("http://httpbin.org/ip")
	if err != nil {
		t.Fatal(err)
	}
	Release(resp)
}

func TestSetEnvProxy(t *testing.T) {
	SetEnvHTTPProxy()

	resp, err := Get("http://httpbin.org/ip")
	if err != nil {
		t.Fatal(err)
	}
	Release(resp)
}
