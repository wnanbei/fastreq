package auth

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"github.com/valyala/fasthttp"
	"github.com/wnanbei/fastreq"
	"math/rand"
	"strconv"
	"time"
)

func MiddlewareOauth1(o *Oauth1) fastreq.Middleware {
	return func(ctx *fastreq.Ctx) error {
		auth := o.GenHeader(ctx.Request)
		ctx.Request.Header.SetBytesV("Authorization", auth)
		return ctx.Next()
	}
}

type Oauth1 struct {
	ConsumerKey    string
	ConsumerSecret string
	AccessToken    string
	AccessSecret   string
}

func (o Oauth1) GenHeader(req *fastreq.Request) []byte {
	args := fasthttp.AcquireArgs()

	args.Add("oauth_consumer_key", o.ConsumerKey)
	args.AddBytesV("oauth_nonce", genNonce())
	args.Add("oauth_signature_method", "HMAC-SHA1")
	args.Add("oauth_timestamp", strconv.FormatInt(time.Now().Unix(), 10))
	args.Add("oauth_token", o.AccessToken)
	args.Add("oauth_version", "1.0")

	req.URI().QueryArgs().VisitAll(func(key, value []byte) {
		args.AddBytesKV(key, value)
	})

	signature := o.signature(req, args)
	return o.header(req, args, signature)
}

func (o Oauth1) signature(req *fastreq.Request, args *fasthttp.Args) []byte {
	args.Sort(bytes.Compare)
	queryString := args.QueryString()

	signatureBase := bytes.Buffer{}
	signatureBase.Grow(len(queryString))
	signatureBase.Write(req.Header.Method())
	signatureBase.WriteString("&")
	signatureBase.Write(queryEscape(req.URI().Scheme()))
	signatureBase.Write(queryEscape([]byte("://")))
	signatureBase.Write(queryEscape(req.URI().Host()))
	signatureBase.Write(queryEscape(req.URI().Path()))
	signatureBase.WriteString("&")
	signatureBase.Write(queryEscape(queryString))

	signatureKey := bytes.Buffer{}
	signatureKey.Grow(len(o.ConsumerSecret) + len(o.AccessSecret) + 1)
	signatureKey.WriteString(o.ConsumerSecret)
	signatureKey.WriteString("&")
	signatureKey.WriteString(o.AccessSecret)

	h := hmac.New(sha1.New, signatureKey.Bytes())
	h.Write(signatureBase.Bytes())
	signature := h.Sum(nil)
	encodedSignature := make([]byte, base64.StdEncoding.EncodedLen(len(signature)))
	base64.StdEncoding.Encode(encodedSignature, signature)
	return encodedSignature
}

func (o Oauth1) header(req *fastreq.Request, args *fasthttp.Args, signature []byte) []byte {
	header := bytes.NewBuffer([]byte(`OAuth oauth_consumer_key="`))
	header.Write(queryEscape(args.Peek("oauth_consumer_key")))
	header.WriteString(`", oauth_nonce="`)
	header.Write(queryEscape(args.Peek("oauth_nonce")))
	header.WriteString(`", oauth_signature="`)
	header.Write(queryEscape(signature))
	header.WriteString(`", oauth_signature_method="`)
	header.Write(queryEscape(args.Peek("oauth_signature_method")))
	header.WriteString(`", oauth_timestamp="`)
	header.Write(queryEscape(args.Peek("oauth_timestamp")))
	header.WriteString(`", oauth_token="`)
	header.Write(queryEscape(args.Peek("oauth_token")))
	header.WriteString(`", oauth_version="`)
	header.Write(queryEscape(args.Peek("oauth_version")))
	header.WriteString(`"`)

	return header.Bytes()
}

const allowed = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func genNonce() []byte {
	b := make([]byte, 48)
	for i := range b {
		b[i] = allowed[rand.Intn(len(allowed))]
	}
	return b
}
