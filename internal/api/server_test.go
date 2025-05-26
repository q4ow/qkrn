package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/q4ow/qkrn/internal/auth"
	"github.com/q4ow/qkrn/internal/store"
	"github.com/q4ow/qkrn/pkg/types"
)

func setupTestServer(authEnabled bool, apiKey string) *Server {
	kvStore := store.NewMemoryStore()
	var authenticator *auth.Authenticator

	if authEnabled {
		authenticator = auth.NewAuthenticator(true, apiKey)
	} else {
		authenticator = auth.NewAuthenticator(false, "")
	}

	return NewServer(kvStore, 8080, authenticator)
}

func TestNewServer(t *testing.T) {
	kvStore := store.NewMemoryStore()
	authenticator := auth.NewAuthenticator(false, "")
	server := NewServer(kvStore, 8080, authenticator)

	if server == nil {
		t.Fatal("Expected server to be created, got nil")
	}

	if server.store != kvStore {
		t.Error("Expected server store to match provided store")
	}

	if server.port != 8080 {
		t.Errorf("Expected server port to be 8080, got %d", server.port)
	}

	if server.auth != authenticator {
		t.Error("Expected server authenticator to match provided authenticator")
	}
}

func TestHandleRoot(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		path           string
		authEnabled    bool
		expectedStatus int
		checkAuth      bool
	}{
		{
			name:           "GET root without auth",
			method:         "GET",
			path:           "/",
			authEnabled:    false,
			expectedStatus: http.StatusOK,
			checkAuth:      false,
		},
		{
			name:           "GET root with auth enabled",
			method:         "GET",
			path:           "/",
			authEnabled:    true,
			expectedStatus: http.StatusOK,
			checkAuth:      true,
		},
		{
			name:           "GET invalid path",
			method:         "GET",
			path:           "/invalid",
			authEnabled:    false,
			expectedStatus: http.StatusNotFound,
			checkAuth:      false,
		},
		{
			name:           "POST root not allowed",
			method:         "POST",
			path:           "/",
			authEnabled:    false,
			expectedStatus: http.StatusMethodNotAllowed,
			checkAuth:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := setupTestServer(tt.authEnabled, "test-key")
			req := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()

			server.handleRoot(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedStatus == http.StatusOK {
				var response map[string]interface{}
				if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}

				if response["service"] != "qkrn" {
					t.Errorf("Expected service to be 'qkrn', got %v", response["service"])
				}

				if response["version"] != "0.1.0" {
					t.Errorf("Expected version to be '0.1.0', got %v", response["version"])
				}

				if response["status"] != "running" {
					t.Errorf("Expected status to be 'running', got %v", response["status"])
				}

				if tt.checkAuth {
					if response["authentication"] != true {
						t.Errorf("Expected authentication to be true, got %v", response["authentication"])
					}
				} else {
					if response["authentication"] != false {
						t.Errorf("Expected authentication to be false, got %v", response["authentication"])
					}
				}
			}
		})
	}
}

func TestHandleHealth(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		expectedStatus int
	}{
		{
			name:           "GET health check",
			method:         "GET",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "POST health not allowed",
			method:         "POST",
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "PUT health not allowed",
			method:         "PUT",
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "DELETE health not allowed",
			method:         "DELETE",
			expectedStatus: http.StatusMethodNotAllowed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := setupTestServer(false, "")
			req := httptest.NewRequest(tt.method, "/health", nil)
			w := httptest.NewRecorder()

			server.handleHealth(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedStatus == http.StatusOK {
				var response map[string]string
				if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}

				if response["status"] != "healthy" {
					t.Errorf("Expected status to be 'healthy', got %s", response["status"])
				}
			}
		})
	}
}

