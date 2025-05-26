package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"net/http"
	"strings"

	"github.com/q4ow/qkrn/pkg/types"
)

type Authenticator struct {
	enabled bool
	apiKey  string
}

func NewAuthenticator(enabled bool, apiKey string) *Authenticator {
	return &Authenticator{
		enabled: enabled,
		apiKey:  apiKey,
	}
}

func GenerateAPIKey() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (a *Authenticator) Middleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !a.enabled {
			next.ServeHTTP(w, r)
			return
		}

		if r.URL.Path == "/health" || r.URL.Path == "/" {
			next.ServeHTTP(w, r)
			return
		}

		token := a.extractToken(r)
		if token == "" {
			a.sendAuthError(w, "Missing authentication token", http.StatusUnauthorized)
			return
		}

		if !a.validateToken(token) {
			a.sendAuthError(w, "Invalid authentication token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	}
}

func (a *Authenticator) extractToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
			return parts[1]
		}
	}

	if apiKey := r.Header.Get("X-API-Key"); apiKey != "" {
		return apiKey
	}

	return r.URL.Query().Get("api_key")
}

func (a *Authenticator) validateToken(token string) bool {
	if a.apiKey == "" {
		return false
	}

	return subtle.ConstantTimeCompare([]byte(token), []byte(a.apiKey)) == 1
}

func (a *Authenticator) sendAuthError(w http.ResponseWriter, message string, statusCode int) {
	response := types.Response{
		Success: false,
		Error:   message,
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("WWW-Authenticate", "Bearer")
	w.WriteHeader(statusCode)

	_ = writeJSON(w, response)
}

func writeJSON(w http.ResponseWriter, v interface{}) error {
	if resp, ok := v.(types.Response); ok {
		if resp.Error != "" {
			_, err := w.Write([]byte(`{"success":false,"error":"` + resp.Error + `"}`))
			return err
		}
		_, err := w.Write([]byte(`{"success":true}`))
		return err
	}
	return nil
}

func (a *Authenticator) IsEnabled() bool {
	return a.enabled
}

func (a *Authenticator) HasValidKey() bool {
	return a.apiKey != ""
}
