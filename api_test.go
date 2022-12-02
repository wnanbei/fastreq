package fastreq

import "testing"

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

	resp, err := Get("http://httpbin.org/ip", nil)
	if err != nil {
		t.Fatal(err)
	}
	Release(resp)
}

func TestSetSocks5Proxy(t *testing.T) {
	SetSocks5Proxy("SOCKS5://localhost:1081")

	resp, err := Get("http://httpbin.org/ip", nil)
	if err != nil {
		t.Fatal(err)
	}
	Release(resp)
}

func TestSetEnvProxy(t *testing.T) {
	SetEnvHTTPProxy()

	resp, err := Get("http://httpbin.org/ip", nil)
	if err != nil {
		t.Fatal(err)
	}
	Release(resp)
}

// func TestDownloadFile(t *testing.T) {
// 	req := NewRequest(GET, "https://cf-markting-1256732272.cos.ap-shanghai.myqcloud.com/W00000012229/material/file/d7f0bf2c-01c6-449a-a291-e8cf5084fce6/h7icq706xp71650363280917.pdf")

// 	err := DownloadFile(req, "./data", "")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// }
