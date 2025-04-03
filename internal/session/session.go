package session

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	chatgpt "chat-tree/internal/chatgpt"
	"chat-tree/internal/config"
)

type FileInfo struct {
	Path string
}

// Save creates a new session directory and writes prompt/response placeholders,
// then processes the session with OpenAI, including original filenames.
func Save(files []FileInfo) (string, error) {
	timestamp := time.Now().Format("20060102_150405")
	sessionDir := filepath.Join(config.DefaultOutputDir, timestamp)

	if err := os.MkdirAll(sessionDir, 0755); err != nil {
		return "", err
	}

	filenameMap := make(map[string]string)

	for i, file := range files {
		content, err := os.ReadFile(file.Path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to read %s: %v\n", file.Path, err)
			continue
		}

		prompt := fmt.Sprintf(config.BuildPrompt+"\n\n%s", string(content))

		model := os.Getenv("OPENAI_MODEL")
		if model == "" {
			model = config.DefaultModel
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

		filenameMap[filepath.Base(promptPath)] = filepath.Base(file.Path)
	}

	err := processSessionWithOpenAI(sessionDir, filenameMap)
	if err != nil {
		return "", fmt.Errorf("failed to process session: %w", err)
	}

	fmt.Println("Session written to:", sessionDir)

	// Build file tree string
	var fileTree strings.Builder
	fileTree.WriteString("# File Tree\n\n")
	for _, file := range files {
		relPath := file.Path
		fileTree.WriteString("- " + relPath + "\n")
	}

	// Save file tree to session directory
	treePath := filepath.Join(sessionDir, "_file-tree.txt")
	if err := os.WriteFile(treePath, []byte(fileTree.String()), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write file tree: %v\n", err)
	}

	return sessionDir, nil
}

// createPrompt constructs a full prompt string from pieces.
func createPrompt(instruction, prevDoc, currentFile string) string {
	var sb strings.Builder
	sb.WriteString(instruction)
	if prevDoc != "" {
		sb.WriteString("\n\n# Previous Documentation\n")
		sb.WriteString(prevDoc)
	}
	sb.WriteString("\n\n# Current File\n```")
	sb.WriteString(currentFile)
	sb.WriteString("\n```")
	return sb.String()
}

// processSessionWithOpenAI reads prompts, sends them to OpenAI, and writes responses.
func processSessionWithOpenAI(sessionDir string, filenameMap map[string]string) error {
	client := chatgpt.NewChatClient()

	promptFiles, err := filepath.Glob(filepath.Join(sessionDir, "*_prompt.txt"))
	if err != nil {
		return err
	}
	if len(promptFiles) == 0 {
		return fmt.Errorf("no prompt files found in %s", sessionDir)
	}

	sort.Strings(promptFiles)

	fmt.Printf("# DOCSTART %s\n\n", time.Now().Format("20060102_150405"))

	previousResponse := ""
	for i, promptPath := range promptFiles {
		content, err := ioutil.ReadFile(promptPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to read %s: %v\n", promptPath, err)
			continue
		}

		originalFilename := filenameMap[filepath.Base(promptPath)]
		currentFile := strings.TrimSpace(string(content))

		prompt := createPrompt(config.BuildPrompt, previousResponse, currentFile)

		fmt.Printf("## PROMPT #%05d\n\nPASSING %s IN FOR REVIEW\n\n%s\n\n", i, originalFilename, prompt)

		response, err := client.SendPrompt(prompt)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to get response for %s: %v\n", promptPath, err)
			continue
		}

		fmt.Printf("RESPONSE #%05d\n\n%s\n\n", i, response)

		base := filepath.Base(promptPath)
		num := strings.Split(base, "_")[0]
		responsePath := filepath.Join(sessionDir, fmt.Sprintf("%s_response.txt", num))
		ioutil.WriteFile(responsePath, []byte(response), 0644)

		previousResponse = response

		time.Sleep(2 * time.Second)
	}

	fmt.Println("DOCEND===================")

	// Write final documentation to a single markdown file
	finalDocPath := filepath.Join(sessionDir, "_documentation.md")
	err = os.WriteFile(finalDocPath, []byte(previousResponse), 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write final documentation file: %v\n", err)
	} else {
		fmt.Printf("✅ Final documentation saved to: %s\n", finalDocPath)
	}

	_ = reviewFinalDocumentation(sessionDir, previousResponse, getSortedValues(filenameMap))

	return nil
}

func reviewFinalDocumentation(sessionDir, finalDoc string, filenames []string) error {
	client := chatgpt.NewChatClient()

	// Build file tree text
	var fileList strings.Builder
	for _, name := range filenames {
		fileList.WriteString("- " + name + "\n")
	}

	reviewPrompt := fmt.Sprintf(
		"%s\n\n# Documentation to Review\n\n%s\n\n# Files Reviewed\n\n%s",
		config.ReviewPrompt,
		finalDoc,
		fileList.String(),
	)

	fmt.Println("## REVIEW PROMPT\n\n", reviewPrompt)

	reviewResponse, err := client.SendPrompt(reviewPrompt)
	if err != nil {
		return fmt.Errorf("failed to get review response: %w", err)
	}

	reviewPath := filepath.Join(sessionDir, "_review_documentation.md")
	err = os.WriteFile(reviewPath, []byte(reviewResponse), 0644)
	if err != nil {
		return fmt.Errorf("failed to write review file: %w", err)
	}

	fmt.Printf("✅ Review documentation saved to: %s\n", reviewPath)
	return nil
}

func getSortedValues(m map[string]string) []string {
	var values []string
	for _, v := range m {
		values = append(values, v)
	}
	sort.Strings(values)
	return values
}
