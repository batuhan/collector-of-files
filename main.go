package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	gitignore "github.com/sabhiram/go-gitignore"
)

func main() {
	// Define command-line flags
	startPath := flag.String("path", "", "Start path for searching files")
	includeExtensions := flag.String("includeExtensions", "", "Comma-separated list of file extensions to include")
	outputFile := flag.String("output", "combined_code.md", "Output file name")
	excludeDirs := flag.String("excludeDirs", "", "Comma-separated list of additional directories to exclude")
	flag.Parse()

	// Check if required flags are provided
	if *startPath == "" || *includeExtensions == "" {
		fmt.Println("Usage: combine-code --path <start_path> --includeExtensions <comma_separated_extensions> [--output <output_file>] [--excludeDirs <comma_separated_dirs>]")
		fmt.Println("Example: combine-code --path /path/to/start --includeExtensions go,js --excludeDirs node_modules,vendor")
		os.Exit(1)
	}

	// Convert comma-separated strings to slices
	extensions := strings.Split(*includeExtensions, ",")
	excludedDirs := strings.Split(*excludeDirs, ",")

	// Create or truncate the output file
	file, err := os.Create(*outputFile)
	if err != nil {
		fmt.Printf("Error creating output file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	// Parse .gitignore file
	gitignoreObj, err := loadGitignore(*startPath)
	if err != nil {
		fmt.Printf("Error loading .gitignore: %v\n", err)
		// Continue without .gitignore rules if there's an error
	}

	// Walk through the directory
	err = filepath.Walk(*startPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Get the relative path
		relPath, err := filepath.Rel(*startPath, path)
		if err != nil {
			return err
		}

		// Check if the file/directory should be ignored based on .gitignore
		if gitignoreObj != nil && gitignoreObj.MatchesPath(relPath) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// Check if the directory should be excluded based on the excludeDirs flag
		if info.IsDir() {
			if containsPath(excludedDirs, relPath) {
				return filepath.SkipDir
			}
			return nil
		}

		// Check if the file has one of the specified extensions
		ext := strings.TrimPrefix(filepath.Ext(path), ".")
		if !contains(extensions, ext) {
			return nil
		}

		// Read file contents
		content, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		// Write file header
		_, err = fmt.Fprintf(file, "===== FILE START: %s =====\n", relPath)
		if err != nil {
			return err
		}

		// Write file contents with line numbers
		lines := strings.Split(string(content), "\n")
		for i, line := range lines {
			_, err = fmt.Fprintf(file, "%4d | %s\n", i+1, line)
			if err != nil {
				return err
			}
		}

		// Write file footer
		_, err = fmt.Fprintf(file, "===== FILE END: %s =====\n\n", relPath)
		return err
	})

	if err != nil {
		fmt.Printf("Error walking the path %v: %v\n", *startPath, err)
		os.Exit(1)
	}

	fmt.Printf("Combined code has been written to %s\n", *outputFile)
}

// Helper function to check if a slice contains a string
func contains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}

// Helper function to check if a path should be excluded
func containsPath(excludedDirs []string, path string) bool {
	for _, dir := range excludedDirs {
		if dir == path || strings.HasPrefix(path, dir+string(os.PathSeparator)) {
			return true
		}
	}
	return false
}

// Helper function to load .gitignore file
func loadGitignore(startPath string) (*gitignore.GitIgnore, error) {
	gitignorePath := filepath.Join(startPath, ".gitignore")
	if _, err := os.Stat(gitignorePath); os.IsNotExist(err) {
		return nil, nil // No .gitignore file found
	}

	return gitignore.CompileIgnoreFile(gitignorePath)
}
