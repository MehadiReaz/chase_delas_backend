package middleware

import (
	"log"
	"net/http"
	"time"
)

// Logging is a middleware that logs HTTP request information
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log request details
		startTime := time.Now()
		log.Printf("Started %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)

		// Call the next handler
		next.ServeHTTP(w, r)

		// Log completion time
		duration := time.Since(startTime)
		log.Printf("Completed %s %s in %v", r.Method, r.URL.Path, duration)
	})
}

// JSONContentType is a middleware that sets the Content-Type header to application/json
func JSONContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set JSON content type for all responses
		w.Header().Set("Content-Type", "application/json")

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}
