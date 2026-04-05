package config

import (
	"fmt"
	"os"

	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	Rules         map[string][]string `toml:"rules"`
	InvertedRules map[string]string
}

func GetConfig(path string) (*Config, error) {
	var cfg Config

	doc, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file %q: %w", path, err)
	}

	err = toml.Unmarshal(doc, &cfg)
	if err != nil {
		return nil, fmt.Errorf("reading toml doc %q: %w", path, err)
	}

	// Inverting config once in first call
	cfg.InvertConfig()

	return &cfg, nil
}

func (cfg *Config) InvertConfig() {
	cfg.InvertedRules = make(map[string]string)

	for folder, exts := range cfg.Rules {
		for _, ext := range exts {
			cfg.InvertedRules[ext] = folder
		}
	}
}

func (cfg *Config) TargetFolderName(fileExt string) (string, error) {
	targetFolder, ok := cfg.InvertedRules[fileExt]
	if ok {
		return targetFolder, nil
	}
	return "", fmt.Errorf("ext isn't in config: %s", fileExt)
}
