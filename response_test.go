package fastreq

import (
	"testing"

	jsoniter "github.com/json-iterator/go"
)

var jsonExamples = []byte(`{
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

func BenchmarkJson(b *testing.B) {
	resp := NewResponse()
	resp.SetBody(jsonExamples)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var v map[string]interface{}
		resp.Json(v)
	}
}

func BenchmarkJsoniter(b *testing.B) {
	SetJsonUnmarshal(jsoniter.Unmarshal)
	resp := NewResponse()
	resp.SetBody(jsonExamples)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var v map[string]interface{}
		resp.Json(v)
	}
}

func BenchmarkJsonPart(b *testing.B) {
	resp := NewResponse()
	resp.SetBody(jsonExamples)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var v map[string]interface{}
		resp.JsonPart("friends.0", v)
	}
}

func BenchmarkJsonGet(b *testing.B) {
	resp := NewResponse()
	resp.SetBody(jsonExamples)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		resp.JsonGet("friends.2.last")
	}
}

func BenchmarkJsonGetMany(b *testing.B) {
	resp := NewResponse()
	resp.SetBody(jsonExamples)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		resp.JsonGetMany("friends.2.last")
	}
}

func TestJson(t *testing.T) {
	resp := NewResponse()
	resp.SetBody(jsonExamples)

	var data map[string]interface{}
	if err := resp.Json(&data); err != nil {
		t.Fatal(err)
	}

	t.Logf("%+v", data)
	Release(resp)
}

func TestJsonGet(t *testing.T) {
	resp := NewResponse()
	resp.SetBody(jsonExamples)

	var data []map[string]interface{}
	if err := resp.JsonPart("friends", &data); err != nil {
		t.Fatal(err)
	}

	t.Logf("%+v", data)
	Release(resp)
}

func TestJsonPart(t *testing.T) {
	resp := NewResponse()
	resp.SetBody(jsonExamples)

	t.Log(resp.JsonGet("friends.2.last").String())

	Release(resp)
}
