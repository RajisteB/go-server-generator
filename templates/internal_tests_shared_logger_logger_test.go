package logger_test

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"strings"
	"testing"
	"time"

	"{{.Module}}/internal/shared/logger"
)

func TestLogger_Info(t *testing.T) {
	var buf bytes.Buffer
	config := &logger.Config{
		Level:     slog.LevelInfo,
		AddSource: false,
		Service:   "test-service",
		Version:   "1.0.0",
		Writer:    &buf,
	}

	log := logger.NewLogger(config)
	log.Info("test message", "key", "value")

	var logEntry map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &logEntry)
	if err != nil {
		t.Fatalf("Failed to unmarshal log entry: %v", err)
	}

	if logEntry["msg"] != "test message" {
		t.Errorf("Expected message 'test message', got %v", logEntry["msg"])
	}

	if logEntry["key"] != "value" {
		t.Errorf("Expected key 'value', got %v", logEntry["key"])
	}

	if logEntry["service"] != "test-service" {
		t.Errorf("Expected service 'test-service', got %v", logEntry["service"])
	}

	if logEntry["version"] != "1.0.0" {
		t.Errorf("Expected version '1.0.0', got %v", logEntry["version"])
	}
}

func TestLogger_Error(t *testing.T) {
	var buf bytes.Buffer
	config := &logger.Config{
		Level:     slog.LevelError,
		AddSource: false,
		Service:   "test-service",
		Version:   "1.0.0",
		Writer:    &buf,
	}

	log := logger.NewLogger(config)
	log.Error("error message", "error", "test error")

	var logEntry map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &logEntry)
	if err != nil {
		t.Fatalf("Failed to unmarshal log entry: %v", err)
	}

	if logEntry["level"] != "ERROR" {
		t.Errorf("Expected level 'ERROR', got %v", logEntry["level"])
	}

	if logEntry["msg"] != "error message" {
		t.Errorf("Expected message 'error message', got %v", logEntry["msg"])
	}
}

func TestLogger_Debug(t *testing.T) {
	var buf bytes.Buffer
	config := &logger.Config{
		Level:     slog.LevelDebug,
		AddSource: false,
		Service:   "test-service",
		Version:   "1.0.0",
		Writer:    &buf,
	}

	log := logger.NewLogger(config)
	log.Debug("debug message", "debug", "info")

	var logEntry map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &logEntry)
	if err != nil {
		t.Fatalf("Failed to unmarshal log entry: %v", err)
	}

	if logEntry["level"] != "DEBUG" {
		t.Errorf("Expected level 'DEBUG', got %v", logEntry["level"])
	}
}

func TestLogger_Warn(t *testing.T) {
	var buf bytes.Buffer
	config := &logger.Config{
		Level:     slog.LevelWarn,
		AddSource: false,
		Service:   "test-service",
		Version:   "1.0.0",
		Writer:    &buf,
	}

	log := logger.NewLogger(config)
	log.Warn("warning message", "warning", "info")

	var logEntry map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &logEntry)
	if err != nil {
		t.Fatalf("Failed to unmarshal log entry: %v", err)
	}

	if logEntry["level"] != "WARN" {
		t.Errorf("Expected level 'WARN', got %v", logEntry["level"])
	}
}

func TestLogger_With(t *testing.T) {
	var buf bytes.Buffer
	config := &logger.Config{
		Level:     slog.LevelInfo,
		AddSource: false,
		Service:   "test-service",
		Version:   "1.0.0",
		Writer:    &buf,
	}

	log := logger.NewLogger(config)
	childLog := log.With("parent", "value")
	childLog.Info("child message", "child", "value")

	var logEntry map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &logEntry)
	if err != nil {
		t.Fatalf("Failed to unmarshal log entry: %v", err)
	}

	if logEntry["parent"] != "value" {
		t.Errorf("Expected parent 'value', got %v", logEntry["parent"])
	}

	if logEntry["child"] != "value" {
		t.Errorf("Expected child 'value', got %v", logEntry["child"])
	}
}

