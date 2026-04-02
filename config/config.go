package config

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/mario-gc/opencode-agent-switcher/models"
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

// LoadGlobalConfig loads the global opencode configuration file.
// Returns the configuration or an error if the file cannot be read or parsed.
func LoadGlobalConfig() (*models.OpencodeConfig, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	configPath := filepath.Join(home, ".config", "opencode", "opencode.json")
	return loadConfigFile(configPath)
}

// LoadProjectConfig loads the project-level opencode configuration file.
// Returns nil if the file does not exist, or an error if it cannot be parsed.
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

// GetAgentsFromConfig extracts agents from an opencode configuration.
// Parameters location and format specify the source metadata for each agent.
// Returns a list of agents defined in the configuration.
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

// GetAvailableModels extracts available models from an opencode configuration.
// Returns a list of model options with provider, ID, and display name.
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

// GetModelsFromCLI retrieves available models by calling the opencode CLI.
// Returns a list of model options or an error if the CLI command fails.
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

// GetGlobalConfigPath returns the path to the global opencode configuration file.
// Returns an error if the user home directory cannot be determined.
func GetGlobalConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "opencode", "opencode.json"), nil
}

// GetProjectConfigPath returns the path to the project-level opencode configuration file.
func GetProjectConfigPath() string {
	return filepath.Join(".", "opencode.json")
}

// UpdateAgentInJSON updates a field for an agent in a JSON configuration file.
// Returns an error if the agent does not exist or the file cannot be updated.
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
