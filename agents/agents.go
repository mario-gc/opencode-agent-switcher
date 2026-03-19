package agents

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"opencode-agent-switcher/config"
	"opencode-agent-switcher/models"

	"gopkg.in/yaml.v3"
)

const maxFrontmatterSize = 64 * 1024

var validSegment = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_\-\.]*$`)

func ValidateModelID(modelID string) error {
	if modelID == "" {
		return fmt.Errorf("model ID cannot be empty")
	}
	if len(modelID) > 256 {
		return fmt.Errorf("model ID exceeds maximum length")
	}

	segments := strings.Split(modelID, "/")
	if len(segments) < 2 {
		return fmt.Errorf("model ID must be in format 'provider/model'")
	}

	for _, segment := range segments {
		if segment == "" {
			return fmt.Errorf("model ID contains empty segment")
		}
		if !validSegment.MatchString(segment) {
			return fmt.Errorf("invalid segment '%s': must start with alphanumeric and contain only alphanumeric, dash, underscore, or dot", segment)
		}
	}

	return nil
}

func isValidMode(mode string) bool {
	return mode == "" || mode == "primary" || mode == "subagent" || mode == "all"
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

func LoadAllAgents() ([]models.Agent, error) {
	agentMap := make(map[string]models.Agent)

	globalAgentsDir, err := getGlobalAgentsDir()
	if err == nil {
		globalMD, err := LoadAgentsFromDir(globalAgentsDir, models.SourceGlobal, models.FormatMarkdown)
		if err == nil {
			for _, agent := range globalMD {
				agentMap[agent.Name] = agent
			}
		}
	}

	globalCfg, err := config.LoadGlobalConfig()
	if err == nil {
		globalJSON := config.GetAgentsFromConfig(globalCfg, models.SourceGlobal, models.FormatJSON)
		for _, agent := range globalJSON {
			agentMap[agent.Name] = agent
		}
	}

	projectMD, err := LoadAgentsFromDir(getProjectAgentsDir(), models.SourceProject, models.FormatMarkdown)
	if err == nil {
		for _, agent := range projectMD {
			agentMap[agent.Name] = agent
		}
	}

	projectCfg, err := config.LoadProjectConfig()
	if err == nil && projectCfg != nil {
		projectJSON := config.GetAgentsFromConfig(projectCfg, models.SourceProject, models.FormatJSON)
		for _, agent := range projectJSON {
			agentMap[agent.Name] = agent
		}
	}

	var agentList []models.Agent
	for _, agent := range agentMap {
		agentList = append(agentList, agent)
	}

	return agentList, nil
}

func getGlobalAgentsDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "opencode", "agents"), nil
}

func getProjectAgentsDir() string {
	return filepath.Join(".", ".opencode", "agents")
}

func LoadAgents(agentsDir string) ([]models.Agent, error) {
	return LoadAgentsFromDir(agentsDir, models.SourceGlobal, models.FormatMarkdown)
}

func LoadAgentsFromDir(agentsDir, location, format string) ([]models.Agent, error) {
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

			mode, ok := frontmatter["mode"].(string)
			if !ok || !isValidMode(mode) {
				mode = ""
			}

			agentList = append(agentList, models.Agent{
				Name:         strings.TrimSuffix(file.Name(), ".md"),
				Path:         path,
				CurrentModel: model,
				Description:  description,
				Mode:         mode,
				Source: models.AgentSource{
					Location: location,
					Format:   format,
				},
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

func UpdateAgentModel(agentPath, agentName, newModel string) error {
	if err := ValidateModelID(newModel); err != nil {
		return fmt.Errorf("invalid model ID: %w", err)
	}

	if strings.HasPrefix(agentPath, "json:") {
		location := strings.TrimPrefix(agentPath, "json:")
		var configPath string
		if location == models.SourceGlobal {
			var err error
			configPath, err = config.GetGlobalConfigPath()
			if err != nil {
				return err
			}
		} else {
			configPath = config.GetProjectConfigPath()
		}
		return updateAgentFieldInJSON(configPath, agentName, "model", newModel)
	}

	return updateAgentFieldInMarkdown(agentPath, "model", newModel)
}

func UpdateAgentMode(agentPath, agentName, newMode string) error {
	if !isValidMode(newMode) {
		return fmt.Errorf("invalid mode: must be 'primary', 'subagent', or 'all'")
	}

	if strings.HasPrefix(agentPath, "json:") {
		location := strings.TrimPrefix(agentPath, "json:")
		var configPath string
		if location == models.SourceGlobal {
			var err error
			configPath, err = config.GetGlobalConfigPath()
			if err != nil {
				return err
			}
		} else {
			configPath = config.GetProjectConfigPath()
		}
		return updateAgentFieldInJSON(configPath, agentName, "mode", newMode)
	}

	return updateAgentFieldInMarkdown(agentPath, "mode", newMode)
}

func updateAgentFieldInJSON(configPath, agentName, field, value string) error {
	return config.UpdateAgentInJSON(configPath, agentName, field, value)
}

func updateAgentFieldInMarkdown(agentPath, field, value string) error {
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

	frontmatter[field] = value

	newFrontmatter, err := yaml.Marshal(frontmatter)
	if err != nil {
		return err
	}

	newContent := fmt.Sprintf("---\n%s---%s", string(newFrontmatter), parts[2])
	return os.WriteFile(agentPath, []byte(newContent), 0600)
}
