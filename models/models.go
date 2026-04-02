package models

// Source constants for agent location and format.
const (
	SourceGlobal   = "global"
	SourceProject  = "project"
	FormatMarkdown = "markdown"
	FormatJSON     = "json"
)

// Sort constants for agent and model sorting options.
const (
	SortAgentAsc  = "agent-asc"
	SortAgentDesc = "agent-desc"
	SortModelAsc  = "model-asc"
	SortModelDesc = "model-desc"
)

// DefaultSort is the default sorting method for agent lists.
var DefaultSort = SortAgentAsc

// OpencodeConfig represents the opencode configuration file structure.
type OpencodeConfig struct {
	Provider map[string]Provider    `json:"provider"`
	Agent    map[string]AgentConfig `json:"agent"`
}

// Provider represents a model provider configuration.
type Provider struct {
	Name   string           `json:"name"`
	Models map[string]Model `json:"models"`
}

// Model represents a model configuration with limits and modalities.
type Model struct {
	Name       string                 `json:"name"`
	Limit      map[string]int         `json:"limit"`
	Modalities map[string][]string    `json:"modalities"`
	Variants   map[string]interface{} `json:"variants"`
}

// AgentConfig represents an agent configuration in opencode.json.
type AgentConfig struct {
	Description string `json:"description"`
	Model       string `json:"model"`
	Mode        string `json:"mode"`
}

// ModelOption represents a selectable model option.
type ModelOption struct {
	ID       string
	Display  string
	Provider string
}

// AgentSource represents the source location and format of an agent.
type AgentSource struct {
	Location string
	Format   string
}

// Agent represents a loaded agent with its configuration.
type Agent struct {
	Name         string
	Path         string
	CurrentModel string
	Description  string
	Mode         string
	Source       AgentSource
}

// AgentAssignment represents a model and mode assignment for an agent.
type AgentAssignment struct {
	Model  string      `json:"model"`
	Mode   string      `json:"mode"`
	Source AgentSource `json:"source"`
}

// Template represents a saved agent configuration template.
type Template struct {
	Name      string                     `json:"name"`
	CreatedAt string                     `json:"created_at"`
	Agents    map[string]AgentAssignment `json:"agents"`
}
