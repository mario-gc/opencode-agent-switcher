# Agentic Coding Guidelines for Agent Switcher

This document provides instructions and guidelines for AI agents operating within the `agent-switcher` codebase.

## 1. Project Overview

`agent-switcher` is a Go CLI tool designed to manage and switch AI models for various agents in the Opencode environment. It interacts with the `opencode` CLI and modifies agent configuration files.

- **Language:** Go 1.22+
- **Entry Point:** `main.go`
- **Module:** `agent-switcher`
- **Dependencies:** `gopkg.in/yaml.v3` (only external dependency)

### Package Structure
| Package | Purpose |
|---------|---------|
| `cli/` | User interface and prompts (`cli/prompt.go`) |
| `config/` | Configuration loading and parsing (`config/config.go`) |
| `agents/` | Agent discovery and modification logic (`agents/agents.go`) |
| `models/` | Shared data structures (`models/models.go`) |

## 2. Build and Test Commands

### Build
```bash
go build -o agent-switcher main.go    # Build binary
go run main.go                         # Run without building
./agent-switcher                       # Run built binary
```

### Testing
```bash
go test ./...                          # Run all tests
go test -v ./...                       # Verbose output
go test -v ./agents -run TestLoadAgents  # Single test function
go test -v ./config                    # Tests for specific package
go test -race ./...                    # Race detector
go test -cover ./...                   # Coverage report
go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out  # HTML coverage
```

### Linting and Formatting
```bash
gofmt -s -w .                          # Format all files
go vet ./...                           # Static analysis
go mod tidy                            # Clean up go.mod
```

## 3. Code Style and Conventions

### Formatting
- **Strictly follow `gofmt` standards**
- All files end with a newline
- Use tabs for indentation (standard Go behavior)
- Max line length follows Go convention (no strict limit, but be reasonable)

### Imports
Group imports with blank lines: Standard Library → Third-party → Local packages
```go
import (
	"bufio"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"

	"agent-switcher/models"
)
```

### Naming Conventions
| Type | Convention | Example |
|------|------------|---------|
| Exported functions | PascalCase | `LoadAgents`, `ParseFrontmatter` |
| Unexported functions | camelCase | `parseConfig` |
| Structs | PascalCase | `OpencodeConfig`, `ModelOption` |
| Interfaces | PascalCase + verb/er | `Reader`, `Writer` |
| Constants | PascalCase or UPPER_CASE | `MaxRetries` |
| Acronyms | Consistent case | `ID` (not `Id`), `HTTP` (not `Http`) |
| Variables | Short, descriptive | `cfg`, `err`, `agentList` |

### Error Handling
Always check errors immediately. Wrap with context when bubbling up:
```go
if err != nil {
    return nil, fmt.Errorf("failed to load config: %w", err)
}
```

- Use `log.Fatalf` **only in `main.go`**
- Return errors from all other packages
- Use `fmt.Errorf` with `%w` for error wrapping
- Silent continuation (with `continue`) is acceptable for non-critical file processing errors

### Functions and Comments
- Add comments to all exported functions explaining purpose
- Comment format: `// FunctionName does X. Returns Y if Z.`
```go
// LoadAgents reads all .md files from agents directory
func LoadAgents(agentsDir string) ([]models.Agent, error) {
```

### Types and Structs
- Define shared structures in `models/models.go`
- Use struct tags for serialization:
```go
type Provider struct {
	Name   string           `json:"name"`
	Models map[string]Model `json:"models"`
}
```

## 4. Architecture Patterns

### Configuration Location
- Global config: `~/.config/opencode/opencode.json`
- Agent configs: `~/.config/opencode/agents/*.md`

### Agent File Format
Agent files are Markdown with YAML frontmatter:
```yaml
---
model: provider/model-id
description: Agent description
---
# Agent instructions here
```

### External CLI Dependency
- Tool calls `opencode models` to fetch available models
- Falls back to parsing `opencode.json` if CLI unavailable

### File Operations
- Use `os` and `path/filepath` for cross-platform compatibility
- Use `filepath.Join` for path construction
- Use `os.ReadFile`/`os.WriteFile` for file I/O

## 5. Agent Behavior Rules

- **No Hallucinations:** Do not invent flags or commands not in the codebase
- **Safety:** Ensure file modifications are safe (read-modify-write pattern)
- **Idempotency:** Operations should be idempotent where possible
- **User Confirmation:** Prompt before making destructive changes
- **Graceful Degradation:** Fall back to alternative methods when primary fails

## 6. Common Patterns in This Codebase

### Function Return Pattern
```go
func LoadAgents(agentsDir string) ([]models.Agent, error) {
    files, err := os.ReadDir(agentsDir)
    if err != nil {
        return nil, err
    }
    // ... processing
    return agents, nil
}
```

### Iterating and Filtering
```go
for _, file := range files {
    if filepath.Ext(file.Name()) == ".md" {
        // process
    }
}
```

### Type Assertion from map[string]interface{}
```go
model, ok := frontmatter["model"].(string)
if !ok {
    continue
}
```

## 7. Adding New Code

When adding new functionality:
1. Place data structures in `models/` if shared
2. Place business logic in appropriate package (`agents/`, `config/`, `cli/`)
3. Export functions that may be used by other packages
4. Add comments to exported functions
5. Run `gofmt -s -w .` and `go vet ./...` before committing

## 8. Cursor/Copilot Rules

No `.cursor/rules/`, `.cursorrules`, or `.github/copilot-instructions.md` files exist. Follow standard Go best practices and this document.