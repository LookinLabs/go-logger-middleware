package logger

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// responseCapture is a custom writer that captures the response body and status code.
type responseCapture struct {
	writer     http.ResponseWriter
	body       *bytes.Buffer
	statusCode int
}

func (resp *responseCapture) Write(body []byte) (int, error) {
	resp.body.Write(body)
	return resp.writer.Write(body)
}

func (resp *responseCapture) WriteHeader(statusCode int) {
	resp.statusCode = statusCode
	resp.writer.WriteHeader(statusCode)
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

// Middleware is the actual middleware function.
func (lm *LoggerMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Start timer
		startTime := time.Now()

		// Generate a unique request ID
		requestID := uuid.New().String()

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

		// Calculate latency in milliseconds
		latency := time.Since(startTime).Seconds() * 1000

		// Log details
		lm.logger.Printf("Request details: client_ip=%s method=%s status_code=%d body_size=%d request_body=%s response_body=%s path=%s user_agent=%s referer=%s request_id=%s host=%s latency_ms=%.4fms",
			clientIP, method, statusCode, bodySize, string(sanitizedRequestBody), string(sanitizedResponseBody), path, userAgent, referer, requestID, host, latency)
	})
}

// sanitizeBody removes or masks sensitive fields from the body.
func (lm *LoggerMiddleware) sanitizeBody(body []byte) []byte {
	var data map[string]interface{}
	if err := Unmarshal(body, &data); err != nil {
		return body
	}
	for _, field := range lm.sensitiveFields {
		if _, exists := data[field]; exists {
			data[field] = "****"
		}
	}
	sanitizedBody, err := Marshal(data)
	if err != nil {
		return body
	}
	return sanitizedBody
}
