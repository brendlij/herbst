package themes

import (
	"os"
	"path/filepath"
	"strings"

	"herbst/internal/util"

	"github.com/pelletier/go-toml/v2"
)

const (
	envConfigDir       = "HERBST_CONFIG_DIR"
	devConfigDir       = "./runtime/config"
	containerConfigDir = "/app/config"
	themesFilename     = "themes.toml"
)

// Theme represents a single theme configuration
type Theme struct {
	Name string            `toml:"name" json:"name"`
	Vars map[string]string `toml:"vars" json:"vars"`
}

// ThemeFile holds all themes
type ThemeFile struct {
	Themes map[string]Theme `toml:"theme" json:"theme"`
}

// EnsureAndLoadThemes loads the themes file, creating it with defaults if it doesn't exist.
// Returns the theme file, the absolute path to the themes file, and any error.
func EnsureAndLoadThemes() (*ThemeFile, string, error) {
	// Determine config directory (themes.toml lives alongside config.toml)
	dir := util.ResolveDir(envConfigDir, devConfigDir, containerConfigDir)

	// Ensure directory exists
	if err := util.EnsureDir(dir); err != nil {
		return nil, "", err
	}

	themesPath := filepath.Join(dir, themesFilename)

	// Check if themes file exists
	if _, err := os.Stat(themesPath); os.IsNotExist(err) {
		// Write default themes
		if err := os.WriteFile(themesPath, []byte(DefaultThemesTOML), 0644); err != nil {
			return nil, "", err
		}
	}

	// Read and parse themes
	data, err := os.ReadFile(themesPath)
	if err != nil {
		return nil, "", err
	}

	var tf ThemeFile
	if err := toml.Unmarshal(data, &tf); err != nil {
		return nil, "", err
	}

	absPath, _ := filepath.Abs(themesPath)
	return &tf, absPath, nil
}

// ActiveTheme returns the theme with the given key or display name, or "default" if not found
func (tf *ThemeFile) ActiveTheme(name string) Theme {
	if name == "" {
		name = "default"
	}

	// First, try to find by key (e.g., "autumn_mist")
	if theme, ok := tf.Themes[name]; ok {
		return theme
	}

	// Second, try to find by display name (e.g., "Autumn Mist")
	for _, theme := range tf.Themes {
		if strings.EqualFold(theme.Name, name) {
			return theme
		}
	}

	// Fall back to default theme
	if theme, ok := tf.Themes["default"]; ok {
		return theme
	}

	// Return empty theme if nothing found
	return Theme{
		Name: "default",
		Vars: make(map[string]string),
	}
}
