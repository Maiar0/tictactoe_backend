package utils

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type APIError struct {
	Error string `json:"error"`
}

func WriteError(w http.ResponseWriter, code int, msg string) {
	log.Printf("[API Error] %d: %s", code, msg)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	_ = json.NewEncoder(w).Encode(APIError{Error: msg})

}

// statusRecorder lets us capture the status code written by a handler.
type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

// LoggingMiddleware logs the method, path, status code, and duration.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		log.Printf("[REQ] %s %s", r.Method, r.URL.Path)

		next.ServeHTTP(rec, r) // call the real handler

		duration := time.Since(start)
		log.Printf("[RES] %s %s -> %d (%v)", r.Method, r.URL.Path, rec.status, duration)
	})
}
