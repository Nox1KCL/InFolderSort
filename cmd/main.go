package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/Nox1KCL/InFolderSort/internal/config"
	"github.com/Nox1KCL/InFolderSort/internal/logger"
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

	levels := map[slog.Level]string{
		slog.LevelInfo:  "logs/info.log",
		slog.LevelError: "logs/error.log",
	}
	handler, err := logger.GetHandler(&cfg.Logger, levels)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "creating logger: %v\n", err)
		os.Exit(1)
	}
	slog.SetDefault(slog.New(handler))

	// Start tui
	err = tui.Core(cfg)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "running application: %v\n", err)
		os.Exit(1)
	}
}
