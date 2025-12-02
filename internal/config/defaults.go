package config

// DefaultConfigTOML is the default configuration file content
const DefaultConfigTOML = `title = "herbst – homelab"
theme = "default"

[ui]
background = ""

[[services]]
name = "Home Assistant"
url  = "https://ha.local"
icon = ""

[[services]]
name = "NAS"
url  = "https://nas.local"
icon = ""
`
