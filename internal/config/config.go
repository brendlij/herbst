package config

import (
	"os"
	"path/filepath"

	"herbst/internal/util"

	"github.com/pelletier/go-toml/v2"
)

const (
	envConfigDir       = "HERBST_CONFIG_DIR"
	devConfigDir       = "./runtime/config"
	containerConfigDir = "/app/config"
	configFilename     = "config.toml"
)

// UI holds UI-related configuration
type UI struct {
	Background string `toml:"background" json:"background"`
}

// Service represents a dashboard service entry
type Service struct {
	Name string `toml:"name" json:"name"`
	URL  string `toml:"url"  json:"url"`
	Icon string `toml:"icon" json:"icon"`
}

// Config is the main configuration structure
type Config struct {
	Title    string    `toml:"title"    json:"title"`
	Theme    string    `toml:"theme"    json:"theme"`
	UI       UI        `toml:"ui"       json:"ui"`
	Services []Service `toml:"services" json:"services"`
}

// EnsureAndLoadConfig loads the config file, creating it with defaults if it doesn't exist.
// Returns the config, the absolute path to the config file, and any error.
func EnsureAndLoadConfig() (*Config, string, error) {
	// Determine config directory
	dir := util.ResolveDir(envConfigDir, devConfigDir, containerConfigDir)

	// Ensure directory exists
	if err := util.EnsureDir(dir); err != nil {
		return nil, "", err
	}

	configPath := filepath.Join(dir, configFilename)

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Write default config
		if err := os.WriteFile(configPath, []byte(DefaultConfigTOML), 0644); err != nil {
			return nil, "", err
		}
	}

	// Read and parse config
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, "", err
	}

	var cfg Config
	if err := toml.Unmarshal(data, &cfg); err != nil {
		return nil, "", err
	}

	absPath, _ := filepath.Abs(configPath)
	return &cfg, absPath, nil
}
