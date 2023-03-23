// Package logm is a Go log management providing facilities to deal with logs and trace context.
package logm

import (
	"context"
	"io"
	"math"
	"runtime/debug"

	"golang.org/x/exp/slog"
)

// DebugLogger returns a new instance of Logger dedicated to debug or test environment.
// Debug messages are included and each message will include the application name and version.
func DebugLogger(name string, w io.Writer) *slog.Logger {
	return NewLogger(name, w, slog.LevelDebug)
}

// DiscardLogger is a logger doing anything. Useful for test purpose and default behavior.
func DiscardLogger() *slog.Logger {
	l := slog.New(slog.NewTextHandler(io.Discard))
	l.Enabled(context.Background(), math.MaxInt)
	return l
}

// NewLogger returns a new instance of Logger with level is the minimum log level to consider.
// Each message will include the application name and version.
func NewLogger(name string, w io.Writer, level slog.Level) *slog.Logger {
	h := slog.HandlerOptions{
		Level: level,
	}
	l := slog.New(h.NewTextHandler(w).WithAttrs([]slog.Attr{
		slog.String(AppNameKey, name),
		slog.String(AppVersionKey, vcsVersion()),
	}))
	return l
}

// vcsVersion returns the VCS version available since go1.18 in build info.
func vcsVersion() string {
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" {
				return setting.Value
			}
		}
	}
	return ""
}

// StdLogger returns a new instance of Logger, ready to be use in production mode.
// Debug messages are ignored and each message will include the application name and version.
func StdLogger(name string, w io.Writer) *slog.Logger {
	return NewLogger(name, w, slog.LevelInfo)
}

// TimeElapsed allows to monitor the execution time of a function.
// Example:
//
//	defer t.TimeElapsed(ctx, logger, slog.LevelDebug, "func")()
//
// It offers a useful interface to be called in defer statement.
func TimeElapsed(ctx context.Context, l *slog.Logger, level slog.Level, msg string, attrs ...slog.Attr) func() {
	t := NewTraceFromContext(ctx)
	t.Start()
	return func() {
		t.End()
		l.LogAttrs(ctx, level, msg, append(attrs, t.LogAttr())...)
	}
}
