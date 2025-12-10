<p align="center">
  <img src="/docs/images/logoherbstfarben.svg" alt="herbst logo" width="128" />
</p>

<h1 align="center">Herbst</h1>

<p align="center"><i>(pronounced "herpst" — German for "autumn")</i></p>

<p align="center">A cozy, pastel-minimal homelab dashboard.<br>
No database. Just TOML, a tiny Go backend, and soft autumn vibes.</p>

<p align="center">
  <img src="/docs/images/image.png" alt="herbst dashboard screenshot" width="800" />
</p>

> [!WARNING]
> herbst is in active development — expect bugs, missing features, and occasional UI teleportation :D
> See [CHANGELOG.md](CHANGELOG.md) for recent changes.

---

## Features

- Service cards with optional health checks
- Local Docker container monitoring (mount the socket)
- Remote Docker agents for multi-node setups
- **System monitoring** — CPU, RAM, disk usage, and uptime
- Weather widget (OpenWeatherMap)
- Multiple themes (earthy, bright, autumn, glass)
- In-browser config editor with live reload
- Background image support with blur
- **Configurable clock format** — 12h/24h time, short/numeric date
- **Version display** in footer (auto-set from Git tag on Docker builds)
- No database, just TOML files

## Quick Start

```yaml
services:
  herbst:
    image: ghcr.io/brendlij/herbst:latest
    container_name: herbst
    restart: unless-stopped
    # Run as your user (UID:GID) so config files are owned by you, not root
    # Find your IDs with: id -u && id -g
    # user: "1000:1000"
    ports:
      - "8080:8080" # External 8088 → Internal 8080
    volumes:
      # Config directory (contains config.toml and themes.toml)
      - ./config:/app/config

      # Static files (for custom icons, backgrounds, etc.)
      - ./static:/app/static

      # Docker socket (required for Docker tab)
      - /var/run/docker.sock:/var/run/docker.sock
    environment:
      - TZ=Europe/Berlin
      # Optional: Add environment variables for config file substitution
      # - OPENWEATHER_API_KEY=your-api-key
      # - HERBST_HOST=192.168.1.100:8080   # For remote docker agents
      # - HERBST_AGENT_PROTOCOL=wss        # Use wss for SSL/TLS (default: ws)
```

## Configuration

Config files live in `/app/config` inside the container (mount a volume to persist):

- `config.toml` — services, weather, docker, system settings
- `themes.toml` — theme color definitions

Edit directly or use the built-in config editor at `/configuration`.

> **Tip:** Use `${ENV_VAR_NAME}` syntax in config values to reference environment variables.

### General Settings

```toml
title = "herbst – homelab"
theme = "autumn"  # Available: autumn, earthy, bright, glass
```

### UI Settings

```toml
[ui]
font = ""  # Custom font name (e.g., "Inter", "Fira Code")

[ui.background]
image = ""  # Filename from /static (e.g., "bg.jpg") or full URL
blur = 0    # Blur amount in pixels

[ui.clock]
time-format = "24h"    # "24h" or "12h"
date-format = "short"  # "short" (3. Dez 2025) or "numeric" (03.12.2025)
```

### Weather (OpenWeatherMap)

```toml
[weather]
enabled = false
api-key = "${OPENWEATHER_API_KEY}"
location = ""       # City "London,GB", zip "10115,DE", or leave empty for lat/lon
lat = 0.0
lon = 0.0
units = "metric"    # metric (°C), imperial (°F), standard (K)
```

### Docker - Local

```toml
[docker.local]
socket-path = "/var/run/docker.sock"
# enabled = true  # Auto-detects if socket exists
```

### Docker - Remote Agents

For monitoring Docker on remote machines, add agents to your config:

```toml
[[docker.agent]]
name = "server-name"
```

The token is auto-generated. Go to the **Configuration** page in the UI to find the ready-to-use `docker run` command with the correct token.

For agents to connect, set these environment variables on the herbst container:

```toml
host = "${HERBST_HOST}"              # e.g., "192.168.1.100:8080"
agent-protocol = "${HERBST_AGENT_PROTOCOL}"  # "ws" (default) or "wss" for SSL
```

### System Monitoring

```toml
[system]
enabled = true
disk-path = "/"  # Path to monitor disk usage (e.g., "/" or "/mnt/data")
```

### Services

Group services into sections:

```toml
[[section]]
title = "Home"

[[section.service]]
name = "Home Assistant"
url = "https://ha.local"
icon = ""              # Icon name or URL (optional)
online-badge = true    # Show online/offline status (optional)

[[section.service]]
name = "NAS"
url = "https://nas.local"
online-badge = true
```

---

## Development

### Prerequisites

- **Go** 1.21+ ([install](https://go.dev/doc/install))
- **Node.js** 20+ ([install](https://nodejs.org/))
- **Bun** (optional, faster than npm) ([install](https://bun.sh/))

### Setup

1. **Clone the repository**

   ```bash
   git clone https://github.com/brendlij/herbst.git
   cd herbst
   ```

2. **Create a `.env` file** in the project root (required for local dev):

   ```env
   HERBST_CONFIG_DIR=./config
   HERBST_STATIC_DIR=./static
   ```

   > Without this, the app will try to use `/app/config` (Docker path) and fail on macOS/Linux.

3. **Create config directories**

   ```bash
   mkdir -p config static
   ```

4. **Install Go dependencies**

   ```bash
   go mod download
   ```

5. **Install frontend dependencies**
   ```bash
   cd web
   bun install  # or: npm install
   cd ..
   ```

### Running locally

You need two terminals — one for the backend, one for the frontend:

**Terminal 1 — Backend (Go)**

```bash
go run ./cmd/herbst
```

The API server runs on `http://localhost:8080`

**Terminal 2 — Frontend (Vite)**

```bash
cd web
bun run dev  # or: npm run dev
```

The frontend dev server runs on `http://localhost:5173` with hot reload.

> **Note:** In development, the frontend proxies API requests to the Go backend (configured in `vite.config.ts`).

### Building

**Build Go binary:**

```bash
go build -o herbst ./cmd/herbst
```

**Build frontend for production:**

```bash
cd web
bun run build  # or: npm run build
```

Output goes to `web/dist/`.

**Build with version (like Docker does):**

```bash
go build -ldflags="-X main.Version=v0.2.7" -o herbst ./cmd/herbst
```

### Project structure

```
herbst/
├── cmd/
│   ├── herbst/              # Main dashboard server
│   └── herbst-docker-agent/ # Remote Docker agent
├── internal/
│   ├── config/              # Config loading & types
│   ├── agents/              # WebSocket agent handling
│   ├── themes/              # Theme loading
│   └── util/                # Utilities
├── web/                     # Vue 3 + Vite frontend
│   └── src/
│       ├── components/      # Vue components
│       ├── views/           # Page views
│       ├── types/           # TypeScript types
│       └── lib/             # Utilities
├── config/                  # Local dev config (gitignored)
├── static/                  # Static files (backgrounds, icons)
└── .env                     # Local environment (gitignored)
```

### Tips

- **Config changes**: The app watches `config.toml` and `themes.toml` for changes and auto-reloads
- **API testing**: All endpoints are under `/api/` (e.g., `/api/config`, `/api/version`, `/api/system/stats`)
- **Themes**: Edit `themes.toml` to customize or add new themes

---

## License

MIT
