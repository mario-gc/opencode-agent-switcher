package main

import (
	"fmt"
	"log"
	"os"

	"opencode-agent-switcher/agents"
	"opencode-agent-switcher/cli"
	"opencode-agent-switcher/config"
	"opencode-agent-switcher/models"
	"opencode-agent-switcher/templates"
)

var (
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
)

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "-v") {
		fmt.Printf("opencode-agent-switcher %s (commit: %s, built: %s)\n", version, commit, date)
		os.Exit(0)
	}

	cfg, err := config.LoadGlobalConfig()
	if err != nil {
		fmt.Printf("Warning: Failed to load global opencode config: %v\n", err)
	}

	agentList, err := agents.LoadAllAgents()
	if err != nil {
		log.Fatalf("Failed to load agents: %v", err)
	}

	if len(agentList) == 0 {
		log.Fatalf("No agents found")
	}

	modelOptions, err := config.GetModelsFromCLI()
	if err != nil {
		fmt.Println("Warning: Failed to get models from CLI, falling back to config file...")
		modelOptions = config.GetAvailableModels(cfg)
	}

	if len(modelOptions) == 0 {
		log.Fatalf("No models found available")
	}

	currentSort := models.DefaultSort
	caseSensitive := true

	for {
		shouldContinue, loopErr := runAgentUpdate(agentList, modelOptions, &currentSort, &caseSensitive)
		if loopErr != nil {
			log.Fatalf("Error: %v", loopErr)
		}

		if !shouldContinue {
			break
		}

		var reloadErr error
		agentList, reloadErr = agents.LoadAllAgents()
		if reloadErr != nil {
			log.Fatalf("Failed to reload agents: %v", reloadErr)
		}
	}

	fmt.Println("\nGoodbye!")
}

func runAgentUpdate(agentList []models.Agent, modelOptions []models.ModelOption, currentSort *string, caseSensitive *bool) (bool, error) {
	sortedAgents := cli.SortAgents(agentList, *currentSort, *caseSensitive)

	selectedName, err := cli.PromptAgentSelection(sortedAgents, *currentSort, *caseSensitive)
	if err != nil {
		return false, err
	}

	if selectedName == cli.ExitChoice {
		return false, nil
	}

	if selectedName == cli.SortChoice {
		newSort, newCaseSensitive, sortErr := cli.PromptSortSelection(*currentSort, *caseSensitive)
		if sortErr != nil {
			return false, sortErr
		}
		*currentSort = newSort
		*caseSensitive = newCaseSensitive
		return true, nil
	}

	if selectedName == cli.TemplatesChoice {
		continueToMenu, templateErr := handleTemplates(agentList)
		if templateErr != nil {
			return false, templateErr
		}
		return continueToMenu, nil
	}

	var selectedAgent models.Agent
	for _, agent := range sortedAgents {
		if agent.Name == selectedName {
			selectedAgent = agent
			break
		}
	}

	for {
		action, err := cli.PromptActionSelection(selectedAgent.CurrentModel, selectedAgent.Mode)
		if err != nil {
			return false, err
		}

		if action == cli.BackChoice {
			return true, nil
		}

		if action == cli.ActionModel {
			continueToMenu, err := handleModelChange(selectedAgent, agentList, modelOptions, *currentSort, *caseSensitive)
			if err != nil {
				return false, err
			}
			return continueToMenu, nil
		}

		if action == cli.ActionMode {
			continueToMenu, err := handleModeChange(selectedAgent)
			if err != nil {
				return false, err
			}
			return continueToMenu, nil
		}
	}
}

