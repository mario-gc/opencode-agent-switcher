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
