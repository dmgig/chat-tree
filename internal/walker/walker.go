package walker

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"chat-tree/internal/session"

	ds "github.com/bmatcuk/doublestar/v4"
)

func isBinaryFile(path string) bool {
	f, err := os.Open(path)
	if err != nil {
		return true
	}
	defer f.Close()

	buf := make([]byte, 8000)
	n, err := f.Read(buf)
	if err != nil && err != io.EOF {
		return true
	}
	return !utf8.Valid(buf[:n])
}

func loadExcludeFile(dir string) []string {
	patterns := []string{}
	excludePath := filepath.Join(dir, ".exclude")

	file, err := os.Open(excludePath)
	if err != nil {
		return patterns
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

	fmt.Fprintln(os.Stderr, patterns)

	return patterns
}

func WalkFiles(paths []string, cliExcludePatterns []string) ([]session.FileInfo, error) {
	var results []session.FileInfo

	projectRoot := findProjectRoot()
	if projectRoot == "" {
		fmt.Fprintln(os.Stderr, "Warning: Could not find project root.")
		projectRoot = "." // fallback
	}

	// Load .exclude from project root
	allExcludePatterns := append([]string{}, cliExcludePatterns...)
	allExcludePatterns = append(allExcludePatterns, loadExcludeFile(projectRoot)...)

	for _, basePath := range paths {
		baseAbs, _ := filepath.Abs(basePath) // capture the base path for relative matching

		err := filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
			return handlePath(path, info, err, baseAbs, allExcludePatterns, &results)
		})

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error walking %s: %v\n", basePath, err)
		}
	}

	return results, nil
}

func handlePath(path string, info os.FileInfo, err error, baseAbs string, excludePatterns []string, results *[]session.FileInfo) error {
	if err != nil {
		return nil
	}

	absPath, _ := filepath.Abs(path)
	relPath, err := filepath.Rel(baseAbs, absPath)
	if err != nil {
		return nil
	}
	relPath = filepath.ToSlash(relPath)

	// Skip "." to avoid matching against itself
	if relPath == "." {
		return nil
	}

	// Directories
	if info.IsDir() {
		for _, pattern := range excludePatterns {
			match, err := ds.Match(pattern, relPath)
			if err == nil && match {
				fmt.Printf("Skipping directory: %s (matched %s)\n", relPath, pattern)
				return filepath.SkipDir
			}
		}
	}

	// Files
	for _, pattern := range excludePatterns {
		match, err := ds.Match(pattern, relPath)
		if err == nil && match {
			fmt.Printf("Excluding file: %s (matched %s)\n", relPath, pattern)
			return nil
		}
	}

	if !info.IsDir() && isBinaryFile(path) {
		return nil
	}

	if !info.IsDir() {
		*results = append(*results, session.FileInfo{Path: path})
	}

	return nil
}

func findProjectRoot() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}

	for {
		if exists(filepath.Join(dir, ".git")) || exists(filepath.Join(dir, "go.mod")) {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return ""
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
