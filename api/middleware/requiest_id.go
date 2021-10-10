package middleware

import (
	"context"
	"net/http"

	"github.com/artem-malko/auth-and-go/api/context_utils"

	"github.com/google/uuid"
)

// RequestID set request ID (UUID V4) into context and headers
func RequestID(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		requestID := r.URL.Query().Get("request_id")

		if requestID == "" {
			requestID = uuid.New().String()
		}

		ctx = context.WithValue(ctx, context_utils.RequestIDKey, requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}
