package agents

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"agent-switcher/models"
	"gopkg.in/yaml.v3"
)

// LoadAgents reads all .md files from agents directory
func LoadAgents(agentsDir string) ([]models.Agent, error) {
	files, err := os.ReadDir(agentsDir)
	if err != nil {
		return nil, err
	}

	var agents []models.Agent
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".md" {
			path := filepath.Join(agentsDir, file.Name())
			content, err := os.ReadFile(path)
			if err != nil {
				continue
			}

			frontmatter, err := ParseFrontmatter(string(content))
			if err != nil {
				continue
			}

			model, ok := frontmatter["model"].(string)
			if !ok {
				continue
			}

			description, _ := frontmatter["description"].(string)

			agents = append(agents, models.Agent{
				Name:         strings.TrimSuffix(file.Name(), ".md"),
				Path:         path,
				CurrentModel: model,
				Description:  description,
			})
		}
	}

	return agents, nil
}

// ParseFrontmatter extracts YAML frontmatter from markdown content
func ParseFrontmatter(content string) (map[string]interface{}, error) {
	if !strings.HasPrefix(content, "---") {
		return nil, fmt.Errorf("no frontmatter found")
	}

	parts := strings.SplitN(content, "---", 3)
	if len(parts) < 3 {
		return nil, fmt.Errorf("invalid frontmatter format")
	}

	var frontmatter map[string]interface{}
	if err := yaml.Unmarshal([]byte(parts[1]), &frontmatter); err != nil {
		return nil, err
	}

	return frontmatter, nil
}

// UpdateAgentModel modifies the model field in agent frontmatter
func UpdateAgentModel(agentPath, newModel string) error {
	content, err := os.ReadFile(agentPath)
	if err != nil {
		return err
	}

	parts := strings.SplitN(string(content), "---", 3)
	if len(parts) < 3 {
		return fmt.Errorf("invalid frontmatter format")
	}

	var frontmatter map[string]interface{}
	if err := yaml.Unmarshal([]byte(parts[1]), &frontmatter); err != nil {
		return err
	}

	frontmatter["model"] = newModel

	newFrontmatter, err := yaml.Marshal(frontmatter)
	if err != nil {
		return err
	}

	newContent := fmt.Sprintf("---\n%s---%s", string(newFrontmatter), parts[2])
	return os.WriteFile(agentPath, []byte(newContent), 0644)
}
