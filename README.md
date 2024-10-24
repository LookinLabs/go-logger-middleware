# go-logger-middleware

Go Logger Middleware is a lightweight, fast and simple HTTP middleware that logs incoming HTTP requests and outgoing HTTP responses. 

It uses only standard Go libraries and is compatible with any Go web framework that supports HTTP middleware.

## Usage

The examples usage can be found under the `examples` directory.

### Gin

```go
package examples

import (
	"log"
	"net/http"

	"go-logger-middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize the logger middleware
	sensitiveFields := []string{"password", "token"}
	loggerMiddleware := logger.NewLoggerMiddleware(sensitiveFields)

	// Create a Gin router
	r := gin.Default()

	// Use the middleware
	r.Use(gin.WrapH(loggerMiddleware.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Your handler logic here
		w.Write([]byte("Hello, World!"))
	}))))

	// Define a simple endpoint
	r.GET("/hello", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, World!")
	})

	// Start the server
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
```

### Chi

```go
package examples

import (
	"log"
	"net/http"

	"go-logger-middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	// Initialize the logger middleware
	sensitiveFields := []string{"password", "token"}
	loggerMiddleware := logger.NewLoggerMiddleware(sensitiveFields)

	// Create a Chi router
	r := chi.NewRouter()

	// Use the built-in Chi middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Use the custom logger middleware
	r.Use(loggerMiddleware.Middleware)

	// Define a simple endpoint
	r.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	// Start the server
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
```

### Default Logger

```go
package main

import (
	"log"
	"net/http"

	"github.com/lookinlabs/go-logger-middleware"
)

func main() {
	// Initialize the logger middleware
	sensitiveFields := []string{"password", "token"}
	loggerMiddleware := logger.NewLoggerMiddleware(sensitiveFields)

	// Create a new HTTP mux (router)
	mux := http.NewServeMux()

	// Define a simple endpoint
	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	// Wrap the mux with the logger middleware
	handler := loggerMiddleware.Middleware(mux)

	// Start the server
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
```

## Benchmarks



## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.