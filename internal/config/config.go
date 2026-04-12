package config

import (
	_ "embed"
	"errors"
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

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("validating config: %w", err)
	}
	// Inverting config once in first call
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
	var conflicts []error

	for folderName, folderRule := range cfg.Rules {
		for _, ext := range folderRule.Extensions {
			ext = strings.ToLower(strings.TrimSpace(ext))

			if ext == "" {
				conflicts = append(conflicts, fmt.Errorf("empty extension in %s", folderName))
				continue
			}

			if !strings.HasPrefix(ext, ".") {
				conflicts = append(conflicts, fmt.Errorf("missing dot in extension %s in %s", ext, folderName))
				continue
			}

			if firstDebut, exists := seenExtensions[ext]; exists {
				conflicts = append(conflicts, fmt.Errorf("duplicate extension: %s | Seen it in %s and %s", ext, firstDebut, folderName))
				continue
			} else {
				seenExtensions[ext] = folderName
			}
		}
	}

	if len(conflicts) != 0 {
		return errors.Join(conflicts...)
	}
	return nil
}
