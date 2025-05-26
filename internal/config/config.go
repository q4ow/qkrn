package config

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Config struct {
	NodeID      string `toml:"node_id"`
	Address     string `toml:"address"`
	Port        int    `toml:"port"`
	LogLevel    string `toml:"log_level"`
	AuthEnabled bool   `toml:"auth_enabled"`
	APIKey      string `toml:"api_key"`
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

func GetConfigPaths() []string {
	paths := []string{}

	paths = append(paths, "./config.toml")

	if homeDir, err := os.UserHomeDir(); err == nil {
		userConfigPath := filepath.Join(homeDir, ".config", "qkrn", "config.toml")
		paths = append(paths, userConfigPath)
	}

	return paths
}

func LoadFromFile(filename string) (*Config, error) {
	cfg := DefaultConfig()

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file does not exist: %s", filename)
	}

	if _, err := toml.DecodeFile(filename, cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file %s: %w", filename, err)
	}

	return cfg, nil
}

func LoadConfig() *Config {
	configFile := ""
	exportConfig := false

	for i, arg := range os.Args[1:] {
		if arg == "--config" || arg == "-config" {
			if i+1 < len(os.Args[1:]) {
				configFile = os.Args[i+2]
			}
		}
		if arg == "--export-config" || arg == "-export-config" {
			exportConfig = true
		}
	}

	cfg := DefaultConfig()

	if configFile != "" {
		if fileCfg, err := LoadFromFile(configFile); err != nil {
			fmt.Fprintf(os.Stderr, "Error loading config file %s: %v\n", configFile, err)
			os.Exit(1)
		} else {
			fmt.Printf("Loaded configuration from: %s\n", configFile)
			cfg = fileCfg
		}
	} else {
		for _, configPath := range GetConfigPaths() {
			if fileCfg, err := LoadFromFile(configPath); err == nil {
				fmt.Printf("Loaded configuration from: %s\n", configPath)
				cfg = fileCfg
				break
			}
		}
	}

	if exportConfig {
		flag.StringVar(&cfg.NodeID, "node-id", cfg.NodeID, "Node ID")
		flag.StringVar(&cfg.Address, "address", cfg.Address, "Node address")
		flag.IntVar(&cfg.Port, "port", cfg.Port, "Node port")
		flag.StringVar(&cfg.LogLevel, "log-level", cfg.LogLevel, "Log level (debug, info, warn, error)")
		flag.BoolVar(&cfg.AuthEnabled, "auth-enabled", cfg.AuthEnabled, "Enable authentication")
		flag.StringVar(&cfg.APIKey, "api-key", cfg.APIKey, "API key for authentication")
		flag.StringVar(&configFile, "config", "", "Path to config file")
		flag.BoolVar(&exportConfig, "export-config", false, "Export current configuration to ./config.toml")

		flag.Parse()

		if err := cfg.ExportToFile("./config.toml"); err != nil {
			fmt.Fprintf(os.Stderr, "Error exporting config: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Configuration exported to ./config.toml")
		os.Exit(0)
	}

	flag.StringVar(&cfg.NodeID, "node-id", cfg.NodeID, "Node ID")
	flag.StringVar(&cfg.Address, "address", cfg.Address, "Node address")
	flag.IntVar(&cfg.Port, "port", cfg.Port, "Node port")
	flag.StringVar(&cfg.LogLevel, "log-level", cfg.LogLevel, "Log level (debug, info, warn, error)")
	flag.BoolVar(&cfg.AuthEnabled, "auth-enabled", cfg.AuthEnabled, "Enable authentication")
	flag.StringVar(&cfg.APIKey, "api-key", cfg.APIKey, "API key for authentication")
	flag.StringVar(&configFile, "config", "", "Path to config file")
	flag.BoolVar(&exportConfig, "export-config", false, "Export current configuration to ./config.toml")

	flag.Parse()

	return cfg
}

func (c *Config) ExportToFile(filename string) error {
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create config file %s: %w", filename, err)
	}
	defer file.Close()

	encoder := toml.NewEncoder(file)
	if err := encoder.Encode(c); err != nil {
		return fmt.Errorf("failed to encode config to TOML: %w", err)
	}

	return nil
}

func (c *Config) String() string {
	return fmt.Sprintf("NodeID=%s, Address=%s:%d, LogLevel=%s, AuthEnabled=%t, APIKey=%s",
		c.NodeID, c.Address, c.Port, c.LogLevel, c.AuthEnabled, c.APIKey)
}
