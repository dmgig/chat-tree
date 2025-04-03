package session

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"chat-tree/internal/chatgpt"
)

// ReviewSession uses the full generated documentation and file tree to perform a review pass
func ReviewSession(sessionDir string) error {
	client := chatgpt.NewChatClient()

	docPath := filepath.Join(sessionDir, "_documentation.md")
	content, err := os.ReadFile(docPath)
	if err != nil {
		return fmt.Errorf("failed to read documentation file: %w", err)
	}

	// Create file tree summary
	var fileList []string
	err = filepath.Walk(sessionDir, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".go") || strings.HasSuffix(path, ".mod") || strings.HasSuffix(path, ".sum") {
			rel, _ := filepath.Rel(sessionDir, path)
			fileList = append(fileList, rel)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to build file tree: %w", err)
	}
	fileTree := strings.Join(fileList, "\n")

	// Build review prompt
	prompt := fmt.Sprintf(`You are reviewing the documentation generated for this project.

Here is the file structure of the reviewed files:
%s

And here is the full documentation to be reviewed:

%s

Please revise it if needed. If it looks good, you can return it unchanged.`, fileTree, string(content))

	// Send it to GPT
	fmt.Println("🔍 Sending review prompt to GPT...")
	response, err := client.SendPrompt(prompt)
	if err != nil {
		return fmt.Errorf("review prompt failed: %w", err)
	}

	reviewPath := filepath.Join(sessionDir, "_review_documentation.md")
	err = os.WriteFile(reviewPath, []byte(response), 0644)
	if err != nil {
		return fmt.Errorf("failed to save review output: %w", err)
	}

	fmt.Printf("✅ Review complete: %s\n", reviewPath)
	return nil
}
