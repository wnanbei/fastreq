package fastreq

import (
	"fmt"
	"time"
)

func ExampleNewClient() {
	client := NewClient()

	resp, err := client.Get("https://www.baidu.com")
	if err != nil {
		panic(err)
	}
	fmt.Println(resp.String())
}

func ExampleNewClient_withConfig() {
	client := NewClient(&ClientConfig{
		Timeout:           time.Second * 10,
		DebugLevel:        DebugDetail,
		MaxRedirectsCount: 3,
		DefaultUserAgent:  "my http client",
	})

	resp, err := client.Get("https://www.baidu.com")
	if err != nil {
		panic(err)
	}
	fmt.Println(resp.String())
}
