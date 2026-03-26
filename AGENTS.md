# Agentic Coding Guidelines for Opencode Agent Switcher

This document provides instructions and guidelines for AI agents operating within the `opencode-agent-switcher` codebase.

## 1. Project Overview

`opencode-agent-switcher` is a Go CLI tool designed to manage and switch AI models and modes for various agents in the Opencode environment. It interacts with the `opencode` CLI and modifies agent configuration files.

- **Language:** Go 1.23+
- **Entry Point:** `main.go`
- **Module:** `opencode-agent-switcher`
- **Dependencies:** 
  - `gopkg.in/yaml.v3` - YAML parsing
  - `github.com/charmbracelet/huh` - Interactive TUI

### Package Structure
| Package | Purpose |
|---------|---------|
| `cli/` | User interface and TUI prompts (`cli/prompt.go`) |
| `config/` | Configuration loading and parsing (`config/config.go`) |
| `agents/` | Agent discovery and modification logic (`agents/agents.go`) |
| `models/` | Shared data structures (`models/models.go`) |

## 2. Build and Test Commands

### Build
```bash
go build -o opencode-agent-switcher main.go    # Build binary
go run main.go                                  # Run without building
./opencode-agent-switcher                       # Run built binary
```

### Testing
```bash
go test ./...                                   # Run all tests
go test -v ./...                                # Verbose output
go test -v ./agents -run TestLoadAgents         # Single test function
go test -v ./config                             # Tests for specific package
go test -race ./...                             # Race detector
go test -cover ./...                            # Coverage report
go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out  # HTML coverage
```

### Linting and Formatting
```bash
gofmt -s -w .                                   # Format all files
go vet ./...                                    # Static analysis
go mod tidy                                     # Clean up go.mod
make lint                                       # Run golangci-lint
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
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
	"gopkg.in/yaml.v3"

	"opencode-agent-switcher/models"
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

### Configuration Sources
The tool loads agents from multiple sources (in order of precedence):
1. **Project JSON:** `./opencode.json` (highest precedence)
2. **Project Markdown:** `.opencode/agents/*.md`
3. **Global JSON:** `~/.config/opencode/opencode.json`
4. **Global Markdown:** `~/.config/opencode/agents/*.md`

Project configurations override global ones for agents with the same name.

### Agent File Format
Agent files are Markdown with YAML frontmatter:
```yaml
---
model: provider/model-id
description: Agent description
mode: primary  # optional: primary, subagent, or all
---
# Agent instructions here
```

### JSON Config Format
Agents can also be defined in `opencode.json`:
```json
{
  "agent": {
    "my-agent": {
      "description": "Agent description",
      "model": "provider/model-id",
      "mode": "subagent"
    }
  }
}
```

### TUI Interaction (Huh? Library)
The tool uses `github.com/charmbracelet/huh` for interactive prompts:
- `PromptAgentSelection()` - Select from agent list with sort option and source indicators
- `PromptSortSelection()` - Select sorting method (Agent A-Z/Z-A, Model A-Z/Z-A)
- `PromptActionSelection()` - Choose action (Change Model / Change Mode / Back)
- `PromptModelSelection()` - Select from model list with custom input option
- `PromptModeSelection()` - Select mode (primary/subagent/all)
- `PromptCustomModelInput()` - Enter custom model ID
- `PromptAddModeField()` - Ask whether to add mode field
- `PromptConfirm()` - Yes/No confirmation
- `PromptUndo()` - Undo changes confirmation

### Sorting Feature
The main menu includes a "Sort by..." option that allows sorting the agent list:
- **Agent name (A-Z)** - Sort alphabetically by agent name (default)
- **Agent name (Z-A)** - Sort reverse-alphabetically by agent name
- **Model name (A-Z)** - Sort alphabetically by model ID
- **Model name (Z-A)** - Sort reverse-alphabetically by model ID

The sort preference persists during the session. Sorting also applies to the model selection menu.

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

### TUI Form Pattern
```go
form := huh.NewForm(
    huh.NewGroup(
        huh.NewSelect[string]().
            Title("Select an option").
            Options(options...).
            Value(&selected),
    ),
)
if err := form.Run(); err != nil {
    return err
}
```

## 7. Adding New Code

When adding new functionality:
1. Place data structures in `models/` if shared
2. Place business logic in appropriate package (`agents/`, `config/`, `cli/`)
3. Export functions that may be used by other packages
4. Add comments to exported functions
5. Add unit tests in `*_test.go` files
6. Run `make check` before committing

## 8. Cursor/Copilot Rules

No `.cursor/rules/`, `.cursorrules`, or `.github/copilot-instructions.md` files exist. Follow standard Go best practices and this document.