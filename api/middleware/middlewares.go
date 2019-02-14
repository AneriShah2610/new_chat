package middleware

import (
	"context"
	"log"
	"net/http"

	"github.com/aneri/new_chat/api/dal"
)

type Middleware func(http.Handler) http.Handler

var ctx context.Context

// Cockroachdb middleware
func CockroachDbMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		crConn, err := dal.DbConnect()
		if err != nil {
			log.Println("db middleware error", err)
		}
		ctx = context.WithValue(request.Context(), "crConn", crConn)
		next.ServeHTTP(writer, request.WithContext(ctx))
	})
}

// Multiple handler middleware
func MultipleMiddleware(h http.Handler, m ...Middleware) http.Handler {
	if len(m) < 1 {
		return h
	}
	wrapped := h
	// loop in reverse to preserve middleware order
	for i := len(m) - 1; i >= 0; i-- {
		wrapped = m[i](wrapped)
	}
	return wrapped
}
