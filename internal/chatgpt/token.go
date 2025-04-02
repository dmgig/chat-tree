package chatgpt

import (
	"fmt"
	"os"

	"chat-tree/internal/config"

	tiktoken "github.com/pkoukk/tiktoken-go"
)

func MaxTokensForModel(model string) (int, error) {
	if max, ok := config.ModelMaxTokens[model]; ok {
		return max, nil
	}
	return 0, fmt.Errorf("unknown max tokens for model: %s", model)
}

// CountTokens estimates the number of tokens in a string for a given model.
func CountTokens(model string, input string) (int, error) {
	tokenizer, err := tiktoken.EncodingForModel(model)
	if err != nil || tokenizer == nil {
		// Fallback to cl100k_base
		fmt.Fprintf(os.Stderr, "Warning: falling back to cl100k_base encoding for model %s\n", model)
		tokenizer, err = tiktoken.GetEncoding("cl100k_base")
		if err != nil {
			return 0, fmt.Errorf("failed to get fallback tokenizer: %w", err)
		}
	}

	tokens := tokenizer.Encode(input, nil, nil)
	return len(tokens), nil
}

// SafeTokenCount wraps CountTokens and exits if it encounters an error.
func SafeTokenCount(model string, input string) int {
	total, err := CountTokens(model, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Token count error: %v\n", err)
		return 0
	}
	return total
}
