package middleware

import (
	"net/http"
	"time"

	"github.com/artem-malko/auth-and-go/api/context_utils"

	"github.com/apex/log"
	"github.com/go-chi/chi/v5/middleware"
)

// CreateAccessLogMiddleware create middleware to log every request
func CreateAccessLogMiddleware(logger log.Interface) Middleware {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "OPTIONS" {
				next.ServeHTTP(w, r)
				return
			}

			requestID := context_utils.GetRequestID(r.Context())
			remoteIP := GetClientIP(r.Context())
			URI := r.URL.Path
			query := r.URL.RawQuery
			referer := r.Referer()
			userAgent := r.UserAgent()
			method := r.Method
			proto := r.Proto

			start := time.Now()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			next.ServeHTTP(ww, r)

			defer func() {
				fieldsCompleted := map[string]interface{}{
					"remote_ip":        remoteIP,
					"http_uri":         URI,
					"http_query":       query,
					"http_referer":     referer,
					"http_user_agent":  userAgent,
					"http_method":      method,
					"http_proto":       proto,
					"request_id":       requestID,
					"lead_time":        int64(time.Since(start) / time.Millisecond),
					"lead_time_ms":     time.Since(start),
					"http_status_code": ww.Status(),
					"http_status_text": http.StatusText(ww.Status()),
					"stage":            "completed",
				}

				logger.WithFields(log.Fields(fieldsCompleted)).Infof("Completed incoming request")
			}()
		}
		return http.HandlerFunc(fn)
	}
}
