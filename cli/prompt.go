package cli

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"

	"agent-switcher/models"
)

// PromptAgentSelection presents interactive agent selection
func PromptAgentSelection(reader *bufio.Reader, agents []models.Agent) (int, error) {
	fmt.Println("\nAvailable Agents:")
	fmt.Println("-----------------")
	for i, agent := range agents {
		fmt.Printf("%d. %s (Current: %s)\n", i+1, agent.Name, agent.CurrentModel)
	}
	fmt.Println("-----------------")

	fmt.Print("Select an agent to update (enter number): ")
	input, err := reader.ReadString('\n')
	if err != nil {
		return -1, err
	}

	input = strings.TrimSpace(input)
	selection, err := strconv.Atoi(input)
	if err != nil || selection < 1 || selection > len(agents) {
		return -1, fmt.Errorf("invalid selection")
	}

	return selection - 1, nil
}

// PromptModelSelection presents interactive model selection
func PromptModelSelection(reader *bufio.Reader, models []models.ModelOption) (int, error) {
	fmt.Println("\nAvailable Models:")
	fmt.Println("-----------------")
	for i, model := range models {
		fmt.Printf("%d. %s\n", i+1, model.Display)
	}
	fmt.Println("-----------------")

	fmt.Print("Select a new model (enter number): ")
	input, err := reader.ReadString('\n')
	if err != nil {
		return -1, err
	}

	input = strings.TrimSpace(input)
	selection, err := strconv.Atoi(input)
	if err != nil || selection < 1 || selection > len(models) {
		return -1, fmt.Errorf("invalid selection")
	}

	return selection - 1, nil
}
