package config

// DefaultConfigTOML is the default configuration file content
// Note: You can use ${ENV_VAR_NAME} syntax to reference environment variables
const DefaultConfigTOML = `# ╔═══════════════════════════════════════════════════════════════════════════╗
# ║  HERBST CONFIGURATION                                                     ║
# ║  Tip: Use ${ENV_VAR_NAME} to reference environment variables              ║
# ╚═══════════════════════════════════════════════════════════════════════════╝

title = "herbst – homelab"
theme = "Autumn"  # Available: Autumn, Aarthy, Bright, Glass, Meadow, Noir, Arctic, Blossom, Ember, Nebula


# ┌───────────────────────────────────────────────────────────────────────────┐
# │  UI SETTINGS                                                              │
# └───────────────────────────────────────────────────────────────────────────┘

[ui]
font = ""

[ui.background]
image = ""  # Filename from /static (e.g. "bg.jpg") or full URL
blur = 0


# ┌───────────────────────────────────────────────────────────────────────────┐
# │  WEATHER (OpenWeatherMap)                                                 │
# └───────────────────────────────────────────────────────────────────────────┘

[weather]
enabled = false
api-key = "${OPENWEATHER_API_KEY}"
location = ""    # City "London,GB", zip "10115,DE", or leave empty for lat/lon
lat = 0.0
lon = 0.0
units = "metric"  # metric (°C), imperial (°F), standard (K)


# ┌───────────────────────────────────────────────────────────────────────────┐
# │  DOCKER - LOCAL                                                           │
# │  Shows containers from the machine where herbst runs                      │
# └───────────────────────────────────────────────────────────────────────────┘

[docker.local]
socket-path = "/var/run/docker.sock"
# enabled = true  # Auto-detects if socket exists


# ┌───────────────────────────────────────────────────────────────────────────┐
# │  DOCKER - REMOTE AGENTS                                                   │
# │  Connect to other machines running herbst-docker-agent                    │
# └───────────────────────────────────────────────────────────────────────────┘

# Set these in your docker-compose environment:
#   HERBST_HOST=192.168.1.100:8080
#   HERBST_AGENT_PROTOCOL=wss  (optional, default: ws)

host = "${HERBST_HOST}"
agent-protocol = "${HERBST_AGENT_PROTOCOL}"

# Add agents by name only - tokens are auto-generated and shown in the UI
# [[docker.agent]]
# name = "server1"
#
# [[docker.agent]]
# name = "raspberry-pi"


# ┌───────────────────────────────────────────────────────────────────────────┐
# │  SYSTEM MONITORING                                                        │
# │  Shows CPU, RAM, disk usage, and uptime for the herbst host               │
# └───────────────────────────────────────────────────────────────────────────┘

[system]
enabled = true
disk-path = "/"  # Path to monitor disk usage (e.g., "/" or "/mnt/data")


# ┌───────────────────────────────────────────────────────────────────────────┐
# │  SERVICES                                                                 │
# │  Group services into sections with [[section]]                            │
# └───────────────────────────────────────────────────────────────────────────┘

[[section]]
title = "Home"

[[section.service]]
name = "Home Assistant"
url = "https://ha.local"
icon = ""
online-badge = true

[[section.service]]
name = "NAS"
url = "https://nas.local"
icon = ""
online-badge = true

# [[section]]
# title = "Media"
#
# [[section.service]]
# name = "Plex"
# url = "https://plex.local"
# online-badge = true
`
