package chatgpt

import (
	"chat-tree/internal/config"
	"context"
	"fmt"
	"os"

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
		model = config.DefaultModel // fallback if unset
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
