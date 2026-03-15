package config

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"agent-switcher/models"
)

// LoadOpencodeConfig reads and parses ~/.config/opencode/opencode.json
func LoadOpencodeConfig() (*models.OpencodeConfig, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	configPath := filepath.Join(home, ".config", "opencode", "opencode.json")
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config models.OpencodeConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// GetAvailableModels extracts all model IDs from providers
func GetAvailableModels(config *models.OpencodeConfig) []models.ModelOption {
	var options []models.ModelOption

	for providerID, provider := range config.Provider {
		for modelID, model := range provider.Models {
			display := fmt.Sprintf("%s/%s", providerID, modelID)
			if model.Name != "" {
				display = fmt.Sprintf("%s (%s)", display, model.Name)
			}
			options = append(options, models.ModelOption{
				ID:       fmt.Sprintf("%s/%s", providerID, modelID),
				Display:  display,
				Provider: providerID,
			})
		}
	}

	return options
}

// GetModelsFromCLI runs `opencode models` and parses output
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

		// Parse line format: provider/model
		parts := strings.Split(line, "/")
		if len(parts) >= 2 {
			provider := parts[0]
			_ = strings.Join(parts[1:], "/") // model part not needed but kept for clarity
			options = append(options, models.ModelOption{
				ID:       line,
				Display:  line,
				Provider: provider,
			})
		}
	}

	return options, nil
}
