package config

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"opencode-agent-switcher/models"
)

var validSegmentPattern = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_\-\.]*$`)

func isValidModelID(modelID string) bool {
	if modelID == "" || len(modelID) > 256 {
		return false
	}

	segments := strings.Split(modelID, "/")
	if len(segments) < 2 {
		return false
	}

	for _, segment := range segments {
		if segment == "" || !validSegmentPattern.MatchString(segment) {
			return false
		}
	}

	return true
}

func LoadGlobalConfig() (*models.OpencodeConfig, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	configPath := filepath.Join(home, ".config", "opencode", "opencode.json")
	return loadConfigFile(configPath)
}

func LoadProjectConfig() (*models.OpencodeConfig, error) {
	configPath := filepath.Join(".", "opencode.json")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, nil
	}
	return loadConfigFile(configPath)
}

func loadConfigFile(configPath string) (*models.OpencodeConfig, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var cfg models.OpencodeConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func GetAgentsFromConfig(cfg *models.OpencodeConfig, location, format string) []models.Agent {
	if cfg == nil || cfg.Agent == nil {
		return nil
	}

	var agentList []models.Agent
	for name, agentCfg := range cfg.Agent {
		agentList = append(agentList, models.Agent{
			Name:         name,
			Path:         fmt.Sprintf("json:%s", location),
			CurrentModel: agentCfg.Model,
			Description:  agentCfg.Description,
			Mode:         agentCfg.Mode,
			Source: models.AgentSource{
				Location: location,
				Format:   format,
			},
		})
	}
	return agentList
}

func GetAvailableModels(cfg *models.OpencodeConfig) []models.ModelOption {
	var options []models.ModelOption

	if cfg == nil || cfg.Provider == nil {
		return options
	}

	for providerID, provider := range cfg.Provider {
		for modelID, model := range provider.Models {
			modelIDFull := fmt.Sprintf("%s/%s", providerID, modelID)
			if !isValidModelID(modelIDFull) {
				continue
			}

			display := modelIDFull
			if model.Name != "" {
				display = fmt.Sprintf("%s (%s)", display, model.Name)
			}
			options = append(options, models.ModelOption{
				ID:       modelIDFull,
				Display:  display,
				Provider: providerID,
			})
		}
	}

	return options
}

func GetModelsFromCLI() ([]models.ModelOption, error) {
	cmd := exec.Command("opencode", "models")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var options []models.ModelOption
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.Split(line, "/")
		if len(parts) >= 2 && isValidModelID(line) {
			provider := parts[0]
			options = append(options, models.ModelOption{
				ID:       line,
				Display:  line,
				Provider: provider,
			})
		}
	}

	return options, nil
}

func GetGlobalConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "opencode", "opencode.json"), nil
}

func GetProjectConfigPath() string {
	return filepath.Join(".", "opencode.json")
}

func UpdateAgentInJSON(configPath, agentName, field, value string) error {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	var cfg map[string]interface{}
	if unmarshalErr := json.Unmarshal(data, &cfg); unmarshalErr != nil {
		return unmarshalErr
	}

	agents, ok := cfg["agent"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("no agent configurations found in %s", configPath)
	}

	agent, ok := agents[agentName].(map[string]interface{})
	if !ok {
		return fmt.Errorf("agent %s not found in %s", agentName, configPath)
	}

	agent[field] = value

	newData, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, newData, 0600)
}
