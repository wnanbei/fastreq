package fastreq

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"net/url"
	"strconv"
	"time"

	"github.com/valyala/fasthttp"
)

type Oauth1 struct {
	ConsumerKey    string
	ConsumerSecret string
	AccessToken    string
	AccessSecret   string
}

func (o Oauth1) GenHeader(req *Request) []byte {
	args := fasthttp.AcquireArgs()

	args.Add("oauth_token", o.AccessToken)
	args.Add("oauth_consumer_key", o.ConsumerKey)
	args.Add("oauth_signature_method", "HMAC-SHA1")
	args.Add("oauth_version", "1.0")
	args.Add("oauth_nonce", "1.0")
	args.Add("oauth_timestamp", strconv.FormatInt(time.Now().Unix(), 10))

	req.req.URI().QueryArgs().VisitAll(func(key, value []byte) {
		args.AddBytesKV(key, value)
	})

	signature := o.signature(req, args)
	return o.header(req, args, signature)
}

func (o Oauth1) signature(req *Request, args *fasthttp.Args) []byte {
	queryString := args.QueryString() //TODO %20

	signatureBase := bytes.Buffer{}
	signatureBase.Grow(len(queryString))
	signatureBase.Write(req.req.Header.Method())
	signatureBase.WriteString("&")
	signatureBase.Write(req.req.URI().Scheme())
	signatureBase.WriteString("://")
	signatureBase.Write(req.req.URI().Host())
	signatureBase.Write(req.req.URI().Path())
	signatureBase.WriteString("&")
	signatureBase.WriteString(url.QueryEscape(string(queryString)))

	signatureKey := bytes.Buffer{}
	signatureKey.Grow(len(o.ConsumerSecret) + len(o.AccessSecret) + 1)
	signatureKey.WriteString(o.ConsumerSecret)
	signatureKey.WriteString("&")
	signatureKey.WriteString(o.AccessSecret)

	h := hmac.New(sha1.New, signatureKey.Bytes())
	signature := h.Sum(signatureBase.Bytes())
	var encodedSignature []byte
	base64.StdEncoding.Encode(encodedSignature, signature)
	return encodedSignature
}

func (o Oauth1) header(req *Request, args *fasthttp.Args, signature []byte) []byte {
	header := bytes.NewBuffer([]byte(`OAuth oauth_consumer_key="`))
	header.Write(queryEscape(args.Peek("oauth_consumer_key"), encodeQueryComponent))
	header.WriteString(`", oauth_nonce="`)
	header.Write(queryEscape(args.Peek("oauth_nonce"), encodeQueryComponent))
	header.WriteString(`", oauth_signature="`)
	header.Write(queryEscape(signature, encodeQueryComponent))
	header.WriteString(`", oauth_signature_method="`)
	header.Write(queryEscape(args.Peek("oauth_signature_method"), encodeQueryComponent))
	header.WriteString(`", oauth_timestamp="`)
	header.Write(queryEscape(args.Peek("oauth_timestamp"), encodeQueryComponent))
	header.WriteString(`", oauth_token="`)
	header.Write(queryEscape(args.Peek("oauth_token"), encodeQueryComponent))
	header.WriteString(`", oauth_version="`)
	header.Write(queryEscape(args.Peek("oauth_version"), encodeQueryComponent))
	header.WriteString(`"`)

	return header.Bytes()
}
