package middleware

import (
	"context"
	"net/http"
)

var trustedHeaderKey = "X-Is-Trusted-Client"

func CreateCheckClientMiddleware(needToCheckClient bool) Middleware {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			if !needToCheckClient {
				ctx = context.WithValue(ctx, trustedHeaderKey, true)
			}

			isTrusted := r.Header.Get(trustedHeaderKey)

			if isTrusted == "trusted" {
				ctx = context.WithValue(ctx, trustedHeaderKey, true)
			} else {
				ctx = context.WithValue(ctx, trustedHeaderKey, false)
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}

func IsTrustedClient(ctx context.Context) bool {
	if isTrusted, ok := ctx.Value(trustedHeaderKey).(bool); ok {
		return isTrusted
	}

	return false
}
