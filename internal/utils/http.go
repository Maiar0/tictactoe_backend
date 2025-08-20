package utils

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"time"
)

type APIError struct {
	Error string `json:"error"`
}

// WriteJSONError writes a JSON error response with the given HTTP status code and error message.
// It automatically logs the error and sets the appropriate Content-Type header.
func WriteJSONError(w http.ResponseWriter, code int, msg string) {
	log.Printf("[API Error] %d: %s", code, msg)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	_ = json.NewEncoder(w).Encode(APIError{Error: msg}) //TODO:: needs to return error

}

// ReadRequestBody reads and decodes the HTTP request body into the target struct.
// It limits the body size to 1MB and logs successful decoding for debugging.
// Returns an error if reading or JSON decoding fails.
func ReadRequestBody(w http.ResponseWriter, r *http.Request, target any) error {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	log.Printf("[ReadRequestBody] Raw body: %s", string(bodyBytes))
	if err := json.Unmarshal(bodyBytes, target); err != nil {
		return err
	}
	log.Printf("[ReadRequestBody] Request body decoded successfully: %+v", target)
	return nil
}

// WriteJSONResponse writes a JSON response with the given HTTP status code and data.
// It automatically sets the Content-Type header and logs successful responses.
// Returns an error if JSON encoding fails.
func WriteJSONResponse(w http.ResponseWriter, code int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	log.Printf("[WriteJSONResponse] Response written successfully: %d: %+v", code, data)
	return json.NewEncoder(w).Encode(data)
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

// Implement http.Hijacker to support WebSocket upgrades
func (r *statusRecorder) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hijacker, ok := r.ResponseWriter.(http.Hijacker); ok {
		return hijacker.Hijack()
	}
	return nil, nil, fmt.Errorf("underlying ResponseWriter does not implement http.Hijacker")
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

// CORSMiddleware allows all origins for development purposes.
// TODO:: This should be configured properly for production.
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow all origins
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// Allow common headers
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")

		// Allow common methods
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
