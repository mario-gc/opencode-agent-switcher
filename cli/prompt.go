package cli

import (
	"fmt"

	"github.com/charmbracelet/huh"

	"opencode-agent-switcher/models"
)

const (
	ExitChoice        = "__EXIT__"
	ContinueChoice    = "__CONTINUE__"
	BackChoice        = "__BACK__"
	CustomModelChoice = "__CUSTOM__"
	ActionModel       = "model"
	ActionMode        = "mode"
)

func PromptAgentSelection(agents []models.Agent) (int, error) {
	var selectedName string
	options := make([]huh.Option[string], len(agents)+1)

	options[0] = huh.NewOption("Exit", ExitChoice)
	for i, agent := range agents {
		sourceTag := formatSourceTag(agent.Source)
		modeDisplay := agent.Mode
		if modeDisplay == "" {
			modeDisplay = "all (default)"
		}
		display := fmt.Sprintf("%s [%s] (Model: %s, Mode: %s)", agent.Name, sourceTag, agent.CurrentModel, modeDisplay)
		options[i+1] = huh.NewOption(display, agent.Name)
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select an agent to update").
				Options(options...).
				Value(&selectedName).
				Height(15),
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

func formatSourceTag(source models.AgentSource) string {
	loc := "g"
	if source.Location == models.SourceProject {
		loc = "p"
	}
	fmtChar := "md"
	if source.Format == models.FormatJSON {
		fmtChar = "json"
	}
	return fmt.Sprintf("%s/%s", loc, fmtChar)
}

func PromptActionSelection(currentModel, currentMode string) (string, error) {
	var action string

	modeDisplay := currentMode
	if modeDisplay == "" {
		modeDisplay = "not set (default: all)"
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("What would you like to change?").
				Options(
					huh.NewOption(fmt.Sprintf("Change Model (current: %s)", currentModel), ActionModel),
					huh.NewOption(fmt.Sprintf("Change Mode (current: %s)", modeDisplay), ActionMode),
					huh.NewOption("Back", BackChoice),
				).
				Value(&action),
		),
	)

	if err := form.Run(); err != nil {
		return "", err
	}

	return action, nil
}

func PromptModeSelection(currentMode string) (string, error) {
	var mode string

	modeDisplay := currentMode
	if modeDisplay == "" {
		modeDisplay = "all (default)"
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title(fmt.Sprintf("Select mode (current: %s)", modeDisplay)).
				Options(
					huh.NewOption("primary", "primary"),
					huh.NewOption("subagent", "subagent"),
					huh.NewOption("all", "all"),
					huh.NewOption("Back", BackChoice),
				).
				Value(&mode),
		),
	)

	if err := form.Run(); err != nil {
		return "", err
	}

	return mode, nil
}

func PromptAddModeField() (bool, error) {
	var choice string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("This agent has no mode set (defaults to 'all'). Add explicit mode field?").
				Options(
					huh.NewOption("Yes, add mode field", "yes"),
					huh.NewOption("No, keep it unset", "no"),
				).
				Value(&choice),
		),
	)

	if err := form.Run(); err != nil {
		return false, err
	}

	return choice == "yes", nil
}

func PromptModelSelection(modelOptions []models.ModelOption) (string, error) {
	var selectedID string
	options := make([]huh.Option[string], len(modelOptions)+1)

	options[0] = huh.NewOption("Enter custom model...", CustomModelChoice)
	for i, model := range modelOptions {
		options[i+1] = huh.NewOption(model.Display, model.ID)
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
		return "", err
	}

	return selectedID, nil
}

func PromptCustomModelInput() (string, error) {
	var modelID string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Enter custom model ID (format: provider/model)").
				Placeholder("openai/gpt-4-turbo").
				Value(&modelID),
		),
	)

	if err := form.Run(); err != nil {
		return "", err
	}

	return modelID, nil
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
					huh.NewOption("Undo - restore previous value", "undo"),
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
