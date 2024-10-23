// package main

// import (
// 	"log"
// 	"net/http"
// 	"os"

// 	"github.com/lookinlabs/go-logger-middleware"

// 	"github.com/gin-gonic/gin"
// )

// func main() {
// 	// Initialize the logger middleware
// 	sensitiveFields := []string{"password", "token"}
// 	appLogger := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
// 	logger := logger.NewLoggerMiddleware(sensitiveFields, appLogger)

// 	// Create a Gin router
// 	r := gin.Default()

// 	// Use the middleware
// 	r.Use(gin.WrapH(logger.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		// Your handler logic here
// 		w.Write([]byte("Hello, World!"))
// 	}))))

// 	// Define a simple endpoint
// 	r.GET("/hello", func(c *gin.Context) {
// 		c.String(http.StatusOK, "Hello, World!")
// 	})

// 	// Start the server
// 	if err := r.Run(":8080"); err != nil {
// 		log.Fatalf("Failed to run server: %v", err)
// 	}
// }
