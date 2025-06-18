package middleware

import (
	"fmt"
	"net/http"
	"strings"
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

type CORSOptions struct {
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string
	MaxAge         int
}

var DefaultCORSOptions = CORSOptions{
	AllowedOrigins: []string{"*"},
	AllowedMethods: []string{"*"},
	AllowedHeaders: []string{"*"},
	MaxAge:         86400,
}

func CORS(opts CORSOptions) func(next http.Handler) http.Handler {
	if opts.AllowedOrigins == nil {
		opts.AllowedOrigins = DefaultCORSOptions.AllowedOrigins
	}
	if opts.AllowedMethods == nil {
		opts.AllowedMethods = DefaultCORSOptions.AllowedMethods
	}
	if opts.AllowedHeaders == nil {
		opts.AllowedHeaders = DefaultCORSOptions.AllowedHeaders
	}
	if opts.MaxAge == 0 {
		opts.MaxAge = DefaultCORSOptions.MaxAge
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Set CORS headers
			w.Header().Set("Access-Control-Allow-Origin", strings.Join(opts.AllowedOrigins, ", "))
			w.Header().Set("Access-Control-Allow-Methods", strings.Join(opts.AllowedMethods, ", "))
			w.Header().Set("Access-Control-Allow-Headers", strings.Join(opts.AllowedHeaders, ", "))
			w.Header().Set("Access-Control-Max-Age", fmt.Sprintf("%d", opts.MaxAge))

			// Handle preflight requests
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
