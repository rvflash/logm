package logm_test

import (
	"context"
	"math"
	"net/http"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/matryer/is"
	"github.com/rvflash/logm"

	"golang.org/x/exp/slog"
)

const (
	traceID = "7300cb05-8323-4dcc-8272-8d6a2c6b7fbc"
	spanID  = "1d889d18-9159-4ff1-9397-fab540cbeb17"
)

func TestNewTrace(t *testing.T) {
	t.Parallel()
	tc := logm.NewTrace()
	is.New(t).True(tc.ID != "") // expected trace ID
}

func TestNewTraceFromContext(t *testing.T) {
	t.Parallel()
	are := is.New(t)

	t.Run("Undefined", func(t *testing.T) {
		t.Parallel()
		tc := logm.NewTraceFromContext(context.Background())
		are.True(tc.ID != "") // expected trace ID
	})

	t.Run("Predefined", func(t *testing.T) {
		t.Parallel()
		t1 := logm.Trace{ID: traceID}
		t2 := logm.NewTraceFromContext(t1.NewContext(context.Background()))
		is.New(t).Equal(traceID, t2.ID) // mismatch trace ID
	})
}

func TestNewTraceFromHTTPRequest(t *testing.T) {
	t.Parallel()
	are := is.New(t)

	t.Run("Undefined", func(t *testing.T) {
		t.Parallel()
		tc := logm.NewTraceFromHTTPRequest(&http.Request{})
		are.True(tc.ID != "") // expected trace ID
	})

	t.Run("Undefined", func(t *testing.T) {
		t.Parallel()
		var (
			req = &http.Request{
				Header: map[string][]string{
					logm.TraceIDHTTPHeader: {traceID},
				},
			}
			tc = logm.NewTraceFromHTTPRequest(req)
		)
		are.Equal(traceID, tc.ID) // mismatch trace ID
	})
}

func TestNewTraceSpan(t *testing.T) {
	t.Parallel()
	are := is.New(t)

	t.Run("Default", func(t *testing.T) {
		t.Parallel()
		tc := logm.NewTraceSpan("")
		are.True(tc.ID != "")     // expected trace ID
		are.True(tc.SpanID == "") // unexpected trace span ID
	})

	t.Run("OK", func(t *testing.T) {
		t.Parallel()
		tc := logm.NewTraceSpan(traceID)
		are.Equal(traceID, tc.ID) // expected trace ID
		are.True(tc.SpanID != "") // expected trace span ID
	})
}

func TestTrace_End(t *testing.T) {
	t.Parallel()
	tc := logm.Trace{StartTime: time.Now().Add(-time.Second)}
	tc.End()
	is.New(t).Equal(time.Second.Milliseconds(), tc.TimeElapsedMs) // mismatch time elapsed
}

func TestTrace_LogAttr(t *testing.T) {
	t.Parallel()

	are := is.New(t)

	for name, tc := range map[string]struct {
		in  logm.Trace
		out []slog.Attr
	}{
		"Default": {out: []slog.Attr{slog.String(logm.TraceIDKey, "")}},
		"ID only": {in: logm.Trace{ID: traceID}, out: []slog.Attr{slog.String(logm.TraceIDKey, traceID)}},
		"With span ID": {
			in: logm.Trace{ID: traceID, SpanID: spanID},
			out: []slog.Attr{
				slog.String(logm.TraceIDKey, traceID),
				slog.String(logm.TraceSpanIDKey, spanID),
			},
		},
		"Unexpected time elapsed (missing start time)": {
			in: logm.Trace{
				TimeElapsedMs: math.MaxUint8,
				ID:            traceID,
				SpanID:        spanID,
			},
			out: []slog.Attr{
				slog.String(logm.TraceIDKey, traceID),
				slog.String(logm.TraceSpanIDKey, spanID),
			},
		},
		"Complete": {
			in: logm.Trace{
				StartTime:     time.Now(), // not used
				TimeElapsedMs: math.MaxUint8,
				ID:            traceID,
				SpanID:        spanID,
			},
			out: []slog.Attr{
				slog.String(logm.TraceIDKey, traceID),
				slog.String(logm.TraceSpanIDKey, spanID),
				slog.Int64(logm.TraceTimeElapsedKey, math.MaxUint8),
			},
		},
	} {
		tt := tc
		t.Run(name, func(t *testing.T) {
			out := tt.in.LogAttr()
			are.Equal(logm.TraceKey, out.Key)                         // mismatch key
			are.Equal("", cmp.Diff(tt.out, out.Value.Group()))        // mismatch attr value
			are.Equal("", cmp.Diff(tt.out, tt.in.LogValue().Group())) // mismatch value
		})
	}
}

func TestTrace_NewContext(t *testing.T) {
	t.Parallel()
	var (
		t1  = logm.Trace{ID: traceID}
		ctx = t1.NewContext(context.Background())
		t2  = logm.NewTraceFromContext(ctx)
	)
	is.New(t).Equal(traceID, t2.ID) // mismatch trace ID
}

func TestTrace_Start(t *testing.T) {
	t.Parallel()
	tc := logm.Trace{}
	tc.Start()
	is.New(t).True(tc.StartTime.Before(time.Now()))
}
