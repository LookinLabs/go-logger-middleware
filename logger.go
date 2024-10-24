package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"time"
)

// responseCapture is a custom writer that captures the response body and status code.
type responseCapture struct {
	writer     http.ResponseWriter
	body       *bytes.Buffer
	statusCode int
}

func (resp *responseCapture) Write(body []byte) (int, error) {
	return resp.body.Write(body)
}

func (resp *responseCapture) WriteHeader(statusCode int) {
	if resp.statusCode == 0 {
		resp.statusCode = statusCode
		resp.writer.WriteHeader(statusCode)
	}
}

func (resp *responseCapture) Header() http.Header {
	return resp.writer.Header()
}

// LoggerMiddleware is a struct that holds the configuration for the middleware.
type LoggerMiddleware struct {
	sensitiveFields []string
	logger          *log.Logger
}

// NewLoggerMiddleware creates a new LoggerMiddleware with the given sensitive fields and logger.
func NewLoggerMiddleware(sensitiveFields []string, logger *log.Logger) *LoggerMiddleware {
	return &LoggerMiddleware{sensitiveFields: sensitiveFields, logger: logger}
}

// generateRequestID generates a unique request ID.
func generateRequestID() string {
	return fmt.Sprintf("%d", rand.Int63())
}

// Middleware is the actual middleware function.
func (lm *LoggerMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Start timer
		startTime := time.Now()

		// Generate a unique request ID
		requestID := generateRequestID()

		// Read the request body
		var requestBody []byte
		if r.Body != nil {
			requestBody, _ = io.ReadAll(r.Body)
			r.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}
		// Sanitize the request body
		sanitizedRequestBody := lm.sanitizeBody(requestBody)

		// Create a custom response writer
		responseWriter := &responseCapture{body: bytes.NewBufferString(""), writer: w}

		// Process request
		next.ServeHTTP(responseWriter, r)

		// Get request details
		clientIP := r.RemoteAddr
		method := r.Method
		path := r.URL.Path
		userAgent := r.UserAgent()
		referer := r.Referer()
		host := r.Host

		// Get response details
		statusCode := responseWriter.statusCode
		if statusCode == 0 {
			statusCode = http.StatusOK
		}

		bodySize := responseWriter.body.Len()
		responseBody := responseWriter.body.String()

		// Sanitize the response body
		sanitizedResponseBody := lm.sanitizeBody([]byte(responseBody))

		// Write the sanitized response body to the response writer
		w.Header().Set("Content-Type", "application/json")
		w.Write(sanitizedResponseBody)

		// Calculate latency in milliseconds
		latency := time.Since(startTime).Seconds() * 1000

		// Log details in JSON format
		logDetails := map[string]interface{}{
			"client_ip":     clientIP,
			"method":        method,
			"status_code":   statusCode,
			"body_size":     bodySize,
			"request_body":  string(sanitizedRequestBody),
			"response_body": string(sanitizedResponseBody),
			"path":          path,
			"user_agent":    userAgent,
			"referer":       referer,
			"request_id":    requestID,
			"host":          host,
			"latency_ms":    fmt.Sprintf("%.4fms", latency),
		}

		logDetailsJSON, err := json.Marshal(logDetails)
		if err != nil {
			lm.logger.Printf("Error marshalling log details: %v", err)
		} else {
			lm.logger.Println(string(logDetailsJSON))
		}
	})
}

// sanitizeBody removes or masks sensitive fields from the body.
func (lm *LoggerMiddleware) sanitizeBody(body []byte) []byte {
	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return body
	}

	for _, field := range lm.sensitiveFields {
		if _, exists := data[field]; exists {
			data[field] = "****"
		}
	}

	sanitizedBody, err := json.Marshal(data)
	if err != nil {
		return body
	}

	return sanitizedBody
}
