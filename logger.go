package logger

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"sync"
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
	bufferPool      sync.Pool
}

// NewLoggerMiddleware creates a new LoggerMiddleware with the given sensitive fields and logger.
func NewLoggerMiddleware(sensitiveFields []string, logger *log.Logger) *LoggerMiddleware {
	return &LoggerMiddleware{
		sensitiveFields: sensitiveFields,
		logger:          logger,
		bufferPool: sync.Pool{
			New: func() interface{} {
				return new(bytes.Buffer)
			},
		},
	}
}

// generateRequestID generates a unique request ID.
func generateRequestID() string {
	return fmt.Sprintf("%d", rand.Int63())
}

// Middleware is the actual middleware function.
func (lm *LoggerMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, req *http.Request) {
		// Create a custom response writer
		responseWriter := &responseCapture{body: lm.bufferPool.Get().(*bytes.Buffer), writer: response}
		defer lm.bufferPool.Put(responseWriter.body)
		responseWriter.body.Reset()

		var (
			startTime    = time.Now()
			requestID    = generateRequestID()
			requestBody  []byte
			bodySize     = responseWriter.body.Len()
			responseBody = responseWriter.body.Bytes()

			// Get request details
			clientIP  = req.RemoteAddr
			method    = req.Method
			path      = req.URL.Path
			userAgent = req.UserAgent()
			host      = req.Host

			sanitizedRequestBody  = lm.sanitizeBody(requestBody)
			sanitizedResponseBody = lm.sanitizeBody(responseBody)
		)

		if req.Body != nil {
			requestBody, _ = io.ReadAll(req.Body)
			req.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// Process request
		next.ServeHTTP(responseWriter, req)

		// Get response details
		statusCode := responseWriter.statusCode
		if statusCode == 0 {
			statusCode = http.StatusOK
		}

		// Write the sanitized response body to the response writer
		responseWriter.WriteHeader(statusCode)
		if _, err := response.Write(sanitizedResponseBody); err != nil {
			lm.logger.Printf("Error writing response: %v", err)
		}

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
			"request_id":    requestID,
			"host":          host,
			"latency_ms":    fmt.Sprintf("%.4fms", latency),
		}

		// Convert logDetails to a slice of KeyValuePair
		logDetailsPairs := MapToKeyValuePairs(logDetails)

		// Marshal logDetails using the custom Marshal function
		logDetailsJSON, err := Marshal(logDetailsPairs)
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
	if err := Unmarshal(body, &data); err != nil {
		return body
	}

	for _, field := range lm.sensitiveFields {
		if _, exists := data[field]; exists {
			data[field] = "****"
		}
	}

	sanitizedBodyPairs := MapToKeyValuePairs(data)
	sanitizedBody, err := Marshal(sanitizedBodyPairs)
	if err != nil {
		return body
	}

	return sanitizedBody
}
