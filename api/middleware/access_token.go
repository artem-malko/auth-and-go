package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/artem-malko/auth-and-go/models"
	"github.com/google/uuid"

	"github.com/artem-malko/auth-and-go/credentials"

	"github.com/artem-malko/auth-and-go/constants"
)

var accessTokenKey = &contextKey{"accessToken"}

type sessionGetter interface {
	GetSessionByAccessToken(accessToken uuid.UUID) (*models.Session, error)
}

// @TODO add logging
func CreateAccessTokenMiddleware(sessionGetter sessionGetter, decodeKey string) Middleware {
	decodeKeyBytes := []byte(decodeKey)

	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			accessTokenCookie, err := r.Cookie(constants.Values.AccessTokenCookieName)

			if err != nil {
				ctx = context.WithValue(ctx, accessTokenKey, nil)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			accessTokenData, err := credentials.ParseRawAccessToken(accessTokenCookie.Value, decodeKeyBytes)

			fmt.Println("accessTokenData: ", accessTokenData)

			if err != nil {
				ctx = context.WithValue(ctx, accessTokenKey, nil)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			if time.Now().Second()%2 == 0 {
				session, err := sessionGetter.GetSessionByAccessToken(accessTokenData.AccessToken)

				if err != nil || session == nil || time.Now().After(session.AccessTokenExpiresDate) {
					ctx = context.WithValue(ctx, accessTokenKey, nil)
					next.ServeHTTP(w, r.WithContext(ctx))
					return
				}
			}

			ctx = context.WithValue(ctx, accessTokenKey, *accessTokenData)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}

func GetAccessTokenInfo(ctx context.Context) *credentials.AccessTokenInfo {
	if accessToken, ok := ctx.Value(accessTokenKey).(credentials.AccessTokenInfo); ok {
		return &accessToken
	}
	return nil
}
