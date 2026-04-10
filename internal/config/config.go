package config

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	
	"github.com/pelletier/go-toml/v2"
)

//go:embed config.toml
var defaultConfig []byte

type FolderRule struct {
	TargetPath string   `toml:"target_path"`
	Extensions []string `toml:"extensions"`
}

type Config struct {
	Rules         map[string]FolderRule `toml:"rules"`
	InvertedRules map[string]string
}

func GetConfig(path string) (*Config, error) {
	var doc []byte
	var err error

	if path != "" {
		doc, err = os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("reading config file %q: %w", path, err)
		}
	} else {
		doc = defaultConfig
	}

	var cfg Config
	if err := toml.Unmarshal(doc, &cfg); err != nil {
		return nil, fmt.Errorf("reading toml doc %q: %w", path, err)
	}

	// Inverting config once in first call

	if err := cfg.Validate(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "invalid config: %v\n", err)
		os.Exit(1)
	}
	cfg.InvertConfig()

	return &cfg, nil
}

func (cfg *Config) InvertConfig() {
	if cfg.InvertedRules != nil {
		return
	}
	cfg.InvertedRules = make(map[string]string)

	for _, folderRule := range cfg.Rules {
		expandedPath := os.ExpandEnv(folderRule.TargetPath)
		finalPath := filepath.Clean(expandedPath)
		for _, ext := range folderRule.Extensions {
			cfg.InvertedRules[ext] = finalPath
		}
	}
}

func (cfg *Config) GetTargetPath(fileExt string) (string, error) {
	targetPath, ok := cfg.InvertedRules[fileExt]
	if ok {
		return targetPath, nil
	}
	return "", fmt.Errorf("ext isn't in config: %s", fileExt)
}

func (cfg *Config) Validate() error {
	seenExtensions := make(map[string]string)
	var conflicts []string

	for folderName, folderRule := range cfg.Rules {
		for _, ext := range folderRule.Extensions {
			ext = strings.ToLower(ext)
			if firstDebut, exists := seenExtensions[ext]; exists {
				conflicts = append(conflicts, fmt.Sprintf("duplicate extension: %s | Seen it in %s and %s", ext, firstDebut, folderName))
			} else {
				seenExtensions[ext] = folderName
			}
		}
	}

	if len(conflicts) != 0 {
		ReportConflicts(conflicts)
		return fmt.Errorf("conflicting extensions: %q", conflicts)
	}
	return nil
}

func ReportConflicts(conflicts []string) {
	fmt.Printf("Conflicts: %d\n", len(conflicts))
	fmt.Println("Conflicting extensions:")
	for _, conflict := range conflicts {
		fmt.Println(conflict)
	}
}
