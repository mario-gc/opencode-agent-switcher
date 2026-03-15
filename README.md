# Agent Switcher

[![Go Report Card](https://goreportcard.com/badge/github.com/mario-gc/opencode-agent-switcher)](https://goreportcard.com/report/github.com/mario-gc/opencode-agent-switcher)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?logo=go)](https://golang.org/)

A CLI tool for managing and switching AI models for Opencode agents.

## Features

- **Agent Discovery:** Automatically detects available agents configured in `~/.config/opencode/agents/`
- **Model Discovery:** Fetches available AI models from the `opencode` CLI or falls back to the configuration file
- **Interactive Selection:** User-friendly CLI to select an agent and a target model
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
go build -o agent-switcher main.go
```

### Using Make

```bash
make build    # Build the binary
make install  # Install to GOPATH/bin
```

## Usage

Run the tool directly from the terminal:

```bash
agent-switcher
```

### Workflow

1. The tool loads your Opencode configuration and available agents
2. It presents a list of agents to choose from
3. After selecting an agent, it displays a list of available AI models
4. Select the new model you wish to assign to the agent
5. If other agents are currently using the same model as the selected agent, the tool will ask if you want to update those agents as well
6. The tool updates the agent configuration files with the new model ID

## Development

### Project Structure

```
.
├── main.go          # Entry point
├── cli/             # User interaction and prompts
├── config/          # Configuration loading and parsing
├── agents/          # Agent discovery and modification
├── models/          # Shared data structures
├── Makefile         # Build automation
├── .golangci.yml    # Linting configuration
└── go.mod           # Go module definition
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