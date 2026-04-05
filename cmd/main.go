package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Nox1KCL/InFolderSort/internal/config"
	"github.com/Nox1KCL/InFolderSort/internal/tui"
)

func main() {
	configPath := flag.String("config", "", "path to config file (uses embedded default if empty)")
	flag.Parse()
	actualConfigPath := *configPath

	cfg, err := config.GetConfig(actualConfigPath)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "get configuration file: %v\n", err)
		os.Exit(1)
	}

	// Start tui
	err = tui.Core(cfg)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "running application: %v\n", err)
		os.Exit(1)
	}
}
