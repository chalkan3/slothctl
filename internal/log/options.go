package log

import (
	"io"
	"log/slog"
	"os"
)

// Format defines the output format for the logger.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Options holds the configuration for the logger.
type Options struct {
	Level  slog.Level
	Format Format
	Output io.Writer
}

// DefaultOptions returns a set of default logger options.
func DefaultOptions() *Options {
	return &Options{
		Level:  slog.LevelInfo,
		Format: FormatText,
		Output: os.Stderr,
	}
}
