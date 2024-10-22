package logger

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func BenchmarkMiddleware(b *testing.B) {
	// Define a simple handler to wrap with the middleware
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	// Initialize the logger middleware
	sensitiveFields := []string{"password", "token"}
	loggerMiddleware := NewLoggerMiddleware(sensitiveFields)

	// Create the middleware
	middleware := loggerMiddleware.Middleware(handler)

	// Create a sample request body
	requestBody := []byte(`{"username": "john_doe", "password": "secret"}`)

	// Create a sample request
	req := httptest.NewRequest("GET", "http://example.com/foo", bytes.NewBuffer(requestBody))

	// Set the number of bytes processed per iteration
	b.SetBytes(int64(len(requestBody)))

	// Run the benchmark
	for i := 0; i < b.N; i++ {
		// Create a response recorder
		rr := httptest.NewRecorder()

		// Serve the request
		middleware.ServeHTTP(rr, req)
	}
}

func TestMain(m *testing.M) {
	// Set a timeout for the tests
	timeout := time.AfterFunc(10*time.Second, func() {
		panic("Tests timed out")
	})
	defer timeout.Stop()

	// Run the tests
	m.Run()
}
