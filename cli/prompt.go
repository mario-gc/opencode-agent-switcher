package cli

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/huh"

	"opencode-agent-switcher/models"
)

const (
	ExitChoice          = "__EXIT__"
	ContinueChoice      = "__CONTINUE__"
	BackChoice          = "__BACK__"
	SortChoice          = "__SORT__"
	CustomModelChoice   = "__CUSTOM__"
	ActionModel         = "model"
	ActionMode          = "mode"
	CaseSensitiveToggle = "__CASE_TOGGLE__"
	TemplatesChoice     = "__TEMPLATES__"
	TemplateSave        = "save"
	TemplateShow        = "show"
	TemplateLoad        = "load"
	TemplateDelete      = "delete"
)

func PromptAgentSelection(agents []models.Agent, currentSort string, caseSensitive bool) (string, error) {
	var selectedName string

	sortDisplay := getSortDisplay(currentSort, caseSensitive)
	sortOption := fmt.Sprintf("Sort by... (%s)", sortDisplay)

	options := make([]huh.Option[string], 0, len(agents)+3)
	options = append(options, huh.NewOption(sortOption, SortChoice))
	options = append(options, huh.NewOption("Templates", TemplatesChoice))

	for _, agent := range agents {
		sourceTag := formatSourceTag(agent.Source)
		modeDisplay := agent.Mode
		if modeDisplay == "" {
			modeDisplay = "all (default)"
		}
		display := fmt.Sprintf("%s [%s] (Model: %s, Mode: %s)", agent.Name, sourceTag, agent.CurrentModel, modeDisplay)
		options = append(options, huh.NewOption(display, agent.Name))
	}

	options = append(options, huh.NewOption("Exit", ExitChoice))

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
		return "", err
	}

	return selectedName, nil
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

func getSortDisplay(sortBy string, caseSensitive bool) string {
	var sortType string
	switch sortBy {
	case models.SortAgentAsc:
		sortType = "Agent A-Z"
	case models.SortAgentDesc:
		sortType = "Agent Z-A"
	case models.SortModelAsc:
		sortType = "Model A-Z"
	case models.SortModelDesc:
		sortType = "Model Z-A"
	default:
		sortType = "Agent A-Z"
	}

	if !caseSensitive {
		sortType += " (case-insensitive)"
	}
	return sortType
}

func PromptSortSelection(currentSort string, caseSensitive bool) (string, bool, error) {
	var selected string

	caseDisplay := "On"
	if !caseSensitive {
		caseDisplay = "Off"
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Sort by").
				Options(
					huh.NewOption("Agent name (A-Z)", models.SortAgentAsc),
					huh.NewOption("Agent name (Z-A)", models.SortAgentDesc),
					huh.NewOption("Model name (A-Z)", models.SortModelAsc),
					huh.NewOption("Model name (Z-A)", models.SortModelDesc),
					huh.NewOption(fmt.Sprintf("Case-sensitive: %s", caseDisplay), CaseSensitiveToggle),
					huh.NewOption("Back", BackChoice),
				).
				Value(&selected),
		),
	)

	if err := form.Run(); err != nil {
		return currentSort, caseSensitive, err
	}

	if selected == BackChoice {
		return currentSort, caseSensitive, nil
	}

	if selected == CaseSensitiveToggle {
		return currentSort, !caseSensitive, nil
	}

	return selected, caseSensitive, nil
}

func SortAgents(agents []models.Agent, sortBy string, caseSensitive bool) []models.Agent {
	result := make([]models.Agent, len(agents))
	copy(result, agents)

	compareStrings := func(a, b string) bool {
		if caseSensitive {
			return a < b
		}
		return strings.ToLower(a) < strings.ToLower(b)
	}

	switch sortBy {
	case models.SortAgentAsc:
		sort.Slice(result, func(i, j int) bool {
			return compareStrings(result[i].Name, result[j].Name)
		})
	case models.SortAgentDesc:
		sort.Slice(result, func(i, j int) bool {
			return compareStrings(result[j].Name, result[i].Name)
		})
	case models.SortModelAsc:
		sort.Slice(result, func(i, j int) bool {
			return compareStrings(result[i].CurrentModel, result[j].CurrentModel)
		})
	case models.SortModelDesc:
		sort.Slice(result, func(i, j int) bool {
			return compareStrings(result[j].CurrentModel, result[i].CurrentModel)
		})
	default:
		sort.Slice(result, func(i, j int) bool {
			return compareStrings(result[i].Name, result[j].Name)
		})
	}

	return result
}

