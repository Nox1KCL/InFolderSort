package files

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Nox1KCL/InFolderSort/internal/config"
)

type SortResult struct {
	Moved   []string
	Skipped []string
	Errors  []error
}

func GetHomeDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to find home directory: %w", err)
	}
	return homeDir, nil
}

func GetDownloadsPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to find home directory: %w", err)
	}

	potentialDirs := []string{"Downloads", "downloads"}

	for _, dirName := range potentialDirs {
		downloadPath := filepath.Join(homeDir, dirName)
		if info, err := os.Stat(downloadPath); err == nil && info.IsDir() {
			return downloadPath, nil
		}
	}

	return "", fmt.Errorf("downloads directory not found")
}

// InDirSorting sorts all files in the target directory according to the provided config
func InDirSorting(targetPath string, cfg *config.Config) (SortResult, error) {
	var report SortResult

	entries, err := os.ReadDir(targetPath)
	if err != nil {
		return report, fmt.Errorf("reading directory %q: %w", targetPath, err)
	}

	for _, entry := range entries {
		fileName := entry.Name()
		if !entry.IsDir() && !strings.HasPrefix(fileName, ".") {
			fileExt := filepath.Ext(fileName)
			targetFolder, err := cfg.TargetFolderName(fileExt) // Get a name of save folder
			if err != nil {
				report.Skipped = append(report.Skipped, entry.Name())
				continue
			}

			finalPath := filepath.Join(targetPath, targetFolder) // Make a save path
			// Make the directories by path
			if err := os.MkdirAll(finalPath, 0755); err != nil {
				report.Errors = append(report.Errors,
					fmt.Errorf("creating dirs by path %q: %w", finalPath, err),
				)
				continue
			}

			oldFilePath := filepath.Join(targetPath, fileName)
			newFilePath := filepath.Join(finalPath, fileName)

			if err := os.Rename(oldFilePath, newFilePath); err != nil {
				report.Errors = append(report.Errors,
					fmt.Errorf("moving file from %q to %q: %w", oldFilePath, newFilePath, err),
				)
				continue
			}
			report.Moved = append(report.Moved, entry.Name())
		}
	}

	return report, nil
}
