<p align="center">
  <img src="/docs/images/logoherbstfarben.svg" alt="herbst logo" width="128" />
</p>

<h1 align="center">Herbst</h1>

<p align="center"><i>(pronounced "herpst" — German for "autumn")</i></p>

<p align="center">A cozy, pastel-minimal homelab dashboard.<br>
No database. Just TOML, a tiny Go backend, and soft autumn vibes.</p>

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

For local development, config files are in `./runtime/config/`.

## Development

```bash
# Backend
go run ./cmd/herbst

# Frontend
cd web && npm install && npm run dev
```

## License

MIT
