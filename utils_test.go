package fastreq

import (
	"testing"
)

func TestQueryEscape(t *testing.T) {
	u := []byte("hello world")

	escaped := queryEscape(u)
	t.Log(string(escaped))
}
