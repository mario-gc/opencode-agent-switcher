package models

// OpencodeConfig represents the parsed opencode.json
type OpencodeConfig struct {
	Provider map[string]Provider `json:"provider"`
}

// Provider represents a model provider (google, ollama, etc.)
type Provider struct {
	Name   string           `json:"name"`
	Models map[string]Model `json:"models"`
}

// Model represents a single AI model
type Model struct {
	Name       string                 `json:"name"`
	Limit      map[string]int         `json:"limit"`
	Modalities map[string][]string    `json:"modalities"`
	Variants   map[string]interface{} `json:"variants"`
}

// ModelOption for user selection
type ModelOption struct {
	ID       string
	Display  string
	Provider string
}

// Agent represents an agent configuration file
type Agent struct {
	Name         string
	Path         string
	CurrentModel string
	Description  string
}
