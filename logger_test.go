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
	logger := log.New(&logBuffer, "", 0) // No log prefix

	// Create a LoggerMiddleware instance
	lm := NewLoggerMiddleware([]string{"password", "token"}, logger)

	// Create a test HTTP handler
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(`{"message": "success", "password": "secret", "token": "12345"}`)); err != nil {
			t.Fatalf("failed to write response: %v", err)
		}
	})

	// Wrap the test handler with the middleware
	wrappedHandler := lm.Middleware(testHandler)

	// Create a test HTTP request
	req := httptest.NewRequest(http.MethodPost, "http://example.com/foo", strings.NewReader(`{"username": "john", "password": "secret", "token": "12345"}`))
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
	logLines := strings.Split(logOutput, "\n")
	for _, line := range logLines {
		if line == "" {
			continue
		}
		var logDetails map[string]interface{}
		if err := json.Unmarshal([]byte(line), &logDetails); err != nil {
			t.Fatalf("failed to parse log output: %v", err)
		}

		// Check log details
		expectedRequestBody := map[string]interface{}{
			"username": "john",
			"password": "****",
			"token":    "****",
		}
		var actualRequestBody map[string]interface{}
		if err := json.Unmarshal([]byte(logDetails["request_body"].(string)), &actualRequestBody); err != nil {
			t.Fatalf("failed to parse request body in log: %v", err)
		}
		if !equal(actualRequestBody, expectedRequestBody) {
			t.Errorf("log output does not mask sensitive fields in request body: %v", logDetails["request_body"])
		}

		expectedResponseBody := map[string]interface{}{
			"message":  "success",
			"password": "****",
			"token":    "****",
		}
		var actualResponseBody map[string]interface{}
		if err := json.Unmarshal([]byte(logDetails["response_body"].(string)), &actualResponseBody); err != nil {
			t.Fatalf("failed to parse response body in log: %v", err)
		}
		if !equal(actualResponseBody, expectedResponseBody) {
			t.Errorf("log output does not mask sensitive fields in response body: %v", logDetails["response_body"])
		}
	}
}
