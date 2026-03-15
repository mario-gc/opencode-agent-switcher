package cli

import (
	"fmt"

	"github.com/charmbracelet/huh"

	"opencode-agent-switcher/models"
)

func PromptAgentSelection(agents []models.Agent) (int, error) {
	var selectedName string
	options := make([]huh.Option[string], len(agents))

	for i, agent := range agents {
		display := fmt.Sprintf("%s (Current: %s)", agent.Name, agent.CurrentModel)
		options[i] = huh.NewOption(display, agent.Name)
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select an agent to update").
				Options(options...).
				Value(&selectedName).
				Height(10),
		),
	)

	if err := form.Run(); err != nil {
		return -1, err
	}

	for i, agent := range agents {
		if agent.Name == selectedName {
			return i, nil
		}
	}

	return -1, fmt.Errorf("selected agent not found")
}

func PromptModelSelection(modelOptions []models.ModelOption) (int, error) {
	var selectedID string
	options := make([]huh.Option[string], len(modelOptions))

	for i, model := range modelOptions {
		options[i] = huh.NewOption(model.Display, model.ID)
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select a new model").
				Options(options...).
				Value(&selectedID).
				Height(15),
		),
	)

	if err := form.Run(); err != nil {
		return -1, err
	}

	for i, model := range modelOptions {
		if model.ID == selectedID {
			return i, nil
		}
	}

	return -1, fmt.Errorf("selected model not found")
}

func PromptConfirm(message string) (bool, error) {
	var confirmed bool

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title(message).
				Value(&confirmed).
				Affirmative("Yes").
				Negative("No"),
		),
	)

	if err := form.Run(); err != nil {
		return false, err
	}

	return confirmed, nil
}
