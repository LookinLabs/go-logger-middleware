package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// LoggerMiddleware is a struct that holds the configuration for the middleware.
type LoggerMiddleware struct {
	sensitiveFields []string
}

// NewLoggerMiddleware creates a new LoggerMiddleware with the given sensitive fields.
func NewLoggerMiddleware(sensitiveFields []string) *LoggerMiddleware {
	return &LoggerMiddleware{sensitiveFields: sensitiveFields}
}

// Middleware is the actual middleware function.
func (lm *LoggerMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Start timer
		startTime := time.Now()

		// Read the request body
		var requestBody []byte
		if r.Body != nil {
			requestBody, _ = io.ReadAll(r.Body)
			r.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}
		// Sanitize the request body
		sanitizedRequestBody := lm.sanitizeBody(requestBody)

		// Create a custom response writer
		customWriter := &responseWriter{body: bytes.NewBufferString(""), ResponseWriter: w}
		w = customWriter

		// Process request
		next.ServeHTTP(customWriter, r)

		// Calculate latency
		latency := time.Since(startTime).Seconds() * 1000 // Convert to milliseconds

		// Get request details
		clientIP := r.RemoteAddr
		method := r.Method
		path := r.URL.Path
		userAgent := r.UserAgent()
		referer := r.Referer()
		requestID := r.Header.Get("X-Request-ID")
		host := r.Host

		// Get response details
		statusCode := customWriter.statusCode
		if statusCode == 0 {
			statusCode = http.StatusOK // Default to 200 OK if not set
		}
		bodySize := customWriter.body.Len()
		responseBody := customWriter.body.String()

		// Sanitize the response body
		sanitizedResponseBody := lm.sanitizeBody([]byte(responseBody))

		// Log details
		log.Printf("Request details: client_ip=%s method=%s status_code=%d body_size=%d request_body=%s response_body=%s path=%s user_agent=%s referer=%s request_id=%s host=%s latency_ms=%.4fms",
			clientIP, method, statusCode, bodySize, string(sanitizedRequestBody), string(sanitizedResponseBody), path, userAgent, referer, requestID, host, latency)
	})
}

// sanitizeBody removes or masks sensitive fields from the body.
func (lm *LoggerMiddleware) sanitizeBody(body []byte) []byte {
	var data map[string]interface{}
	if err := customUnmarshal(body, &data); err != nil {
		return body
	}
	for _, field := range lm.sensitiveFields {
		if _, exists := data[field]; exists {
			data[field] = "****"
		}
	}
	sanitizedBody, err := customMarshal(data)
	if err != nil {
		return body
	}
	return sanitizedBody
}

// customMarshal is a custom JSON marshalling function.
func customMarshal(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteByte('{')
	first := true
	for key, value := range v.(map[string]interface{}) {
		if !first {
			buf.WriteByte(',')
		}
		first = false
		buf.WriteString(`"` + key + `":`)
		switch value := value.(type) {
		case string:
			buf.WriteString(`"` + value + `"`)
		case int:
			buf.WriteString(fmt.Sprintf("%d", value))
		case float64:
			buf.WriteString(fmt.Sprintf("%f", value))
		default:
			buf.WriteString(`null`)
		}
	}
	buf.WriteByte('}')
	return buf.Bytes(), nil
}

// customUnmarshal is a custom JSON unmarshalling function.
func customUnmarshal(data []byte, v interface{}) error {
	var result map[string]interface{}
	err := json.Unmarshal(data, &result)
	if err != nil {
		return err
	}
	*v.(*map[string]interface{}) = result
	return nil
}

// responseWriter is a custom writer that captures the response body and status code.
type responseWriter struct {
	http.ResponseWriter
	body       *bytes.Buffer
	statusCode int
}

func (rw *responseWriter) Write(body []byte) (int, error) {
	rw.body.Write(body)
	return rw.ResponseWriter.Write(body)
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}
