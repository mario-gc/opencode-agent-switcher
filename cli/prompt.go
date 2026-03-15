package cli

import (
	"fmt"

	"github.com/charmbracelet/huh"

	"opencode-agent-switcher/models"
)

const ExitChoice = "__EXIT__"
const ContinueChoice = "__CONTINUE__"

func PromptAgentSelection(agents []models.Agent) (int, error) {
	var selectedName string
	options := make([]huh.Option[string], len(agents)+1)

	options[0] = huh.NewOption("Exit", ExitChoice)
	for i, agent := range agents {
		display := fmt.Sprintf("%s (Current: %s)", agent.Name, agent.CurrentModel)
		options[i+1] = huh.NewOption(display, agent.Name)
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

	if selectedName == ExitChoice {
		return -2, nil
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

func PromptContinueOrExit() (bool, error) {
	var choice string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("What would you like to do?").
				Options(
					huh.NewOption("Continue (select another agent)", ContinueChoice),
					huh.NewOption("Exit", ExitChoice),
				).
				Value(&choice),
		),
	)

	if err := form.Run(); err != nil {
		return false, err
	}

	return choice == ContinueChoice, nil
}

func PromptUndo(message string) (bool, error) {
	var choice string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title(message).
				Options(
					huh.NewOption("Undo - restore previous model", "undo"),
					huh.NewOption("Keep changes", "keep"),
				).
				Value(&choice),
		),
	)

	if err := form.Run(); err != nil {
		return false, err
	}

	return choice == "undo", nil
}
