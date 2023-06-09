package logm

import (
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"

	"golang.org/x/exp/slog"
)

// Middleware provides some standard HTTP handlers to deal with logs.
type Middleware struct {
	Logger       *slog.Logger
	ErrorMessage string
}

// LogHandler is an HTTP middleware designed to log every request and response.
func (m Middleware) LogHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := NewTraceFromHTTPRequest(r)
		t.Start()
		wh := newHTTPResponseWriter(w)
		next.ServeHTTP(wh, r)
		t.End()
		m.Logger.Info(
			fmt.Sprintf("%d %s %s", wh.statusCode, r.Method, r.URL.Path),
			logHTTPRequest(r),
			logHTTPResponse(wh),
			t.LogAttr(),
		)
	})
}

func logHTTPRequest(r *http.Request) slog.Attr {
	return slog.Group(HTTPRequestKey,
		slog.String(HTTPPathKey, r.URL.Path),
		slog.String(HTTPMethodKey, r.Method),
		slog.String(HTTPRemoteAddrKey, r.RemoteAddr),
		slog.String(HTTPQueryKey, r.URL.RawQuery),
	)
}

func logHTTPResponse(w *httpResponseWriter) slog.Attr {
	return slog.Group(HTTPResponseKey,
		slog.Int(HTTPStatusKey, w.statusCode),
		slog.Int(HTTPSizeKey, w.size),
	)
}

func newHTTPResponseWriter(w http.ResponseWriter) *httpResponseWriter {
	return &httpResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
}

type httpResponseWriter struct {
	http.ResponseWriter
	size       int
	statusCode int
}

// Write wraps the response writer to follow the response size.
func (w *httpResponseWriter) Write(data []byte) (n int, err error) {
	n, err = w.ResponseWriter.Write(data)
	w.size += n
	return
}

// WriteHeader captures the HTTP response header.
func (w *httpResponseWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

// RecoverHandler is an HTTP middleware designed to recover on panic, log the error and debug the stack trace.
// If no error message is provided, we used the default internal error message.
func (m Middleware) RecoverHandler(next http.Handler) http.Handler {
	msg := m.ErrorMessage
	if msg == "" {
		msg = http.StatusText(http.StatusInternalServerError)
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			pr := recover()
			if pr != nil {
				var err error
				switch t := pr.(type) {
				case string:
					err = errors.New(t)
				case error:
					err = t
				default:
					err = fmt.Errorf("unsupported panic type: %#v", t)
				}
				t := NewTraceFromHTTPRequest(r)
				m.Logger.Error(err.Error(), PanicKey, t, logHTTPRequest(r))
				m.Logger.Debug(string(debug.Stack()), PanicKey, t, logHTTPRequest(r))
				http.Error(w, msg, http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// LogHandler is an HTTP middleware designed to log every request and response.
func LogHandler(l *slog.Logger, next http.Handler) http.Handler {
	return Middleware{Logger: l}.LogHandler(next)
}

// RecoverHandler is an HTTP middleware designed to recover on panic, log the error and debug the stack trace.
func RecoverHandler(msg string, l *slog.Logger, next http.Handler) http.Handler {
	return Middleware{ErrorMessage: msg, Logger: l}.RecoverHandler(next)
}

// TraceHandler is an HTTP middleware designed to share the trace context in the request context.
func TraceHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := NewTraceFromHTTPRequest(r)
		r = r.WithContext(t.NewContext(r.Context()))
		r.Header.Set(TraceIDHTTPHeader, t.ID)
		next.ServeHTTP(w, r)
	})
}