func TestLogger_WithTraceID(t *testing.T) {
	var buf bytes.Buffer
	config := &logger.Config{
		Level:     slog.LevelInfo,
		AddSource: false,
		Service:   "test-service",
		Version:   "1.0.0",
		Writer:    &buf,
	}

	log := logger.NewLogger(config)
	traceLog := log.WithTraceID("trace-123")
	traceLog.Info("trace message")

	var logEntry map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &logEntry)
	if err != nil {
		t.Fatalf("Failed to unmarshal log entry: %v", err)
	}

	if logEntry["trace_id"] != "trace-123" {
		t.Errorf("Expected trace_id 'trace-123', got %v", logEntry["trace_id"])
	}
}

func TestLogger_WithRequestID(t *testing.T) {
	var buf bytes.Buffer
	config := &logger.Config{
		Level:     slog.LevelInfo,
		AddSource: false,
		Service:   "test-service",
		Version:   "1.0.0",
		Writer:    &buf,
	}

	log := logger.NewLogger(config)
	requestLog := log.WithRequestID("req-456")
	requestLog.Info("request message")

	var logEntry map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &logEntry)
	if err != nil {
		t.Fatalf("Failed to unmarshal log entry: %v", err)
	}

	if logEntry["request_id"] != "req-456" {
		t.Errorf("Expected request_id 'req-456', got %v", logEntry["request_id"])
	}
}

func TestLogger_WithUserID(t *testing.T) {
	var buf bytes.Buffer
	config := &logger.Config{
		Level:     slog.LevelInfo,
		AddSource: false,
		Service:   "test-service",
		Version:   "1.0.0",
		Writer:    &buf,
	}

	log := logger.NewLogger(config)
	userLog := log.WithUserID("user-789")
	userLog.Info("user message")

	var logEntry map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &logEntry)
	if err != nil {
		t.Fatalf("Failed to unmarshal log entry: %v", err)
	}

	if logEntry["user_id"] != "user-789" {
		t.Errorf("Expected user_id 'user-789', got %v", logEntry["user_id"])
	}
}

func TestLogger_WithFields(t *testing.T) {
	var buf bytes.Buffer
	config := &logger.Config{
		Level:     slog.LevelInfo,
		AddSource: false,
		Service:   "test-service",
		Version:   "1.0.0",
		Writer:    &buf,
	}

	log := logger.NewLogger(config)
	fields := map[string]interface{}{
		"field1": "value1",
		"field2": "value2",
	}
	fieldsLog := log.WithFields(fields)
	fieldsLog.Info("fields message")

	var logEntry map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &logEntry)
	if err != nil {
		t.Fatalf("Failed to unmarshal log entry: %v", err)
	}

	if logEntry["field1"] != "value1" {
		t.Errorf("Expected field1 'value1', got %v", logEntry["field1"])
	}

	if logEntry["field2"] != "value2" {
		t.Errorf("Expected field2 'value2', got %v", logEntry["field2"])
	}
}

func TestLogger_WithContext(t *testing.T) {
	var buf bytes.Buffer
	config := &logger.Config{
		Level:     slog.LevelInfo,
		AddSource: false,
		Service:   "test-service",
		Version:   "1.0.0",
		Writer:    &buf,
	}

	log := logger.NewLogger(config)
	ctx := context.Background()
	ctx = logger.ContextWithTraceID(ctx, "trace-ctx")
	ctx = logger.ContextWithRequestID(ctx, "req-ctx")
	ctx = logger.ContextWithUserID(ctx, "user-ctx")

	contextLog := log.WithContext(ctx)
	contextLog.Info("context message")

	var logEntry map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &logEntry)
	if err != nil {
		t.Fatalf("Failed to unmarshal log entry: %v", err)
	}

	if logEntry["trace_id"] != "trace-ctx" {
		t.Errorf("Expected trace_id 'trace-ctx', got %v", logEntry["trace_id"])
	}

	if logEntry["request_id"] != "req-ctx" {
		t.Errorf("Expected request_id 'req-ctx', got %v", logEntry["request_id"])
	}

	if logEntry["user_id"] != "user-ctx" {
		t.Errorf("Expected user_id 'user-ctx', got %v", logEntry["user_id"])
	}
}

