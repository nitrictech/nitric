package middleware

import (
	"fmt"
	"net/http"
	"time"
)

type wrappedResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *wrappedResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func ProxyLogging(source string, target string, logResult bool) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			wrapped := &wrappedResponseWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			fmt.Printf("%s -> %s %s %s\n", source, target, r.Method, r.URL.Path)
			next.ServeHTTP(wrapped, r)
			if logResult {
				fmt.Printf("%s <- %s %d %s\n", source, target, wrapped.statusCode, time.Since(start))
			}
		})
	}
}
