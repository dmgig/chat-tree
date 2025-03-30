package chatgpt

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	openai "github.com/sashabaranov/go-openai"
)

type ChatClient struct {
	client *openai.Client
}

func NewChatClient(apiKey string) *ChatClient {
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

// Temporary test runner
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

	prompt := args[1]
	client := NewChatClient(apiKey)
	resp, err := client.SendPrompt(prompt)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	fmt.Println("Response:\n", resp)
}
