package tui

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Nox1KCL/InFolderSort/internal/config"
	"github.com/Nox1KCL/InFolderSort/internal/files"
)

func Core(cfg *config.Config) error {
	userChoice := askChoice("basic sort or manual?(b/m): ", "b", "m")

	var targetPath string
	var err error

	switch userChoice {
	case "b":
		targetPath, err = files.GetDownloadsPath()
	case "m":
		targetPath, err = getManualPath()
	}

	if err != nil {
		return fmt.Errorf("failed to get target path: %w", err)
	}

	return performSort(targetPath, cfg)
}

func performSort(targetDir string, cfg *config.Config) error {
	fileInfo, err := os.Stat(targetDir)
	if os.IsNotExist(err) {
		return fmt.Errorf("path doesn't exist: %w", err)
	}
	if !fileInfo.IsDir() {
		return fmt.Errorf("path is not a directory: %q", fileInfo.Name())
	}

	if err := files.InDirSorting(targetDir, cfg); err != nil {
		return fmt.Errorf("directory sorting error: %w", err)
	}
	return nil
}

func getManualPath() (string, error) {
	homeDir, err := files.GetHomeDir()
	if err != nil {
		return "", fmt.Errorf("getting home directory: %w", err)
	}

	fmt.Printf("We suppose your home dir is: %s\n", homeDir)
	useHome := askChoice("Use it for base of path?(y/n): ", "y", "n")
	userInput := askInput("Enter folder's path you want to sort: ")

	if useHome == "y" {
		return filepath.Join(homeDir, userInput), nil
	}
	return userInput, nil
}

func askChoice(prompt string, validOptions ...string) string {
	var input string
	for {
		fmt.Print(prompt)
		_, _ = fmt.Scanln(&input)
		input = strings.TrimSpace(input)

		for _, opt := range validOptions {
			if input == opt {
				return input
			}
		}
		fmt.Println("Invalid input, try again.")
	}
}

func askInput(prompt string) string {
	fmt.Print(prompt)
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		return strings.TrimSpace(scanner.Text())
	}
	return ""
}
