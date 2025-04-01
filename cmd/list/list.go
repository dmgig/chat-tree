package list

import (
	"fmt"
	"os"

	"chat-tree/internal/walker"
)

func ListFiles(args []string) {
	if len(args) == 0 {
		fmt.Println("Usage: chat-tree list <path> [--exclude pattern ...]")
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

	for _, file := range files {
		fmt.Println(file.Path)
	}
}
