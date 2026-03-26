package main

import (
	"fmt"
	"log"

	"opencode-agent-switcher/agents"
	"opencode-agent-switcher/cli"
	"opencode-agent-switcher/config"
	"opencode-agent-switcher/models"
)

func main() {
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

	for {
		shouldContinue, loopErr := runAgentUpdate(agentList, modelOptions, &currentSort)
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

func runAgentUpdate(agentList []models.Agent, modelOptions []models.ModelOption, currentSort *string) (bool, error) {
	sortedAgents := cli.SortAgents(agentList, *currentSort)

	selectedName, err := cli.PromptAgentSelection(sortedAgents, *currentSort)
	if err != nil {
		return false, err
	}

	if selectedName == cli.ExitChoice {
		return false, nil
	}

	if selectedName == cli.SortChoice {
		newSort, sortErr := cli.PromptSortSelection(*currentSort)
		if sortErr != nil {
			return false, sortErr
		}
		*currentSort = newSort
		return true, nil
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
			continueToMenu, err := handleModelChange(selectedAgent, agentList, modelOptions, *currentSort)
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

func handleModelChange(selectedAgent models.Agent, agentList []models.Agent, modelOptions []models.ModelOption, currentSort string) (bool, error) {
	sortedModels := cli.SortModels(modelOptions, currentSort)
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
