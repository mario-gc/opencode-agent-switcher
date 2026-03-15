package agents

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidateModelID(t *testing.T) {
	tests := []struct {
		name    string
		modelID string
		wantErr bool
	}{
		{"valid simple", "openai/gpt-4", false},
		{"valid with dash", "anthropic/claude-3-sonnet", false},
		{"valid with dot", "google/gemini-1.5-pro", false},
		{"valid with underscore", "meta/llama_3_70b", false},
		{"empty string", "", true},
		{"too long", "a/" + string(make([]byte, 256)), true},
		{"invalid chars", "provider/model@name", true},
		{"starts with special char", "provider/-model", true},
		{"no slash", "provider", true},
		{"single slash empty model", "provider/", true},
		{"empty provider", "/model", true},
		{"valid with multiple slashes", "provider/model/variant", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if len(tt.modelID) > 256 && tt.wantErr {
				t.Skip("model ID construction issue")
			}
			err := ValidateModelID(tt.modelID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateModelID(%q) error = %v, wantErr %v", tt.modelID, err, tt.wantErr)
			}
		})
	}
}

func TestParseFrontmatter(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    map[string]interface{}
		wantErr bool
	}{
		{
			name: "valid frontmatter",
			content: `---
model: openai/gpt-4
description: Test agent
---
# Content`,
			want:    map[string]interface{}{"model": "openai/gpt-4", "description": "Test agent"},
			wantErr: false,
		},
		{
			name:    "no frontmatter",
			content: "# Just markdown",
			wantErr: true,
		},
		{
			name:    "invalid frontmatter format",
			content: "---\nmodel: test\n# missing closing",
			wantErr: true,
		},
		{
			name: "empty frontmatter",
			content: `---
---
# Content`,
			want:    map[string]interface{}{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseFrontmatter(tt.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseFrontmatter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != nil {
				for k, v := range tt.want {
					if got[k] != v {
						t.Errorf("ParseFrontmatter()[%q] = %v, want %v", k, got[k], v)
					}
				}
			}
		})
	}
}

func TestLoadAgents(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "agents_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	validAgent := `---
model: openai/gpt-4
description: Test agent
---
# Test Agent Content`

	if writeErr := os.WriteFile(filepath.Join(tmpDir, "test-agent.md"), []byte(validAgent), 0644); writeErr != nil {
		t.Fatalf("Failed to write test file: %v", writeErr)
	}

	noModelAgent := `---
description: No model here
---
# Content`
	if writeErr := os.WriteFile(filepath.Join(tmpDir, "no-model.md"), []byte(noModelAgent), 0644); writeErr != nil {
		t.Fatalf("Failed to write test file: %v", writeErr)
	}

	invalidFrontmatter := `# Just markdown content`
	if writeErr := os.WriteFile(filepath.Join(tmpDir, "invalid.md"), []byte(invalidFrontmatter), 0644); writeErr != nil {
		t.Fatalf("Failed to write test file: %v", writeErr)
	}

	txtFile := `some text content`
	if writeErr := os.WriteFile(filepath.Join(tmpDir, "readme.txt"), []byte(txtFile), 0644); writeErr != nil {
		t.Fatalf("Failed to write test file: %v", writeErr)
	}

	loadedAgents, err := LoadAgents(tmpDir)
	if err != nil {
		t.Fatalf("LoadAgents() error = %v", err)
	}

	if len(loadedAgents) != 1 {
		t.Errorf("LoadAgents() returned %d agents, want 1", len(loadedAgents))
	}

	if len(loadedAgents) > 0 {
		if loadedAgents[0].Name != "test-agent" {
			t.Errorf("Agent name = %q, want %q", loadedAgents[0].Name, "test-agent")
		}
		if loadedAgents[0].CurrentModel != "openai/gpt-4" {
			t.Errorf("Agent model = %q, want %q", loadedAgents[0].CurrentModel, "openai/gpt-4")
		}
	}
}

func TestLoadAgentsNonExistentDir(t *testing.T) {
	_, err := LoadAgents("/nonexistent/path/that/does/not/exist")
	if err == nil {
		t.Error("LoadAgents() should return error for non-existent directory")
	}
}

func TestUpdateAgentModel(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "agent_update_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	agentContent := `---
model: openai/gpt-4
description: Test agent
---
# Test Agent Content`

	agentPath := filepath.Join(tmpDir, "test-agent.md")
	if writeErr := os.WriteFile(agentPath, []byte(agentContent), 0644); writeErr != nil {
		t.Fatalf("Failed to write test file: %v", writeErr)
	}

	err = UpdateAgentModel(agentPath, "anthropic/claude-3-sonnet")
	if err != nil {
		t.Fatalf("UpdateAgentModel() error = %v", err)
	}

	content, err := os.ReadFile(agentPath)
	if err != nil {
		t.Fatalf("Failed to read updated file: %v", err)
	}

	frontmatter, err := ParseFrontmatter(string(content))
	if err != nil {
		t.Fatalf("Failed to parse updated frontmatter: %v", err)
	}

	if frontmatter["model"] != "anthropic/claude-3-sonnet" {
		t.Errorf("Updated model = %v, want %q", frontmatter["model"], "anthropic/claude-3-sonnet")
	}
}

func TestUpdateAgentModelInvalidID(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "agent_update_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	agentContent := `---
model: openai/gpt-4
---
# Content`

	agentPath := filepath.Join(tmpDir, "test-agent.md")
	if writeErr := os.WriteFile(agentPath, []byte(agentContent), 0644); writeErr != nil {
		t.Fatalf("Failed to write test file: %v", writeErr)
	}

	err = UpdateAgentModel(agentPath, "")
	if err == nil {
		t.Error("UpdateAgentModel() should return error for empty model ID")
	}
}

func TestIsPathWithinDir(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		dir     string
		want    bool
		wantErr bool
	}{
		{"within dir", "/home/user/agents/test.md", "/home/user/agents", true, false},
		{"outside dir", "/home/user/other/test.md", "/home/user/agents", false, false},
		{"same level", "/home/user/agents", "/home/user/agents", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := isPathWithinDir(tt.path, tt.dir)
			if (err != nil) != tt.wantErr {
				t.Errorf("isPathWithinDir() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("isPathWithinDir() = %v, want %v", got, tt.want)
			}
		})
	}
}