func TestLogger_LogHTTPRequest(t *testing.T) {
	var buf bytes.Buffer
	config := &logger.Config{
		Level:     slog.LevelInfo,
		AddSource: false,
		Service:   "test-service",
		Version:   "1.0.0",
		Writer:    &buf,
	}

	log := logger.NewLogger(config)
	log.LogHTTPRequest("GET", "/api/users", "Mozilla/5.0", "192.168.1.1", "application/json")

	var logEntry map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &logEntry)
	if err != nil {
		t.Fatalf("Failed to unmarshal log entry: %v", err)
	}

	if logEntry["method"] != "GET" {
		t.Errorf("Expected method 'GET', got %v", logEntry["method"])
	}

	if logEntry["path"] != "/api/users" {
		t.Errorf("Expected path '/api/users', got %v", logEntry["path"])
	}

	if logEntry["user_agent"] != "Mozilla/5.0" {
		t.Errorf("Expected user_agent 'Mozilla/5.0', got %v", logEntry["user_agent"])
	}

	if logEntry["client_ip"] != "192.168.1.1" {
		t.Errorf("Expected client_ip '192.168.1.1', got %v", logEntry["client_ip"])
	}

	if logEntry["content_type"] != "application/json" {
		t.Errorf("Expected content_type 'application/json', got %v", logEntry["content_type"])
	}
}

func TestLogger_LogDBOperation(t *testing.T) {
	var buf bytes.Buffer
	config := &logger.Config{
		Level:     slog.LevelDebug,
		AddSource: false,
		Service:   "test-service",
		Version:   "1.0.0",
		Writer:    &buf,
	}

	log := logger.NewLogger(config)
	log.LogDBOperation("SELECT", "users", 100*time.Millisecond, 5, nil)

	var logEntry map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &logEntry)
	if err != nil {
		t.Fatalf("Failed to unmarshal log entry: %v", err)
	}

	if logEntry["operation"] != "SELECT" {
		t.Errorf("Expected operation 'SELECT', got %v", logEntry["operation"])
	}

	if logEntry["table"] != "users" {
		t.Errorf("Expected table 'users', got %v", logEntry["table"])
	}

	if logEntry["duration_ms"] != float64(100) {
		t.Errorf("Expected duration_ms 100, got %v", logEntry["duration_ms"])
	}

	if logEntry["rows_affected"] != float64(5) {
		t.Errorf("Expected rows_affected 5, got %v", logEntry["rows_affected"])
	}
}

func TestLogger_LogDBOperation_WithError(t *testing.T) {
	var buf bytes.Buffer
	config := &logger.Config{
		Level:     slog.LevelError,
		AddSource: false,
		Service:   "test-service",
		Version:   "1.0.0",
		Writer:    &buf,
	}

	log := logger.NewLogger(config)
	log.LogDBOperation("INSERT", "users", 50*time.Millisecond, 0, &TestError{"database error"})

	var logEntry map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &logEntry)
	if err != nil {
		t.Fatalf("Failed to unmarshal log entry: %v", err)
	}

	if logEntry["level"] != "ERROR" {
		t.Errorf("Expected level 'ERROR', got %v", logEntry["level"])
	}

	if logEntry["error"] != "database error" {
		t.Errorf("Expected error 'database error', got %v", logEntry["error"])
	}
}

func TestLogger_LogAPICall(t *testing.T) {
	var buf bytes.Buffer
	config := &logger.Config{
		Level:     slog.LevelInfo,
		AddSource: false,
		Service:   "test-service",
		Version:   "1.0.0",
		Writer:    &buf,
	}

	log := logger.NewLogger(config)
	log.LogAPICall("external-api", "/users", "GET", 200, 150*time.Millisecond, nil)

	var logEntry map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &logEntry)
	if err != nil {
		t.Fatalf("Failed to unmarshal log entry: %v", err)
	}

	if logEntry["external_service"] != "external-api" {
		t.Errorf("Expected external_service 'external-api', got %v", logEntry["external_service"])
	}

	if logEntry["endpoint"] != "/users" {
		t.Errorf("Expected endpoint '/users', got %v", logEntry["endpoint"])
	}

	if logEntry["method"] != "GET" {
		t.Errorf("Expected method 'GET', got %v", logEntry["method"])
	}

	if logEntry["status_code"] != float64(200) {
		t.Errorf("Expected status_code 200, got %v", logEntry["status_code"])
	}

	if logEntry["duration_ms"] != float64(150) {
		t.Errorf("Expected duration_ms 150, got %v", logEntry["duration_ms"])
	}
}

