package logger

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestLoggerMiddleware(t *testing.T) {
	// Initialize the logger middleware
	sensitiveFields := []string{"password", "token"}
	appLogger := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	loggerMiddleware := NewLoggerMiddleware(sensitiveFields, appLogger)

	// Create a new HTTP mux (router)
	mux := http.NewServeMux()

	// Define a simple endpoint
	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	// Define another endpoint
	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"username":"john_doe","password":"secret"}`))
	})

	// Define an endpoint with a token
	mux.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"token":"1234567890"}`))
	})

	// Wrap the mux with the logger middleware
	handler := loggerMiddleware.Middleware(mux)

	// Create a test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// Define test cases
	tests := []struct {
		name       string
		endpoint   string
		wantStatus int
		wantBody   string
	}{
		{"HelloEndpoint", "/hello", http.StatusOK, "Hello, World!"},
		{"LoginEndpoint", "/login", http.StatusOK, `{"username":"john_doe","password":"****"}`},
		{"TokenEndpoint", "/token", http.StatusOK, `{"token":"****"}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := http.Get(server.URL + tt.endpoint)
			if err != nil {
				t.Fatalf("Failed to send request: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.wantStatus {
				t.Errorf("Expected status %d, got %d", tt.wantStatus, resp.StatusCode)
			}

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Failed to read response body: %v", err)
			}

			if !strings.Contains(string(body), tt.wantBody) {
				t.Errorf("Expected body to contain %q, got %q", tt.wantBody, string(body))
			}
		})
	}
}
