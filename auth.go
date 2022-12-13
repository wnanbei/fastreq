package fastreq

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"math/rand"
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
	return o.header(args, signature)
}

func (o Oauth1) signature(req *Request, args *fasthttp.Args) []byte {
	args.Sort(bytes.Compare)
	queryString := args.String()

	signatureBase := bytes.Buffer{}
	signatureBase.Grow(len(queryString))
	signatureBase.Write(req.Header.Method())
	signatureBase.WriteString("&")
	signatureBase.WriteString(url.QueryEscape(string(req.URI().Scheme())))
	signatureBase.WriteString(url.QueryEscape("://"))
	signatureBase.WriteString(url.QueryEscape(string(req.URI().Host())))
	signatureBase.WriteString(url.QueryEscape(string(req.URI().Path())))
	signatureBase.WriteString("&")
	signatureBase.WriteString(url.QueryEscape(queryString))

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

func (o Oauth1) header(args *fasthttp.Args, signature []byte) []byte {
	header := bytes.NewBuffer([]byte(`OAuth oauth_consumer_key="`))
	header.WriteString(url.QueryEscape(unsafeB2S(args.Peek("oauth_consumer_key"))))
	header.WriteString(`", oauth_nonce="`)
	header.WriteString(url.QueryEscape(unsafeB2S(args.Peek("oauth_nonce"))))
	header.WriteString(`", oauth_signature="`)
	header.WriteString(url.QueryEscape(unsafeB2S(signature)))
	header.WriteString(`", oauth_signature_method="`)
	header.WriteString(url.QueryEscape(unsafeB2S(args.Peek("oauth_signature_method"))))
	header.WriteString(`", oauth_timestamp="`)
	header.WriteString(url.QueryEscape(unsafeB2S(args.Peek("oauth_timestamp"))))
	header.WriteString(`", oauth_token="`)
	header.WriteString(url.QueryEscape(unsafeB2S(args.Peek("oauth_token"))))
	header.WriteString(`", oauth_version="`)
	header.WriteString(url.QueryEscape(unsafeB2S(args.Peek("oauth_version"))))
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
