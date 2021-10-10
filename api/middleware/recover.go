package middleware

import (
	"fmt"
	"net/http"
	"runtime"

	"github.com/artem-malko/auth-and-go/api/context_utils"

	"github.com/apex/log"
	"github.com/artem-malko/auth-and-go/api/response"
)

const stackSize = 1024 * 2

// This middleware have to be used only for unexpected errors like
// "Index is out of range" and so on.
// Do not panic inside app
func CreateRecoverMiddleware(logger log.Interface) Middleware {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil && err != http.ErrAbortHandler {
					stack := make([]byte, stackSize)
					stack = stack[:runtime.Stack(stack, false)]

					logger.
						WithField("request_id", context_utils.GetRequestID(r.Context())).
						WithField("error_type", "panic").
						WithField("stack", fmt.Sprintf("%s", stack)).
						Errorf("%v", err)

					response.InternalServerError(w)
				}
			}()
			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}
