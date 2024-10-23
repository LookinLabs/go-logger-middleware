package logger

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLoggerMiddleware(t *testing.T) {
	// Create a buffer to capture logs
	var logBuffer bytes.Buffer
	logger := log.New(&logBuffer, "", log.LstdFlags)

	// Create a LoggerMiddleware instance
	lm := NewLoggerMiddleware([]string{"password", "token"}, logger)

	// Create a test HTTP handler
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "success", "password": "secret", "token": "12345"}`))
	})

	// Wrap the test handler with the middleware
	wrappedHandler := lm.Middleware(testHandler)

	// Create a test HTTP request
	req := httptest.NewRequest("POST", "http://example.com/foo", strings.NewReader(`{"username": "john", "password": "secret", "token": "12345"}`))
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Send the request to the wrapped handler
	wrappedHandler.ServeHTTP(rr, req)

	// Verify the response
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Parse the response body
	var responseBody map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &responseBody); err != nil {
		t.Fatalf("failed to parse response body: %v", err)
	}

	// Verify the sanitized fields in the response body
	if responseBody["password"] != "****" {
		t.Errorf("handler returned unexpected password: got %v want %v", responseBody["password"], "****")
	}
	if responseBody["token"] != "****" {
		t.Errorf("handler returned unexpected token: got %v want %v", responseBody["token"], "****")
	}

	// Verify the logs
	logOutput := logBuffer.String()
	if !strings.Contains(logOutput, "Request details:") {
		t.Errorf("log output does not contain expected log entry: %v", logOutput)
	}
	if !strings.Contains(logOutput, `"password":"****"`) {
		t.Errorf("log output does not mask sensitive fields: %v", logOutput)
	}
	if !strings.Contains(logOutput, `"token":"****"`) {
		t.Errorf("log output does not mask sensitive fields: %v", logOutput)
	}
}
