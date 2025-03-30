/*
Copyright © 2025 Dave M. Giglio dave.m.giglio@gmail.com
*/
package main

import (
	"os"
	"fmt"
	"chat-tree/cmd"
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
		cmd.Execute(os.Args[2:])
	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}
