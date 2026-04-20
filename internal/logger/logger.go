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

type LeveledHandler struct {
	infoHandler  slog.Handler
	errorHandler slog.Handler
}

func (h *LeveledHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return true
}

func (h *LeveledHandler) Handle(ctx context.Context, r slog.Record) error {
	if r.Level >= slog.LevelError {
		return h.errorHandler.Handle(ctx, r)
	}
	return h.infoHandler.Handle(ctx, r)
}

func (h *LeveledHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &LeveledHandler{
		infoHandler:  h.infoHandler.WithAttrs(attrs),
		errorHandler: h.errorHandler.WithAttrs(attrs),
	}
}
func (h *LeveledHandler) WithGroup(name string) slog.Handler {
	return &LeveledHandler{
		infoHandler:  h.infoHandler.WithGroup(name),
		errorHandler: h.errorHandler.WithGroup(name),
	}
}

func GetLogger(cfg *LumberConfig) (*LeveledHandler, error) {
	levels := map[slog.Level]string{
		slog.LevelInfo:  "logs/info.log",
		slog.LevelError: "logs/error.log",
	}

	var handlers []slog.Handler
	for level, path := range levels {
		absPath, err := filepath.Abs(path)
		if err != nil {
			return nil, fmt.Errorf("getting abs path %q: %w", path, err)
		}
		handler, err := cfg.HandlerConveyor(absPath, level)
		if err != nil {
			return nil, fmt.Errorf("creating logger %q: %w", absPath, err)
		}
		handlers = append(handlers, handler)

	}

	router := &LeveledHandler{
		infoHandler:  handlers[0],
		errorHandler: handlers[1],
	}
	// TODO Доробити нормальне присвоєння логерів

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