func handleModelChange(selectedAgent models.Agent, agentList []models.Agent, modelOptions []models.ModelOption, currentSort string, caseSensitive bool) (bool, error) {
	sortedModels := cli.SortModels(modelOptions, currentSort, caseSensitive)
	selectedModelID, err := cli.PromptModelSelection(sortedModels)
	if err != nil {
		return false, err
	}

	if selectedModelID == cli.BackChoice {
		return true, nil
	}

	var selectedModel models.ModelOption
	if selectedModelID == cli.CustomModelChoice {
		customID, err := cli.PromptCustomModelInput()
		if err != nil {
			return false, err
		}
		if err := agents.ValidateModelID(customID); err != nil {
			fmt.Printf("Invalid model ID: %v\n", err)
			return true, nil
		}
		selectedModel = models.ModelOption{ID: customID}
	} else {
		for _, m := range modelOptions {
			if m.ID == selectedModelID {
				selectedModel = m
				break
			}
		}
	}

	agentsToUpdate := []models.Agent{selectedAgent}
	previousModels := make(map[string]string)
	previousModels[selectedAgent.Name] = selectedAgent.CurrentModel

	var otherAgents []models.Agent
	for _, a := range agentList {
		if a.Name != selectedAgent.Name && a.CurrentModel == selectedAgent.CurrentModel {
			otherAgents = append(otherAgents, a)
		}
	}

	if len(otherAgents) > 0 {
		message := fmt.Sprintf("%d other agent(s) use the same model. Update all?", len(otherAgents))
		confirmed, confirmErr := cli.PromptConfirm(message)
		if confirmErr != nil {
			return false, confirmErr
		}
		if confirmed {
			for _, a := range otherAgents {
				previousModels[a.Name] = a.CurrentModel
			}
			agentsToUpdate = append(agentsToUpdate, otherAgents...)
		}
	}

	fmt.Printf("\nUpdating %d agent(s) to model '%s'...\n", len(agentsToUpdate), selectedModel.ID)

	updatedAgents := []string{}
	for _, agent := range agentsToUpdate {
		updateErr := agents.UpdateAgentModel(agent.Path, agent.Name, selectedModel.ID)
		if updateErr != nil {
			log.Printf("Failed to update agent %s: %v", agent.Name, updateErr)
		} else {
			fmt.Printf("✓ Updated %s\n", agent.Name)
			updatedAgents = append(updatedAgents, agent.Name)
		}
	}

	if len(updatedAgents) > 0 {
		undoMessage := fmt.Sprintf("Updated %d agent(s). Undo changes?", len(updatedAgents))
		wantUndo, undoErr := cli.PromptUndo(undoMessage)
		if undoErr != nil {
			return false, undoErr
		}

		if wantUndo {
			fmt.Println("\nUndoing changes...")
			for _, agentName := range updatedAgents {
				for _, agent := range agentsToUpdate {
					if agent.Name == agentName {
						previousModel := previousModels[agentName]
						restoreErr := agents.UpdateAgentModel(agent.Path, agent.Name, previousModel)
						if restoreErr != nil {
							log.Printf("Failed to undo agent %s: %v", agentName, restoreErr)
						} else {
							fmt.Printf("✓ Restored %s to %s\n", agentName, previousModel)
						}
						break
					}
				}
			}
		}
	}

	return promptContinue()
}

func handleModeChange(selectedAgent models.Agent) (bool, error) {
	selectedMode, err := cli.PromptModeSelection(selectedAgent.Mode)
	if err != nil {
		return false, err
	}

	if selectedMode == cli.BackChoice {
		return true, nil
	}

	previousMode := selectedAgent.Mode
	shouldAddField := false

	if previousMode == "" {
		shouldAddField, err = cli.PromptAddModeField()
		if err != nil {
			return false, err
		}
		if !shouldAddField {
			fmt.Println("Keeping mode unset (defaults to 'all')")
			return promptContinue()
		}
	}

	fmt.Printf("\nUpdating agent '%s' mode to '%s'...\n", selectedAgent.Name, selectedMode)

	updateErr := agents.UpdateAgentMode(selectedAgent.Path, selectedAgent.Name, selectedMode)
	if updateErr != nil {
		log.Printf("Failed to update agent %s: %v", selectedAgent.Name, updateErr)
		return promptContinue()
	}

	fmt.Printf("✓ Updated %s mode to '%s'\n", selectedAgent.Name, selectedMode)

	undoMessage := "Updated agent mode. Undo changes?"
	wantUndo, undoErr := cli.PromptUndo(undoMessage)
	if undoErr != nil {
		return false, undoErr
	}

	if wantUndo {
		fmt.Println("\nUndoing changes...")
		if previousMode == "" {
			restoreErr := agents.UpdateAgentMode(selectedAgent.Path, selectedAgent.Name, "")
			if restoreErr != nil {
				log.Printf("Failed to undo agent %s: %v", selectedAgent.Name, restoreErr)
			} else {
				fmt.Printf("✓ Restored %s mode to unset (default: all)\n", selectedAgent.Name)
			}
		} else {
			restoreErr := agents.UpdateAgentMode(selectedAgent.Path, selectedAgent.Name, previousMode)
			if restoreErr != nil {
				log.Printf("Failed to undo agent %s: %v", selectedAgent.Name, restoreErr)
			} else {
				fmt.Printf("✓ Restored %s mode to '%s'\n", selectedAgent.Name, previousMode)
			}
		}
	}

	return promptContinue()
}

