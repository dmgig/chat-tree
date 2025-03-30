package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"chat-tree/internal/session"
)

// isBinaryFile returns true if the file contains non-text (binary) data.
func isBinaryFile(path string) bool {
	f, err := os.Open(path)
	if err != nil {
		return true // assume binary if we can't open
	}
	defer f.Close()

	buf := make([]byte, 8000)
	n, err := f.Read(buf)
	if err != nil && err != io.EOF {
		return true
	}

	// Check for null bytes or invalid UTF-8
	if !utf8.Valid(buf[:n]) {
		return true
	}
	return false
}

// loadExcludeFile loads patterns from a .exclude file in the given directory.
func loadExcludeFile(dir string) []string {
	patterns := []string{}
	excludePath := filepath.Join(dir, ".exclude")
	file, err := os.Open(excludePath)
	if err != nil {
		return patterns // no file = no extra patterns
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		patterns = append(patterns, line)
	}
	return patterns
}

// walkFiles returns a list of all non-binary files under the given paths, excluding matching patterns.
func walkFiles(paths []string, cliExcludePatterns []string) ([]session.FileInfo, error) {
	var results []session.FileInfo

	for _, basePath := range paths {
		allExcludePatterns := append([]string{}, cliExcludePatterns...)
		allExcludePatterns = append(allExcludePatterns, loadExcludeFile(basePath)...)

		err := filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil // skip unreadable files
			}
			if info.IsDir() {
				return nil // skip dirs
			}

			relPath, _ := filepath.Rel(basePath, path)

			for _, pattern := range allExcludePatterns {
				// Try matching against base name
				if match, _ := filepath.Match(pattern, filepath.Base(relPath)); match {
					return nil
				}
				// Try matching against relative path
				if match, _ := filepath.Match(pattern, relPath); match {
					return nil
				}
			}

			if isBinaryFile(path) {
				return nil // skip binaries
			}

			results = append(results, session.FileInfo{Path: path})
			return nil
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error walking %s: %v\n", basePath, err)
		}
	}
	return results, nil
}

// Execute is the main CLI entry point
func Execute(args []string) {
	if len(args) == 0 {
		fmt.Println("Usage: chat-tree <path> [<path> ...] --exclude <pattern>")
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

	files, err := walkFiles(paths, excludePatterns)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	if len(files) == 0 {
		fmt.Println("No files found after exclusions.")
		return
	}

	err = session.Save(files)
	if err != nil {
		fmt.Println("Failed to write session:", err)
		os.Exit(1)
	}
}
