package main

import (
	"flag"
	"fmt"
	"main/internal/config"
	"main/internal/tui"
	"os"
)

func main() {
	configPath := flag.String("config", "internal/config/config.toml", "Path to configuration file")
	flag.Parse()
	actualPath := *configPath

	cfg, err := config.GetConfig(actualPath)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "get configuration file: %v\n", err)
		os.Exit(1)
	}

	// Start tui
	err = tui.Core(cfg)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// Зупинився на 7
