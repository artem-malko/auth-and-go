package middleware

import (
	"net/http"

	"github.com/artem-malko/auth-and-go/api/response"
)

func RequireAuth(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		accessTokenInfo := GetAccessTokenInfo(ctx)

		if accessTokenInfo == nil {
			response.Error(w, http.StatusForbidden, "You are not logged in")
			return
		}

		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
