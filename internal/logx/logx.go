package logx

import (
	"log/slog"
	"os"

	"github.com/bafgion/scenaria-golang/internal/secret"
)

var logger = slog.Default()

func Init() {
	opts := &slog.HandlerOptions{Level: slog.LevelInfo}
	var handler slog.Handler
	switch os.Getenv("SCENARIA_LOG") {
	case "json":
		handler = slog.NewJSONHandler(os.Stderr, opts)
	case "debug":
		opts.Level = slog.LevelDebug
		handler = slog.NewTextHandler(os.Stderr, opts)
	default:
		handler = slog.NewTextHandler(os.Stderr, opts)
	}
	logger = slog.New(handler)
	slog.SetDefault(logger)
}

func Info(msg string, args ...any) {
	logger.Info(secret.RedactString(msg), secret.RedactKeyValues(args)...)
}
func Warn(msg string, args ...any) {
	logger.Warn(secret.RedactString(msg), secret.RedactKeyValues(args)...)
}
func Error(msg string, args ...any) {
	logger.Error(secret.RedactString(msg), secret.RedactKeyValues(args)...)
}
func Debug(msg string, args ...any) {
	logger.Debug(secret.RedactString(msg), secret.RedactKeyValues(args)...)
}

func Logger() *slog.Logger { return logger }
