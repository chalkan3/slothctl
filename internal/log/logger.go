package log

import (
	"io"
	"log/slog"
	"math/rand"
	"os"
	"time"
)

// Logger is the interface for our structured logger.
type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
	Fatal(msg string, args ...any)

	With(args ...any) Logger
}

// logger implements the Logger interface using slog.
type logger struct {
	sl *slog.Logger
}

// NewLogger creates a new Logger instance with the given options.
func NewLogger(opts *Options) Logger {
	if opts == nil {
		opts = DefaultOptions()
	}

	var handler slog.Handler
	handlerOpts := &slog.HandlerOptions{Level: opts.Level, AddSource: true}

	switch opts.Format {
	case FormatJSON:
		handler = slog.NewJSONHandler(opts.Output, handlerOpts)
	case FormatText:
		handler = NewTextHandler(opts.Output, handlerOpts)
	default:
		handler = NewTextHandler(opts.Output, handlerOpts) // Default to text
	}

	return &logger{sl: slog.New(handler)}
}

func (l *logger) Debug(msg string, args ...any) {
	l.sl.Debug(msg, args...)
}

func (l *logger) Info(msg string, args ...any) {
	l.sl.Info(msg, args...)
}

func (l *logger) Warn(msg string, args ...any) {
	l.sl.Warn(msg, args...)
}

func (l *logger) Error(msg string, args ...any) {
	l.sl.Error(msg, args...)
}

func (l *logger) Fatal(msg string, args ...any) {
	l.sl.Error(msg, args...)
	os.Exit(1)
}

func (l *logger) With(args ...any) Logger {
	return &logger{sl: l.sl.With(args...)}
}

// Global logger instance
var defaultLogger = NewLogger(DefaultOptions())

// Exported functions for convenience
func Debug(msg string, args ...any) {
	defaultLogger.Debug(msg, args...)
}

func Info(msg string, args ...any) {
	defaultLogger.Info(msg, args...)
}

func Warn(msg string, args ...any) {
	defaultLogger.Warn(msg, args...)
}

func Error(msg string, args ...any) {
	defaultLogger.Error(msg, args...)
}

func Fatal(msg string, args ...any) {
	defaultLogger.Fatal(msg, args...)
}

// SetDefaultLogger allows changing the global logger instance.
func SetDefaultLogger(l Logger) {
	defaultLogger = l
}

// writer is an io.Writer that writes to a slog.Logger at a specific level.
type writer struct {
	logFunc func(msg string, args ...any)
}

// Write implements the io.Writer interface.
func (w *writer) Write(p []byte) (n int, err error) {
	msg := string(p)
	if len(msg) > 0 && msg[len(msg)-1] == '\n' {
		msg = msg[:len(msg)-1] // Remove trailing newline
	}
	w.logFunc(msg)
	return len(p), nil
}

// NewWriter creates a new io.Writer that logs to the specified log function.
func NewWriter(logFunc func(msg string, args ...any)) io.Writer {
	return &writer{logFunc: logFunc}
}

var emojis = []string{"🦥", "😴", "🌳", "🌿", "💚", "✨", "💖", "🌟"}

// GetRandomSlothEmoji returns a random cute sloth emoji.
func GetRandomSlothEmoji() string {
	rand.Seed(time.Now().UnixNano())
	return emojis[rand.Intn(len(emojis))]
}
