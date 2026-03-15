package config

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"agent-switcher/models"
)

var validModelIDPattern = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_\-\.\/]*$`)

func isValidModelID(modelID string) bool {
	if modelID == "" || len(modelID) > 256 {
		return false
	}
	return validModelIDPattern.MatchString(modelID)
}

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

	var cfg models.OpencodeConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func GetAvailableModels(cfg *models.OpencodeConfig) []models.ModelOption {
	var options []models.ModelOption

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
