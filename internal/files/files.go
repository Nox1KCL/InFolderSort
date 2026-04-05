package files

import (
	"fmt"
	"main/internal/config"
	"os"
	"path/filepath"
	"strings"
)

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

	return "", fmt.Errorf("downloads directory didn't found")
}

// InDirSorting Func is sorting all files in a dir with base config
func InDirSorting(targetPath string, cfg *config.Config) error {
	entries, _ := os.ReadDir(targetPath) // Scan all dir

	for _, entry := range entries {
		fileName := entry.Name()
		if !entry.IsDir() && !strings.HasPrefix(fileName, ".") {
			fileExt := filepath.Ext(fileName)
			targetFolder, err := TargetFolderName(cfg, fileExt) // Get a name of save folder
			if err != nil {
				fmt.Printf("getting target folder name: %v", err)
				continue
			}

			finalPath := filepath.Join(targetPath, targetFolder) // Make a save path
			err = os.MkdirAll(finalPath, 0755)                   // Make a directories by path
			if err != nil {
				return fmt.Errorf("creating dirs by path %q: %w", finalPath, err)
			}

			oldFilePath := filepath.Join(targetPath, fileName)
			newFilePath := filepath.Join(finalPath, fileName)

			err = os.Rename(oldFilePath, newFilePath)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func TargetFolderName(cfg *config.Config, fileExt string) (string, error) {
	targetFolder, ok := cfg.InvertedRules[fileExt]
	if ok {
		return targetFolder, nil
	}
	return "", fmt.Errorf("ext isn't in config: %s", fileExt)
}
