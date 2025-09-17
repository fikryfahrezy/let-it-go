package logger

import (
	"io"
	"log/slog"
	"os"
)

type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
)

func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "debug"
	case LevelInfo:
		return "info"
	case LevelWarn:
		return "warn"
	case LevelError:
		return "error"
	default:
		return "info"
	}
}

func ParseLevel(s string) Level {
	switch s {
	case "debug":
		return LevelDebug
	case "info":
		return LevelInfo
	case "warn":
		return LevelWarn
	case "error":
		return LevelError
	default:
		return LevelInfo
	}
}

type Format int

const (
	FormatText Format = iota
	FormatJSON
)

func (f Format) String() string {
	switch f {
	case FormatJSON:
		return "json"
	default:
		return "text"
	}
}

func ParseFormat(s string) Format {
	switch s {
	case "json":
		return FormatJSON
	default:
		return FormatText
	}
}

type Config struct {
	Level         Level
	Format        Format
	DisableOutput bool // For tests - sends output to io.Discard
}

func NewLogger(config Config) *slog.Logger {
	var level slog.Level
	switch config.Level {
	case LevelDebug:
		level = slog.LevelDebug
	case LevelInfo:
		level = slog.LevelInfo
	case LevelWarn:
		level = slog.LevelWarn
	case LevelError:
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level:     level,
		AddSource: true,
	}

	// Choose output destination
	var output io.Writer = os.Stdout
	if config.DisableOutput {
		output = io.Discard
	}

	var handler slog.Handler
	switch config.Format {
	case FormatJSON:
		handler = slog.NewJSONHandler(output, opts)
	default:
		handler = slog.NewTextHandler(output, opts)
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)

	return logger
}

// NewDiscardLogger creates a logger that discards all output for testing
func NewDiscardLogger() *slog.Logger {
	return NewLogger(Config{
		Level:         LevelInfo,
		Format:        FormatText,
		DisableOutput: true,
	})
}
