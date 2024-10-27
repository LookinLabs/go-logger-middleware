package main

import (
	"log"
	"net/http"
	"os"

	"github.com/lookinlabs/go-logger-middleware"
)

func main() {
	// Initialize the logger middleware
	sensitiveFields := []string{"password", "token"}
	appLogger := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	loggerMiddleware := logger.NewLoggerMiddleware(sensitiveFields, appLogger)

	// Create a new HTTP mux (router)
	mux := http.NewServeMux()

	// Define a simple endpoint
	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	// Define another endpoint
	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"username": "john_doe", "password": "secret"}`))
	})

	// Wrap the mux with the logger middleware
	handler := loggerMiddleware.Middleware(mux)

	// Start the server
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
