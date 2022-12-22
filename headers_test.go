package fastreq

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Headers_Add(t *testing.T) {
	m := make(map[string]struct{})
	var h = NewHeaders()
	h.Add("aaa", "bbb")
	h.Add("content-type", "xxx")
	m["bbb"] = struct{}{}
	m["xxx"] = struct{}{}
	for i := 0; i < 10; i++ {
		v := fmt.Sprintf("%d", i)
		h.Add("Foo-Bar", v)
		m[v] = struct{}{}
	}

	h.VisitAll(func(k, v []byte) {
		require.Contains(t, []string{"Aaa", "Foo-Bar", "Content-Type"}, string(k))

		_, ok := m[string(v)]
		require.True(t, ok)
		delete(m, string(v))
	})
	require.Emptyf(t, m, "%d headers are missed", len(m))
}

func Test_Headers_Empty_Value(t *testing.T) {
	h := NewHeaders()
	h.Set("EmptyValue1", "")
	h.Set("EmptyValue2", " ")

	require.Empty(t, h.Peek("EmptyValue1"))
	require.NotEmpty(t, h.Peek("EmptyValue2"))
}

func Test_Headers_Del(t *testing.T) {
	var h = NewHeaders(
		"Host", "aaa",
		"User-agent", "ccc",
		"cookie", "foobar=baz",
	)
	h.Set("Foo-Bar", "baz")
	h.SetBytesK([]byte("aaa"), "bbb")
	h.SetBytesV("Connection", []byte("keep-alive"))
	h.SetBytesKV([]byte("content-Type"), []byte("aaa"))
	h.Set("Content-Length", "1123")

	h.Del("foo-bar")
	h.Del("connection")
	h.DelBytes([]byte("content-type"))
	h.Del("Host")
	h.Del("user-agent")
	h.Del("content-length")
	h.Del("cookie")

	require.Equal(t, "bbb", string(h.Peek("aaa")))
	require.Empty(t, h.Peek("Foo-Bar"))
	require.Empty(t, h.PeekBytes([]byte("Connection")))
	require.Empty(t, h.Peek("Content-Length"))
	require.Empty(t, h.Peek("Host"))
	require.Empty(t, h.Peek("User-Agent"))
	require.Empty(t, h.Peek("Content-Length"))
	require.Empty(t, h.Peek("Cookie"))
}
