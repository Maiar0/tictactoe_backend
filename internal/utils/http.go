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
	if err := json.NewDecoder(r.Body).Decode(target); err != nil {
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
