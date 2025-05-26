package types

import "errors"

var (
	ErrKeyNotFound  = errors.New("key not found")
	ErrEmptyKey     = errors.New("key cannot be empty")
	ErrUnauthorized = errors.New("unauthorized")
	ErrInvalidToken = errors.New("invalid authentication token")
	ErrMissingToken = errors.New("missing authentication token")
)

type Store interface {
	Get(key string) (string, error)
	Set(key, value string) error
	Delete(key string) error
	Keys() []string
}

type Node struct {
	ID      string `json:"id"`
	Address string `json:"address"`
	Port    int    `json:"port"`
}

type Request struct {
	Key   string `json:"key"`
	Value string `json:"value,omitempty"`
}

type Response struct {
	Success bool   `json:"success"`
	Value   string `json:"value,omitempty"`
	Error   string `json:"error,omitempty"`
}

type AuthRequest struct {
	Token string `json:"token,omitempty"`
}

type AuthResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}
