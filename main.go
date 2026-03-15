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

var (
	version = "dev"
	commit  = "none"
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

	agentIndex, err := cli.PromptAgentSelection(agentList)
	if err != nil {
		log.Fatalf("Agent selection failed: %v", err)
	}
	selectedAgent := agentList[agentIndex]

	modelIndex, err := cli.PromptModelSelection(modelOptions)
	if err != nil {
		log.Fatalf("Model selection failed: %v", err)
	}
	selectedModel := modelOptions[modelIndex]

	agentsToUpdate := []models.Agent{selectedAgent}

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
			log.Fatalf("Confirmation failed: %v", err)
		}
		if confirmed {
			agentsToUpdate = append(agentsToUpdate, otherAgents...)
		}
	}

	fmt.Printf("\nUpdating %d agent(s) to model '%s'...\n", len(agentsToUpdate), selectedModel.ID)

	for _, agent := range agentsToUpdate {
		err := agents.UpdateAgentModel(agent.Path, selectedModel.ID)
		if err != nil {
			log.Printf("Failed to update agent %s: %v", agent.Name, err)
		} else {
			fmt.Printf("✓ Updated %s\n", agent.Name)
		}
	}

	fmt.Println("\nDone!")
}
