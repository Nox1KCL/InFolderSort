package files

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Nox1KCL/InFolderSort/internal/config"
)

func TestInDirSorting(t *testing.T) {
	dir := t.TempDir()
	_ = os.WriteFile(filepath.Join(dir, "photo.jpg"), []byte("data"), 0644)
	_ = os.WriteFile(filepath.Join(dir, "doc.pdf"), []byte("data"), 0644)

	cfg := &config.Config{
		Rules: map[string][]string{
			"Images": {".jpg", ".png", ".bmp"},
			"Docs":   {".pdf", ".docx", ".doc"},
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

	t.Logf("Sorting report: %v", report)
}