func TestLogger_LogBusinessEvent(t *testing.T) {
	var buf bytes.Buffer
	config := &logger.Config{
		Level:     slog.LevelInfo,
		AddSource: false,
		Service:   "test-service",
		Version:   "1.0.0",
		Writer:    &buf,
	}

	log := logger.NewLogger(config)
	fields := map[string]interface{}{
		"user_id": "123",
		"action":  "login",
	}
	log.LogBusinessEvent("user_login", fields)

	var logEntry map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &logEntry)
	if err != nil {
		t.Fatalf("Failed to unmarshal log entry: %v", err)
	}

	if logEntry["event"] != "user_login" {
		t.Errorf("Expected event 'user_login', got %v", logEntry["event"])
	}

	if logEntry["user_id"] != "123" {
		t.Errorf("Expected user_id '123', got %v", logEntry["user_id"])
	}

	if logEntry["action"] != "login" {
		t.Errorf("Expected action 'login', got %v", logEntry["action"])
	}
}

func TestLogger_LogSecurityEvent(t *testing.T) {
	var buf bytes.Buffer
	config := &logger.Config{
		Level:     slog.LevelWarn,
		AddSource: false,
		Service:   "test-service",
		Version:   "1.0.0",
		Writer:    &buf,
	}

	log := logger.NewLogger(config)
	fields := map[string]interface{}{
		"ip_address": "192.168.1.100",
		"attempts":   3,
	}
	log.LogSecurityEvent("failed_login", "too many attempts", "high", fields)

	var logEntry map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &logEntry)
	if err != nil {
		t.Fatalf("Failed to unmarshal log entry: %v", err)
	}

	if logEntry["level"] != "WARN" {
		t.Errorf("Expected level 'WARN', got %v", logEntry["level"])
	}

	if logEntry["security_event"] != "failed_login" {
		t.Errorf("Expected security_event 'failed_login', got %v", logEntry["security_event"])
	}

	if logEntry["reason"] != "too many attempts" {
		t.Errorf("Expected reason 'too many attempts', got %v", logEntry["reason"])
	}

	if logEntry["severity"] != "high" {
		t.Errorf("Expected severity 'high', got %v", logEntry["severity"])
	}
}

func TestLogger_LogPerformance(t *testing.T) {
	var buf bytes.Buffer
	config := &logger.Config{
		Level:     slog.LevelDebug,
		AddSource: false,
		Service:   "test-service",
		Version:   "1.0.0",
		Writer:    &buf,
	}

	log := logger.NewLogger(config)
	fields := map[string]interface{}{
		"operation_type": "database_query",
		"table":          "users",
	}
	log.LogPerformance("slow_query", 6*time.Second, fields)

	var logEntry map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &logEntry)
	if err != nil {
		t.Fatalf("Failed to unmarshal log entry: %v", err)
	}

	if logEntry["level"] != "WARN" {
		t.Errorf("Expected level 'WARN' for slow operation, got %v", logEntry["level"])
	}

	if logEntry["performance_metric"] != "slow_query" {
		t.Errorf("Expected performance_metric 'slow_query', got %v", logEntry["performance_metric"])
	}

	if logEntry["duration_ms"] != float64(6000) {
		t.Errorf("Expected duration_ms 6000, got %v", logEntry["duration_ms"])
	}
}

func TestLogger_TimeOperation(t *testing.T) {
	var buf bytes.Buffer
	config := &logger.Config{
		Level:     slog.LevelDebug,
		AddSource: false,
		Service:   "test-service",
		Version:   "1.0.0",
		Writer:    &buf,
	}

	log := logger.NewLogger(config)
	timer := log.TimeOperation("test_operation", "param", "value")
	time.Sleep(10 * time.Millisecond) // Simulate work
	timer()

	// Should have two log entries: start and end
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("Expected 2 log entries, got %d", len(lines))
	}

	// Check start log
	var startEntry map[string]interface{}
	err := json.Unmarshal([]byte(lines[0]), &startEntry)
	if err != nil {
		t.Fatalf("Failed to unmarshal start log entry: %v", err)
	}

	if startEntry["operation"] != "test_operation" {
		t.Errorf("Expected operation 'test_operation', got %v", startEntry["operation"])
	}

	// Check end log
	var endEntry map[string]interface{}
	err = json.Unmarshal([]byte(lines[1]), &endEntry)
	if err != nil {
		t.Fatalf("Failed to unmarshal end log entry: %v", err)
	}

	if endEntry["performance_metric"] != "test_operation" {
		t.Errorf("Expected performance_metric 'test_operation', got %v", endEntry["performance_metric"])
	}
}

