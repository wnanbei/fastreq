package fastreq

import "testing"

func TestGet(t *testing.T) {
	params := NewArgs()
	params.Add("hello", "world")
	params.Add("params", "2")
	resp, err := Get("http://httpbin.org/get", params)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(resp.BodyString())
}

func TestPost(t *testing.T) {
	formBody := NewArgs()
	formBody.Add("hello", "world")
	formBody.Add("params", "2")
	resp, err := Post("http://httpbin.org/post", formBody)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(resp.BodyString())
}
