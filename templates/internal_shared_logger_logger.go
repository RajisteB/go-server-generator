package logger

import (
	"context"
	"io"
	"log/slog"
	"os"
	"runtime"
	"strings"
	"time"
)

type Logger struct {
	*slog.Logger
	config *Config
}

// Info logs an info message
func (l *Logger) Info(msg string, args ...interface{}) {
	l.Logger.Info(msg, args...)
}

// Error logs an error message
func (l *Logger) Error(msg string, args ...interface{}) {
	l.Logger.Error(msg, args...)
}

// Debug logs a debug message
func (l *Logger) Debug(msg string, args ...interface{}) {
	l.Logger.Debug(msg, args...)
}

// Warn logs a warning message
func (l *Logger) Warn(msg string, args ...interface{}) {
	l.Logger.Warn(msg, args...)
}

// With creates a new logger with additional context
func (l *Logger) With(args ...interface{}) *Logger {
	return &Logger{
		Logger: l.Logger.With(args...),
		config: l.config,
	}
}

// withSlog creates a new logger with additional context and returns the underlying slog.Logger
func (l *Logger) withSlog(args ...interface{}) *slog.Logger {
	return l.Logger.With(args...)
}

type Config struct {
	Level     slog.Level
	AddSource bool
	Service   string
	Version   string
	Writer    io.Writer
}

func DefaultConfig() *Config {
	return &Config{
		Level:     slog.LevelDebug,
		AddSource: true,
		Service:   "localhost",
		Version:   "1.0.0",
		Writer:    os.Stdout,
	}
}

func DevelopmentConfig(service string, version string) *Config {
	return &Config{
		Level:     slog.LevelDebug,
		AddSource: false,
		Service:   service,
		Version:   version,
		Writer:    os.Stdout,
	}
}

var defaultLogger *Logger

func Initialize(config *Config) {
	if config == nil {
		config = DefaultConfig()
	}

	opts := &slog.HandlerOptions{
		Level:     config.Level,
		AddSource: config.AddSource,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {

			if a.Key == slog.SourceKey {
				source := a.Value.Any().(*slog.Source)

				parts := strings.Split(source.File, "/")
				source.File = parts[len(parts)-1]
			}

			if a.Key == slog.TimeKey {
				return slog.Attr{
					Key:   "timestamp",
					Value: slog.StringValue(a.Value.Time().Format(time.RFC3339)),
				}
			}
			return a
		},
	}

	handler := slog.NewJSONHandler(config.Writer, opts)

	logger := slog.New(handler).With(
		"service", config.Service,
		"version", config.Version,
	)

	defaultLogger = &Logger{
		Logger: logger,
		config: config,
	}

	slog.SetDefault(logger)
}

func GetLogger() *Logger {
	if defaultLogger == nil {
		Initialize(DefaultConfig())
	}
	return defaultLogger
}

func NewLogger(config *Config) *Logger {
	if config == nil {
		config = DefaultConfig()
	}

	opts := &slog.HandlerOptions{
		Level:     config.Level,
		AddSource: config.AddSource,
	}

	handler := slog.NewJSONHandler(config.Writer, opts)
	logger := slog.New(handler).With(
		"service", config.Service,
		"version", config.Version,
	)

	return &Logger{
		Logger: logger,
		config: config,
	}
}

type contextKey string

const (
	TraceIDKey   contextKey = "trace_id"
	RequestIDKey contextKey = "request_id"
	UserIDKey    contextKey = "user_id"
	SessionIDKey contextKey = "session_id"
)

func (l *Logger) WithTraceID(traceID string) *Logger {
	return &Logger{
		Logger: l.withSlog("trace_id", traceID),
		config: l.config,
	}
}

func (l *Logger) WithRequestID(requestID string) *Logger {
	return &Logger{
		Logger: l.withSlog("request_id", requestID),
		config: l.config,
	}
}

func (l *Logger) WithUserID(userID interface{}) *Logger {
	return &Logger{
		Logger: l.withSlog("user_id", userID),
		config: l.config,
	}
}

