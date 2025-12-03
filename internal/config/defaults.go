package config

// DefaultConfigTOML is the default configuration file content
const DefaultConfigTOML = `title = "herbst – homelab"
theme = "default"

[ui]
font = ""

[ui.background]
image = ""
blur = 0

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
