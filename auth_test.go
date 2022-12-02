package fastreq

import (
	"testing"
)

func TestOauth1(t *testing.T) {
	auth1 := Oauth1{
		ConsumerKey:    "fzhgOqIhjxWIN6QSzMCvju5QU",
		ConsumerSecret: "nlUAGEVxSQA8EiboGue6tU5VmqKkEwP6t56Wv6hOrfwFIkCtxS",
		AccessToken:    "1005370915331272704-KguBzyEU0NC2uaiU0AL3shqdzSITvX",
		AccessSecret:   "0NcqdCcvTRj1wUhckAbQOoufrWo8nm1Z5uUNLp1YItmtO",
	}

	SetOauth1(&auth1)
	SetHTTPProxy("localhost:8001")
	SetDebugLevel(DebugDetail)

	url := "https://api.twitter.com/2/users/1443522425690288140/tweets"
	args := NewQueryParams()
	args.Add("max_results", "1")
	ctx, err := Get(url, args)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(ctx.Response.String())
}