func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	args := make([]interface{}, 0, len(fields)*2)
	for k, v := range fields {
		args = append(args, k, v)
	}
	return &Logger{
		Logger: l.withSlog(args...),
		config: l.config,
	}
}

func (l *Logger) WithComponent(component string) *Logger {
	return &Logger{
		Logger: l.withSlog("component", component),
		config: l.config,
	}
}

func (l *Logger) WithContext(ctx context.Context) *Logger {
	logger := l.Logger

	if traceID, ok := ctx.Value(TraceIDKey).(string); ok && traceID != "" {
		logger = logger.With("trace_id", traceID)
	}

	if requestID, ok := ctx.Value(RequestIDKey).(string); ok && requestID != "" {
		logger = logger.With("request_id", requestID)
	}

	if userID := ctx.Value(UserIDKey); userID != nil {
		logger = logger.With("user_id", userID)
	}

	if sessionID, ok := ctx.Value(SessionIDKey).(string); ok && sessionID != "" {
		logger = logger.With("session_id", sessionID)
	}

	return &Logger{
		Logger: logger,
		config: l.config,
	}
}

func (l *Logger) LogHTTPRequest(method, path, userAgent, clientIP string, contentType string) {
	l.Info("HTTP request",
		"method", method,
		"path", path,
		"user_agent", userAgent,
		"client_ip", clientIP,
		"content_type", contentType,
	)
}

func (l *Logger) LogDBOperation(operation, table string, duration time.Duration, rowsAffected int64, err error) {
	fields := []interface{}{
		"operation", operation,
		"table", table,
		"duration_ms", duration.Milliseconds(),
		"rows_affected", rowsAffected,
	}

	if err != nil {
		fields = append(fields, "error", err.Error())
		l.Error("Database operation failed", fields...)
	} else {
		l.Debug("Database operation completed", fields...)
	}
}

func (l *Logger) LogAPICall(service, endpoint, method string, statusCode int, duration time.Duration, err error) {
	fields := []interface{}{
		"external_service", service,
		"endpoint", endpoint,
		"method", method,
		"status_code", statusCode,
		"duration_ms", duration.Milliseconds(),
	}

	if err != nil {
		fields = append(fields, "error", err.Error())
		l.Error("External API call failed", fields...)
	} else {
		l.Info("External API call completed", fields...)
	}
}

func (l *Logger) LogBusinessEvent(event string, fields map[string]interface{}) {
	args := []interface{}{"event", event}
	for k, v := range fields {
		args = append(args, k, v)
	}
	l.Info("Business event", args...)
}

func (l *Logger) LogSecurityEvent(event, reason string, severity string, fields map[string]interface{}) {
	args := []interface{}{
		"security_event", event,
		"reason", reason,
		"severity", severity,
	}
	for k, v := range fields {
		args = append(args, k, v)
	}
	l.Warn("Security event", args...)
}

func (l *Logger) LogPerformance(operation string, duration time.Duration, fields map[string]interface{}) {
	args := []interface{}{
		"performance_metric", operation,
		"duration_ms", duration.Milliseconds(),
	}
	for k, v := range fields {
		args = append(args, k, v)
	}

	if duration > 5*time.Second {
		l.Warn("Slow operation detected", args...)
	} else {
		l.Debug("Performance metric", args...)
	}
}

func (l *Logger) ErrorWithStack(msg string, err error, fields ...interface{}) {
	args := []interface{}{"error", err.Error()}
	args = append(args, fields...)

	if l.config.Level <= slog.LevelDebug {
		stack := make([]byte, 4096)
		length := runtime.Stack(stack, false)
		args = append(args, "stack_trace", string(stack[:length]))
	}

	l.Error(msg, args...)
}

func (l *Logger) LogPanic(recovered interface{}, fields ...interface{}) {
	args := []interface{}{"panic", recovered}
	args = append(args, fields...)

	stack := make([]byte, 4096)
	length := runtime.Stack(stack, false)
	args = append(args, "stack_trace", string(stack[:length]))

	l.Error("Panic recovered", args...)
}

