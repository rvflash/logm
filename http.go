package logm

import (
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"

	"golang.org/x/exp/slog"
)

// LogHandler is an HTTP middleware designed to log every request and response.
func LogHandler(l *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := NewTraceFromHTTPRequest(r)
		t.Start()
		wh := newHTTPResponseWriter(w)
		next.ServeHTTP(wh, r)
		t.End()
		l.Info(
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
func RecoverHandler(msg string, l *slog.Logger, next http.Handler) http.Handler {
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
				t := NewTrace()
				l.Error(err.Error(), PanicKey, t, logHTTPRequest(r))
				l.Debug(string(debug.Stack()), PanicKey, t, logHTTPRequest(r))
				http.Error(w, msg, http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// TraceHandler is an HTTP middleware designed to share the trace context in the request context.
func TraceHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := NewTraceFromHTTPRequest(r)
		next.ServeHTTP(w, r.WithContext(t.NewContext(r.Context())))
	})
}
