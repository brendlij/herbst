# Changelog

## [0.2.7] - 2025-12-10

### Added

- **Version display**: Shows app version in footer (automatically set from Git tag during Docker build)
- **Clock format config**: Configure time (12h/24h) and date format (short/numeric) via `[ui.clock]` in config.toml

### Changed

- **Settings icon**: Replaced "Configuration" tab text with minimal cog icon, moved to header right side
- **Improved padding**: Better spacing for all views with consistent left/right padding
- **Section spacing**: Increased gap between service sections for clearer visual separation

### Fixed

- **System tab loading speed**: CPU stats now cached in background - tab loads instantly instead of waiting 500ms

## [0.2.6] - 2025-12-09

### Changed

- **Config editor**: Replaced custom TOML editor with [simple-code-editor](https://github.com/justcaliturner/simple-code-editor) for better syntax highlighting, also because own one was glitched out.
- **Keyboard shortcut**: Added Ctrl+S / Cmd+S to save config (same as clicking "Save & Reload")

## [0.2.5] - 2025-12-09

### Added

- **System monitoring tab**: View CPU, RAM, disk usage, and uptime for the herbst host
  - Configurable via `[system]` section in config.toml
  - Custom disk path support for monitoring specific mounts
  - Real-time stats with auto-refresh every 3 seconds
- **Comprehensive README**: Full configuration reference for all settings

## [0.2.4] - 2025-12-07

### Fixed

- **Agent token reload**: Agent tokens now properly reload when config changes - no more server restart needed after adding agents
- **WebSocket keepalive**: Improved ping timing and context handling for more stable connections

## [0.2.3] - 2025-12-07

### Added

- **LICENSE**: Added MIT license

### Changed

- **Favicon**: Custom Herbst leaf icon with beige background for better visibility
- Removed unused Vite default assets (vue.svg, vite.svg)

## [0.2.2] - 2025-12-06

### Fixed

- **Agent connection stability**: Added WebSocket ping/pong keepalive to prevent proxy timeouts and broken pipe errors
- **Agent connection tracking**: Properly track connected/disconnected state for remote Docker agents
- **Docker build**: Fixed platform-specific dependency issues and TypeScript compilation in CI

## [0.2.0] - 2025-12-06

### Added

- **Service sections**: Group services with `[[section]]` and `title` in config.toml
- **Line numbers** in the config editor
- **Environment variables**: `HERBST_HOST` and `HERBST_AGENT_PROTOCOL` for flexible deployment
- **Auto-generated agent tokens**: No need to manually configure tokens anymore

### Changed

- **Docker config restructure** (breaking): `[docker]` split into `[docker.local]` for local containers and `[[docker.agent]]` for remote agents
- Improved default `config.toml` with visual section headers and better comments

### Fixed

- Config editor cursor position now stays in sync with text
- Config editor scroll sync at bottom of file

### Removed

- Search bar (temporarily disabled)

## [0.1.1] - 2025-12-05

### Fixed

- Docker image tag configuration in CI workflow

## [0.1.0] - 2025-12-05

### Added

- Initial Herbst dashboard with service cards, local docker monitoring
- Docker agent for container monitoring on different machines
- Config editor UI for config.toml and themes.toml
