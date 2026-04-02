package templates

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"opencode-agent-switcher/models"
)

func TestValidateTemplateName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid simple", "my-template", false},
		{"valid with underscore", "my_template", false},
		{"valid alphanumeric", "Template123", false},
		{"empty string", "", true},
		{"too long", "a_very_long_template_name_that_exceeds_sixty_four_characters_limit", true},
		{"starts with dash", "-template", true},
		{"starts with underscore", "_template", true},
		{"contains space", "my template", true},
		{"contains special char", "template@name", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTemplateName(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateTemplateName(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestSaveAndLoadTemplate(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "templates_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	oldDir, dirErr := GetTemplatesDir()
	if dirErr == nil {
		os.RemoveAll(oldDir)
	}

	homeDir := filepath.Join(tmpDir, "home")
	if mkdirErr := os.MkdirAll(homeDir, 0755); mkdirErr != nil {
		t.Fatalf("Failed to create home dir: %v", mkdirErr)
	}

	os.Setenv("HOME", homeDir)
	defer os.Unsetenv("HOME")

	agents := []models.Agent{
		{
			Name:         "architect",
			Path:         "/test/path",
			CurrentModel: "openai/gpt-4",
			Mode:         "primary",
			Source:       models.AgentSource{Location: models.SourceGlobal, Format: models.FormatMarkdown},
		},
		{
			Name:         "backend-dev",
			Path:         "json:global",
			CurrentModel: "anthropic/claude-3",
			Mode:         "subagent",
			Source:       models.AgentSource{Location: models.SourceGlobal, Format: models.FormatJSON},
		},
	}

	templateName := "test-template"
	if saveErr := SaveTemplate(templateName, agents); saveErr != nil {
		t.Fatalf("SaveTemplate() error = %v", saveErr)
	}

	templates, loadErr := LoadTemplates()
	if loadErr != nil {
		t.Fatalf("LoadTemplates() error = %v", loadErr)
	}

	if len(templates) != 1 {
		t.Errorf("LoadTemplates() returned %d templates, want 1", len(templates))
	}

	if templates[0].Name != templateName {
		t.Errorf("Template name = %q, want %q", templates[0].Name, templateName)
	}

	if len(templates[0].Agents) != 2 {
		t.Errorf("Template agents count = %d, want 2", len(templates[0].Agents))
	}

	architect, ok := templates[0].Agents["architect"]
	if !ok {
		t.Error("Template missing 'architect' agent")
	} else {
		if architect.Model != "openai/gpt-4" {
			t.Errorf("architect model = %q, want %q", architect.Model, "openai/gpt-4")
		}
		if architect.Mode != "primary" {
			t.Errorf("architect mode = %q, want %q", architect.Mode, "primary")
		}
	}
}

func TestTemplateExists(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "templates_exists_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	os.Setenv("HOME", tmpDir)
	defer os.Unsetenv("HOME")

	exists, err := TemplateExists("nonexistent")
	if err != nil {
		t.Fatalf("TemplateExists() error = %v", err)
	}
	if exists {
		t.Error("TemplateExists() returned true for nonexistent template")
	}

	agents := []models.Agent{
		{Name: "test", CurrentModel: "test/model", Source: models.AgentSource{}},
	}
	if saveErr := SaveTemplate("existing", agents); saveErr != nil {
		t.Fatalf("SaveTemplate() error = %v", saveErr)
	}

	exists, err = TemplateExists("existing")
	if err != nil {
		t.Fatalf("TemplateExists() error = %v", err)
	}
	if !exists {
		t.Error("TemplateExists() returned false for existing template")
	}
}

func TestDeleteTemplate(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "templates_delete_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	os.Setenv("HOME", tmpDir)
	defer os.Unsetenv("HOME")

	agents := []models.Agent{
		{Name: "test", CurrentModel: "test/model", Source: models.AgentSource{}},
	}
	if saveErr := SaveTemplate("to-delete", agents); saveErr != nil {
		t.Fatalf("SaveTemplate() error = %v", saveErr)
	}

	if delErr := DeleteTemplate("to-delete"); delErr != nil {
		t.Fatalf("DeleteTemplate() error = %v", delErr)
	}

	exists, checkErr := TemplateExists("to-delete")
	if checkErr != nil {
		t.Fatalf("TemplateExists() error = %v", checkErr)
	}
	if exists {
		t.Error("Template still exists after deletion")
	}
}

func TestLoadTemplateByName(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "templates_byname_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	os.Setenv("HOME", tmpDir)
	defer os.Unsetenv("HOME")

	agents := []models.Agent{
		{Name: "test", CurrentModel: "test/model", Mode: "primary", Source: models.AgentSource{}},
	}
	if saveErr := SaveTemplate("my-template", agents); saveErr != nil {
		t.Fatalf("SaveTemplate() error = %v", saveErr)
	}

	template, loadErr := LoadTemplateByName("my-template")
	if loadErr != nil {
		t.Fatalf("LoadTemplateByName() error = %v", loadErr)
	}

	if template.Name != "my-template" {
		t.Errorf("Template name = %q, want %q", template.Name, "my-template")
	}

	_, err = LoadTemplateByName("nonexistent")
	if err == nil {
		t.Error("LoadTemplateByName() should return error for nonexistent template")
	}
}

func TestMatchAgents(t *testing.T) {
	template := models.Template{
		Name:      "test",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
		Agents: map[string]models.AgentAssignment{
			"architect": {
				Model:  "openai/gpt-4",
				Mode:   "primary",
				Source: models.AgentSource{Location: models.SourceGlobal, Format: models.FormatMarkdown},
			},
			"backend-dev": {
				Model:  "anthropic/claude-3",
				Mode:   "subagent",
				Source: models.AgentSource{Location: models.SourceProject, Format: models.FormatJSON},
			},
			"missing": {
				Model:  "test/model",
				Mode:   "all",
				Source: models.AgentSource{Location: models.SourceGlobal, Format: models.FormatMarkdown},
			},
		},
	}

	currentAgents := []models.Agent{
		{
			Name:         "architect",
			Path:         "/path/to/architect.md",
			CurrentModel: "old/model",
			Mode:         "all",
			Source:       models.AgentSource{Location: models.SourceGlobal, Format: models.FormatMarkdown},
		},
		{
			Name:         "backend-dev",
			Path:         "json:project",
			CurrentModel: "another/model",
			Mode:         "primary",
			Source:       models.AgentSource{Location: models.SourceProject, Format: models.FormatJSON},
		},
	}

	matched, unmatched := MatchAgents(template, currentAgents)

	if len(matched) != 2 {
		t.Errorf("matched count = %d, want 2", len(matched))
	}

	if len(unmatched) != 1 {
		t.Errorf("unmatched count = %d, want 1", len(unmatched))
	}

	for _, agent := range matched {
		if agent.Name == "architect" {
			if agent.CurrentModel != "openai/gpt-4" {
				t.Errorf("architect model = %q, want %q", agent.CurrentModel, "openai/gpt-4")
			}
		}
		if agent.Name == "backend-dev" {
			if agent.CurrentModel != "anthropic/claude-3" {
				t.Errorf("backend-dev model = %q, want %q", agent.CurrentModel, "anthropic/claude-3")
			}
		}
	}
}

func TestMatchAgentsSourceMismatch(t *testing.T) {
	template := models.Template{
		Name: "test",
		Agents: map[string]models.AgentAssignment{
			"agent": {
				Model:  "new/model",
				Mode:   "primary",
				Source: models.AgentSource{Location: models.SourceGlobal, Format: models.FormatMarkdown},
			},
		},
	}

	currentAgents := []models.Agent{
		{
			Name:         "agent",
			Path:         "json:project",
			CurrentModel: "old/model",
			Source:       models.AgentSource{Location: models.SourceProject, Format: models.FormatJSON},
		},
	}

	matched, unmatched := MatchAgents(template, currentAgents)

	if len(matched) != 0 {
		t.Errorf("matched count = %d, want 0 (source mismatch)", len(matched))
	}

	if len(unmatched) != 1 {
		t.Errorf("unmatched count = %d, want 1", len(unmatched))
	}
}
