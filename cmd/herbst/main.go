package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
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

// SSEBroker manages Server-Sent Events connections
type SSEBroker struct {
	clients    map[chan string]bool
	register   chan chan string
	unregister chan chan string
	broadcast  chan string
	mu         sync.RWMutex
}

func NewSSEBroker() *SSEBroker {
	return &SSEBroker{
		clients:    make(map[chan string]bool),
		register:   make(chan chan string),
		unregister: make(chan chan string),
		broadcast:  make(chan string),
	}
}

func (b *SSEBroker) Run() {
	for {
		select {
		case client := <-b.register:
			b.mu.Lock()
			b.clients[client] = true
			b.mu.Unlock()
			log.Printf("SSE client connected (%d total)", len(b.clients))
		case client := <-b.unregister:
			b.mu.Lock()
			if _, ok := b.clients[client]; ok {
				delete(b.clients, client)
				close(client)
			}
			b.mu.Unlock()
			log.Printf("SSE client disconnected (%d total)", len(b.clients))
		case msg := <-b.broadcast:
			b.mu.RLock()
			for client := range b.clients {
				select {
				case client <- msg:
				default:
					// Client buffer full, skip
				}
			}
			b.mu.RUnlock()
		}
	}
}

func (b *SSEBroker) Notify(event string) {
	b.broadcast <- event
}

// ConfigStore holds the current config with thread-safe access
type ConfigStore struct {
	mu         sync.RWMutex
	apiConfig  APIConfig
	configPath string
	themesPath string
	broker     *SSEBroker
}

func (cs *ConfigStore) Get() APIConfig {
	cs.mu.RLock()
	defer cs.mu.RUnlock()
	return cs.apiConfig
}

func (cs *ConfigStore) Reload() error {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	// Reload config
	cfg, _, err := config.EnsureAndLoadConfig()
	if err != nil {
		return err
	}

	// Reload themes
	themeFile, _, err := themes.EnsureAndLoadThemes()
	if err != nil {
		return err
	}

	// Get active theme
	activeTheme := themeFile.ActiveTheme(cfg.Theme)

	// Update API config
	cs.apiConfig = APIConfig{
		Title:     cfg.Title,
		UI:        cfg.UI,
		Services:  cfg.Services,
		Theme:     cfg.Theme,
		ThemeVars: activeTheme.Vars,
	}

	log.Printf("Config reloaded - Theme: %s", activeTheme.Name)

	// Notify connected clients
	if cs.broker != nil {
		cs.broker.Notify("reload")
	}

	return nil
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

	// Initialize SSE broker for live reload
	broker := NewSSEBroker()
	go broker.Run()

	// Initialize config store
	store := &ConfigStore{
		apiConfig: APIConfig{
			Title:     cfg.Title,
			UI:        cfg.UI,
			Services:  cfg.Services,
			Theme:     cfg.Theme,
			ThemeVars: activeTheme.Vars,
		},
		configPath: configPath,
		themesPath: themesPath,
		broker:     broker,
	}

	// Start file watcher
	go watchFiles(store, configPath, themesPath)

	mux := http.NewServeMux()

	// API endpoint: GET /api/config
	mux.HandleFunc("/api/config", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(store.Get()); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	})

	// SSE endpoint for live reload
	mux.HandleFunc("/api/events", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// Create client channel
		client := make(chan string, 10)
		broker.register <- client

		// Remove client on disconnect
		defer func() {
			broker.unregister <- client
		}()

		// Send initial connection message
		fmt.Fprintf(w, "event: connected\ndata: ok\n\n")
		w.(http.Flusher).Flush()

		// Listen for events
		for {
			select {
			case msg := <-client:
				fmt.Fprintf(w, "event: %s\ndata: %s\n\n", msg, msg)
				w.(http.Flusher).Flush()
			case <-r.Context().Done():
				return
			}
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
	log.Println("Watching for config changes...")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func watchFiles(store *ConfigStore, configPath, themesPath string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Printf("Failed to create file watcher: %v", err)
		return
	}
	defer watcher.Close()

	// Watch the directories containing the files (more reliable than watching files directly)
	configDir := filepath.Dir(configPath)
	themesDir := filepath.Dir(themesPath)

	if err := watcher.Add(configDir); err != nil {
		log.Printf("Failed to watch config directory: %v", err)
	} else {
		log.Printf("Watching config directory: %s", configDir)
	}

	if configDir != themesDir {
		if err := watcher.Add(themesDir); err != nil {
			log.Printf("Failed to watch themes directory: %v", err)
		} else {
			log.Printf("Watching themes directory: %s", themesDir)
		}
	}

	configFile := filepath.Base(configPath)
	themesFile := filepath.Base(themesPath)

	// Debounce timer to avoid reloading on every keystroke
	var debounceTimer *time.Timer
	const debounceDelay = 500 * time.Millisecond

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			// Check if the changed file is our config or themes file
			changedFile := filepath.Base(event.Name)
			if event.Op&(fsnotify.Write|fsnotify.Create) != 0 {
				if changedFile == configFile || changedFile == themesFile {
					// Reset debounce timer
					if debounceTimer != nil {
						debounceTimer.Stop()
					}
					debounceTimer = time.AfterFunc(debounceDelay, func() {
						log.Printf("Detected change in: %s", changedFile)
						if err := store.Reload(); err != nil {
							log.Printf("Failed to reload config: %v", err)
						}
					})
				}
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Printf("File watcher error: %v", err)
		}
	}
}
