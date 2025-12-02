package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"

	"herbst/internal/config"
	"herbst/internal/themes"
	"herbst/internal/util"
)

const (
	envAssetsDir       = "HERBST_ASSETS_DIR"
	devAssetsDir       = "./runtime/assets"
	containerAssetsDir = "/app/assets"
)

// APIConfig is the response structure for /api/config
type APIConfig struct {
	Title     string            `json:"title"`
	UI        config.UI         `json:"ui"`
	Services  []config.Service  `json:"services"`
	Theme     string            `json:"theme"`
	ThemeVars map[string]string `json:"themeVars"`
}

func main() {
	// Load .env file if it exists (won't override existing env vars)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	} else {
		log.Println("Loaded .env file")
	}

	// Load configuration
	cfg, configPath, err := config.EnsureAndLoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	log.Printf("Config loaded from: %s", configPath)

	// Load themes
	themeFile, themesPath, err := themes.EnsureAndLoadThemes()
	if err != nil {
		log.Fatalf("Failed to load themes: %v", err)
	}
	log.Printf("Themes loaded from: %s", themesPath)

	// Get active theme
	activeTheme := themeFile.ActiveTheme(cfg.Theme)
	log.Printf("Active theme: %s", activeTheme.Name)

	// Build API config response
	apiConfig := APIConfig{
		Title:     cfg.Title,
		UI:        cfg.UI,
		Services:  cfg.Services,
		Theme:     cfg.Theme,
		ThemeVars: activeTheme.Vars,
	}

	mux := http.NewServeMux()

	// API endpoint: GET /api/config
	mux.HandleFunc("/api/config", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(apiConfig); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	})

	// Serve user assets (if directory exists) under /static/
	assetsDir := util.ResolveDir(envAssetsDir, devAssetsDir, containerAssetsDir)
	if _, err := os.Stat(assetsDir); err == nil {
		log.Printf("Serving user assets from: %s at /static/", assetsDir)
		mux.Handle("/static/", http.StripPrefix("/static/",
			http.FileServer(http.Dir(assetsDir)),
		))
	} else {
		log.Printf("User assets directory not found, skipping: %s", assetsDir)
	}

	// Serve frontend (Vue + Vite build)
	distPath := filepath.Join("web", "dist")
	if _, err := os.Stat(distPath); err != nil {
		log.Printf("Warning: dist folder missing: %s - frontend won't be served", distPath)
	} else {
		log.Printf("Serving frontend from: %s", distPath)
		mux.Handle("/", http.FileServer(http.Dir(distPath)))
	}

	log.Println("herbst running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
