# Opencode Agent Switcher

[![Go Report Card](https://goreportcard.com/badge/github.com/mario-gc/opencode-agent-switcher)](https://goreportcard.com/report/github.com/mario-gc/opencode-agent-switcher)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?logo=go)](https://golang.org/)

A CLI tool for managing and switching AI models and modes for Opencode agents.

## Features

- **Agent Discovery:** Automatically detects available agents from multiple sources:
  - Global markdown: `~/.config/opencode/agents/*.md`
  - Global JSON: `~/.config/opencode/opencode.json`
  - Project markdown: `.opencode/agents/*.md`
  - Project JSON: `./opencode.json`
- **Model Switching:** Change the AI model assigned to any agent
- **Mode Switching:** Change agent mode (primary/subagent/all)
- **Sorting Options:** Sort agents and models alphabetically (A-Z/Z-A) with case-sensitivity toggle
- **Custom Model Input:** Enter custom model IDs directly (format: `provider/model`)
- **Interactive TUI:** Beautiful terminal user interface using [Huh?](https://github.com/charmbracelet/huh)
- **Batch Updates:** Detects if multiple agents use the same model and offers to update them all
- **Undo Support:** Restore previous settings after updates
- **Source Indicators:** See where each agent is defined (global/project, markdown/JSON)

## Prerequisites

- **Go:** Version 1.23 or higher
- **Opencode:** The `opencode` CLI tool must be installed and configured

## Installation

### Download Binary

Download the latest release for your platform from the [Releases page](https://github.com/mario-gc/opencode-agent-switcher/releases):

```bash
# Linux (amd64)
curl -sL https://github.com/mario-gc/opencode-agent-switcher/releases/latest/download/opencode-agent-switcher_0.6.0_linux_amd64.tar.gz | tar xz

# Linux (arm64)
curl -sL https://github.com/mario-gc/opencode-agent-switcher/releases/latest/download/opencode-agent-switcher_0.6.0_linux_arm64.tar.gz | tar xz

# Make executable
chmod +x opencode-agent-switcher
```

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

### Command Line Options

| Option | Description |
|--------|-------------|
| `-v`, `--version` | Show version information |

```bash
opencode-agent-switcher --version
# Output: opencode-agent-switcher 0.6.0 (commit: abc1234, built: 2026-03-27)
```

### Workflow

1. The tool loads your Opencode configuration and available agents from all sources
2. An interactive menu appears with:
   - **Sort by...** - Change how agents are sorted (Agent A-Z/Z-A, Model A-Z/Z-A)
   - All available agents with their current model, mode, and source
   - An "Exit" option to quit the application
3. Select an agent to modify
4. Choose an action:
   - **Change Model** - Select a new AI model or enter a custom one
   - **Change Mode** - Switch between primary/subagent/all modes
   - **Back** - Return to agent selection
5. If changing mode and the agent has no mode set, choose whether to add the field
6. If other agents use the same model, you'll be asked to update them all
7. After updating, you can undo changes or continue

### Source Indicators

Agents are tagged with their source location:
- `[g/md]` - Global markdown file
- `[g/json]` - Global JSON config
- `[p/md]` - Project markdown file
- `[p/json]` - Project JSON config

### Configuration Precedence

When agents have the same name in different sources, project-level configurations take precedence over global ones.

## Development

### Project Structure

```
.
├── main.go              # Entry point with main loop
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