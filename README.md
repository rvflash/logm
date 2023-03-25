# LogM

[![GoDoc](https://godoc.org/github.com/rvflash/logm?status.svg)](https://godoc.org/github.com/rvflash/logm)
[![Build Status](https://github.com/rvflash/logm/workflows/build/badge.svg)](https://github.com/rvflash/logm/actions?workflow=build)
[![Code Coverage](https://codecov.io/gh/rvflash/logm/branch/master/graph/badge.svg)](https://codecov.io/gh/rvflash/logm)
[![Go Report Card](https://goreportcard.com/badge/github.com/rvflash/logm?)](https://goreportcard.com/report/github.com/rvflash/logm)


LogM is a Go log management package providing facilities to deal with logs and trace context.

## Features

1. Provides simple methods to expose a [slog.Logger](https://pkg.go.dev/golang.org/x/exp/slog) using `logfmt` as format:
   - `DefaultLogger`: A logger ready for production.
   - `DebugLogger`: A logger exposing debug record for development.
   - `DiscardLogger`: Another to discard any logs (test purposes or no space left on disk).
2. Provides a `File` with automatic rotating, maximum file size, zip archives, etc. Thanks to [lumberjack](https://github.com/natefinch/lumberjack).
3. Provides a `Trace` structure to uniquely identified actions, like an HTTP request. See `NewTraceFromContext` to easily propagate or retrieve trace context.  
4. Exposes HTTP middlewares to handle log and tracing:to create a trace context on each request.
   - `LogHandler`: a logging middleware to log detail about the request and the response.
   - `RecoverHandler`: a middleware to recover on panic, log the message as ERROR and the stack trace as DEBUG. 
   - `TraceHandler`: a middleware to retrieve request header `X-Trace-Id` (see `NewTraceFromHTTPRequest`) and propagate its value through the request context.
5. Provides `TimeElapsed` to log in defer the time elapsed of a function. 


### Installation

```bash
$ go get -u github.com/rvflash/logm
```

### Prerequisite

`logm` uses the Go modules, `debug.BuildSetting` and the `slog` package that required Go 1.18 or later.


## Exemples

### Create a logger for an application named `app` that store records in a file named `app.log`, ignoring DEBUG ones.

By default, the `NewFile` function returns a self rolling file as soon as it reaches the size of 100 Mo.
The rotated log files is compressed using gzip and retained forever. Each log is prefixed by the local time.
Customization is available by using `File` directly.
By default, each record has the name and the current version of the application as attributes.

> The version is provided on build by the `debug.ReadBuildInfo` package.
> It's the revision identifier for the current commit or checkout.

```go
log := logm.DefaultLogger("app", logm.NewFile("app.log"))
log.Info("hello")
```
```bash
$ cat app.log
time=2023-03-25T00:04:35.287+01:00 level=INFO msg=hello app=app version=d1da844711730f2f5cbd08be93e62e71475f7d4e
```

### Create a logger to debug on standard output.

```go
log := logm.DebugLogger("app", os.Stdout)
log.Debug("hello")
```
```bash
time=2023-03-25T10:57:51.772+01:00 level=DEBUG msg=hello app=app version=d1da844711730f2f5cbd08be93e62e71475f7d4e
```

### Propagate a trace identifier through the context.

`NewTrace` can create a new trace context with an UUID v4 as identifier.
It's also possible to create a custom one by directly  using `Trace` and propagate it through a `context.Context`.

```go
var (
    t   = logm.Trace{ID: "myID"}
    ctx = t.NewContext(context.Background())
)
```

### Add the trace context on each log.

```go
var (
    l = logm.DefaultLogger("app", os.Stdout)
    t = logm.NewTrace()
)
log := l.With(t.LogAttr())
log.Info("hello")
log.Warn("world")
```
```bash
time=2023-03-25T13:06:37.322+01:00 level=INFO msg=hello app=app version=d1da844711730f2f5cbd08be93e62e71475f7d4e trace.id=0a02e16c-7418-4558-9dcc-718c007162b6
time=2023-03-25T13:06:37.322+01:00 level=WARN msg=world app=app version=d1da844711730f2f5cbd08be93e62e71475f7d4e trace.id=0a02e16c-7418-4558-9dcc-718c007162b6
```

### Monitor the time elapsed by a function on `defer`.

```go
var (
    log = logm.DefaultLogger("app", os.Stdout)
    ctx = context.Background()
)
func(ctx context.Context, log *slog.Logger) {
    defer logm.TimeElapsed(ctx, log, slog.LevelInfo, "example")()
    time.Sleep(time.Millisecond)
}(ctx, log)
```
```bash
time=2023-03-25T12:06:17.605+01:00 level=INFO msg=example app=app version=d1da844711730f2f5cbd08be93e62e71475f7d4e trace.id=ccc05db1-68d2-4442-9353-0789e0b8ca55 trace.time_elapsed_ms=1
```