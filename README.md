# LogM

[![GoDoc](https://godoc.org/github.com/rvflash/logm?status.svg)](https://godoc.org/github.com/rvflash/logm)
[![Build Status](https://github.com/rvflash/logm/workflows/build/badge.svg)](https://github.com/rvflash/logm/actions?workflow=build)
[![Code Coverage](https://codecov.io/gh/rvflash/logm/branch/master/graph/badge.svg)](https://codecov.io/gh/rvflash/logm)
[![Go Report Card](https://goreportcard.com/badge/github.com/rvflash/logm?)](https://goreportcard.com/report/github.com/rvflash/logm)


LogM is a Go log management providing facilities to deal with logs and trace context.

## Features

1. Provides simple methods to expose a [slog.Logger](https://pkg.go.dev/golang.org/x/exp/slog) ready for production, another to debug or even to discard all logs. 
2. Provides a log file with automatic rotating, maximum file size, zip archives, etc. thanks to [lumberjack](https://github.com/natefinch/lumberjack).
3. Exposes an HTTP middleware to create a trace context on each request. 
4. Exposes a middleware to log details about each HTTP request and response.
5. Supplies an HTTP middleware to recover on panic, log the event in error, the stack trace in debug and write an internal server error as response header.
6. Provide a method to log in defer the time elapsed of a function. 


### Installation

```bash
$ go get -u github.com/rvflash/logm
```