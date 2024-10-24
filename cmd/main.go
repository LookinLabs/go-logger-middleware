package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte("Hello, World!")); err != nil {
			log.Printf("Failed to write response: %v", err)
		}
	})

	r.Get("/json", func(w http.ResponseWriter, r *http.Request) {
		response := map[string]string{"message": "Hello, JSON!"}
		responseBody, err := json.Marshal(response)
		if err != nil {
			http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write(responseBody); err != nil {
			log.Printf("Failed to write response: %v", err)
		}
	})

	// Create a custom server with timeouts
	server := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	// Start the server
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
