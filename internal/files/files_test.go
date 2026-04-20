package files

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Nox1KCL/InFolderSort/internal/config"
)

func TestInDirSorting_Basic(t *testing.T) {
	dir := t.TempDir()
	_ = os.WriteFile(filepath.Join(dir, "photo.jpg"), []byte("data"), 0644)
	_ = os.WriteFile(filepath.Join(dir, "doc.pdf"), []byte("data"), 0644)

	cfg := &config.Config{
		Rules: map[string]config.FolderRule{
			"Images": {
				TargetPath: "Images",
				Extensions: []string{".jpg"},
			},
			"Docs": {
				TargetPath: "Docs",
				Extensions: []string{".pdf"},
			},
		},
	}
	cfg.InvertConfig()

	report, err := InDirSorting(dir, cfg)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "Images", "photo.jpg")); err != nil {
		t.Error("photo.jpg was not moved to Images/")
	}
	if _, err := os.Stat(filepath.Join(dir, "Docs", "doc.pdf")); err != nil {
		t.Error("doc.pdf was not moved to Docs/")
	}

	if len(report.Moved) != 2 {
		t.Errorf("expected 2 moved files, got %d", len(report.Moved))
	}
}

func TestSorter_ConflictResolution(t *testing.T) {
	dir := t.TempDir()
	
	// Створюємо цільову папку і файл, який вже там існує (Конфлікт!)
	targetSubdir := filepath.Join(dir, "Images")
	_ = os.MkdirAll(targetSubdir, 0755)
	_ = os.WriteFile(filepath.Join(targetSubdir, "photo.jpg"), []byte("old data"), 0644)

	// Створюємо новий файл у папці для сканування
	_ = os.WriteFile(filepath.Join(dir, "photo.jpg"), []byte("new data"), 0644)

	cfg := &config.Config{
		Rules: map[string]config.FolderRule{
			"Images": {
				TargetPath: "Images",
				Extensions: []string{".jpg"},
			},
		},
	}
	cfg.InvertConfig()

	sorter := NewSorter(dir, cfg)
	
	if err := sorter.Scan(); err != nil {
		t.Fatalf("Scan failed: %v", err)
	}
	
	if err := sorter.Plan(); err != nil {
		t.Fatalf("Plan failed: %v", err)
	}

	// Перевіряємо, чи план дійсно запропонував нове ім'я
	if len(sorter.Tasks) != 1 {
		t.Fatalf("expected 1 task, got %d", len(sorter.Tasks))
	}

	destPath := sorter.Tasks[0].DestPath
	if !strings.Contains(destPath, "photo_") {
		t.Errorf("expected renamed file with timestamp, got %s", destPath)
	}

	// Виконуємо сортування
	report, err := sorter.Execute()
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if len(report.Moved) != 1 {
		t.Errorf("expected 1 moved file, got %d", len(report.Moved))
	}

	// Перевіряємо, чи обидва файли тепер існують у цільовій папці
	entries, _ := os.ReadDir(targetSubdir)
	if len(entries) != 2 {
		t.Errorf("expected 2 files in target subdir, got %d", len(entries))
	}
}
