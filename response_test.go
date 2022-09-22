package fastreq

import (
	"encoding/json"
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
