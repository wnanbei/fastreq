package fastreq

import (
	"testing"
)

func TestOauth1(t *testing.T) {
	auth1 := Oauth1{
		ConsumerKey:    "",
		ConsumerSecret: "",
		AccessToken:    "",
		AccessSecret:   "",
	}

	SetOauth1(&auth1)

	url := "https://api.twitter.com/2/users/1443522425690288140/tweets"
	resp, err := Get(url, nil)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(resp.BodyString())
}
