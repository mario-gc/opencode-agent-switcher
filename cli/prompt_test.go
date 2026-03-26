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
	if CaseSensitiveToggle != "__CASE_TOGGLE__" {
		t.Errorf("CaseSensitiveToggle = %q, want %q", CaseSensitiveToggle, "__CASE_TOGGLE__")
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

func TestGetSortDisplay(t *testing.T) {
	tests := []struct {
		sortBy        string
		caseSensitive bool
		expected      string
	}{
		{models.SortAgentAsc, true, "Agent A-Z"},
		{models.SortAgentDesc, true, "Agent Z-A"},
		{models.SortModelAsc, true, "Model A-Z"},
		{models.SortModelDesc, true, "Model Z-A"},
		{"unknown", true, "Agent A-Z"},
		{models.SortAgentAsc, false, "Agent A-Z (case-insensitive)"},
		{models.SortModelDesc, false, "Model Z-A (case-insensitive)"},
	}

	for _, tt := range tests {
		name := tt.sortBy
		if !tt.caseSensitive {
			name += "_case-insensitive"
		}
		t.Run(name, func(t *testing.T) {
			got := getSortDisplay(tt.sortBy, tt.caseSensitive)
			if got != tt.expected {
				t.Errorf("getSortDisplay(%q, %v) = %q, want %q", tt.sortBy, tt.caseSensitive, got, tt.expected)
			}
		})
	}
}

func TestSortAgents(t *testing.T) {
	agents := []models.Agent{
		{Name: "zebra", CurrentModel: "openai/gpt-3"},
		{Name: "alpha", CurrentModel: "openai/gpt-4"},
		{Name: "beta", CurrentModel: "anthropic/claude"},
	}

	t.Run("Agent A-Z", func(t *testing.T) {
		sorted := SortAgents(agents, models.SortAgentAsc, true)
		if sorted[0].Name != "alpha" || sorted[1].Name != "beta" || sorted[2].Name != "zebra" {
			t.Errorf("SortAgents Agent A-Z failed: got %v", []string{sorted[0].Name, sorted[1].Name, sorted[2].Name})
		}
	})

	t.Run("Agent Z-A", func(t *testing.T) {
		sorted := SortAgents(agents, models.SortAgentDesc, true)
		if sorted[0].Name != "zebra" || sorted[1].Name != "beta" || sorted[2].Name != "alpha" {
			t.Errorf("SortAgents Agent Z-A failed: got %v", []string{sorted[0].Name, sorted[1].Name, sorted[2].Name})
		}
	})

	t.Run("Model A-Z", func(t *testing.T) {
		sorted := SortAgents(agents, models.SortModelAsc, true)
		if sorted[0].CurrentModel != "anthropic/claude" {
			t.Errorf("SortAgents Model A-Z failed: first should be anthropic/claude, got %s", sorted[0].CurrentModel)
		}
	})

	t.Run("Model Z-A", func(t *testing.T) {
		sorted := SortAgents(agents, models.SortModelDesc, true)
		if sorted[0].CurrentModel != "openai/gpt-4" {
			t.Errorf("SortAgents Model Z-A failed: first should be openai/gpt-4, got %s", sorted[0].CurrentModel)
		}
	})

	t.Run("Does not modify original", func(t *testing.T) {
		original := []models.Agent{
			{Name: "zebra", CurrentModel: "z"},
			{Name: "alpha", CurrentModel: "a"},
		}
		_ = SortAgents(original, models.SortAgentAsc, true)
		if original[0].Name != "zebra" || original[1].Name != "alpha" {
			t.Error("SortAgents modified the original slice")
		}
	})

	t.Run("Empty slice", func(t *testing.T) {
		sorted := SortAgents([]models.Agent{}, models.SortAgentAsc, true)
		if len(sorted) != 0 {
			t.Errorf("SortAgents empty slice should return empty, got %d", len(sorted))
		}
	})

	t.Run("Case-insensitive Agent A-Z", func(t *testing.T) {
		mixedCase := []models.Agent{
			{Name: "Zebra", CurrentModel: "z"},
			{Name: "alpha", CurrentModel: "a"},
			{Name: "Beta", CurrentModel: "b"},
		}
		sorted := SortAgents(mixedCase, models.SortAgentAsc, false)
		if sorted[0].Name != "alpha" || sorted[1].Name != "Beta" || sorted[2].Name != "Zebra" {
			t.Errorf("SortAgents case-insensitive A-Z failed: got %v", []string{sorted[0].Name, sorted[1].Name, sorted[2].Name})
		}
	})

	t.Run("Case-sensitive Agent A-Z", func(t *testing.T) {
		mixedCase := []models.Agent{
			{Name: "zebra", CurrentModel: "z"},
			{Name: "Alpha", CurrentModel: "a"},
			{Name: "beta", CurrentModel: "b"},
		}
		sorted := SortAgents(mixedCase, models.SortAgentAsc, true)
		if sorted[0].Name != "Alpha" || sorted[1].Name != "beta" || sorted[2].Name != "zebra" {
			t.Errorf("SortAgents case-sensitive A-Z failed: got %v", []string{sorted[0].Name, sorted[1].Name, sorted[2].Name})
		}
	})
}

func TestSortModels(t *testing.T) {
	opts := []models.ModelOption{
		{ID: "openai/gpt-4", Display: "GPT-4"},
		{ID: "anthropic/claude", Display: "Claude"},
		{ID: "google/gemini", Display: "Gemini"},
	}

	t.Run("Model A-Z", func(t *testing.T) {
		sorted := SortModels(opts, models.SortModelAsc, true)
		if sorted[0].ID != "anthropic/claude" || sorted[1].ID != "google/gemini" || sorted[2].ID != "openai/gpt-4" {
			t.Errorf("SortModels A-Z failed: got %v", []string{sorted[0].ID, sorted[1].ID, sorted[2].ID})
		}
	})

	t.Run("Model Z-A", func(t *testing.T) {
		sorted := SortModels(opts, models.SortModelDesc, true)
		if sorted[0].ID != "openai/gpt-4" || sorted[1].ID != "google/gemini" || sorted[2].ID != "anthropic/claude" {
			t.Errorf("SortModels Z-A failed: got %v", []string{sorted[0].ID, sorted[1].ID, sorted[2].ID})
		}
	})

	t.Run("Does not modify original", func(t *testing.T) {
		original := []models.ModelOption{
			{ID: "z/model", Display: "Z"},
			{ID: "a/model", Display: "A"},
		}
		_ = SortModels(original, models.SortModelAsc, true)
		if original[0].ID != "z/model" || original[1].ID != "a/model" {
			t.Error("SortModels modified the original slice")
		}
	})

	t.Run("Empty slice", func(t *testing.T) {
		sorted := SortModels([]models.ModelOption{}, models.SortModelAsc, true)
		if len(sorted) != 0 {
			t.Errorf("SortModels empty slice should return empty, got %d", len(sorted))
		}
	})

	t.Run("Case-insensitive Model A-Z", func(t *testing.T) {
		mixedCase := []models.ModelOption{
			{ID: "Openai/GPT-4", Display: "GPT-4"},
			{ID: "anthropic/claude", Display: "Claude"},
		}
		sorted := SortModels(mixedCase, models.SortModelAsc, false)
		if sorted[0].ID != "anthropic/claude" {
			t.Errorf("SortModels case-insensitive A-Z failed: first should be anthropic/claude, got %s", sorted[0].ID)
		}
	})
}
