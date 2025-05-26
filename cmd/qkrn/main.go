package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/q4ow/qkrn/internal/api"
	"github.com/q4ow/qkrn/internal/config"
	"github.com/q4ow/qkrn/internal/store"
)

func main() {
	cfg := config.ParseFlags()
	log.Printf("Starting qkrn with config: %s", cfg)

	kvStore := store.NewMemoryStore()
	log.Printf("Initialized memory store")

	server := api.NewServer(kvStore, cfg.Port)

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