func promptContinue() (bool, error) {
	continueChoice, err := cli.PromptContinueOrExit()
	if err != nil {
		return false, err
	}
	return continueChoice, nil
}

func promptTemplateContinue() (bool, error) {
	continueChoice, err := cli.PromptTemplateContinueOrExit()
	if err != nil {
		return false, err
	}
	return continueChoice, nil
}

func handleTemplates(agentList []models.Agent) (bool, error) {
	for {
		choice, err := cli.PromptTemplateMenu()
		if err != nil {
			return false, err
		}

		if choice == cli.BackChoice {
			return true, nil
		}

		if choice == cli.TemplateSave {
			continueToMenu, saveErr := handleTemplateSave(agentList)
			if saveErr != nil {
				return false, saveErr
			}
			if !continueToMenu {
				return false, nil
			}
			continue
		}

		if choice == cli.TemplateShow {
			continueToMenu, showErr := handleTemplateShow(agentList)
			if showErr != nil {
				return false, showErr
			}
			if !continueToMenu {
				return false, nil
			}
			continue
		}
	}
}

func handleTemplateSave(agentList []models.Agent) (bool, error) {
	name, err := cli.PromptTemplateName()
	if err != nil {
		return false, err
	}

	if err := templates.ValidateTemplateName(name); err != nil {
		fmt.Printf("Invalid template name: %v\n", err)
		return true, nil
	}

	exists, err := templates.TemplateExists(name)
	if err != nil {
		return false, err
	}

	if exists {
		overwrite, confirmErr := cli.PromptTemplateOverwrite(name)
		if confirmErr != nil {
			return false, confirmErr
		}
		if !overwrite {
			return true, nil
		}
	}

	if err := templates.SaveTemplate(name, agentList); err != nil {
		fmt.Printf("Failed to save template: %v\n", err)
		return true, nil
	}

	fmt.Printf("\n✓ Template '%s' saved with %d agents\n", name, len(agentList))
	return promptTemplateContinue()
}

func handleTemplateShow(agentList []models.Agent) (bool, error) {
	templateList, err := templates.LoadTemplates()
	if err != nil {
		return false, err
	}

	if len(templateList) == 0 {
		fmt.Println("\nThere are no templates saved")
		return promptTemplateContinue()
	}

	for {
		selectedName, err := cli.PromptTemplateSelection(templateList)
		if err != nil {
			return false, err
		}

		if selectedName == cli.BackChoice || selectedName == "" {
			return true, nil
		}

		action, err := cli.PromptTemplateAction(selectedName)
		if err != nil {
			return false, err
		}

		if action == cli.BackChoice {
			continue
		}

		if action == cli.TemplateInspect {
			handleTemplateInspect(selectedName)
			continue
		}

		if action == cli.TemplateLoad {
			continueToMenu, loadErr := handleTemplateLoad(selectedName, agentList)
			if loadErr != nil {
				return false, loadErr
			}
			return continueToMenu, nil
		}

		if action == cli.TemplateDelete {
			continueToMenu, deleteErr := handleTemplateDelete(selectedName)
			if deleteErr != nil {
				return false, deleteErr
			}
			if !continueToMenu {
				return false, nil
			}
			templateList, err = templates.LoadTemplates()
			if err != nil {
				return false, err
			}
			continue
		}
	}
}

