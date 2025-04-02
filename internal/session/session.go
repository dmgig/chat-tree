package session

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	chatgpt "chat-tree/internal/chatgpt"
)

type FileInfo struct {
	Path string
}

// Save creates a new session directory and writes prompt/response placeholders.
// It returns the session directory path and an error, if any.
func Save(files []FileInfo) (string, error) {
	timestamp := time.Now().Format("20060102_150405")
	sessionDir := filepath.Join("output", timestamp)

	if err := os.MkdirAll(sessionDir, 0755); err != nil {
		return "", err
	}

	for i, file := range files {
		content, err := os.ReadFile(file.Path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to read %s: %v\n", file.Path, err)
			continue
		}

		prompt := fmt.Sprintf(
			"Explain the following file in detail, including what each function does and how it fits into the program:\n\n%s",
			string(content),
		)

		model := os.Getenv("OPENAI_MODEL")
		if model == "" {
			model = "gpt-3.5-turbo" // fallback
		}

		tokenCount, err := chatgpt.CountTokens(model, prompt)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error counting tokens for %s: %v\n", file.Path, err)
			tokenCount = -1
		}

		maxTokens, err := chatgpt.MaxTokensForModel(model)
		if err == nil && tokenCount > maxTokens {
			fmt.Fprintf(os.Stderr, "WARNING: %s exceeds token limit (%d > %d). Consider splitting.\n", file.Path, tokenCount, maxTokens)
		}

		promptPath := filepath.Join(sessionDir, fmt.Sprintf("%05d_prompt.txt", i))
		responsePath := filepath.Join(sessionDir, fmt.Sprintf("%05d_response.txt", i))

		_ = os.WriteFile(promptPath, []byte(prompt), 0644)
		_ = os.WriteFile(responsePath, []byte("[placeholder for response]"), 0644)
	}

	fmt.Println("Session written to:", sessionDir)
	return sessionDir, nil
}
