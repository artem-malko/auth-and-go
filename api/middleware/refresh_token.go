package middleware

import (
	"context"
	"net/http"

	"github.com/artem-malko/auth-and-go/constants"
	"github.com/artem-malko/auth-and-go/credentials"

	"github.com/google/uuid"
)

var refreshTokenKey = &contextKey{"accessToken"}

func CreateRefreshTokenMiddleware(decodeKey string) Middleware {
	decodeKeyBytes := []byte(decodeKey)

	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			refreshTokenCookie, err := r.Cookie(constants.Values.RefreshTokenCookieName)

			if err != nil {
				ctx = context.WithValue(ctx, refreshTokenKey, nil)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			refreshTokenInfo, err := credentials.ParseRawRefreshToken(refreshTokenCookie.Value, decodeKeyBytes)

			if err != nil {
				ctx = context.WithValue(ctx, refreshTokenKey, nil)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			ctx = context.WithValue(ctx, refreshTokenKey, refreshTokenInfo.RefreshToken)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}

func GetRefreshToken(ctx context.Context) *uuid.UUID {
	if refreshToken, ok := ctx.Value(refreshTokenKey).(uuid.UUID); ok {
		return &refreshToken
	}
	return nil
}