func (l *Logger) TimeOperation(operation string, fields ...interface{}) func() {
	start := time.Now()
	l.Debug("Operation started", append([]interface{}{"operation", operation}, fields...)...)

	return func() {
		duration := time.Since(start)
		l.LogPerformance(operation, duration, nil)
	}
}

func (l *Logger) InfoIf(condition bool, msg string, fields ...interface{}) {
	if condition {
		l.Info(msg, fields...)
	}
}

func (l *Logger) WarnIf(condition bool, msg string, fields ...interface{}) {
	if condition {
		l.Warn(msg, fields...)
	}
}

func (l *Logger) ErrorIf(condition bool, msg string, fields ...interface{}) {
	if condition {
		l.Error(msg, fields...)
	}
}

type SamplingLogger struct {
	*Logger
	sampleRate int
	counter    int
}

func (l *Logger) NewSamplingLogger(sampleRate int) *SamplingLogger {
	return &SamplingLogger{
		Logger:     l,
		sampleRate: sampleRate,
		counter:    0,
	}
}

func (sl *SamplingLogger) Info(msg string, fields ...interface{}) {
	sl.counter++
	if sl.counter%sl.sampleRate == 0 {
		sl.Logger.Info(msg, fields...)
	}
}

func Debug(msg string, fields ...interface{}) {
	GetLogger().Debug(msg, fields...)
}

func Info(msg string, fields ...interface{}) {
	GetLogger().Info(msg, fields...)
}

func Warn(msg string, fields ...interface{}) {
	GetLogger().Warn(msg, fields...)
}

func Error(msg string, fields ...interface{}) {
	GetLogger().Error(msg, fields...)
}

func WithTraceID(traceID string) *Logger {
	return GetLogger().WithTraceID(traceID)
}

func WithRequestID(requestID string) *Logger {
	return GetLogger().WithRequestID(requestID)
}

func WithContext(ctx context.Context) *Logger {
	return GetLogger().WithContext(ctx)
}

func WithFields(fields map[string]interface{}) *Logger {
	return GetLogger().WithFields(fields)
}

func LogHTTPRequest(method, path, userAgent, clientIP string, contentType string) {
	GetLogger().LogHTTPRequest(method, path, userAgent, clientIP, contentType)
}

func LogBusinessEvent(event string, fields map[string]interface{}) {
	GetLogger().LogBusinessEvent(event, fields)
}

func TimeOperation(operation string, fields ...interface{}) func() {
	return GetLogger().TimeOperation(operation, fields...)
}

func ErrorWithStack(msg string, err error, fields ...interface{}) {
	GetLogger().ErrorWithStack(msg, err, fields...)
}

func ContextWithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, TraceIDKey, traceID)
}

func ContextWithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, RequestIDKey, requestID)
}

func ContextWithUserID(ctx context.Context, userID interface{}) context.Context {
	return context.WithValue(ctx, UserIDKey, userID)
}

type ErrorDetails struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

func (l *Logger) LogError(err error, details ErrorDetails, fields ...interface{}) {
	args := []interface{}{
		"error", err.Error(),
		"error_code", details.Code,
		"error_message", details.Message,
	}

	if details.Details != nil {
		args = append(args, "error_details", details.Details)
	}

	args = append(args, fields...)
	l.Error("Structured error", args...)
}

func (l *Logger) LogHealthCheck(service string, status string, duration time.Duration, details map[string]interface{}) {
	args := []interface{}{
		"health_check", service,
		"status", status,
		"duration_ms", duration.Milliseconds(),
	}

	for k, v := range details {
		args = append(args, k, v)
	}

	if status == "healthy" {
		l.Debug("Health check passed", args...)
	} else {
		l.Warn("Health check failed", args...)
	}
}

func (l *Logger) LogAudit(action, resource string, userID interface{}, result string, fields map[string]interface{}) {
	args := []interface{}{
		"audit_action", action,
		"resource", resource,
		"user_id", userID,
		"result", result,
		"timestamp", time.Now().UTC().Format(time.RFC3339),
	}

	for k, v := range fields {
		args = append(args, k, v)
	}

	l.Info("Audit event", args...)
}
