package context_utils

import (
	"context"
)

var RequestIDKey = "RequestID"

func GetRequestID(ctx context.Context) string {
	if reqID, ok := ctx.Value(RequestIDKey).(string); ok {
		return reqID
	}
	return ""
}
