package middleware

import (
	"net/http"
)

// contextKey is a value for use with context.WithValue. It's used as
// a pointer so it fits in an interface{} without allocation. This technique
// for defining context keys was copied from Go 1.7's new use of context in net/http.
type contextKey struct {
	name string
}

func (k *contextKey) String() string {
	return k.name
}

// Middleware is chi router middleware type
type Middleware func(next http.Handler) http.Handler
