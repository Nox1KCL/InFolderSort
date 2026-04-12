package files

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Nox1KCL/InFolderSort/internal/config"
)

type SortResult struct {
	Moved   []string
	Skipped []string
	Errors  []error
}

type MoveTask struct {
	FileName   string
	SourcePath string
	DestPath   string
}

type Sorter struct {
	Config  *config.Config
	ScanDir string
	Files   []os.DirEntry
	Tasks   []MoveTask
	Errors  []error
}

func NewSorter(targetDir string, cfg *config.Config) *Sorter {
	return &Sorter{
		Config:  cfg,
		ScanDir: targetDir,
		Files:   make([]os.DirEntry, 0),
		Tasks:   make([]MoveTask, 0),
		Errors:  make([]error, 0),
	}
}

func (s *Sorter) Scan() error {
	entries, err := os.ReadDir(s.ScanDir)
	if err != nil {
		return fmt.Errorf("reading directory %q: %w", s.ScanDir, err)
	}

	for _, entry := range entries {
		fileName := entry.Name()
		if !entry.IsDir() && !strings.HasPrefix(fileName, ".") {
			s.Files = append(s.Files, entry)
		}
	}
	return nil
}

func (s *Sorter) Plan() error {
	if len(s.Files) == 0 {
		return fmt.Errorf("no files found in %q", s.ScanDir)
	}
	if s.Config == nil {
		return fmt.Errorf("config is empty")
	}

	for _, file := range s.Files {
		fileName := file.Name()
		fileExt := filepath.Ext(fileName)
		targetPath, err := s.Config.GetTargetPath(fileExt)
		if err != nil {
			s.Errors = append(s.Errors, err)
			continue
		}

		var savePath string
		if filepath.IsAbs(targetPath) {
			savePath = targetPath
		} else {
			savePath = filepath.Join(s.ScanDir, targetPath)
		}

		exist, err := IsFileExist(filepath.Join(savePath, fileName))
		if err != nil {
			s.Errors = append(s.Errors, err)
			continue
		}
		if exist {
			fileName = RenameFile(fileName)
		}
		destPath := filepath.Join(savePath, fileName)
		s.Tasks = append(s.Tasks, MoveTask{fileName, s.ScanDir, destPath})
	}
	return nil
}

func IsFileExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	} else if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, fmt.Errorf("checking file %q: %w", path, err)
}

func RenameFile(file string) string {
	ext := filepath.Ext(file)
	name := strings.TrimSuffix(file, ext)
	timestamp := time.Now().Format("20060102_150405")
	newName := fmt.Sprintf("%s_%s%s", name, timestamp, ext)
	return newName
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
