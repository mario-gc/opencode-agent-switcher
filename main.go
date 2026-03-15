package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"opencode-agent-switcher/agents"
	"opencode-agent-switcher/cli"
	"opencode-agent-switcher/config"
	"opencode-agent-switcher/models"
)

func main() {
	cfg, err := config.LoadOpencodeConfig()
	if err != nil {
		log.Fatalf("Failed to load opencode config: %v", err)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Failed to get user home dir: %v", err)
	}
	agentsDir := filepath.Join(home, ".config", "opencode", "agents")

	agentList, err := agents.LoadAgents(agentsDir)
	if err != nil {
		log.Fatalf("Failed to load agents: %v", err)
	}

	if len(agentList) == 0 {
		log.Fatalf("No agents found in %s", agentsDir)
	}

	modelOptions, err := config.GetModelsFromCLI()
	if err != nil {
		fmt.Println("Warning: Failed to get models from CLI, falling back to config file...")
		modelOptions = config.GetAvailableModels(cfg)
	}

	if len(modelOptions) == 0 {
		log.Fatalf("No models found available")
	}

	for {
		shouldContinue, err := runAgentUpdate(agentList, modelOptions)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}

		if !shouldContinue {
			break
		}

		agentList, err = agents.LoadAgents(agentsDir)
		if err != nil {
			log.Fatalf("Failed to reload agents: %v", err)
		}
	}

	fmt.Println("\nGoodbye!")
}

func runAgentUpdate(agentList []models.Agent, modelOptions []models.ModelOption) (bool, error) {
	agentIndex, err := cli.PromptAgentSelection(agentList)
	if err != nil {
		return false, err
	}

	if agentIndex == -2 {
		return false, nil
	}

	selectedAgent := agentList[agentIndex]

	modelIndex, err := cli.PromptModelSelection(modelOptions)
	if err != nil {
		return false, err
	}
	selectedModel := modelOptions[modelIndex]

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
		confirmed, err := cli.PromptConfirm(message)
		if err != nil {
			return false, err
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
		err := agents.UpdateAgentModel(agent.Path, selectedModel.ID)
		if err != nil {
			log.Printf("Failed to update agent %s: %v", agent.Name, err)
		} else {
			fmt.Printf("✓ Updated %s\n", agent.Name)
			updatedAgents = append(updatedAgents, agent.Name)
		}
	}

	if len(updatedAgents) > 0 {
		undoMessage := fmt.Sprintf("Updated %d agent(s). Undo changes?", len(updatedAgents))
		wantUndo, err := cli.PromptUndo(undoMessage)
		if err != nil {
			return false, err
		}

		if wantUndo {
			fmt.Println("\nUndoing changes...")
			for _, agentName := range updatedAgents {
				for _, agent := range agentsToUpdate {
					if agent.Name == agentName {
						previousModel := previousModels[agentName]
						err := agents.UpdateAgentModel(agent.Path, previousModel)
						if err != nil {
							log.Printf("Failed to undo agent %s: %v", agentName, err)
						} else {
							fmt.Printf("✓ Restored %s to %s\n", agentName, previousModel)
						}
						break
					}
				}
			}
		}
	}

	continueChoice, err := cli.PromptContinueOrExit()
	if err != nil {
		return false, err
	}

	return continueChoice, nil
}
