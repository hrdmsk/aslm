package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Config represents application configuration
type Config struct {
	HomePath     string `json:"homePath"`
	GeminiAPIKey string `json:"geminiApiKey"`
}

var configPath string

func init() {
	// Get user config directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// Fallback to current directory
		configPath = "config.json"
		return
	}

	configDir := filepath.Join(homeDir, ".aslm")
	os.MkdirAll(configDir, 0755)
	configPath = filepath.Join(configDir, "config.json")
}

// LoadConfig loads configuration from file
func LoadConfig() (*Config, error) {
	// Default config
	defaultConfig := &Config{
		HomePath:     "D:/VRChatAssetPack",
		GeminiAPIKey: "",
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		// If file doesn't exist, return default config
		if os.IsNotExist(err) {
			return defaultConfig, nil
		}
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// SaveConfig saves configuration to file
func SaveConfig(cfg *Config) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}
