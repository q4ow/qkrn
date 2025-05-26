package config

import (
	"flag"
	"fmt"
	"os"
)

type Config struct {
	NodeID   string
	Address  string
	Port     int
	LogLevel string
}

func DefaultConfig() *Config {
	hostname, _ := os.Hostname()
	return &Config{
		NodeID:   hostname,
		Address:  "localhost",
		Port:     8080,
		LogLevel: "info",
	}
}

func ParseFlags() *Config {
	cfg := DefaultConfig()

	flag.StringVar(&cfg.NodeID, "node-id", cfg.NodeID, "Node ID")
	flag.StringVar(&cfg.Address, "address", cfg.Address, "Node address")
	flag.IntVar(&cfg.Port, "port", cfg.Port, "Node port")
	flag.StringVar(&cfg.LogLevel, "log-level", cfg.LogLevel, "Log level (debug, info, warn, error)")
	flag.Parse()

	return cfg
}

func (c *Config) String() string {
	return fmt.Sprintf("NodeID=%s, Address=%s:%d, LogLevel=%s",
		c.NodeID, c.Address, c.Port, c.LogLevel)
}