func TestHandleKeys(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		setupKeys      map[string]string
		expectedStatus int
		expectedKeys   []string
	}{
		{
			name:           "GET empty keys list",
			method:         "GET",
			setupKeys:      map[string]string{},
			expectedStatus: http.StatusOK,
			expectedKeys:   []string{},
		},
		{
			name:   "GET keys with data",
			method: "GET",
			setupKeys: map[string]string{
				"key1": "value1",
				"key2": "value2",
			},
			expectedStatus: http.StatusOK,
			expectedKeys:   []string{"key1", "key2"},
		},
		{
			name:           "POST keys not allowed",
			method:         "POST",
			setupKeys:      map[string]string{},
			expectedStatus: http.StatusMethodNotAllowed,
			expectedKeys:   nil,
		},
		{
			name:           "PUT keys not allowed",
			method:         "PUT",
			setupKeys:      map[string]string{},
			expectedStatus: http.StatusMethodNotAllowed,
			expectedKeys:   nil,
		},
		{
			name:           "DELETE keys not allowed",
			method:         "DELETE",
			setupKeys:      map[string]string{},
			expectedStatus: http.StatusMethodNotAllowed,
			expectedKeys:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := setupTestServer(false, "")

			for key, value := range tt.setupKeys {
				server.store.Set(key, value)
			}

			req := httptest.NewRequest(tt.method, "/keys", nil)
			w := httptest.NewRecorder()

			server.handleKeys(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedStatus == http.StatusOK {
				var response map[string][]string
				if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}

				keys := response["keys"]
				if len(keys) != len(tt.expectedKeys) {
					t.Errorf("Expected %d keys, got %d", len(tt.expectedKeys), len(keys))
				}

				for _, expectedKey := range tt.expectedKeys {
					found := false
					for _, key := range keys {
						if key == expectedKey {
							found = true
							break
						}
					}
					if !found {
						t.Errorf("Expected key %s not found in response", expectedKey)
					}
				}
			}
		})
	}
}

func TestHandleKeyValue(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		path           string
		body           interface{}
		setupData      map[string]string
		expectedStatus int
		expectedValue  string
		expectedError  string
	}{
		{
			name:           "GET existing key",
			method:         "GET",
			path:           "/kv/testkey",
			setupData:      map[string]string{"testkey": "testvalue"},
			expectedStatus: http.StatusOK,
			expectedValue:  "testvalue",
		},
		{
			name:           "GET non-existent key",
			method:         "GET",
			path:           "/kv/nonexistent",
			setupData:      map[string]string{},
			expectedStatus: http.StatusNotFound,
			expectedError:  "Key not found",
		},

		{
			name:           "PUT new key",
			method:         "PUT",
			path:           "/kv/newkey",
			body:           types.Request{Value: "newvalue"},
			setupData:      map[string]string{},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "PUT update existing key",
			method:         "PUT",
			path:           "/kv/existing",
			body:           types.Request{Value: "updatedvalue"},
			setupData:      map[string]string{"existing": "oldvalue"},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "PUT with invalid JSON",
			method:         "PUT",
			path:           "/kv/testkey",
			body:           "invalid json",
			setupData:      map[string]string{},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid JSON",
		},

		{
			name:           "POST new key",
			method:         "POST",
			path:           "/kv/postkey",
			body:           types.Request{Value: "postvalue"},
			setupData:      map[string]string{},
			expectedStatus: http.StatusCreated,
		},

		{
			name:           "DELETE existing key",
			method:         "DELETE",
			path:           "/kv/deletekey",
			setupData:      map[string]string{"deletekey": "deletevalue"},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "DELETE non-existent key",
			method:         "DELETE",
			path:           "/kv/nonexistent",
			setupData:      map[string]string{},
			expectedStatus: http.StatusNotFound,
			expectedError:  "Key not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := setupTestServer(false, "")

			for key, value := range tt.setupData {
				server.store.Set(key, value)
			}

			var reqBody *bytes.Buffer
			if tt.body != nil {
				if str, ok := tt.body.(string); ok {
					reqBody = bytes.NewBufferString(str)
				} else {
					bodyBytes, _ := json.Marshal(tt.body)
					reqBody = bytes.NewBuffer(bodyBytes)
				}
			} else {
				reqBody = bytes.NewBuffer([]byte{})
			}

			req := httptest.NewRequest(tt.method, tt.path, reqBody)
			if tt.body != nil && tt.method != "GET" && tt.method != "DELETE" {
				req.Header.Set("Content-Type", "application/json")
			}
			w := httptest.NewRecorder()

			server.handleKeyValue(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			var response types.Response
			if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
				t.Fatalf("Failed to decode response: %v", err)
			}

			if tt.expectedValue != "" {
				if !response.Success {
					t.Errorf("Expected success response, got failure: %s", response.Error)
				}
				if response.Value != tt.expectedValue {
					t.Errorf("Expected value %s, got %s", tt.expectedValue, response.Value)
				}
			}

			if tt.expectedError != "" {
				if response.Success {
					t.Errorf("Expected error response, got success")
				}
				if response.Error != tt.expectedError {
					t.Errorf("Expected error %s, got %s", tt.expectedError, response.Error)
				}
			}

			if tt.expectedStatus == http.StatusCreated {
				if !response.Success {
					t.Errorf("Expected success response, got failure: %s", response.Error)
				}
			}
		})
	}
}

