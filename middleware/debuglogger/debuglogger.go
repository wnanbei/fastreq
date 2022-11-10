package debuglogger

import (
	"fmt"
	"time"

	"github.com/wnanbei/fastreq"
)

func MiddlewareLogger() fastreq.Middleware {
	return func(ctx *fastreq.Ctx) error {
		start := time.Now()
		fmt.Printf(
			"REQUEST[%s]: %s %s\n%s",
			start.Format(time.RFC3339),
			ctx.Request.Header.Method(),
			ctx.Request.URI().FullURI(),
			ctx.Request.String(),
		)

		if err := ctx.Next(); err != nil {
			return err
		}

		end := time.Now()
		fmt.Printf(
			"RESPONSE[%s]: %d %dms %s %s\n%s",
			end.Format(time.RFC3339),
			ctx.Response.StatusCode(),
			end.Sub(start).Milliseconds(),
			ctx.Request.Header.Method(),
			ctx.Request.URI().FullURI(),
			ctx.Response.String(),
		)

		return nil
	}
}
