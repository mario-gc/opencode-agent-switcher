# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.4.0] - 2026-03-22

### Added
- **Mode switching**: Change agent mode (primary/subagent/all) in addition to model
- **Custom model input**: Enter custom model ID directly instead of selecting from list
- **Multi-source agent discovery**: Load agents from all configuration sources
  - Global markdown: `~/.config/opencode/agents/*.md`
  - Global JSON: `~/.config/opencode/opencode.json`
  - Project markdown: `.opencode/agents/*.md`
  - Project JSON: `./opencode.json`
- **Source indicators**: Agents display their source location (e.g., `[g/md]`, `[p/json]`)
- **Action selection menu**: Choose between "Change Model" or "Change Mode" after selecting agent
- **Mode field prompt**: Ask user whether to add mode field when agent has none set

### Changed
- Refactored agent loading to support multiple configuration sources
- Updated agent display to show mode information in selection list
- `UpdateAgentModel` now requires agent name as parameter for JSON config support
- Project configs take precedence over global configs for same agent name

### Fixed
- Proper handling of agents without mode field (defaults to "all")
- Return to agent selection menu after completing changes (instead of action menu)

## [0.3.0] - 2026-03-15

### Added
- Exit option in main menu - users can now exit directly from agent selection
- Undo functionality - after updating agents, users can undo changes and restore previous models
- Loop to main menu - after completing changes, users can choose to continue or exit
- Basic tests for cli package constants

### Changed
- Refactored main.go into modular `runAgentUpdate` function for better loop handling
- Agent list is now reloaded after each iteration when continuing

## [0.2.0] - 2026-03-15

### Added
- Interactive TUI using Huh? library for agent and model selection
- Confirmation dialog for batch updates using TUI
- CI/CD pipeline with GitHub Actions (build, test, lint, security scan)
- Comprehensive `.gitignore` for Go projects
- `Makefile` for build automation with useful targets
- `.golangci.yml` for linting configuration
- `CONTRIBUTING.md` with development guidelines and GitFlow workflow
- Unit tests for agents package (LoadAgents, ParseFrontmatter, UpdateAgentModel, ValidateModelID)
- Unit tests for config package (GetAvailableModels, isValidModelID)

### Changed
- Renamed project from "Agent Switcher" to "Opencode Agent Switcher"
- Updated module name from `agent-switcher` to `opencode-agent-switcher`
- Replaced text-based prompts with interactive TUI selection
- Improved README with badges and clearer installation instructions

### Security
- Added symlink protection to prevent symlink attacks in agents directory
- Added path traversal validation to ensure files stay within expected directories
- Added model ID validation to prevent injection attacks (each segment validated)
- Changed file permissions from `0644` to `0600` for written agent files
- Added frontmatter size limit to prevent memory exhaustion

## [0.1.0] - 2026-03-15

### Added
- Initial release of opencode-agent-switcher
- CLI tool for managing AI agent model configurations
- Agent discovery from `~/.config/opencode/agents/`
- Model discovery via `opencode models` CLI or config file fallback
- Interactive agent and model selection
- Batch update for multiple agents using the same model
- YAML frontmatter parsing for agent configuration files
- MIT License

[Unreleased]: https://github.com/mario-gc/opencode-agent-switcher/compare/v0.4.0...HEAD
[0.4.0]: https://github.com/mario-gc/opencode-agent-switcher/compare/v0.3.0...v0.4.0
[0.3.0]: https://github.com/mario-gc/opencode-agent-switcher/compare/v0.2.0...v0.3.0
[0.2.0]: https://github.com/mario-gc/opencode-agent-switcher/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/mario-gc/opencode-agent-switcher/releases/tag/v0.1.0