func SortModels(modelOptions []models.ModelOption, sortBy string, caseSensitive bool) []models.ModelOption {
	result := make([]models.ModelOption, len(modelOptions))
	copy(result, modelOptions)

	compareStrings := func(a, b string) bool {
		if caseSensitive {
			return a < b
		}
		return strings.ToLower(a) < strings.ToLower(b)
	}

	switch sortBy {
	case models.SortModelAsc:
		sort.Slice(result, func(i, j int) bool {
			return compareStrings(result[i].ID, result[j].ID)
		})
	case models.SortModelDesc:
		sort.Slice(result, func(i, j int) bool {
			return compareStrings(result[j].ID, result[i].ID)
		})
	default:
		sort.Slice(result, func(i, j int) bool {
			return compareStrings(result[i].ID, result[j].ID)
		})
	}

	return result
}

func PromptTemplateMenu() (string, error) {
	var choice string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Templates").
				Options(
					huh.NewOption("Save current configuration as template", TemplateSave),
					huh.NewOption("Show existing templates", TemplateShow),
					huh.NewOption("Return to main menu", BackChoice),
				).
				Value(&choice),
		),
	)

	if err := form.Run(); err != nil {
		return "", err
	}

	return choice, nil
}

func PromptTemplateName() (string, error) {
	var name string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Enter template name").
				Placeholder("my-template").
				Value(&name),
		),
	)

	if err := form.Run(); err != nil {
		return "", err
	}

	return name, nil
}

func PromptTemplateSelection(templates []models.Template) (string, error) {
	if len(templates) == 0 {
		return "", nil
	}

	var selected string
	options := make([]huh.Option[string], len(templates)+1)

	for i, template := range templates {
		display := fmt.Sprintf("%s (created: %s, agents: %d)", template.Name, formatDate(template.CreatedAt), len(template.Agents))
		options[i] = huh.NewOption(display, template.Name)
	}
	options[len(templates)] = huh.NewOption("Back", BackChoice)

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select a template").
				Options(options...).
				Value(&selected).
				Height(15),
		),
	)

	if err := form.Run(); err != nil {
		return "", err
	}

	return selected, nil
}

func formatDate(timestamp string) string {
	if timestamp == "" {
		return "unknown"
	}
	t, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		return timestamp
	}
	return t.Format("2006-01-02")
}

func PromptTemplateAction(templateName string) (string, error) {
	var action string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title(fmt.Sprintf("Template: %s - What would you like to do?", templateName)).
				Options(
					huh.NewOption("Load this template", TemplateLoad),
					huh.NewOption("Delete this template", TemplateDelete),
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

func PromptTemplateOverwrite(templateName string) (bool, error) {
	var choice string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title(fmt.Sprintf("Template '%s' already exists. Overwrite?", templateName)).
				Options(
					huh.NewOption("Yes, overwrite", "yes"),
					huh.NewOption("No, use different name", "no"),
				).
				Value(&choice),
		),
	)

	if err := form.Run(); err != nil {
		return false, err
	}

	return choice == "yes", nil
}

func PromptTemplateLoadConfirm(matchedCount, unmatchedCount int) (bool, error) {
	var choice string

	message := fmt.Sprintf("Load template? (%d agents will be updated, %d unmatched)", matchedCount, unmatchedCount)

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title(message).
				Options(
					huh.NewOption("Yes, load template", "yes"),
					huh.NewOption("No, go back", "no"),
				).
				Value(&choice),
		),
	)

	if err := form.Run(); err != nil {
		return false, err
	}

	return choice == "yes", nil
}

func PromptTemplateDeleteConfirm(templateName string) (bool, error) {
	var choice string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title(fmt.Sprintf("Delete template '%s'?", templateName)).
				Options(
					huh.NewOption("Yes, delete", "yes"),
					huh.NewOption("No, keep it", "no"),
				).
				Value(&choice),
		),
	)

	if err := form.Run(); err != nil {
		return false, err
	}

	return choice == "yes", nil
}
