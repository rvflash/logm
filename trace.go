package logm

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"

	"golang.org/x/exp/slog"
)

// TraceIDHTTPHeader is the name of the HTTP header used to share a trace context ID.
const TraceIDHTTPHeader = "X-Trace-Id"

// NewTrace creates a new Trace with a new generated UUID v4 as identifier.
func NewTrace() *Trace {
	return &Trace{ID: newUUID()}
}

type contextual string

const ctxTraceID contextual = "traceID"

// NewTraceFromContext returns a new Trace based on the context.Context.
// If the trace ID value is not found or blank, a new one is created.
// Otherwise, we create a trace span with this trace identifier as parent identifier.
func NewTraceFromContext(ctx context.Context) *Trace {
	return NewTraceSpan(contextValue(ctx, ctxTraceID))
}

func contextValue(ctx context.Context, key contextual) string {
	if v, ok := ctx.Value(key).(string); ok {
		return v
	}
	return ""
}

// NewTraceFromHTTPRequest returns a new Trace based on the http.Request.
// If the trace ID value is not found or blank, a new one is created.
// Otherwise, we create a trace span with this trace identifier as parent identifier.
func NewTraceFromHTTPRequest(req *http.Request) *Trace {
	return NewTraceSpan(req.Header.Get(TraceIDHTTPHeader))
}

// NewTraceSpan creates a trace span based on the Trace.
func NewTraceSpan(parentID string) *Trace {
	if parentID == "" {
		return NewTrace()
	}
	return &Trace{
		ID:     parentID,
		SpanID: newUUID(),
	}
}

func newUUID() string {
	id, err := uuid.NewRandom()
	if err != nil {
		return ""
	}
	return id.String()
}

// Trace represents a trace context.
type Trace struct {
	TimeElapsedMs int64
	StartTime     time.Time
	ID            string
	SpanID        string
}

// End ends the context trace and calculates the time elapsed since its starting.
func (t *Trace) End() {
	t.TimeElapsedMs = time.Since(t.StartTime).Milliseconds()
}

// LogAttr returns the trace as a slog.Attr.
func (t *Trace) LogAttr() slog.Attr {
	return slog.Group(TraceKey, t.logAttrs()...)
}

// LogValue implements the slog.logValuer interface.
func (t *Trace) LogValue() slog.Value {
	return slog.GroupValue(t.logAttrs()...)
}

func (t *Trace) logAttrs() []slog.Attr {
	res := []slog.Attr{
		slog.String(TraceIDKey, t.ID),
	}
	if t.SpanID != "" {
		res = append(res, slog.String(TraceSpanIDKey, t.SpanID))
	}
	if !t.StartTime.IsZero() {
		res = append(res, slog.Int64(TraceTimeElapsedKey, t.TimeElapsedMs))
	}
	return res
}

// NewContext creates a new trace context.Context to carry the trace identifier.
func (t *Trace) NewContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, ctxTraceID, t.ID)
}

// Start adds a start time to the trace.
func (t *Trace) Start() {
	t.StartTime = time.Now()
}
