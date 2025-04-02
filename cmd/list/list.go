package list

import (
	"fmt"
	"os"
	"sort"

	"chat-tree/internal/chatgpt"
	"chat-tree/internal/config"
	"chat-tree/internal/walker"
)

func ListFiles(args []string) {
	if len(args) == 0 {
		fmt.Println("Usage: chat-tree list <path> [--exclude pattern ...]")
		os.Exit(1)
	}

	var paths []string
	excludePatterns := []string{}
	for i := 0; i < len(args); i++ {
		if args[i] == "--exclude" && i+1 < len(args) {
			i++
			excludePatterns = append(excludePatterns, args[i])
		} else {
			paths = append(paths, args[i])
		}
	}

	files, err := walker.WalkFiles(paths, excludePatterns)
	if err != nil {
		fmt.Println("Error walking files:", err)
		os.Exit(1)
	}

	for _, file := range files {
		fmt.Println(file.Path)
	}
}

func ListModels() {
	client := chatgpt.NewChatClient()

	models, err := client.GetAvailableModels()
	if err != nil {
		fmt.Printf("Error retrieving models: %v\n", err)
		os.Exit(1)
	}

	// Sort models by ID
	sort.Slice(models, func(i, j int) bool {
		return models[i].ID < models[j].ID
	})

	unknownSeen := false
	for _, model := range models {
		if max, ok := config.ModelMaxTokens[model.ID]; ok {
			fmt.Printf("Model ID: %s (Max tokens: %d)\n", model.ID, max)
		} else {
			fmt.Printf("Model ID: %s (Max tokens: unknown)\n", model.ID)
			unknownSeen = true
		}
	}

	if unknownSeen {
		fmt.Println("\nSome models have unknown token limits. Check https://platform.openai.com/docs/models for details.")
	}
}