func TestLogger_InfoIf(t *testing.T) {
	var buf bytes.Buffer
	config := &logger.Config{
		Level:     slog.LevelInfo,
		AddSource: false,
		Service:   "test-service",
		Version:   "1.0.0",
		Writer:    &buf,
	}

	log := logger.NewLogger(config)

	// Test with true condition
	log.InfoIf(true, "conditional message", "key", "value")
	if buf.Len() == 0 {
		t.Error("Expected log entry when condition is true")
	}

	buf.Reset()

	// Test with false condition
	log.InfoIf(false, "conditional message", "key", "value")
	if buf.Len() != 0 {
		t.Error("Expected no log entry when condition is false")
	}
}

func TestLogger_SamplingLogger(t *testing.T) {
	var buf bytes.Buffer
	config := &logger.Config{
		Level:     slog.LevelInfo,
		AddSource: false,
		Service:   "test-service",
		Version:   "1.0.0",
		Writer:    &buf,
	}

	log := logger.NewLogger(config)
	samplingLog := log.NewSamplingLogger(3) // Sample every 3rd message

	// Log 5 messages
	for i := 0; i < 5; i++ {
		samplingLog.Info("sampled message", "index", i)
	}

	// Should have 1 log entry (3rd message only)
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 1 {
		t.Fatalf("Expected 1 log entry with sampling rate 3, got %d", len(lines))
	}
}

func TestLogger_LogError(t *testing.T) {
	var buf bytes.Buffer
	config := &logger.Config{
		Level:     slog.LevelError,
		AddSource: false,
		Service:   "test-service",
		Version:   "1.0.0",
		Writer:    &buf,
	}

	log := logger.NewLogger(config)
	details := logger.ErrorDetails{
		Code:    "VALIDATION_ERROR",
		Message: "Invalid input provided",
		Details: map[string]interface{}{
			"field": "email",
			"value": "invalid-email",
		},
	}

	log.LogError(&TestError{"validation failed"}, details, "request_id", "req-123")

	var logEntry map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &logEntry)
	if err != nil {
		t.Fatalf("Failed to unmarshal log entry: %v", err)
	}

	if logEntry["error"] != "validation failed" {
		t.Errorf("Expected error 'validation failed', got %v", logEntry["error"])
	}

	if logEntry["error_code"] != "VALIDATION_ERROR" {
		t.Errorf("Expected error_code 'VALIDATION_ERROR', got %v", logEntry["error_code"])
	}

	if logEntry["error_message"] != "Invalid input provided" {
		t.Errorf("Expected error_message 'Invalid input provided', got %v", logEntry["error_message"])
	}

	if logEntry["request_id"] != "req-123" {
		t.Errorf("Expected request_id 'req-123', got %v", logEntry["request_id"])
	}
}

func TestLogger_LogHealthCheck(t *testing.T) {
	var buf bytes.Buffer
	config := &logger.Config{
		Level:     slog.LevelDebug,
		AddSource: false,
		Service:   "test-service",
		Version:   "1.0.0",
		Writer:    &buf,
	}

	log := logger.NewLogger(config)
	details := map[string]interface{}{
		"database": "connected",
		"redis":    "connected",
	}

	log.LogHealthCheck("database", "healthy", 50*time.Millisecond, details)

	var logEntry map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &logEntry)
	if err != nil {
		t.Fatalf("Failed to unmarshal log entry: %v", err)
	}

	if logEntry["health_check"] != "database" {
		t.Errorf("Expected health_check 'database', got %v", logEntry["health_check"])
	}

	if logEntry["status"] != "healthy" {
		t.Errorf("Expected status 'healthy', got %v", logEntry["status"])
	}

	if logEntry["duration_ms"] != float64(50) {
		t.Errorf("Expected duration_ms 50, got %v", logEntry["duration_ms"])
	}
}

