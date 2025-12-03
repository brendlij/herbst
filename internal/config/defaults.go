package config

// DefaultConfigTOML is the default configuration file content
// Note: You can use ${ENV_VAR_NAME} syntax to reference environment variables
const DefaultConfigTOML = `title = "herbst – homelab"
theme = "default"

[ui]
font = ""

[ui.background]
image = ""
blur = 0

# Weather configuration (OpenWeatherMap)
# Tip: Use ${ENV_VAR_NAME} to load values from environment variables
# Note: Config changes hot-reload, but .env changes require a restart
[weather]
enabled = false
api-key = ""     # Direct key or ${ENV_VAR_NAME}
location = ""    # City (e.g. "London,GB"), zip code (e.g. "79650,DE"), or empty for lat/lon
lat = 0.0        # Latitude (only used if location is empty)
lon = 0.0        # Longitude (only used if location is empty)
units = "metric" # metric (°C), imperial (°F), or standard (K)

# Docker integration - shows container status
# Mount Docker socket: -v /var/run/docker.sock:/var/run/docker.sock
[docker]
enabled = false
socket-path = "/var/run/docker.sock"  # Default Docker socket path

[[services]]
name = "Home Assistant"
url  = "https://ha.local"
icon = ""
online-badge = true

[[services]]
name = "NAS"
url  = "https://nas.local"
icon = ""
online-badge = true
`
