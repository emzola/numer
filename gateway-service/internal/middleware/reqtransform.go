package middleware

import (
	"io"
	"net/http"
	"strings"
)

// RequestTransformationMiddleware transforms the request body to lowercase for POST requests
func RequestTransformationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			body, _ := io.ReadAll(r.Body)
			r.Body = io.NopCloser(strings.NewReader(strings.ToLower(string(body))))
		}

		next.ServeHTTP(w, r)
	})
}
