// package main

// import (
// 	"log"
// 	"net/http"

// 	"github.com/lookinlabs/go-logger-middleware"
// )

// func main() {
// 	// Initialize the logger middleware
// 	sensitiveFields := []string{"password", "token"}
// 	loggerMiddleware := logger.NewLoggerMiddleware(sensitiveFields)

// 	// Create a new HTTP mux (router)
// 	mux := http.NewServeMux()

// 	// Define a simple endpoint
// 	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
// 		w.Write([]byte("Hello, World!"))
// 	})

// 	// Wrap the mux with the logger middleware
// 	handler := loggerMiddleware.Middleware(mux)

// 	// Start the server
// 	if err := http.ListenAndServe(":8080", handler); err != nil {
// 		log.Fatalf("Failed to run server: %v", err)
// 	}
// }
