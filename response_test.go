package fastreq

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/tidwall/gjson"
)

func BenchmarkJsonGet(b *testing.B) {
	s := `{
  "name": {"first": "Tom", "last": "Anderson"},
  "age":37,
  "children": ["Sara","Alex","Jack"],
  "fav.movie": "Deer Hunter",
  "friends": [
    {"first": "Dale", "last": "Murphy", "age": 44, "nets": ["ig", "fb", "tw"]},
    {"first": "Roger", "last": "Craig", "age": 68, "nets": ["fb", "tw"]},
    {"first": "Jane", "last": "Murphy", "age": 47, "nets": ["ig", "tw"]}
  ]
}`

	for i := 0; i < b.N; i++ {
		gjson.Get(s, "friends.3.nets").String()
	}
}

func BenchmarkJsonGetBytes(b *testing.B) {
	s := []byte(`{
  "name": {"first": "Tom", "last": "Anderson"},
  "age":37,
  "children": ["Sara","Alex","Jack"],
  "fav.movie": "Deer Hunter",
  "friends": [
    {"first": "Dale", "last": "Murphy", "age": 44, "nets": ["ig", "fb", "tw"]},
    {"first": "Roger", "last": "Craig", "age": 68, "nets": ["fb", "tw"]},
    {"first": "Jane", "last": "Murphy", "age": 47, "nets": ["ig", "tw"]}
  ]
}`)

	for i := 0; i < b.N; i++ {
		gjson.GetBytes(s, "friends.3.nets").String()
	}
}

func BenchmarkJsonParseGet(b *testing.B) {
	s := `{
  "name": {"first": "Tom", "last": "Anderson"},
  "age":37,
  "children": ["Sara","Alex","Jack"],
  "fav.movie": "Deer Hunter",
  "friends": [
    {"first": "Dale", "last": "Murphy", "age": 44, "nets": ["ig", "fb", "tw"]},
    {"first": "Roger", "last": "Craig", "age": 68, "nets": ["fb", "tw"]},
    {"first": "Jane", "last": "Murphy", "age": 47, "nets": ["ig", "tw"]}
  ]
}`

	for i := 0; i < b.N; i++ {
		r := gjson.Parse(s)
		r.Get("friends.3.nets").String()
	}
}

func BenchmarkJsonUnmarshl(b *testing.B) {
	s := []byte(`{
  "name": {"first": "Tom", "last": "Anderson"},
  "age":37,
  "children": ["Sara","Alex","Jack"],
  "fav.movie": "Deer Hunter",
  "friends": [
    {"first": "Dale", "last": "Murphy", "age": 44, "nets": ["ig", "fb", "tw"]},
    {"first": "Roger", "last": "Craig", "age": 68, "nets": ["fb", "tw"]},
    {"first": "Jane", "last": "Murphy", "age": 47, "nets": ["ig", "tw"]}
  ]
}`)

	for i := 0; i < b.N; i++ {
		var v map[string]interface{}
		json.Unmarshal(s, v)
	}
}

func TestJsonPart(t *testing.T) {
	client := NewClient()
	client.AddMiddleware(MiddlewareLogger())

	params := NewArgs()
	params.Add("hello", "world")
	params.Add("params", "2")
	resp, err := client.Get("http://httpbin.org/get", params)
	if err != nil {
		t.Fatal(err)
	}

	var data map[string]interface{}
	if err := resp.Response.JsonPart("headers", &data); err != nil {
		t.Fatal(err)
	}

	fmt.Printf("\n%+v", data)
	Release(resp)
}
