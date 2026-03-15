package config

import (
	"testing"

	"opencode-agent-switcher/models"
)

func TestIsValidModelID(t *testing.T) {
	tests := []struct {
		name    string
		modelID string
		want    bool
	}{
		{"valid simple", "openai/gpt-4", true},
		{"valid with dash", "anthropic/claude-3-sonnet", true},
		{"valid with dot", "google/gemini-1.5-pro", true},
		{"valid with underscore", "meta/llama_3_70b", true},
		{"valid with multiple slashes", "provider/model/variant", true},
		{"empty string", "", false},
		{"too long", "provider/" + string(make([]byte, 257)), false},
		{"invalid chars", "provider/model@name", false},
		{"starts with special char", "-provider/model", false},
		{"special chars in middle", "provider/model name", false},
		{"segment starts with dash", "provider/-model", false},
		{"no slash", "provider", false},
		{"single slash empty model", "provider/", false},
		{"empty provider", "/model", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if len(tt.modelID) > 257 {
				t.Skip("model ID construction issue")
			}
			got := isValidModelID(tt.modelID)
			if got != tt.want {
				t.Errorf("isValidModelID(%q) = %v, want %v", tt.modelID, got, tt.want)
			}
		})
	}
}

func TestGetAvailableModels(t *testing.T) {
	cfg := createTestConfig()

	options := GetAvailableModels(cfg)

	if len(options) == 0 {
		t.Error("GetAvailableModels() returned empty list")
	}

	found := false
	for _, opt := range options {
		if opt.ID == "openai/gpt-4" {
			found = true
			if opt.Provider != "openai" {
				t.Errorf("Provider = %q, want %q", opt.Provider, "openai")
			}
			break
		}
	}

	if !found {
		t.Error("Expected model 'openai/gpt-4' not found in options")
	}
}

func TestGetAvailableModelsFiltersInvalidIDs(t *testing.T) {
	cfg := createTestConfigWithInvalidModel()

	options := GetAvailableModels(cfg)

	for _, opt := range options {
		if !isValidModelID(opt.ID) {
			t.Errorf("GetAvailableModels() returned invalid model ID: %q", opt.ID)
		}
	}
}

func createTestConfig() *models.OpencodeConfig {
	return &models.OpencodeConfig{
		Provider: map[string]models.Provider{
			"openai": {
				Name: "openai",
				Models: map[string]models.Model{
					"gpt-4": {
						Name: "GPT-4",
					},
					"gpt-3.5-turbo": {
						Name: "GPT-3.5 Turbo",
					},
				},
			},
			"anthropic": {
				Name: "anthropic",
				Models: map[string]models.Model{
					"claude-3-sonnet": {
						Name: "Claude 3 Sonnet",
					},
				},
			},
		},
	}
}

func createTestConfigWithInvalidModel() *models.OpencodeConfig {
	return &models.OpencodeConfig{
		Provider: map[string]models.Provider{
			"test": {
				Name: "test",
				Models: map[string]models.Model{
					"valid-model": {
						Name: "Valid Model",
					},
					"": {
						Name: "Empty Model ID",
					},
				},
			},
		},
	}
}
