package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/q4ow/qkrn/internal/api"
	"github.com/q4ow/qkrn/internal/auth"
	"github.com/q4ow/qkrn/internal/config"
	"github.com/q4ow/qkrn/internal/store"
)

func main() {
	cfg := config.LoadConfig()
	log.Printf("Starting qkrn with config: %s", cfg)

	if cfg.AuthEnabled && cfg.APIKey == "" {
		generatedKey, err := auth.GenerateAPIKey()
		if err != nil {
			log.Fatalf("Failed to generate API key: %v", err)
		}
		cfg.APIKey = generatedKey
		log.Printf("Generated API key: %s", generatedKey)
		log.Printf("IMPORTANT: Save this API key - it will be required for all API requests")
	}

	kvStore := store.NewMemoryStore()
	log.Printf("Initialized memory store")

	authenticator := auth.NewAuthenticator(cfg.AuthEnabled, cfg.APIKey)
	if cfg.AuthEnabled {
		log.Printf("Authentication enabled")
		if !authenticator.HasValidKey() {
			log.Printf("WARNING: Authentication enabled but no valid API key configured")
		}
	} else {
		log.Printf("Authentication disabled")
	}

	server := api.NewServer(kvStore, cfg.Port, authenticator)

	go func() {
		if err := server.Start(); err != nil {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
}
