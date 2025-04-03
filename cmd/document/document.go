package document

import (
	"fmt"
	"os"

	"chat-tree/internal/session"
	"chat-tree/internal/walker"
)

// CreateDocumentation is the main entry point for the 'document' command.
func CreateDocumentation(args []string) {
	if len(args) == 0 {
		fmt.Println("Usage: chat-tree document <path> [--exclude pattern ...]")
		os.Exit(1)
	}

	var paths []string
	excludePatterns := []string{}
	for i := 0; i < len(args); i++ {
		if args[i] == "--exclude" && i+1 < len(args) {
			i++
			excludePatterns = append(excludePatterns, args[i])
		} else {
			paths = append(paths, args[i])
		}
	}

	files, err := walker.WalkFiles(paths, excludePatterns)
	if err != nil {
		fmt.Println("Error walking files:", err)
		os.Exit(1)
	}

	if len(files) == 0 {
		fmt.Println("No files found after exclusions.")
		return
	}

	if _, err := session.Save(files); err != nil {
		fmt.Println("Failed to write session:", err)
		os.Exit(1)
	}

	fmt.Println("Documentation generation complete.")
}
