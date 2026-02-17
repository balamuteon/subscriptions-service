package logger

// NoopLogger is a logger implementation that discards all logs.
type NoopLogger struct{}

func (NoopLogger) Debug(msg string, args ...any) {}
func (NoopLogger) Info(msg string, args ...any)  {}
func (NoopLogger) Warn(msg string, args ...any)  {}
func (NoopLogger) Error(msg string, args ...any) {}
func (NoopLogger) With(args ...any) Logger       { return NoopLogger{} }

// NewNoop returns a Logger that does nothing.
func NewNoop() Logger {
	return NoopLogger{}
}
