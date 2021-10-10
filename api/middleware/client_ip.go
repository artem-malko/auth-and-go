package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/apex/log"
	"github.com/artem-malko/auth-and-go/api/response"
)

var clientIPKey = &contextKey{"client ip"}

const clientIPErrorMessage = "clientIP can not be empty"

func CreateClientIPMiddleware(logger log.Interface) Middleware {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			xForwardedFor := r.Header.Get("X-Forwarded-For")
			xRealIP := r.Header.Get("X-Real-IP")
			remoteAddr := r.RemoteAddr
			remoteIP := xForwardedFor

			if xForwardedFor == "" {
				remoteIP = xRealIP
			}

			if xRealIP == "" {
				remoteIP = remoteAddr
			}

			if remoteIP == "" {
				logger.Debugf("%s", errors.New(clientIPErrorMessage))
				response.Error(w, http.StatusBadRequest, clientIPErrorMessage)
				return
			}

			ctx = context.WithValue(ctx, clientIPKey, remoteIP)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}

func GetClientIP(ctx context.Context) string {
	if clientIP, ok := ctx.Value(clientIPKey).(string); ok {
		return clientIP
	}
	return ""
}