func handleTemplateLoad(templateName string, agentList []models.Agent) (bool, error) {
	template, err := templates.LoadTemplateByName(templateName)
	if err != nil {
		return false, err
	}

	matched, unmatched := templates.MatchAgents(template, agentList)

	fmt.Printf("\nTemplate '%s' summary:\n", templateName)
	fmt.Printf("  %d agents will be updated\n", len(matched))
	if len(unmatched) > 0 {
		fmt.Printf("  %d agents in template not matched (different source):\n", len(unmatched))
		for _, name := range unmatched {
			fmt.Printf("    - %s\n", name)
		}
	}

	if len(matched) == 0 {
		fmt.Println("\nNo agents match this template. Cannot apply.")
		return promptTemplateContinue()
	}

	confirmed, confirmErr := cli.PromptTemplateLoadConfirm(len(matched), len(unmatched))
	if confirmErr != nil {
		return false, confirmErr
	}

	if !confirmed {
		return true, nil
	}

	fmt.Printf("\nApplying template '%s'...\n", templateName)

	previousStates := make(map[string]models.AgentAssignment)
	for _, agent := range matched {
		for _, current := range agentList {
			if current.Name == agent.Name &&
				current.Source.Location == agent.Source.Location &&
				current.Source.Format == agent.Source.Format {
				previousStates[agent.Name] = models.AgentAssignment{
					Model:  current.CurrentModel,
					Mode:   current.Mode,
					Source: current.Source,
				}
				break
			}
		}
	}

	updatedAgents := []string{}
	for _, agent := range matched {
		if agent.CurrentModel != "" {
			updateErr := agents.UpdateAgentModel(agent.Path, agent.Name, agent.CurrentModel)
			if updateErr != nil {
				log.Printf("Failed to update model for %s: %v", agent.Name, updateErr)
			} else {
				fmt.Printf("✓ Updated %s model to '%s'\n", agent.Name, agent.CurrentModel)
				updatedAgents = append(updatedAgents, agent.Name)
			}
		}

		if agent.Mode != "" {
			modeErr := agents.UpdateAgentMode(agent.Path, agent.Name, agent.Mode)
			if modeErr != nil {
				log.Printf("Failed to update mode for %s: %v", agent.Name, modeErr)
			} else {
				fmt.Printf("✓ Updated %s mode to '%s'\n", agent.Name, agent.Mode)
			}
		}
	}

	if len(updatedAgents) > 0 {
		undoMessage := fmt.Sprintf("Template applied to %d agents. Undo changes?", len(updatedAgents))
		wantUndo, undoErr := cli.PromptUndo(undoMessage)
		if undoErr != nil {
			return false, undoErr
		}

		if wantUndo {
			fmt.Println("\nUndoing changes...")
			for _, agentName := range updatedAgents {
				previous := previousStates[agentName]
				var agent models.Agent
				for _, m := range matched {
					if m.Name == agentName {
						agent = m
						break
					}
				}

				if previous.Model != "" {
					restoreErr := agents.UpdateAgentModel(agent.Path, agentName, previous.Model)
					if restoreErr != nil {
						log.Printf("Failed to undo model for %s: %v", agentName, restoreErr)
					} else {
						fmt.Printf("✓ Restored %s model to '%s'\n", agentName, previous.Model)
					}
				}

				if previous.Mode != "" {
					modeErr := agents.UpdateAgentMode(agent.Path, agentName, previous.Mode)
					if modeErr != nil {
						log.Printf("Failed to undo mode for %s: %v", agentName, modeErr)
					} else {
						fmt.Printf("✓ Restored %s mode to '%s'\n", agentName, previous.Mode)
					}
				}
			}
		}
	}

	return promptTemplateContinue()
}

func handleTemplateInspect(templateName string) {
	template, err := templates.LoadTemplateByName(templateName)
	if err != nil {
		fmt.Printf("Failed to load template: %v\n", err)
		return
	}

	fmt.Printf("\n%s\n", cli.FormatTemplateInspect(template))
}

func handleTemplateDelete(templateName string) (bool, error) {
	confirmed, err := cli.PromptTemplateDeleteConfirm(templateName)
	if err != nil {
		return false, err
	}

	if !confirmed {
		return true, nil
	}

	if err := templates.DeleteTemplate(templateName); err != nil {
		fmt.Printf("Failed to delete template: %v\n", err)
		return true, nil
	}

	fmt.Printf("\n✓ Template '%s' deleted\n", templateName)
	return promptTemplateContinue()
}
