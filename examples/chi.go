// package examples

// import (
// 	"log"
// 	"net/http"

// 	"github.com/lookinlabs/go-logger-middleware"

// 	"github.com/go-chi/chi/v5"
// 	"github.com/go-chi/chi/v5/middleware"
// )

// func main() {
// 	// Initialize the logger middleware
// 	sensitiveFields := []string{"password", "token"}
// 	loggerMiddleware := logger.NewLoggerMiddleware(sensitiveFields)

// 	// Create a Chi router
// 	r := chi.NewRouter()

// 	// Use the built-in Chi middleware
// 	r.Use(middleware.RequestID)
// 	r.Use(middleware.RealIP)
// 	r.Use(middleware.Logger)
// 	r.Use(middleware.Recoverer)

// 	// Use the custom logger middleware
// 	r.Use(loggerMiddleware.Middleware)

// 	// Define a simple endpoint
// 	r.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
// 		w.Write([]byte("Hello, World!"))
// 	})

// 	// Start the server
// 	if err := http.ListenAndServe(":8080", r); err != nil {
// 		log.Fatalf("Failed to run server: %v", err)
// 	}
// }