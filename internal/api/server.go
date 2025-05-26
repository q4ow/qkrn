package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/q4ow/qkrn/internal/auth"
	"github.com/q4ow/qkrn/pkg/types"
)

type Server struct {
	store  types.Store
	port   int
	server *http.ServeMux
	auth   *auth.Authenticator
}

func NewServer(store types.Store, port int, authenticator *auth.Authenticator) *Server {
	s := &Server{
		store:  store,
		port:   port,
		server: http.NewServeMux(),
		auth:   authenticator,
	}

	s.setupRoutes()
	return s
}

func (s *Server) setupRoutes() {
	s.server.HandleFunc("/", s.handleRoot)
	s.server.HandleFunc("/health", s.handleHealth)
	s.server.HandleFunc("/keys", s.auth.Middleware(s.handleKeys))
	s.server.HandleFunc("/kv/", s.auth.Middleware(s.handleKeyValue))
}

func (s *Server) Start() error {
	addr := fmt.Sprintf(":%d", s.port)
	log.Printf("Starting server on %s", addr)
	return http.ListenAndServe(addr, s.server)
}

func (s *Server) handleRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	response := map[string]interface{}{
		"service":        "qkrn",
		"version":        "0.1.0",
		"status":         "running",
		"authentication": s.auth.IsEnabled(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	response := map[string]string{
		"status": "healthy",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *Server) handleKeys(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	keys := s.store.Keys()
	response := map[string][]string{
		"keys": keys,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *Server) handleKeyValue(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/kv/")
	if path == "" {
		http.Error(w, "Key is required", http.StatusBadRequest)
		return
	}

	key := path

	switch r.Method {
	case http.MethodGet:
		s.handleGet(w, r, key)
	case http.MethodPut, http.MethodPost:
		s.handleSet(w, r, key)
	case http.MethodDelete:
		s.handleDelete(w, r, key)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleGet(w http.ResponseWriter, r *http.Request, key string) {
	value, err := s.store.Get(key)
	if err != nil {
		if err == types.ErrKeyNotFound {
			s.sendErrorResponse(w, "Key not found", http.StatusNotFound)
			return
		}
		s.sendErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := types.Response{
		Success: true,
		Value:   value,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *Server) handleSet(w http.ResponseWriter, r *http.Request, key string) {
	var req types.Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.sendErrorResponse(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := s.store.Set(key, req.Value); err != nil {
		s.sendErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := types.Response{
		Success: true,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (s *Server) handleDelete(w http.ResponseWriter, r *http.Request, key string) {
	if err := s.store.Delete(key); err != nil {
		if err == types.ErrKeyNotFound {
			s.sendErrorResponse(w, "Key not found", http.StatusNotFound)
			return
		}
		s.sendErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := types.Response{
		Success: true,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *Server) sendErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	response := types.Response{
		Success: false,
		Error:   message,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}