func TestHandleKeyValueUnsupportedMethod(t *testing.T) {
	server := setupTestServer(false, "")
	req := httptest.NewRequest("PATCH", "/kv/testkey", nil)
	w := httptest.NewRecorder()

	server.handleKeyValue(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
	}
}

func TestAuthenticationIntegration(t *testing.T) {
	apiKey := "test-api-key-123"
	server := setupTestServer(true, apiKey)

	tests := []struct {
		name           string
		endpoint       string
		method         string
		authHeader     string
		expectedStatus int
	}{
		{
			name:           "Access protected endpoint without auth",
			endpoint:       "/keys",
			method:         "GET",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Access protected endpoint with valid auth",
			endpoint:       "/keys",
			method:         "GET",
			authHeader:     "Bearer " + apiKey,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Access protected endpoint with invalid auth",
			endpoint:       "/keys",
			method:         "GET",
			authHeader:     "Bearer invalid-key",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Access root endpoint without auth (should work)",
			endpoint:       "/",
			method:         "GET",
			authHeader:     "",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Access health endpoint without auth (should work)",
			endpoint:       "/health",
			method:         "GET",
			authHeader:     "",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.endpoint, nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			w := httptest.NewRecorder()

			switch tt.endpoint {
			case "/":
				server.handleRoot(w, req)
			case "/health":
				server.handleHealth(w, req)
			case "/keys":
				server.auth.Middleware(server.handleKeys)(w, req)
			}

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestSendErrorResponse(t *testing.T) {
	server := setupTestServer(false, "")
	w := httptest.NewRecorder()

	server.sendErrorResponse(w, "Test error message", http.StatusBadRequest)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	var response types.Response
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response.Success {
		t.Error("Expected success to be false")
	}

	if response.Error != "Test error message" {
		t.Errorf("Expected error message 'Test error message', got '%s'", response.Error)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type 'application/json', got '%s'", contentType)
	}
}

func TestContentTypeHeaders(t *testing.T) {
	server := setupTestServer(false, "")

	endpoints := []struct {
		name    string
		path    string
		method  string
		handler func(http.ResponseWriter, *http.Request)
	}{
		{"root", "/", "GET", server.handleRoot},
		{"health", "/health", "GET", server.handleHealth},
		{"keys", "/keys", "GET", server.handleKeys},
	}

	for _, endpoint := range endpoints {
		t.Run(endpoint.name, func(t *testing.T) {
			req := httptest.NewRequest(endpoint.method, endpoint.path, nil)
			w := httptest.NewRecorder()

			endpoint.handler(w, req)

			contentType := w.Header().Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("Expected Content-Type 'application/json', got '%s'", contentType)
			}
		})
	}
}
