// package examples

// import (
// 	"log"
// 	"net/http"

// 	"go-logger-middleware"

// 	"github.com/gin-gonic/gin"
// )

// func main() {
// 	// Initialize the logger middleware
// 	sensitiveFields := []string{"password", "token"}
// 	loggerMiddleware := logger.NewLoggerMiddleware(sensitiveFields)

// 	// Create a Gin router
// 	r := gin.Default()

// 	// Use the middleware
// 	r.Use(gin.WrapH(loggerMiddleware.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
