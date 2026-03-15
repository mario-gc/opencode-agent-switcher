package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"agent-switcher/agents"
	"agent-switcher/cli"
	"agent-switcher/config"
	"agent-switcher/models"
)

func main() {
	// Create a single reader for all stdin input
	reader := bufio.NewReader(os.Stdin)

	// 1. Load Opencode Config
	cfg, err := config.LoadOpencodeConfig()
	if err != nil {
		log.Fatalf("Failed to load opencode config: %v", err)
	}

	// 2. Load Agents
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

	// 3. Get Available Models
	// Try CLI first
	modelOptions, err := config.GetModelsFromCLI()
	if err != nil {
		fmt.Println("Warning: Failed to get models from CLI, falling back to config file...")
		modelOptions = config.GetAvailableModels(cfg)
	}

	if len(modelOptions) == 0 {
		log.Fatalf("No models found available")
	}

	// 4. Prompt for Agent Selection
	agentIndex, err := cli.PromptAgentSelection(reader, agentList)
	if err != nil {
		log.Fatalf("Selection failed: %v", err)
	}
	selectedAgent := agentList[agentIndex]

	// 5. Prompt for Model Selection
	fmt.Printf("\nUpdating agent '%s' (Current: %s)\n", selectedAgent.Name, selectedAgent.CurrentModel)
	modelIndex, err := cli.PromptModelSelection(reader, modelOptions)
	if err != nil {
		log.Fatalf("Selection failed: %v", err)
	}
	selectedModel := modelOptions[modelIndex]

	// 6. Update Agent
	// Find all agents with the same model if user wants to update all?
	// For now, just update the selected agent as per original plan,
	// but the user mentioned "if several agents use the same model, we change it in all those agents"
	// Let's ask the user if they want to update all agents with this model.

	agentsToUpdate := []models.Agent{selectedAgent}

	// Check for other agents with same model
	var otherAgents []models.Agent
	for _, a := range agentList {
		if a.Name != selectedAgent.Name && a.CurrentModel == selectedAgent.CurrentModel {
			otherAgents = append(otherAgents, a)
		}
	}

	if len(otherAgents) > 0 {
		fmt.Printf("\nThe following agents also use model '%s':\n", selectedAgent.CurrentModel)
		for _, a := range otherAgents {
			fmt.Printf("- %s\n", a.Name)
		}

		fmt.Print("\nDo you want to update these agents as well? (y/n): ")
		response, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("Failed to read response: %v", err)
		}
		response = strings.TrimSpace(response)
		if response == "y" || response == "Y" {
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