func TestLogger_LogAudit(t *testing.T) {
	var buf bytes.Buffer
	config := &logger.Config{
		Level:     slog.LevelInfo,
		AddSource: false,
		Service:   "test-service",
		Version:   "1.0.0",
		Writer:    &buf,
	}

	log := logger.NewLogger(config)
	fields := map[string]interface{}{
		"ip_address": "192.168.1.1",
		"user_agent": "Mozilla/5.0",
	}

	log.LogAudit("CREATE", "user", "user-123", "success", fields)

	var logEntry map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &logEntry)
	if err != nil {
		t.Fatalf("Failed to unmarshal log entry: %v", err)
	}

	if logEntry["audit_action"] != "CREATE" {
		t.Errorf("Expected audit_action 'CREATE', got %v", logEntry["audit_action"])
	}

	if logEntry["resource"] != "user" {
		t.Errorf("Expected resource 'user', got %v", logEntry["resource"])
	}

	if logEntry["user_id"] != "user-123" {
		t.Errorf("Expected user_id 'user-123', got %v", logEntry["user_id"])
	}

	if logEntry["result"] != "success" {
		t.Errorf("Expected result 'success', got %v", logEntry["result"])
	}
}

func TestDefaultConfig(t *testing.T) {
	config := logger.DefaultConfig()

	if config.Level != slog.LevelDebug {
		t.Errorf("Expected default level DEBUG, got %v", config.Level)
	}

	if !config.AddSource {
		t.Error("Expected AddSource to be true by default")
	}

	if config.Service != "localhost" {
		t.Errorf("Expected default service 'localhost', got %v", config.Service)
	}

	if config.Version != "1.0.0" {
		t.Errorf("Expected default version '1.0.0', got %v", config.Version)
	}
}

func TestDevelopmentConfig(t *testing.T) {
	config := logger.DevelopmentConfig("test-service", "2.0.0")

	if config.Level != slog.LevelDebug {
		t.Errorf("Expected development level DEBUG, got %v", config.Level)
	}

	if config.AddSource {
		t.Error("Expected AddSource to be false in development")
	}

	if config.Service != "test-service" {
		t.Errorf("Expected service 'test-service', got %v", config.Service)
	}

	if config.Version != "2.0.0" {
		t.Errorf("Expected version '2.0.0', got %v", config.Version)
	}
}

func TestGetLogger(t *testing.T) {
	// Initialize with custom config
	var buf bytes.Buffer
	config := &logger.Config{
		Level:     slog.LevelInfo,
		AddSource: false,
		Service:   "test-service",
		Version:   "1.0.0",
		Writer:    &buf,
	}

	logger.Initialize(config)
	log := logger.GetLogger()

	if log == nil {
		t.Error("Expected GetLogger to return a logger instance")
	}

	log.Info("test message")
	if buf.Len() == 0 {
		t.Error("Expected log output from GetLogger")
	}
}

// Helper types for testing
type TestError struct {
	message string
}

func (e *TestError) Error() string {
	return e.message
}

// Benchmark tests
func BenchmarkLogger_Info(b *testing.B) {
	var buf bytes.Buffer
	config := &logger.Config{
		Level:     slog.LevelInfo,
		AddSource: false,
		Service:   "test-service",
		Version:   "1.0.0",
		Writer:    &buf,
	}

	log := logger.NewLogger(config)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		log.Info("benchmark message", "iteration", i)
	}
}

func BenchmarkLogger_WithFields(b *testing.B) {
	var buf bytes.Buffer
	config := &logger.Config{
		Level:     slog.LevelInfo,
		AddSource: false,
		Service:   "test-service",
		Version:   "1.0.0",
		Writer:    &buf,
	}

	log := logger.NewLogger(config)
	fields := map[string]interface{}{
		"field1": "value1",
		"field2": "value2",
		"field3": "value3",
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		log.WithFields(fields).Info("benchmark message")
	}
}

func BenchmarkLogger_LogHTTPRequest(b *testing.B) {
	var buf bytes.Buffer
	config := &logger.Config{
		Level:     slog.LevelInfo,
		AddSource: false,
		Service:   "test-service",
		Version:   "1.0.0",
		Writer:    &buf,
	}

	log := logger.NewLogger(config)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		log.LogHTTPRequest("GET", "/api/users", "Mozilla/5.0", "192.168.1.1", "application/json")
	}
}
