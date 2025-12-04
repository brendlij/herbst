package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/joho/godotenv"

	"herbst/internal/agents"
	"herbst/internal/config"
	"herbst/internal/themes"
	"herbst/internal/util"
)

const (
	envStaticDir       = "HERBST_STATIC_DIR"
	devStaticDir       = "./runtime/static"
	containerStaticDir = "/app/static"
)

// DockerAPIConfig is the resolved Docker config for API responses
type DockerAPIConfig struct {
	Enabled    bool   `json:"enabled"`
	SocketPath string `json:"socketPath"`
}

// APIConfig is the response structure for /api/config
type APIConfig struct {
	Title     string            `json:"title"`
	UI        config.UI         `json:"ui"`
	Weather   config.Weather    `json:"weather"`
	Docker    DockerAPIConfig   `json:"docker"`
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
		Title:   cfg.Title,
		UI:      cfg.UI,
		Weather: cfg.Weather,
		Docker: DockerAPIConfig{
			Enabled:    cfg.Docker.IsEnabled(),
			SocketPath: cfg.Docker.SocketPath,
		},
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
			Title:   cfg.Title,
			UI:      cfg.UI,
			Weather: cfg.Weather,
			Docker: DockerAPIConfig{
				Enabled:    cfg.Docker.IsEnabled(),
				SocketPath: cfg.Docker.SocketPath,
			},
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

	registry := agents.NewRegistry()
	agentServer := agents.NewServer(cfg, registry)

	mux := http.NewServeMux()

		// WebSocket für Agents: /api/agents/ws
	mux.HandleFunc("/api/agents/ws", agentServer.HandleWS)

	// Remote-Docker-Nodes: /api/docker/nodes
	mux.HandleFunc("/api/docker/nodes", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		snapshot := registry.Snapshot()

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(snapshot); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	})

	// API endpoint: POST /api/reload
	// Reloads the configuration files and notifies all connected clients
	mux.HandleFunc("/api/reload", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		log.Println("Config reload requested via API")

		if err := store.Reload(); err != nil {
			log.Printf("Failed to reload config: %v", err)
			http.Error(w, "Failed to reload configuration", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "Configuration reloaded successfully",
		})
	})

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

	// API endpoint: GET /api/health?url=<service-url>
	// Checks if a service URL is reachable
	mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		targetURL := r.URL.Query().Get("url")
		if targetURL == "" {
			http.Error(w, "Missing 'url' parameter", http.StatusBadRequest)
			return
		}

		// Create a client with timeout and skip TLS verification (for self-signed certs)
		client := &http.Client{
			Timeout: 5 * time.Second,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		}

		// Try HEAD first, fall back to GET if it fails
		req, err := http.NewRequest(http.MethodHead, targetURL, nil)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{"online": false, "error": err.Error()})
			return
		}

		resp, err := client.Do(req)
		
		// If HEAD fails or returns 405 (Method Not Allowed), try GET
		if err != nil || (resp != nil && resp.StatusCode == 405) {
			if resp != nil {
				resp.Body.Close()
			}
			req, _ = http.NewRequest(http.MethodGet, targetURL, nil)
			resp, err = client.Do(req)
		}

		online := err == nil && resp != nil && resp.StatusCode < 500
		if resp != nil {
			resp.Body.Close()
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"online": online})
	})

	// API endpoint: GET /api/weather
	// Fetches current weather from OpenWeatherMap
	mux.HandleFunc("/api/weather", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		weatherCfg := store.Get().Weather
		if !weatherCfg.Enabled || weatherCfg.APIKey == "" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"enabled": false,
				"error":   "Weather not configured",
			})
			return
		}

		client := &http.Client{Timeout: 10 * time.Second}
		var lat, lon float64
		var locationName string

		// Determine coordinates based on location config
		if weatherCfg.Location != "" {
			location := strings.TrimSpace(weatherCfg.Location)
			
			// Auto-detect if it's a zip code: starts with digits or has "zip:" prefix
			isZipCode := strings.HasPrefix(strings.ToLower(location), "zip:")
			if !isZipCode && len(location) > 0 {
				// Check if it starts with a digit (likely a zip code like "79650,DE" or "10001,US")
				firstChar := location[0]
				isZipCode = firstChar >= '0' && firstChar <= '9'
			}
			
			if isZipCode {
				// Remove "zip:" prefix if present
				zipPart := strings.TrimPrefix(location, "zip:")
				zipPart = strings.TrimPrefix(zipPart, "ZIP:")
				zipPart = strings.TrimSpace(zipPart)
				
				geoURL := fmt.Sprintf(
					"http://api.openweathermap.org/geo/1.0/zip?zip=%s&appid=%s",
					zipPart,
					weatherCfg.APIKey,
				)

				resp, err := client.Get(geoURL)
				if err != nil {
					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(map[string]interface{}{
						"enabled": true,
						"error":   "Failed to geocode zip code",
					})
					return
				}
				defer resp.Body.Close()

				var zipResp struct {
					Lat  float64 `json:"lat"`
					Lon  float64 `json:"lon"`
					Name string  `json:"name"`
				}
				if err := json.NewDecoder(resp.Body).Decode(&zipResp); err != nil {
					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(map[string]interface{}{
						"enabled": true,
						"error":   "Failed to parse zip code response",
					})
					return
				}
				lat, lon = zipResp.Lat, zipResp.Lon
				locationName = zipResp.Name
			} else {
				// It's a city name
				geoURL := fmt.Sprintf(
					"http://api.openweathermap.org/geo/1.0/direct?q=%s&limit=1&appid=%s",
					weatherCfg.Location,
					weatherCfg.APIKey,
				)

				resp, err := client.Get(geoURL)
				if err != nil {
					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(map[string]interface{}{
						"enabled": true,
						"error":   "Failed to geocode city name",
					})
					return
				}
				defer resp.Body.Close()

				var geoResp []struct {
					Lat     float64 `json:"lat"`
					Lon     float64 `json:"lon"`
					Name    string  `json:"name"`
					Country string  `json:"country"`
				}
				if err := json.NewDecoder(resp.Body).Decode(&geoResp); err != nil || len(geoResp) == 0 {
					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(map[string]interface{}{
						"enabled": true,
						"error":   "City not found",
					})
					return
				}
				lat, lon = geoResp[0].Lat, geoResp[0].Lon
				locationName = geoResp[0].Name
			}
		} else {
			// Use direct coordinates from config
			lat, lon = weatherCfg.Lat, weatherCfg.Lon
		}

		// Build OpenWeatherMap API URL
		apiURL := fmt.Sprintf(
			"https://api.openweathermap.org/data/2.5/weather?lat=%f&lon=%f&appid=%s&units=%s",
			lat,
			lon,
			weatherCfg.APIKey,
			weatherCfg.Units,
		)

		resp, err := client.Get(apiURL)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"enabled": true,
				"error":   "Failed to fetch weather data",
			})
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"enabled": true,
				"error":   fmt.Sprintf("Weather API returned status %d", resp.StatusCode),
			})
			return
		}

		// Parse OpenWeatherMap response
		var owmResp struct {
			Main struct {
				Temp      float64 `json:"temp"`
				FeelsLike float64 `json:"feels_like"`
				Humidity  int     `json:"humidity"`
			} `json:"main"`
			Weather []struct {
				Main        string `json:"main"`
				Description string `json:"description"`
				Icon        string `json:"icon"`
			} `json:"weather"`
			Name string `json:"name"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&owmResp); err != nil {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"enabled": true,
				"error":   "Failed to parse weather data",
			})
			return
		}

		// Use geocoded location name, fall back to API response
		cityName := locationName
		if cityName == "" {
			cityName = owmResp.Name
		}

		description := ""
		icon := ""
		if len(owmResp.Weather) > 0 {
			description = owmResp.Weather[0].Description
			icon = owmResp.Weather[0].Icon
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"enabled":     true,
			"temp":        owmResp.Main.Temp,
			"feelsLike":   owmResp.Main.FeelsLike,
			"humidity":    owmResp.Main.Humidity,
			"description": description,
			"icon":        icon,
			"city":        cityName,
			"units":       weatherCfg.Units,
		})
	})

	// API endpoint: GET /api/docker/containers
	// Lists all Docker containers via Docker socket
	mux.HandleFunc("/api/docker/containers", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		dockerCfg := store.Get().Docker
		if !dockerCfg.Enabled {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"enabled": false,
				"error":   "Docker integration not enabled",
			})
			return
		}

		socketPath := dockerCfg.SocketPath
		if socketPath == "" {
			socketPath = "/var/run/docker.sock"
		}

		// Create HTTP client that connects via Unix socket
		client := &http.Client{
			Transport: &http.Transport{
				DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
					return net.Dial("unix", socketPath)
				},
			},
			Timeout: 10 * time.Second,
		}

		// Call Docker API
		req, err := http.NewRequest("GET", "http://localhost/containers/json?all=true", nil)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"enabled": true,
				"error":   "Failed to create request",
			})
			return
		}

		resp, err := client.Do(req)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"enabled": true,
				"error":   fmt.Sprintf("Failed to connect to Docker: %v", err),
			})
			return
		}
		defer resp.Body.Close()

		// Parse Docker API response
		var containers []struct {
			ID      string   `json:"Id"`
			Names   []string `json:"Names"`
			Image   string   `json:"Image"`
			State   string   `json:"State"`
			Status  string   `json:"Status"`
			Created int64    `json:"Created"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&containers); err != nil {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"enabled": true,
				"error":   "Failed to parse Docker response",
			})
			return
		}

		// Transform to our format
		result := make([]map[string]interface{}, len(containers))
		for i, c := range containers {
			name := c.ID[:12]
			if len(c.Names) > 0 {
				name = strings.TrimPrefix(c.Names[0], "/")
			}
			result[i] = map[string]interface{}{
				"id":      c.ID[:12],
				"name":    name,
				"image":   c.Image,
				"state":   c.State,
				"status":  c.Status,
				"created": c.Created,
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"enabled":    true,
			"containers": result,
		})
	})

	// API endpoint: GET /api/docker/agents
	// Lists all configured docker agents with their connection status
	mux.HandleFunc("/api/docker/agents", func(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// UI-Config aus dem Store (ist immer "valid", weil nur bei erfolgreichem Reload ersetzt wird)
	currentCfg := store.Get()

	// Verbundene Nodes aus Registry
	connectedNodes := registry.Snapshot()

	// Original-Config für Tokens / Agents laden
	cfg, _, err := config.EnsureAndLoadConfig()
	if err != nil {
		log.Printf("Failed to load config in /api/docker/agents: %v", err)

		// Kein Panic, sondern sauberer Fehler-Response
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"enabled":    currentCfg.Docker.Enabled,
			"agents":     []interface{}{},
			"serverHost": r.Host,
			"error":      "failed to load config",
		})
		return
	}

	type AgentResponse struct {
		Name       string      `json:"name"`
		Token      string      `json:"token"`
		Connected  bool        `json:"connected"`
		LastSeen   *string     `json:"lastSeen"`
		Containers interface{} `json:"containers"`
	}

	agents := make([]AgentResponse, 0, len(cfg.Docker.Agents))

	for _, agentCfg := range cfg.Docker.Agents {
		agent := AgentResponse{
			Name:       agentCfg.Name,
			Token:      agentCfg.Token,
			Connected:  false,
			LastSeen:   nil,
			Containers: []interface{}{},
		}

		if node, exists := connectedNodes[agentCfg.Name]; exists {
			agent.Connected = true
			lastSeen := node.LastSeen.Format(time.RFC3339)
			agent.LastSeen = &lastSeen
			agent.Containers = node.Containers
		}

		agents = append(agents, agent)
	}

	// Host für den Docker-Run-Command
	hostURL := cfg.Docker.Host
	if hostURL == "" {
		hostURL = r.Host
	}

	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"enabled":    currentCfg.Docker.Enabled,
		"agents":     agents,
		"serverHost": hostURL,
	})
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

	// Serve static files (if directory exists) under /static/
	staticDir := util.ResolveDir(envStaticDir, devStaticDir, containerStaticDir)
	if _, err := os.Stat(staticDir); err == nil {
		log.Printf("Serving static files from: %s at /static/", staticDir)
		mux.Handle("/static/", http.StripPrefix("/static/",
			http.FileServer(http.Dir(staticDir)),
		))
	} else {
		log.Printf("Static directory not found, skipping: %s", staticDir)
	}

	// Serve frontend (Vue + Vite build)
	// Serve frontend (Vue + Vite build)
	distPath := filepath.Join("web", "dist")
	if _, err := os.Stat(distPath); err != nil {
		log.Printf("Warning: dist folder missing: %s - frontend won't be served", distPath)
	} else {
		log.Printf("Serving frontend from: %s", distPath)

		fs := http.FileServer(http.Dir(distPath))

		// SPA-Handler mit Fallback auf index.html
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			// API & static sollen NICHT vom SPA gehandlet werden
			if strings.HasPrefix(r.URL.Path, "/api/") || strings.HasPrefix(r.URL.Path, "/static/") {
				http.NotFound(w, r)
				return
			}

			// Pfad im dist-Ordner prüfen
			requestPath := filepath.Join(distPath, filepath.Clean(r.URL.Path))
			if info, err := os.Stat(requestPath); err == nil && !info.IsDir() {
				// Datei existiert -> normal ausliefern
				fs.ServeHTTP(w, r)
				return
			}

			// Fallback: index.html für SPA-Routen wie /docker, /docker-nodes, ...
			http.ServeFile(w, r, filepath.Join(distPath, "index.html"))
		})
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
