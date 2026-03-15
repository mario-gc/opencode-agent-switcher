package agents

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"agent-switcher/models"

	"gopkg.in/yaml.v3"
)

const maxFrontmatterSize = 64 * 1024

var validModelID = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_\-\.\/]*$`)

func ValidateModelID(modelID string) error {
	if modelID == "" {
		return fmt.Errorf("model ID cannot be empty")
	}
	if len(modelID) > 256 {
		return fmt.Errorf("model ID exceeds maximum length")
	}
	if !validModelID.MatchString(modelID) {
		return fmt.Errorf("invalid model ID format: contains disallowed characters")
	}
	return nil
}

func isPathWithinDir(path, dir string) (bool, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return false, err
	}
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return false, err
	}
	return strings.HasPrefix(absPath, absDir+string(filepath.Separator)), nil
}

func LoadAgents(agentsDir string) ([]models.Agent, error) {
	absAgentsDir, err := filepath.Abs(agentsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve agents directory: %w", err)
	}

	files, err := os.ReadDir(agentsDir)
	if err != nil {
		return nil, err
	}

	var agentList []models.Agent
	for _, file := range files {
		if file.Type()&os.ModeSymlink != 0 {
			continue
		}

		if filepath.Ext(file.Name()) == ".md" {
			path := filepath.Join(agentsDir, file.Name())

			isWithin, err := isPathWithinDir(path, absAgentsDir)
			if err != nil || !isWithin {
				continue
			}

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

			description, ok := frontmatter["description"].(string)
			if !ok {
				description = ""
			}

			agentList = append(agentList, models.Agent{
				Name:         strings.TrimSuffix(file.Name(), ".md"),
				Path:         path,
				CurrentModel: model,
				Description:  description,
			})
		}
	}

	return agentList, nil
}

func ParseFrontmatter(content string) (map[string]interface{}, error) {
	if !strings.HasPrefix(content, "---") {
		return nil, fmt.Errorf("no frontmatter found")
	}

	parts := strings.SplitN(content, "---", 3)
	if len(parts) < 3 {
		return nil, fmt.Errorf("invalid frontmatter format")
	}

	frontmatterStr := parts[1]
	if len(frontmatterStr) > maxFrontmatterSize {
		return nil, fmt.Errorf("frontmatter exceeds maximum allowed size")
	}

	var frontmatter map[string]interface{}
	if err := yaml.Unmarshal([]byte(frontmatterStr), &frontmatter); err != nil {
		return nil, err
	}

	return frontmatter, nil
}

func UpdateAgentModel(agentPath, newModel string) error {
	if err := ValidateModelID(newModel); err != nil {
		return fmt.Errorf("invalid model ID: %w", err)
	}

	content, err := os.ReadFile(agentPath)
	if err != nil {
		return err
	}

	parts := strings.SplitN(string(content), "---", 3)
	if len(parts) < 3 {
		return fmt.Errorf("invalid frontmatter format")
	}

	var frontmatter map[string]interface{}
	if err = yaml.Unmarshal([]byte(parts[1]), &frontmatter); err != nil {
		return err
	}

	frontmatter["model"] = newModel

	newFrontmatter, err := yaml.Marshal(frontmatter)
	if err != nil {
		return err
	}

	newContent := fmt.Sprintf("---\n%s---%s", string(newFrontmatter), parts[2])
	return os.WriteFile(agentPath, []byte(newContent), 0600)
}