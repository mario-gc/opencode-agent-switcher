# Opencode Agent Switcher

[![Go Report Card](https://goreportcard.com/badge/github.com/mario-gc/opencode-agent-switcher)](https://goreportcard.com/report/github.com/mario-gc/opencode-agent-switcher)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?logo=go)](https://golang.org/)

A CLI tool for managing and switching AI models for Opencode agents.

## Features

- **Agent Discovery:** Automatically detects available agents configured in `~/.config/opencode/agents/`
- **Model Discovery:** Fetches available AI models from the `opencode` CLI or falls back to the configuration file
- **Interactive TUI:** Beautiful terminal user interface using [Huh?](https://github.com/charmbracelet/huh) for selection
- **Batch Updates:** Detects if multiple agents are using the same model and offers to update them all simultaneously
- **Configuration Management:** Safely updates the YAML configuration files for the agents

## Prerequisites

- **Go:** Version 1.23 or higher
- **Opencode:** The `opencode` CLI tool must be installed and configured
- **Configuration:** Expects `~/.config/opencode/opencode.json` and agent configurations in `~/.config/opencode/agents/`

## Installation

### From Source

```bash
# Clone the repository
git clone https://github.com/mario-gc/opencode-agent-switcher.git
cd opencode-agent-switcher

# Build and install
go install

# Or build locally
go build -o opencode-agent-switcher main.go
```

### Using Make

```bash
make build    # Build the binary
make install  # Install to GOPATH/bin
```

## Usage

Run the tool directly from the terminal:

```bash
opencode-agent-switcher
```

### Workflow

1. The tool loads your Opencode configuration and available agents
2. An interactive menu appears showing all available agents with their current models
3. Use arrow keys to navigate and Enter to select an agent
4. A second menu appears with all available models
5. Select the new model you wish to assign to the agent
6. If other agents use the same model, you'll be asked if you want to update them too
7. The tool updates the agent configuration files and shows confirmation

## Development

### Project Structure

```
.
├── main.go              # Entry point
├── cli/                 # User interaction and TUI prompts
│   └── prompt.go        # Huh? based interactive prompts
├── config/              # Configuration loading and parsing
│   └── config.go        # Opencode config handling
├── agents/              # Agent discovery and modification
│   └── agents.go        # Agent file operations
├── models/              # Shared data structures
├── Makefile             # Build automation
├── .golangci.yml        # Linting configuration
└── go.mod               # Go module definition
```

### Commands

```bash
make build         # Build the binary
make test          # Run tests with race detection
make test-coverage # Generate coverage report
make lint          # Run golangci-lint
make fmt           # Format code
make vet           # Run go vet
make check         # Run all checks (fmt, vet, lint, test)
make clean         # Remove build artifacts
```

### Running Tests

```bash
go test -v -race ./...
```

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## Changelog

See [CHANGELOG.md](CHANGELOG.md) for version history.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.