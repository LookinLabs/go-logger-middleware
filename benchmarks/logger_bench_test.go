package main

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/handlers"
	"github.com/lookinlabs/go-logger-middleware"
	"github.com/urfave/negroni"
)

// No-op logger for benchmarks
type noopLogger struct{}

func (l *noopLogger) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func BenchmarkCustomLoggerMiddleware(b *testing.B) {
	// Define a simple handler to wrap with the middleware
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	// Initialize the logger middleware with a no-op logger
	sensitiveFields := []string{"password", "token"}
	noOpLogger := log.New(&noopLogger{}, "", 0)
	loggerMiddleware := logger.NewLoggerMiddleware(sensitiveFields, noOpLogger)

	// Create the middleware
	middleware := loggerMiddleware.Middleware(handler)

	// Create a sample request body
	requestBody := []byte(`{"username": "john_doe", "password": "secret"}`)

	// Create a request
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

func BenchmarkChiLoggerMiddleware(b *testing.B) {
	// Define a simple handler to wrap with the middleware
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	// Create the middleware with a no-op logger
	noOpLogger := log.New(&noopLogger{}, "", 0)
	middleware := middleware.RequestLogger(&middleware.DefaultLogFormatter{Logger: noOpLogger})(handler)

	// Create a sample request body
	requestBody := []byte(`{"username": "john_doe", "password": "secret"}`)

	// Create a request
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

func BenchmarkNegroniLoggerMiddleware(b *testing.B) {
	// Define a simple handler to wrap with the middleware
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	// Create the middleware with a no-op logger
	noOpLogger := log.New(&noopLogger{}, "", 0)
	n := negroni.New()
	n.Use(negroni.NewLogger())
	n.UseHandler(handler)

	// Override the logger for the Negroni logger middleware
	for _, middleware := range n.Handlers() {
		if logger, ok := middleware.(*negroni.Logger); ok {
			logger.ALogger = noOpLogger
		}
	}

	// Create a sample request body
	requestBody := []byte(`{"username": "john_doe", "password": "secret"}`)

	// Create a request
	req := httptest.NewRequest("GET", "http://example.com/foo", bytes.NewBuffer(requestBody))

	// Set the number of bytes processed per iteration
	b.SetBytes(int64(len(requestBody)))

	// Run the benchmark
	for i := 0; i < b.N; i++ {
		// Create a response recorder
		rr := httptest.NewRecorder()

		// Serve the request
		n.ServeHTTP(rr, req)
	}
}

func BenchmarkGorillaLoggerMiddleware(b *testing.B) {
	// Define a simple handler to wrap with the middleware
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	// Create the middleware with a no-op logger
	noOpLogger := log.New(&noopLogger{}, "", 0)
	middleware := handlers.LoggingHandler(noOpLogger.Writer(), handler)

	// Create a sample request body
	requestBody := []byte(`{"username": "john_doe", "password": "secret"}`)

	// Create a request
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
