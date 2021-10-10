package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"

	"github.com/artem-malko/auth-and-go/api/response"

	"github.com/apex/log"
	"github.com/go-chi/chi/v5"
)

func CreateCorrectUUIDInParamMiddleware(logger log.Interface, contextKey string, paramName string) Middleware {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			rawEntityID := chi.URLParam(r, paramName)

			if rawEntityID == "" {
				response.Error(w, http.StatusBadRequest, paramName+" is not passed")
				return
			}

			entityUUID, err := uuid.Parse(rawEntityID)

			if err != nil {
				logger.Debugf("%s", err)
				response.Error(w, http.StatusBadRequest, paramName+" is not valid UUID")
				return
			}

			ctx = context.WithValue(ctx, contextKey, entityUUID)

			next.ServeHTTP(w, r.WithContext(ctx))
		}

		return http.HandlerFunc(fn)
	}
}
