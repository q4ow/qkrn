package config

import (
	"flag"
	"fmt"
	"os"
)

type Config struct {
	NodeID      string
	Address     string
	Port        int
	LogLevel    string
	AuthEnabled bool
	APIKey      string
}

func DefaultConfig() *Config {
	hostname, _ := os.Hostname()
	return &Config{
		NodeID:      hostname,
		Address:     "localhost",
		Port:        8080,
		LogLevel:    "info",
		AuthEnabled: false,
		APIKey:      "",
	}
}

func ParseFlags() *Config {
	cfg := DefaultConfig()

	flag.StringVar(&cfg.NodeID, "node-id", cfg.NodeID, "Node ID")
	flag.StringVar(&cfg.Address, "address", cfg.Address, "Node address")
	flag.IntVar(&cfg.Port, "port", cfg.Port, "Node port")
	flag.StringVar(&cfg.LogLevel, "log-level", cfg.LogLevel, "Log level (debug, info, warn, error)")
	flag.BoolVar(&cfg.AuthEnabled, "auth-enabled", cfg.AuthEnabled, "Enable authentication")
	flag.StringVar(&cfg.APIKey, "api-key", cfg.APIKey, "API key for authentication")
	flag.Parse()

	return cfg
}

func (c *Config) String() string {
	return fmt.Sprintf("NodeID=%s, Address=%s:%d, LogLevel=%s, AuthEnabled=%t, APIKey=%s",
		c.NodeID, c.Address, c.Port, c.LogLevel, c.AuthEnabled, c.APIKey)
}
