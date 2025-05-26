package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGenerateAPIKey(t *testing.T) {
	key1, err := GenerateAPIKey()
	if err != nil {
		t.Fatalf("GenerateAPIKey failed: %v", err)
	}

	if len(key1) != 64 {
		t.Errorf("Expected key length 64, got %d", len(key1))
	}

	key2, err := GenerateAPIKey()
	if err != nil {
		t.Fatalf("GenerateAPIKey failed: %v", err)
	}

	if key1 == key2 {
		t.Error("Generated keys should be unique")
	}
}

func TestAuthenticator_Middleware(t *testing.T) {
	testAPIKey := "test-api-key-123"

	tests := []struct {
		name           string
		enabled        bool
		apiKey         string
		requestPath    string
		authHeader     string
		expectedStatus int
	}{
		{
			name:           "auth disabled - should pass",
			enabled:        false,
			apiKey:         testAPIKey,
			requestPath:    "/kv/test",
			authHeader:     "",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "health endpoint - should pass without auth",
			enabled:        true,
			apiKey:         testAPIKey,
			requestPath:    "/health",
			authHeader:     "",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "root endpoint - should pass without auth",
			enabled:        true,
			apiKey:         testAPIKey,
			requestPath:    "/",
			authHeader:     "",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "missing token - should fail",
			enabled:        true,
			apiKey:         testAPIKey,
			requestPath:    "/kv/test",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "invalid token - should fail",
			enabled:        true,
			apiKey:         testAPIKey,
			requestPath:    "/kv/test",
			authHeader:     "Bearer invalid-token",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "valid bearer token - should pass",
			enabled:        true,
			apiKey:         testAPIKey,
			requestPath:    "/kv/test",
			authHeader:     "Bearer " + testAPIKey,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "valid x-api-key header - should pass",
			enabled:        true,
			apiKey:         testAPIKey,
			requestPath:    "/kv/test",
			authHeader:     "",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			auth := NewAuthenticator(tt.enabled, tt.apiKey)

			handler := auth.Middleware(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			req := httptest.NewRequest("GET", tt.requestPath, nil)

			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			if tt.name == "valid x-api-key header - should pass" {
				req.Header.Set("X-API-Key", tt.apiKey)
			}

			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestAuthenticator_extractToken(t *testing.T) {
	auth := NewAuthenticator(true, "test-key")

	tests := []struct {
		name          string
		authHeader    string
		apiKeyHeader  string
		queryParam    string
		expectedToken string
	}{
		{
			name:          "bearer token",
			authHeader:    "Bearer test-token-123",
			expectedToken: "test-token-123",
		},
		{
			name:          "bearer token case insensitive",
			authHeader:    "bearer test-token-123",
			expectedToken: "test-token-123",
		},
		{
			name:          "x-api-key header",
			apiKeyHeader:  "test-token-123",
			expectedToken: "test-token-123",
		},
		{
			name:          "query parameter",
			queryParam:    "test-token-123",
			expectedToken: "test-token-123",
		},
		{
			name:          "no token",
			expectedToken: "",
		},
		{
			name:          "invalid auth header format",
			authHeader:    "InvalidFormat",
			expectedToken: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/test", nil)

			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			if tt.apiKeyHeader != "" {
				req.Header.Set("X-API-Key", tt.apiKeyHeader)
			}
			if tt.queryParam != "" {
				q := req.URL.Query()
				q.Add("api_key", tt.queryParam)
				req.URL.RawQuery = q.Encode()
			}

			token := auth.extractToken(req)
			if token != tt.expectedToken {
				t.Errorf("Expected token '%s', got '%s'", tt.expectedToken, token)
			}
		})
	}
}

func TestAuthenticator_validateToken(t *testing.T) {
	validKey := "valid-api-key-123"
	auth := NewAuthenticator(true, validKey)

	tests := []struct {
		name     string
		token    string
		expected bool
	}{
		{
			name:     "valid token",
			token:    validKey,
			expected: true,
		},
		{
			name:     "invalid token",
			token:    "invalid-token",
			expected: false,
		},
		{
			name:     "empty token",
			token:    "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := auth.validateToken(tt.token)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestAuthenticator_NoAPIKey(t *testing.T) {
	auth := NewAuthenticator(true, "")

	if auth.validateToken("any-token") {
		t.Error("Should not validate any token when no API key is configured")
	}

	if auth.HasValidKey() {
		t.Error("Should report no valid key when API key is empty")
	}
}
