package models

const (
	SourceGlobal   = "global"
	SourceProject  = "project"
	FormatMarkdown = "markdown"
	FormatJSON     = "json"
)

type OpencodeConfig struct {
	Provider map[string]Provider    `json:"provider"`
	Agent    map[string]AgentConfig `json:"agent"`
}

type Provider struct {
	Name   string           `json:"name"`
	Models map[string]Model `json:"models"`
}

type Model struct {
	Name       string                 `json:"name"`
	Limit      map[string]int         `json:"limit"`
	Modalities map[string][]string    `json:"modalities"`
	Variants   map[string]interface{} `json:"variants"`
}

type AgentConfig struct {
	Description string `json:"description"`
	Model       string `json:"model"`
	Mode        string `json:"mode"`
}

type ModelOption struct {
	ID       string
	Display  string
	Provider string
}

type AgentSource struct {
	Location string
	Format   string
}

type Agent struct {
	Name         string
	Path         string
	CurrentModel string
	Description  string
	Mode         string
	Source       AgentSource
}
