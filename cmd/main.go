package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/lookinlabs/go-logger-middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	// Initialize the logger middleware
	sensitiveFields := []string{"password", "token"}
	appLogger := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	loggerMiddleware := logger.NewLoggerMiddleware(sensitiveFields, appLogger)

	// Create a Chi router
	r := chi.NewRouter()

	// Use the built-in Chi middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)

	// Use the custom logger middleware
	r.Use(loggerMiddleware.Middleware)

	// Define a simple GET endpoint
	r.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	// Define a POST endpoint to test sanitization
	r.Post("/test", func(w http.ResponseWriter, r *http.Request) {
		var requestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		responseBody, err := json.Marshal(requestBody)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(responseBody)
	})

	// Start the server
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
