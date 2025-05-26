package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.Address != "localhost" {
		t.Errorf("Expected default address to be 'localhost', got '%s'", cfg.Address)
	}

	if cfg.Port != 8080 {
		t.Errorf("Expected default port to be 8080, got %d", cfg.Port)
	}

	if cfg.LogLevel != "info" {
		t.Errorf("Expected default log level to be 'info', got '%s'", cfg.LogLevel)
	}

	if cfg.AuthEnabled != false {
		t.Errorf("Expected default auth enabled to be false, got %t", cfg.AuthEnabled)
	}

	if cfg.APIKey != "" {
		t.Errorf("Expected default API key to be empty, got '%s'", cfg.APIKey)
	}
}

func TestGetConfigPaths(t *testing.T) {
	paths := GetConfigPaths()

	if len(paths) < 2 {
		t.Errorf("Expected at least 2 config paths, got %d", len(paths))
	}

	if paths[0] != "./config.toml" {
		t.Errorf("Expected first path to be './config.toml', got '%s'", paths[0])
	}

	homeDir, err := os.UserHomeDir()
	if err == nil {
		expectedUserPath := filepath.Join(homeDir, ".config", "qkrn", "config.toml")
		if paths[1] != expectedUserPath {
			t.Errorf("Expected second path to be '%s', got '%s'", expectedUserPath, paths[1])
		}
	}
}

func TestLoadFromFile(t *testing.T) {
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "test-config.toml")

	configContent := `node_id = "test-node"
address = "0.0.0.0"
port = 9999
log_level = "debug"
auth_enabled = true
api_key = "test-api-key"`

	if err := os.WriteFile(configFile, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	cfg, err := LoadFromFile(configFile)
	if err != nil {
		t.Fatalf("Failed to load config from file: %v", err)
	}

	if cfg.NodeID != "test-node" {
		t.Errorf("Expected NodeID to be 'test-node', got '%s'", cfg.NodeID)
	}

	if cfg.Address != "0.0.0.0" {
		t.Errorf("Expected Address to be '0.0.0.0', got '%s'", cfg.Address)
	}

	if cfg.Port != 9999 {
		t.Errorf("Expected Port to be 9999, got %d", cfg.Port)
	}

	if cfg.LogLevel != "debug" {
		t.Errorf("Expected LogLevel to be 'debug', got '%s'", cfg.LogLevel)
	}

	if cfg.AuthEnabled != true {
		t.Errorf("Expected AuthEnabled to be true, got %t", cfg.AuthEnabled)
	}

	if cfg.APIKey != "test-api-key" {
		t.Errorf("Expected APIKey to be 'test-api-key', got '%s'", cfg.APIKey)
	}
}

func TestLoadFromFile_NotExists(t *testing.T) {
	_, err := LoadFromFile("/nonexistent/config.toml")
	if err == nil {
		t.Error("Expected error when loading non-existent config file")
	}
}

func TestLoadFromFile_InvalidTOML(t *testing.T) {
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "invalid-config.toml")

	invalidContent := `node_id = "test
invalid toml content`

	if err := os.WriteFile(configFile, []byte(invalidContent), 0644); err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	_, err := LoadFromFile(configFile)
	if err == nil {
		t.Error("Expected error when loading invalid TOML config file")
	}
}

func TestExportToFile(t *testing.T) {
	cfg := &Config{
		NodeID:      "export-test",
		Address:     "127.0.0.1",
		Port:        8888,
		LogLevel:    "warn",
		AuthEnabled: true,
		APIKey:      "export-key",
	}

	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "exported-config.toml")

	if err := cfg.ExportToFile(configFile); err != nil {
		t.Fatalf("Failed to export config: %v", err)
	}

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		t.Error("Exported config file does not exist")
	}

	loadedCfg, err := LoadFromFile(configFile)
	if err != nil {
		t.Fatalf("Failed to load exported config: %v", err)
	}

	if loadedCfg.NodeID != cfg.NodeID {
		t.Errorf("Expected NodeID to be '%s', got '%s'", cfg.NodeID, loadedCfg.NodeID)
	}

	if loadedCfg.Address != cfg.Address {
		t.Errorf("Expected Address to be '%s', got '%s'", cfg.Address, loadedCfg.Address)
	}

	if loadedCfg.Port != cfg.Port {
		t.Errorf("Expected Port to be %d, got %d", cfg.Port, loadedCfg.Port)
	}

	if loadedCfg.LogLevel != cfg.LogLevel {
		t.Errorf("Expected LogLevel to be '%s', got '%s'", cfg.LogLevel, loadedCfg.LogLevel)
	}

	if loadedCfg.AuthEnabled != cfg.AuthEnabled {
		t.Errorf("Expected AuthEnabled to be %t, got %t", cfg.AuthEnabled, loadedCfg.AuthEnabled)
	}

	if loadedCfg.APIKey != cfg.APIKey {
		t.Errorf("Expected APIKey to be '%s', got '%s'", cfg.APIKey, loadedCfg.APIKey)
	}
}

func TestExportToFile_CreateDirectory(t *testing.T) {
	cfg := DefaultConfig()

	tempDir := t.TempDir()
	nestedDir := filepath.Join(tempDir, "nested", "directory")
	configFile := filepath.Join(nestedDir, "config.toml")

	if err := cfg.ExportToFile(configFile); err != nil {
		t.Fatalf("Failed to export config to nested directory: %v", err)
	}

	if _, err := os.Stat(nestedDir); os.IsNotExist(err) {
		t.Error("Nested directory was not created")
	}

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		t.Error("Config file was not created in nested directory")
	}
}

func TestConfigString(t *testing.T) {
	cfg := &Config{
		NodeID:      "string-test",
		Address:     "test.example.com",
		Port:        3000,
		LogLevel:    "debug",
		AuthEnabled: true,
		APIKey:      "test-key-123",
	}

	expected := "NodeID=string-test, Address=test.example.com:3000, LogLevel=debug, AuthEnabled=true, APIKey=test-key-123"
	actual := cfg.String()

	if actual != expected {
		t.Errorf("Expected string representation to be '%s', got '%s'", expected, actual)
	}
}
