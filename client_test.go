package fastreq

import "testing"

func TestDebug(t *testing.T) {
	SetDebugLevel(DebugDetail)
	params := NewArgs()
	params.Add("hello", "world")
	params.Add("params", "2")
	resp, err := Get("http://httpbin.org/get", params)
	if err != nil {
		t.Fatal(err)
	}
	Release(resp)
}
