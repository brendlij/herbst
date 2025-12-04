package config

import (
	"os"
	"path/filepath"
	"regexp"

	"herbst/internal/util"

	"github.com/pelletier/go-toml/v2"
)

const (
	envConfigDir       = "HERBST_CONFIG_DIR"
	devConfigDir       = "./runtime/config"
	containerConfigDir = "/app/config"
	configFilename     = "config.toml"
)

// Background holds background-related configuration
type Background struct {
	Image string  `toml:"image" json:"image"`
	Blur  float64 `toml:"blur"  json:"blur"`
}

// Weather holds weather-related configuration
type Weather struct {
	Enabled  bool    `toml:"enabled"  json:"enabled"`
	APIKey   string  `toml:"api-key"  json:"apiKey"`
	Location string  `toml:"location" json:"location"` // City name, "zip:CODE,COUNTRY", or empty for lat/lon
	Lat      float64 `toml:"lat"      json:"lat"`
	Lon      float64 `toml:"lon"      json:"lon"`
	Units    string  `toml:"units"    json:"units"`
}

// Docker holds Docker integration configuration
type Docker struct {
	Enabled    *bool               `toml:"enabled"     json:"enabled"`    // Pointer to detect if explicitly set
	Host       string              `toml:"host"        json:"host"`       // External host URL for agents (e.g. "192.168.1.100:8080")
	SocketPath string              `toml:"socket-path" json:"socketPath"`
	Agents     []DockerAgentConfig `toml:"agents"      json:"agents"`     // [[docker.agents]]
}

// DockerAgentConfig represents a remote docker agent node
type DockerAgentConfig struct {
	Name  string `toml:"name"  json:"name"`
	Token string `toml:"token" json:"token"`
}



// IsEnabled returns true if Docker is enabled (auto-detects if not explicitly set)
func (d *Docker) IsEnabled() bool {
	// If explicitly set in config, use that value
	if d.Enabled != nil {
		return *d.Enabled
	}
	// Auto-detect: check if socket exists
	if d.SocketPath == "" {
		d.SocketPath = "/var/run/docker.sock"
	}
	_, err := os.Stat(d.SocketPath)
	return err == nil
}

// UI holds UI-related configuration
type UI struct {
	Background Background `toml:"background" json:"background"`
	Font       string     `toml:"font"       json:"font"`
}

// Service represents a dashboard service entry
type Service struct {
	Name        string `toml:"name"         json:"name"`
	URL         string `toml:"url"          json:"url"`
	Icon        string `toml:"icon"         json:"icon"`
	OnlineBadge bool   `toml:"online-badge" json:"onlineBadge"`
}

// Config is the main configuration structure
type Config struct {
	Title    string    `toml:"title"    json:"title"`
	Theme    string    `toml:"theme"    json:"theme"`
	UI       UI        `toml:"ui"       json:"ui"`
	Weather  Weather   `toml:"weather"  json:"weather"`
	Docker   Docker    `toml:"docker"   json:"docker"`
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

	// Expand environment variables in config values
	expandEnvVars(&cfg)

	absPath, _ := filepath.Abs(configPath)
	return &cfg, absPath, nil
}

// expandEnvVars expands ${VAR_NAME} references in config string values
func expandEnvVars(cfg *Config) {
	envVarRegex := regexp.MustCompile(`\$\{([^}]+)\}`)

	expand := func(s string) string {
		return envVarRegex.ReplaceAllStringFunc(s, func(match string) string {
			// Extract variable name from ${VAR_NAME}
			varName := envVarRegex.FindStringSubmatch(match)[1]
			if value, exists := os.LookupEnv(varName); exists {
				return value
			}
			// Return original if env var not found
			return match
		})
	}

	// Expand in Weather config
	cfg.Weather.APIKey = expand(cfg.Weather.APIKey)
	cfg.Weather.Location = expand(cfg.Weather.Location)

	// Expand in Docker config
	cfg.Docker.SocketPath = expand(cfg.Docker.SocketPath)

	// Expand in Docker agent configs
	for i := range cfg.Docker.Agents {
		cfg.Docker.Agents[i].Name  = expand(cfg.Docker.Agents[i].Name)
		cfg.Docker.Agents[i].Token = expand(cfg.Docker.Agents[i].Token)
	}

	// Expand in UI config
	cfg.UI.Background.Image = expand(cfg.UI.Background.Image)
	cfg.UI.Font = expand(cfg.UI.Font)

	// Expand in Services
	for i := range cfg.Services {
		cfg.Services[i].Name = expand(cfg.Services[i].Name)
		cfg.Services[i].URL = expand(cfg.Services[i].URL)
		cfg.Services[i].Icon = expand(cfg.Services[i].Icon)
	}

	// Expand title and theme
	cfg.Title = expand(cfg.Title)
	cfg.Theme = expand(cfg.Theme)
}
