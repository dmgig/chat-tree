package cmd

import (
	"fmt"
    "bufio"
	"github.com/gdamore/tcell/v2"
    "github.com/rivo/tview"
    "github.com/spf13/cobra"
    "os"
    "path/filepath"
    "strings"
)

var excludePatterns []string // Stores ignore patterns

var selectedFiles = make(map[string]bool) // Track selected file paths

var treeCmd = &cobra.Command{
    Use:   "tree [directory]",
    Short: "Displays a file tree with checkboxes",
    Args:  cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        directory := args[0]
        loadIgnorePatterns(directory) // Load ignore rules
        displayFileTreeWithCheckboxes(directory)
    },
}

func init() {
    rootCmd.AddCommand(treeCmd)
}

// Load `.chat-tree-ignore` file if it exists
func loadIgnorePatterns(rootPath string) {
    ignoreFile := filepath.Join(rootPath, ".exclude")
    excludePatterns = []string{} // Reset ignore list

    file, err := os.Open(ignoreFile)
    if err != nil {
        return // No ignore file found, do nothing
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())
        if line != "" && !strings.HasPrefix(line, "#") { // Ignore empty lines & comments
            excludePatterns = append(excludePatterns, line)
        }
    }
}

func displayFileTreeWithCheckboxes(rootPath string) {
    app := tview.NewApplication()
    treeView := tview.NewTreeView()

    rootNode := tview.NewTreeNode("[ ] " + rootPath).SetReference(rootPath).SetExpanded(true)
    treeView.SetRoot(rootNode).SetCurrentNode(rootNode)

    populateTreeWithCheckboxes(rootNode, rootPath)

    // Add a special "Concatenate & Exit" option at the bottom
    concatNode := tview.NewTreeNode("[▶] Concatenate & Exit").SetReference("concat")
    rootNode.AddChild(concatNode)

    // Handle checkbox toggling & concatenation selection
    treeView.SetSelectedFunc(func(node *tview.TreeNode) {
        text := node.GetText()
        filePath, _ := node.GetReference().(string)

        if filePath == "concat" { // User selected "Concatenate & Exit"
            app.Stop() // Exit UI, return to terminal
            concatenateFiles() // Run concatenation
            return
        }

        // Toggle checkboxes for files
        if strings.HasPrefix(text, "[ ] ") {
            node.SetText("[✔] " + strings.TrimPrefix(text, "[ ] "))
            selectedFiles[filePath] = true
        } else if strings.HasPrefix(text, "[✔] ") {
            node.SetText("[ ] " + strings.TrimPrefix(text, "[✔] "))
            delete(selectedFiles, filePath)
        }

        saveSelectedFiles() // Save selection after toggle
    })

    // Handle expand/collapse with arrow keys
    treeView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
        node := treeView.GetCurrentNode()
        if node == nil {
            return event
        }

        switch event.Key() {
        case tcell.KeyRight: // Expand folder
            if len(node.GetChildren()) > 0 {
                node.SetExpanded(true)
            }
        case tcell.KeyLeft: // Collapse folder
            if len(node.GetChildren()) > 0 {
                node.SetExpanded(false)
            }
        }
        return event
    })

    app.SetRoot(treeView, true)

    if err := app.Run(); err != nil {
        panic(err)
    }
}

func saveSelectedFiles() {
    file, err := os.Create("selected_files.txt")
    if err != nil {
        fmt.Println("Error saving selected files:", err)
        return
    }
    defer file.Close()

    for filePath := range selectedFiles {
        file.WriteString(filePath + "\n")
    }
}

func loadSelectedFiles() []string {
    var selectedFiles []string

    file, err := os.Open("selected_files.txt")
    if err != nil {
        return selectedFiles // No file = no selections
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        selectedFiles = append(selectedFiles, scanner.Text())
    }

    return selectedFiles
}

// Recursively populates the tree
func populateTreeWithCheckboxes(parentNode *tview.TreeNode, path string) {
    files, err := os.ReadDir(path)
    if err != nil {
        return
    }

    for _, file := range files {
        fullPath := filepath.Join(path, file.Name())

        if isExcluded(file.Name()) {
            continue // Skip ignored files
        }

        childNode := tview.NewTreeNode("[ ] " + file.Name()).SetReference(fullPath)

        if file.IsDir() {
            childNode.SetExpanded(false) // Expandable folders
            populateTreeWithCheckboxes(childNode, fullPath)
        }

        parentNode.AddChild(childNode)
    }
}

// Check if a file matches the ignore patterns
func isExcluded(filename string) bool {
    for _, pattern := range excludePatterns {
        match, _ := filepath.Match(pattern, filename)
        if match {
            return true
        }
    }
    return false
}
