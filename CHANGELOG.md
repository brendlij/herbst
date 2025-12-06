# Changelog

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

### Added

- Initial Herbst dashboard with service cards, local docker monitoring.
- Docker agent for container monitoring on different machine.
- user can edit config.toml and themes.toml over UI.
