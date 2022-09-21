package fastreq

import (
	"testing"

	oauth1 "github.com/klaidas/go-oauth1"
)

func TestOauth1(t *testing.T) {
	auth1 := Oauth1{}

	SetOauth1(&auth1)
	SetHTTPProxy("localhost:8001")

	url := "https://api.twitter.com/2/users/1443522425690288140/tweets"
	resp, err := Get(url, nil)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(resp.BodyString())
}

func TestOauth2(t *testing.T) {
	SetHTTPProxy("localhost:8001")

	auth := oauth1.OAuth1{}

	url := "https://api.twitter.com/2/users/1443522425690288140/tweets"
	authHeader := auth.BuildOAuth1Header("POST", url, map[string]string{
		"include_entities": "true",
	})
	t.Log(authHeader)

	req := NewRequest()
	req.SetRequestURI(url)
	req.SetMethod(GET)
	req.SetHeader("Authorization", authHeader)

	resp, err := Do(req)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(resp.BodyString())
}
