package logm_test

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"

	"github.com/matryer/is"
	"github.com/rvflash/logm"

	"golang.org/x/exp/slog"
)

const (
	name  = "app"
	info  = "hello"
	debug = "world"
	warn  = "earth"
)

func TestDebugLogger(t *testing.T) {
	t.Parallel()
	var (
		are = is.New(t)
		buf = new(bytes.Buffer)
		l   = logm.DebugLogger(name, buf)
	)
	l.Info(info)
	l.Debug(debug)
	l.Warn(warn)

	out := buf.String()
	are.True(strings.Contains(out, name))  // missing app message
	are.True(strings.Contains(out, info))  // missing log message
	are.True(strings.Contains(out, debug)) // missing debug message
	are.True(strings.Contains(out, warn))  // missing warn message
}

func TestDiscardLogger(t *testing.T) {
	t.Parallel()
	l := logm.DiscardLogger()
	is.New(t).True(l != nil)
}

func TestNewLogger(t *testing.T) {
	t.Parallel()
	var (
		are = is.New(t)
		buf = new(bytes.Buffer)
		l   = logm.NewLogger(name, buf, slog.LevelWarn)
	)
	l.Info(info)
	l.Debug(debug)
	l.Warn(warn)

	out := buf.String()
	are.True(strings.Contains(out, name))   // missing app message
	are.True(!strings.Contains(out, info))  // unexpected log message
	are.True(!strings.Contains(out, debug)) // unexpected debug message
	are.True(strings.Contains(out, warn))   // missing warn message
}

func TestDefaultLogger(t *testing.T) {
	t.Parallel()
	var (
		are = is.New(t)
		buf = new(bytes.Buffer)
		l   = logm.DefaultLogger(name, buf)
	)
	l.Info(info)
	l.Debug(debug)
	l.Warn(warn)

	out := buf.String()
	are.True(strings.Contains(out, name))   // missing app message
	are.True(strings.Contains(out, info))   // missing log message
	are.True(!strings.Contains(out, debug)) // unexpected debug message
	are.True(strings.Contains(out, warn))   // missing warn message
}

func TestTimeElapsed(t *testing.T) {
	t.Parallel()
	var (
		tc  = logm.Trace{ID: traceID}
		buf = new(bytes.Buffer)
		l   = logm.DebugLogger(name, buf)
		are = is.New(t)
	)
	t.Cleanup(func() {
		out := buf.String()
		are.True(strings.Contains(out, info))                     // missing log message
		are.True(strings.Contains(out, "trace.time_elapsed_ms=")) // missing time elapsed
	})
	defer logm.TimeElapsed(tc.NewContext(context.Background()), l, slog.LevelInfo, info)()
	time.Sleep(time.Millisecond)
}
