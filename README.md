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

- `config.toml` — services, weather, docker settings
- `themes.toml` — theme color definitions

Edit directly or use the built-in config editor at `/configuration`.

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
