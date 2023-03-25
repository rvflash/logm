package logm_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/matryer/is"
	"github.com/rvflash/logm"
)

const (
	target = "http://testing/"
	intErr = "oops I did it again"
)

func TestLogHandler(t *testing.T) {
	t.Parallel()
	var (
		are  = is.New(t)
		next = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(info))
		})
		buf = new(bytes.Buffer)
		log = logm.DefaultLogger(name, buf)
		hdl = logm.LogHandler(log, next)
		req = httptest.NewRequest(http.MethodGet, target, nil)
		res = httptest.NewRecorder()
	)
	hdl.ServeHTTP(res, req)
	out := buf.String()
	are.True(strings.Contains(out, "200 GET /"))        // message expected
	are.True(strings.Contains(out, "app=app"))          // application expected
	are.True(strings.Contains(out, "resp.size=5"))      // response size expected
	are.True(strings.Contains(out, "req.path=/"))       // request path expected
	are.True(strings.Contains(out, "req.method=GET"))   // request method expected
	are.True(strings.Contains(out, "req.remote_addr=")) // remote addr path expected
	are.True(strings.Contains(out, "req.query="))       // request query expected
	are.True(strings.Contains(out, "trace.id="))        // request trace id expected
}

func TestTraceHandler(t *testing.T) {
	t.Parallel()
	var (
		are  = is.New(t)
		next = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tc := logm.NewTraceFromContext(r.Context())
			if tc.ID == traceID {
				w.WriteHeader(http.StatusOK)
			} else {
				w.WriteHeader(http.StatusBadRequest)
				_, _ = w.Write([]byte(tc.ID))
			}
		})
		hdl = logm.TraceHandler(next)
		req = httptest.NewRequest(http.MethodGet, target, nil)
		res = httptest.NewRecorder()
	)
	req.Header.Add(logm.TraceIDHTTPHeader, traceID)
	hdl.ServeHTTP(res, req)
	are.Equal(http.StatusOK, res.Code) // unexpected response code
	are.Equal("", res.Body.String())   // unexpected response content
}

func TestRecoverHandler(t *testing.T) {
	t.Parallel()
	var (
		are  = is.New(t)
		next = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			panic(warn)
		})
		buf = new(bytes.Buffer)
		log = logm.DebugLogger(name, buf)
		hdl = logm.RecoverHandler(intErr, log, next)
		req = httptest.NewRequest(http.MethodGet, target, nil)
		res = httptest.NewRecorder()
	)
	hdl.ServeHTTP(res, req)
	are.Equal(http.StatusInternalServerError, res.Code) // unexpected response code
	out := buf.String()
	are.True(strings.Contains(out, "level=ERROR msg=earth"))      // unexpected error message
	are.True(strings.Contains(out, `level=DEBUG msg="goroutine`)) // unexpected debug message
}
