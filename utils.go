package fastreq

import (
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
