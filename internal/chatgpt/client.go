package chatgpt

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/joho/godotenv"
	openai "github.com/sashabaranov/go-openai"
)

type ChatClient struct {
	client *openai.Client
}

func NewChatClient() *ChatClient {
	_ = godotenv.Load()
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		fmt.Println("Missing OPENAI_API_KEY environment variable")
		os.Exit(1)
	}
	client := openai.NewClient(apiKey)
	return &ChatClient{client: client}
}

func (c *ChatClient) SendPrompt(prompt string) (string, error) {
	resp, err := c.client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleUser, Content: prompt},
		},
	})
	if err != nil {
		return "", err
	}
	if len(resp.Choices) > 0 {
		return resp.Choices[0].Message.Content, nil
	}
	return "", fmt.Errorf("no response from OpenAI")
}

// ProcessSession reads prompts from a session directory and writes OpenAI responses
func (c *ChatClient) ProcessSession(sessionDir string) error {
	promptFiles, err := filepath.Glob(filepath.Join(sessionDir, "*_prompt.txt"))
	if err != nil {
		return err
	}
	if len(promptFiles) == 0 {
		return fmt.Errorf("no prompt files found in %s", sessionDir)
	}

	sort.Strings(promptFiles) // ensure deterministic order

	var previousContext string
	for _, promptPath := range promptFiles {
		content, err := ioutil.ReadFile(promptPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to read %s: %v\n", promptPath, err)
			continue
		}

		prompt := strings.TrimSpace(string(content))
		if previousContext != "" {
			prompt = previousContext + "\n\n" + prompt
		}

		fmt.Printf("Sending prompt from %s\n", promptPath)
		fmt.Printf("Prompt:\n%s\n\n", prompt) // Add this to see actual prompt content

		response, err := c.SendPrompt(prompt)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to get response for %s: %v\n", promptPath, err)
			continue
		}

		fmt.Printf("Response received for %s\n", promptPath)

		base := filepath.Base(promptPath)
		num := strings.Split(base, "_")[0] // e.g., "00001"
		responsePath := filepath.Join(sessionDir, fmt.Sprintf("%s_response.txt", num))
		ioutil.WriteFile(responsePath, []byte(response), 0644)

		previousContext += "\n\n" + prompt + "\n" + response

		time.Sleep(2 * time.Second) // Delay between requests to avoid rate limiting
	}

	return nil
}

func RunFromCLI(args []string) {
	if len(args) < 2 || args[0] != "--prompt" {
		fmt.Println("Usage: chat-tree openai --prompt \"your message here\"")
		os.Exit(1)
	}

	_ = godotenv.Load() // loads .env if it exists

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		fmt.Println("Missing OPENAI_API_KEY environment variable")
		os.Exit(1)
	}

	model := os.Getenv("OPENAI_MODEL")
	if model == "" {
		model = "gpt-3.5-turbo" // fallback if unset
	}

	prompt := args[1]

	// Token counting
	tokenCount, err := CountTokens(model, prompt)
	if err != nil {
		fmt.Printf("Failed to count tokens: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Using model: %s\n", model)
	fmt.Printf("Token count: %d\n", tokenCount)

	client := NewChatClient()
	resp, err := client.SendPrompt(prompt)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	fmt.Println("Response:\n", resp)
}

func (c *ChatClient) GetAvailableModels() ([]openai.Model, error) {
	resp, err := c.client.ListModels(context.Background())
	if err != nil {
		return nil, err
	}
	return resp.Models, nil
}
