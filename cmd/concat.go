package cmd

import (
    "fmt"
    "os"
    "github.com/spf13/cobra"
)

var concatCmd = &cobra.Command{
    Use:   "concat",
    Short: "Concatenates selected files with their relative paths",
    Run: func(cmd *cobra.Command, args []string) {
        concatenateFiles()
    },
}

func init() {
    rootCmd.AddCommand(concatCmd)
}

func concatenateFiles() {
    selectedFiles := loadSelectedFiles() // Load selected files

    if len(selectedFiles) == 0 {
        fmt.Println("No files selected.")
        return
    }

    fmt.Println("\n--- Concatenated Output ---\n")
    for _, filePath := range selectedFiles {
        content, err := os.ReadFile(filePath)
        if err != nil {
            fmt.Printf("Error reading %s: %v\n", filePath, err)
            continue
        }

        fmt.Printf("%s %s\n", filePath, string(content))
    }
    fmt.Println("\n--- Done! ---\n")
}
