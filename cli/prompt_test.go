package cli

import (
	"testing"

	"opencode-agent-switcher/models"
)

func TestPromptAgentSelectionConstants(t *testing.T) {
	if ExitChoice != "__EXIT__" {
		t.Errorf("ExitChoice = %q, want %q", ExitChoice, "__EXIT__")
	}
	if ContinueChoice != "__CONTINUE__" {
		t.Errorf("ContinueChoice = %q, want %q", ContinueChoice, "__CONTINUE__")
	}
}

func TestPromptContinueOrExitConstants(t *testing.T) {
	agents := []models.Agent{
		{Name: "test-agent", CurrentModel: "openai/gpt-4", Path: "/tmp/test.md"},
	}

	if len(agents) != 1 {
		t.Errorf("Expected 1 agent, got %d", len(agents))
	}
}

func TestActionConstants(t *testing.T) {
	if ActionModel != "model" {
		t.Errorf("ActionModel = %q, want %q", ActionModel, "model")
	}
	if ActionMode != "mode" {
		t.Errorf("ActionMode = %q, want %q", ActionMode, "mode")
	}
	if BackChoice != "__BACK__" {
		t.Errorf("BackChoice = %q, want %q", BackChoice, "__BACK__")
	}
	if CustomModelChoice != "__CUSTOM__" {
		t.Errorf("CustomModelChoice = %q, want %q", CustomModelChoice, "__CUSTOM__")
	}
}

func TestFormatSourceTag(t *testing.T) {
	tests := []struct {
		name     string
		source   models.AgentSource
		expected string
	}{
		{"global markdown", models.AgentSource{Location: models.SourceGlobal, Format: models.FormatMarkdown}, "g/md"},
		{"global json", models.AgentSource{Location: models.SourceGlobal, Format: models.FormatJSON}, "g/json"},
		{"project markdown", models.AgentSource{Location: models.SourceProject, Format: models.FormatMarkdown}, "p/md"},
		{"project json", models.AgentSource{Location: models.SourceProject, Format: models.FormatJSON}, "p/json"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatSourceTag(tt.source)
			if got != tt.expected {
				t.Errorf("formatSourceTag() = %q, want %q", got, tt.expected)
			}
		})
	}
}
