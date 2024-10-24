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

The benchmarks results can be found under the `benchmarks` directory.

![Benchmark Results](benchmarks/benchmark_graph.png)

### Benchmark Results

| Benchmark                              | Iterations | Time per Operation (ns/op) | Throughput (MB/s) |
|----------------------------------------|------------|----------------------------|-------------------|
| BenchmarkMarshal_Logger_Small-11       | 2,428,286  | 471.0                      | -                 |
| BenchmarkMarshal_Logger_Medium-11      | 1,000,000  | 1055                       | -                 |
| BenchmarkMarshal_Logger_Large-11       | 749,445    | 1611                       | -                 |
| BenchmarkUnmarshal_Logger_Small-11     | 647,900    | 1734                       | -                 |
| BenchmarkUnmarshal_Logger_Medium-11    | 250,807    | 4780                       | -                 |
| BenchmarkUnmarshal_Logger_Large-11     | 148,404    | 8062                       | -                 |
| BenchmarkMarshal_Stdlib_Small-11       | 1,000,000  | 1056                       | -                 |
| BenchmarkMarshal_Stdlib_Medium-11      | 363,956    | 3273                       | -                 |
| BenchmarkMarshal_Stdlib_Large-11       | 171,868    | 6893                       | -                 |
| BenchmarkUnmarshal_Stdlib_Small-11     | 691,006    | 1762                       | -                 |
| BenchmarkUnmarshal_Stdlib_Medium-11    | 212,148    | 5613                       | -                 |
| BenchmarkUnmarshal_Stdlib_Large-11     | 92,278     | 12845                      | -                 |
| BenchmarkMarshal_Goccy_Small-11        | 1,884,061  | 638.1                      | -                 |
| BenchmarkMarshal_Goccy_Medium-11       | 548,515    | 2184                       | -                 |
| BenchmarkMarshal_Goccy_Large-11        | 255,254    | 4653                       | -                 |
| BenchmarkUnmarshal_Goccy_Small-11      | 1,599,242  | 751.9                      | -                 |
| BenchmarkUnmarshal_Goccy_Medium-11     | 438,148    | 2733                       | -                 |
| BenchmarkUnmarshal_Goccy_Large-11      | 170,019    | 7037                       | -                 |
| BenchmarkMarshal_Jsoniter_Small-11     | 2,425,772  | 494.9                      | -                 |
| BenchmarkMarshal_Jsoniter_Medium-11    | 839,830    | 1437                       | -                 |
| BenchmarkMarshal_Jsoniter_Large-11     | 381,164    | 3117                       | -                 |
| BenchmarkUnmarshal_Jsoniter_Small-11   | 1,341,801  | 896.5                      | -                 |
| BenchmarkUnmarshal_Jsoniter_Medium-11  | 366,684    | 3316                       | -                 |
| BenchmarkUnmarshal_Jsoniter_Large-11   | 140,568    | 8607                       | -                 |
| BenchmarkCustomLoggerMiddleware-11     | 566,791    | 2084                       | 22.07             |
| BenchmarkChiLoggerMiddleware-11        | 1,389,481  | 862.3                      | 53.35             |
| BenchmarkNegroniLoggerMiddleware-11    | 887,360    | 1358                       | 33.87             |
| BenchmarkGorillaLoggerMiddleware-11    | 1,275,544  | 940.1                      | -                 |


## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.