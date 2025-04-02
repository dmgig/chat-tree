package config

// ModelMaxTokens maps known model IDs to their maximum token limits.
// Values are sourced from official OpenAI documentation.
var ModelMaxTokens = map[string]int{
	"gpt-4":              8192,
	"gpt-4-32k":          32768,
	"gpt-4-turbo":        128000,
	"gpt-3.5-turbo":      4096,
	"gpt-3.5-turbo-16k":  16384,
	"gpt-4-0125-preview": 128000,
	"gpt-3.5-turbo-1106": 16384,
	"text-davinci-003":   4097,
	"text-curie-001":     2048,
	"text-babbage-001":   2048,
	"text-ada-001":       2048,
}
