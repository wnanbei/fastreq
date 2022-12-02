package fastreq

import (
	"reflect"
	"unsafe"

	"github.com/valyala/fasthttp"
)

type Args struct {
	*fasthttp.Args
}

func NewArgs() *Args {
	return &Args{fasthttp.AcquireArgs()}
}

func (a *Args) Release() {
	fasthttp.ReleaseArgs(a.Args)
}

// unsafeB2S converts byte slice to a string without memory allocation.
// See https://groups.google.com/forum/#!msg/Golang-Nuts/ENgbUzYvCuU/90yGx7GUAgAJ .
//
// Note it may break if string and/or slice header will change
// in the future go versions.
func unsafeB2S(b []byte) string {
	/* #nosec G103 */
	return *(*string)(unsafe.Pointer(&b))
}

// unsafeS2B converts string to a byte slice without memory allocation.
//
// Note it may break if string and/or slice header will change
// in the future go versions.
func unsafeS2B(s string) (b []byte) {
	/* #nosec G103 */
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	/* #nosec G103 */
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh.Data = sh.Data
	bh.Cap = sh.Len
	bh.Len = sh.Len
	return b
}
