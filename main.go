package main

import (
	"fmt"
	"os"

	"chat-tree/cmd/document"
	"chat-tree/cmd/list"
	"chat-tree/internal/chatgpt"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: chat-tree <command> [options]")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "openai":
		chatgpt.RunFromCLI(os.Args[2:])
	case "document":
		document.CreateDocumentation(os.Args[2:])
	case "list":
		list.ListFiles(os.Args[2:])
	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}
