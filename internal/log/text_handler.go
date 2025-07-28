package log

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"runtime"
	"strconv"
	"sync"
)

// ANSI color codes
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorMagenta = "\033[35m"
	colorCyan   = "\033[36m"
	colorWhite  = "\033[37m"
)

// textHandler is a slog.Handler that writes log records to an io.Writer
// in a human-readable, elegant format.
type textHandler struct {
	slog.Handler
	out io.Writer
	mu  *sync.Mutex
}

// NewTextHandler creates a new textHandler.
func NewTextHandler(out io.Writer, opts *slog.HandlerOptions) slog.Handler {
	return &textHandler{
		Handler: slog.NewTextHandler(out, opts),
		out:	 out,
		mu:		&sync.Mutex{},
	}
}

// Handle formats the log record and writes it to the output.
func (h *textHandler) Handle(ctx context.Context, r slog.Record) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	buf := make([]byte, 0, 1024)

	// Add color based on level
	switch r.Level {
	case slog.LevelDebug:
		buf = append(buf, colorBlue...)
	case slog.LevelInfo:
		buf = append(buf, colorGreen...)
	case slog.LevelWarn:
		buf = append(buf, colorYellow...)
	case slog.LevelError:
		buf = append(buf, colorRed...)
	}

	// Time
	buf = r.Time.AppendFormat(buf, "2006-01-02 15:04:05.000")
	buf = append(buf, ' ')

	// Level
	buf = append(buf, '[')
	buf = append(buf, r.Level.String()...)
	buf = append(buf, ']')
	buf = append(buf, ' ')

	// Source (file:line)
	if r.PC != 0 {
		fs := runtime.CallersFrames([]uintptr{r.PC})
		f, _ := fs.Next()
		file := f.File
		line := f.Line
		shortFile := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				shortFile = file[i+1:]
				break
			}
		}
		buf = append(buf, '(')
		buf = append(buf, shortFile...)
		buf = append(buf, ':')
		buf = strconv.AppendInt(buf, int64(line), 10)
		buf = append(buf, []byte(") ")...)
	}

	// Message
	buf = append(buf, r.Message...)

	// Attributes
	r.Attrs(func(attr slog.Attr) bool {
		buf = append(buf, []byte("  ")...)
		buf = append(buf, attr.Key...)
		buf = append(buf, '=')
		buf = append(buf, []byte(fmt.Sprintf("%v", attr.Value.Any()))...)
		return true
	})

	buf = append(buf, colorReset...) // Reset color
	buf = append(buf, '\n')

	_, err := h.out.Write(buf)
	return err
}

// WithAttrs returns a new textHandler that includes the given attributes.
func (h *textHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &textHandler{
		Handler: h.Handler.WithAttrs(attrs),
		out:	 h.out,
		mu:		 h.mu,
	}
}

// WithGroup returns a new textHandler that starts a new group.
func (h *textHandler) WithGroup(name string) slog.Handler {
	return &textHandler{
		Handler: h.Handler.WithGroup(name),
		out:	 h.out,
		mu:		 h.mu,
	}
}
