package middleware

import (
	"context"
	"net/http"
)

var trustedSearchBotKey = "X-Is-Trusted-Search-Bot"

func CreateCheckSearchBotMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		isTrustedSearchBotCookie, isTrustedSearchBotCookieErr := r.Cookie("_is_trusted_search_bot")

		if isTrustedSearchBotCookieErr == nil && isTrustedSearchBotCookie.Value == "trusted" {
			ctx = context.WithValue(ctx, trustedSearchBotKey, true)
		} else {
			ctx = context.WithValue(ctx, trustedSearchBotKey, false)
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

func IsTrustedSearchBot(ctx context.Context) bool {
	if isTrusted, ok := ctx.Value(trustedSearchBotKey).(bool); ok {
		return isTrusted
	}

	return false
}
