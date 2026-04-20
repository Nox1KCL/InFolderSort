package logger

import (
	"context"
	"fmt"
	"log/slog"
	"path/filepath"

	"gopkg.in/natefinch/lumberjack.v2"
)

type LumberConfig struct {
	MaxSize    int
	MaxAge     int
	MaxBackups int
	Compress   bool
}

type handlerEntry struct {
	level   slog.Level
	handler slog.Handler
}

type LeveledHandler struct {
	handlers []handlerEntry
}

func (h *LeveledHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return true
}

func (h *LeveledHandler) Handle(ctx context.Context, r slog.Record) error {
	for _, entry := range h.handlers {
		if r.Level >= entry.level {
			return entry.handler.Handle(ctx, r)
		}
	}
	return nil
}

func (h *LeveledHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	var newHandles []handlerEntry
	for _, entry := range h.handlers {
		newHandles = append(newHandles, handlerEntry{entry.level, entry.handler.WithAttrs(attrs)})
	}
	return &LeveledHandler{handlers: newHandles}
}
func (h *LeveledHandler) WithGroup(name string) slog.Handler {
	var newHandles []handlerEntry
	for _, entry := range h.handlers {
		newHandles = append(newHandles, handlerEntry{entry.level, entry.handler.WithGroup(name)})
	}
	return &LeveledHandler{handlers: newHandles}
}

func GetLogger(cfg *LumberConfig) (*LeveledHandler, error) {
	levels := map[slog.Level]string{
		slog.LevelInfo:  "logs/info.log",
		slog.LevelError: "logs/error.log",
	}
	absInfoPath, err := filepath.Abs(levels[slog.LevelInfo])
	if err != nil {
		return nil, fmt.Errorf("getting abs path %q: %w", absInfoPath, err)
	}
	infoHandler, err := cfg.HandlerConveyor(absInfoPath, slog.LevelInfo)
	if err != nil {
		return nil, fmt.Errorf("creating logger %q: %w", absInfoPath, err)
	}

	absErrorPath, err := filepath.Abs(levels[slog.LevelError])
	if err != nil {
		return nil, fmt.Errorf("getting abs path %q: %w", absErrorPath, err)
	}
	errorHandler, err := cfg.HandlerConveyor(absErrorPath, slog.LevelError)
	if err != nil {
		return nil, fmt.Errorf("creating logger %q: %w", absErrorPath, err)
	}

	router := &LeveledHandler{
		handlers: []handlerEntry{
			{slog.LevelInfo, infoHandler},
			{slog.LevelError, errorHandler},
		},
	}

	return router, nil
}

func (cfg *LumberConfig) HandlerConveyor(filename string, level slog.Level) (slog.Handler, error) {
	handler := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    cfg.MaxSize,
		MaxAge:     cfg.MaxAge,
		MaxBackups: cfg.MaxBackups,
		Compress:   cfg.Compress,
	}
	logger := slog.NewJSONHandler(handler, &slog.HandlerOptions{Level: level})
	return logger, nil
}
