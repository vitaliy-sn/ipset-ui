package logger

import (
	"log/slog"
	"os"
	"sync"
)

var (
	instance *slog.Logger
	once     sync.Once
)

// Init initializes the global slog.Logger (once).
func initialize() {
	once.Do(func() {
		instance = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))
	})
}

// Get returns the global slog.Logger.
func get() *slog.Logger {
	if instance == nil {
		initialize()
	}
	return instance
}

// Info logs an Info level message without context.
func Info(msg string, args ...any) {
	get().Info(msg, args...)
}

// Error logs an Error level message without context.
func Error(msg string, args ...any) {
	get().Error(msg, args...)
}

// Warn logs a Warn level message without context.
func Warn(msg string, args ...any) {
	get().Warn(msg, args...)
}

// Debug logs a Debug level message without context.
func Debug(msg string, args ...any) {
	get().Debug(msg, args...)
}

// Fatal logs an Error level message and terminates the program.
func Fatal(msg string, args ...any) {
	get().Error(msg, args...)
	os.Exit(1)
}
