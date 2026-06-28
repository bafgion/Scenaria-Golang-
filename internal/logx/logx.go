package logx

import (
	"log/slog"
	"os"
)

var logger = slog.Default()

func Init() {
	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo})
	logger = slog.New(handler)
	slog.SetDefault(logger)
}

func Info(msg string, args ...any)  { logger.Info(msg, args...) }
func Warn(msg string, args ...any)  { logger.Warn(msg, args...) }
func Error(msg string, args ...any) { logger.Error(msg, args...) }
func Debug(msg string, args ...any) { logger.Debug(msg, args...) }

func Logger() *slog.Logger { return logger }
