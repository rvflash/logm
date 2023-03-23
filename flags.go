package logm

// List of standard key to use to qualify structured log.
const (
	// AppNameKey is the name of the application in structured log.
	AppNameKey = "app"
	// AppVersionKey is the version of the application in structured log.
	AppVersionKey = "version"
	// HTTPRequestKey is the HTTP request name in structured log.
	HTTPRequestKey = "req"
	// HTTPPathKey is the HTTP request path in structured log.
	HTTPPathKey = "path"
	// HTTPMethodKey is the HTTP request method in structured log.
	HTTPMethodKey = "method"
	// HTTPRemoteAddrKey is the HTTP remote address in structured log.
	HTTPRemoteAddrKey = "remote_addr"
	// HTTPQueryKey is the HTTP request query in structured log.
	HTTPQueryKey = "query"
	// HTTPResponseKey is the HTTP response name in structured log.
	HTTPResponseKey = "resp"
	// HTTPStatusKey is the HTTP response status in structured log.
	HTTPStatusKey = "status"
	// HTTPSizeKey is the HTTP response size in structured log.
	HTTPSizeKey = "size"
	// PanicKey is a panic in structured log.
	PanicKey = "panic"
	// TraceKey is the name of the trace in structured log.
	TraceKey = "trace"
	// TraceIDKey is the name of the trace ID in structured log.
	TraceIDKey = "id"
	// TraceSpanIDKey is the name of the trace span ID in structured log.
	TraceSpanIDKey = "span_id"
	// TraceTimeElapsedKey is the name of the trace time in structured log.
	TraceTimeElapsedKey = "time_elapsed_ms"
)
