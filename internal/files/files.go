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

type MoveTask struct {
	File       string
	SourcePath string
	DestPath   string
}

type Sorter struct {
	Config    *config.Config
	TargetDir string
	Tasks     []MoveTask
	Errors    []error
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
func InDirSorting(sortPath string, cfg *config.Config) (SortResult, error) {
	var report SortResult

	entries, err := os.ReadDir(sortPath)
	if err != nil {
		return report, fmt.Errorf("reading directory %q: %w", sortPath, err)
	}

	for _, entry := range entries {
		fileName := entry.Name()
		if !entry.IsDir() && !strings.HasPrefix(fileName, ".") {
			fileExt := filepath.Ext(fileName)
			targetPath, err := cfg.GetTargetPath(fileExt)
			if err != nil {
				report.Skipped = append(report.Skipped, entry.Name())
				continue
			}

			var savePath string
			if filepath.IsAbs(targetPath) {
				savePath = targetPath
			} else {
				savePath = filepath.Join(sortPath, targetPath)
			}

			// Make the directories by path
			if err := os.MkdirAll(savePath, 0755); err != nil {
				report.Errors = append(report.Errors,
					fmt.Errorf("creating dirs by path %q: %w", savePath, err),
				)
				continue
			}

			oldFilePath := filepath.Join(sortPath, fileName)
			newFilePath := filepath.Join(savePath, fileName)

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
