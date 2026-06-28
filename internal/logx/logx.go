package logx

import (
	"log/slog"
	"os"
)

var logger = slog.Default()

func Init() {
	opts := &slog.HandlerOptions{Level: slog.LevelInfo}
	var handler slog.Handler
	switch os.Getenv("SCENARIA_LOG") {
	case "json":
		handler = slog.NewJSONHandler(os.Stderr, opts)
	default:
		handler = slog.NewTextHandler(os.Stderr, opts)
	}
	logger = slog.New(handler)
	slog.SetDefault(logger)
}

func Info(msg string, args ...any)  { logger.Info(msg, args...) }
func Warn(msg string, args ...any)  { logger.Warn(msg, args...) }
func Error(msg string, args ...any) { logger.Error(msg, args...) }
func Debug(msg string, args ...any) { logger.Debug(msg, args...) }

func Logger() *slog.Logger { return logger }
