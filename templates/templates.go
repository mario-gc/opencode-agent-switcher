package templates

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/mario-gc/opencode-agent-switcher/models"
)

var validTemplateName = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_\-]*$`)

const templatesDirName = "opencode-agent-switcher"

func GetTemplatesDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	dir := filepath.Join(home, ".config", templatesDirName, "templates")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("failed to create templates directory: %w", err)
	}
	return dir, nil
}

func ValidateTemplateName(name string) error {
	if name == "" {
		return fmt.Errorf("template name cannot be empty")
	}
	if len(name) > 64 {
		return fmt.Errorf("template name exceeds maximum length (64 characters)")
	}
	if !validTemplateName.MatchString(name) {
		return fmt.Errorf("template name must start with alphanumeric and contain only alphanumeric, underscore, or dash")
	}
	return nil
}

func TemplateExists(name string) (bool, error) {
	dir, err := GetTemplatesDir()
	if err != nil {
		return false, err
	}
	path := filepath.Join(dir, name+".json")
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	return err == nil, err
}

func SaveTemplate(name string, agents []models.Agent) error {
	if err := ValidateTemplateName(name); err != nil {
		return err
	}

	agentMap := make(map[string]models.AgentAssignment)
	for _, agent := range agents {
		agentMap[agent.Name] = models.AgentAssignment{
			Model:  agent.CurrentModel,
			Mode:   agent.Mode,
			Source: agent.Source,
		}
	}

	template := models.Template{
		Name:      name,
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
		Agents:    agentMap,
	}

	dir, err := GetTemplatesDir()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(template, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal template: %w", err)
	}

	path := filepath.Join(dir, name+".json")
	return os.WriteFile(path, data, 0600)
}

func LoadTemplates() ([]models.Template, error) {
	dir, err := GetTemplatesDir()
	if err != nil {
		return nil, err
	}

	files, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var templates []models.Template
	for _, file := range files {
		if filepath.Ext(file.Name()) != ".json" {
			continue
		}

		path := filepath.Join(dir, file.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}

		var template models.Template
		if err := json.Unmarshal(data, &template); err != nil {
			continue
		}

		templates = append(templates, template)
	}

	sort.Slice(templates, func(i, j int) bool {
		return strings.ToLower(templates[i].Name) < strings.ToLower(templates[j].Name)
	})

	return templates, nil
}

func LoadTemplateByName(name string) (models.Template, error) {
	dir, err := GetTemplatesDir()
	if err != nil {
		return models.Template{}, err
	}

	path := filepath.Join(dir, name+".json")
	data, err := os.ReadFile(path)
	if err != nil {
		return models.Template{}, fmt.Errorf("template '%s' not found", name)
	}

	var template models.Template
	if err := json.Unmarshal(data, &template); err != nil {
		return models.Template{}, fmt.Errorf("failed to parse template: %w", err)
	}

	return template, nil
}

func DeleteTemplate(name string) error {
	dir, err := GetTemplatesDir()
	if err != nil {
		return err
	}

	path := filepath.Join(dir, name+".json")
	if err := os.Remove(path); err != nil {
		return fmt.Errorf("failed to delete template: %w", err)
	}

	return nil
}

func MatchAgents(template models.Template, currentAgents []models.Agent) ([]models.Agent, []string) {
	var matched []models.Agent
	var unmatched []string

	for templateName, assignment := range template.Agents {
		found := false
		for _, agent := range currentAgents {
			if agent.Name == templateName &&
				agent.Source.Location == assignment.Source.Location &&
				agent.Source.Format == assignment.Source.Format {
				matched = append(matched, models.Agent{
					Name:         agent.Name,
					Path:         agent.Path,
					CurrentModel: assignment.Model,
					Mode:         assignment.Mode,
					Source:       agent.Source,
				})
				found = true
				break
			}
		}
		if !found {
			sourceTag := fmt.Sprintf("%s/%s", assignment.Source.Location, assignment.Source.Format)
			unmatched = append(unmatched, fmt.Sprintf("%s [%s]", templateName, sourceTag))
		}
	}

	sort.Slice(matched, func(i, j int) bool {
		return strings.ToLower(matched[i].Name) < strings.ToLower(matched[j].Name)
	})
	sort.Strings(unmatched)

	return matched, unmatched
}
