package constants_test

import (
	"testing"
	"time"

	"{{.Module}}/internal/shared/constants"
)

func TestAllowedMethod_Values(t *testing.T) {
	tests := []struct {
		name     string
		method   constants.AllowedMethod
		expected string
	}{
		{
			name:     "GET method",
			method:   constants.AllowedMethodGET,
			expected: "GET",
		},
		{
			name:     "POST method",
			method:   constants.AllowedMethodPOST,
			expected: "POST",
		},
		{
			name:     "PUT method",
			method:   constants.AllowedMethodPUT,
			expected: "PUT",
		},
		{
			name:     "PATCH method",
			method:   constants.AllowedMethodPATCH,
			expected: "PATCH",
		},
		{
			name:     "DELETE method",
			method:   constants.AllowedMethodDELETE,
			expected: "DELETE",
		},
		{
			name:     "OPTIONS method",
			method:   constants.AllowedMethodOPTIONS,
			expected: "OPTIONS",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.method) != tt.expected {
				t.Errorf("Expected method %s, got %s", tt.expected, string(tt.method))
			}
		})
	}
}

func TestAllowedMethod_StringConversion(t *testing.T) {
	method := constants.AllowedMethodGET
	str := string(method)
	if str != "GET" {
		t.Errorf("Expected string conversion to return 'GET', got %s", str)
	}
}

func TestAllowedMethod_Comparison(t *testing.T) {
	if constants.AllowedMethodGET != "GET" {
		t.Error("AllowedMethodGET should equal 'GET'")
	}

	if constants.AllowedMethodPOST != "POST" {
		t.Error("AllowedMethodPOST should equal 'POST'")
	}

	if constants.AllowedMethodGET == constants.AllowedMethodPOST {
		t.Error("GET and POST methods should not be equal")
	}
}

func TestServiceAPIPrefix(t *testing.T) {
	if constants.SERVICE_API_PREFIX != "api/v1" {
		t.Errorf("Expected SERVICE_API_PREFIX to be 'api/v1', got %s", constants.SERVICE_API_PREFIX)
	}
}

func TestServerAllowedOrigins(t *testing.T) {
	tests := []struct {
		name     string
		origin   string
		expected string
	}{
		{
			name:     "Local origin",
			origin:   constants.ServerAllowedOriginLocal,
			expected: "http://localhost:3000",
		},
		{
			name:     "Vite origin",
			origin:   constants.ServerAllowedOriginVite,
			expected: "http://localhost:5173",
		},
		{
			name:     "React origin",
			origin:   constants.ServerAllowedOriginReact,
			expected: "http://localhost:3001",
		},
		{
			name:     "React Native origin",
			origin:   constants.ServerAllowedOriginReactNative,
			expected: "http://localhost:8081",
		},
		{
			name:     "Postman origin",
			origin:   constants.ServerAllowedOriginPostman,
			expected: "https://www.postman.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.origin != tt.expected {
				t.Errorf("Expected origin %s, got %s", tt.expected, tt.origin)
			}
		})
	}
}

func TestTimeouts(t *testing.T) {
	tests := []struct {
		name     string
		timeout  time.Duration
		expected time.Duration
	}{
		{
			name:     "Write timeout",
			timeout:  constants.WriteTimeout,
			expected: 15 * time.Second,
		},
		{
			name:     "Read timeout",
			timeout:  constants.ReadTimeout,
			expected: 15 * time.Second,
		},
		{
			name:     "Idle timeout",
			timeout:  constants.IdleTimeout,
			expected: 60 * time.Second,
		},
		{
			name:     "Shutdown grace period",
			timeout:  constants.ShutdownGracePeriod,
			expected: 30 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.timeout != tt.expected {
				t.Errorf("Expected timeout %v, got %v", tt.expected, tt.timeout)
			}
		})
	}
}

func TestTimeouts_ArePositive(t *testing.T) {
	if constants.WriteTimeout <= 0 {
		t.Error("WriteTimeout should be positive")
	}

	if constants.ReadTimeout <= 0 {
		t.Error("ReadTimeout should be positive")
	}

	if constants.IdleTimeout <= 0 {
		t.Error("IdleTimeout should be positive")
	}

	if constants.ShutdownGracePeriod <= 0 {
		t.Error("ShutdownGracePeriod should be positive")
	}
}

func TestJSONMaxSize(t *testing.T) {
	expected := 10 * 1024 * 1024 // 10MB
	if constants.JSONMaxSize != expected {
		t.Errorf("Expected JSONMaxSize to be %d bytes (10MB), got %d", expected, constants.JSONMaxSize)
	}
}

func TestJSONMaxSize_IsPositive(t *testing.T) {
	if constants.JSONMaxSize <= 0 {
		t.Error("JSONMaxSize should be positive")
	}
}

func TestWebhookEventType_Values(t *testing.T) {
	tests := []struct {
		name      string
		eventType constants.WebhookEventType
		expected  string
	}{
		{
			name:      "User created event",
			eventType: constants.WebhookEventUserCreated,
			expected:  "user.created",
		},
		{
			name:      "User updated event",
			eventType: constants.WebhookEventUserUpdated,
			expected:  "user.updated",
		},
		{
			name:      "User deleted event",
			eventType: constants.WebhookEventUserDeleted,
			expected:  "user.deleted",
		},
		{
			name:      "Organization created event",
			eventType: constants.WebhookEventOrganizationCreated,
			expected:  "organization.created",
		},
		{
			name:      "Organization updated event",
			eventType: constants.WebhookEventOrganizationUpdated,
			expected:  "organization.updated",
		},
		{
			name:      "Organization deleted event",
			eventType: constants.WebhookEventOrganizationDeleted,
			expected:  "organization.deleted",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.eventType) != tt.expected {
				t.Errorf("Expected event type %s, got %s", tt.expected, string(tt.eventType))
			}
		})
	}
}

func TestWebhookEventType_StringConversion(t *testing.T) {
	eventType := constants.WebhookEventUserCreated
	str := string(eventType)
	if str != "user.created" {
		t.Errorf("Expected string conversion to return 'user.created', got %s", str)
	}
}

func TestWebhookEventType_Comparison(t *testing.T) {
	if constants.WebhookEventUserCreated != "user.created" {
		t.Error("WebhookEventUserCreated should equal 'user.created'")
	}

	if constants.WebhookEventUserCreated == constants.WebhookEventUserUpdated {
		t.Error("User created and updated events should not be equal")
	}
}

func TestConstants_Immutability(t *testing.T) {
	// Test that constants cannot be modified (compile-time check)
	// This test ensures the constants are properly defined as const

	// These should compile without issues
	var method constants.AllowedMethod = constants.AllowedMethodGET
	var eventType constants.WebhookEventType = constants.WebhookEventUserCreated

	// Test that we can use them in comparisons
	if method != constants.AllowedMethodGET {
		t.Error("Method constant should be immutable")
	}

	if eventType != constants.WebhookEventUserCreated {
		t.Error("Event type constant should be immutable")
	}
}

func TestConstants_UsageInMaps(t *testing.T) {
	// Test that constants can be used as map keys
	methods := map[constants.AllowedMethod]bool{
		constants.AllowedMethodGET:     true,
		constants.AllowedMethodPOST:    true,
		constants.AllowedMethodPUT:     true,
		constants.AllowedMethodPATCH:   true,
		constants.AllowedMethodDELETE:  true,
		constants.AllowedMethodOPTIONS: true,
	}

	if !methods[constants.AllowedMethodGET] {
		t.Error("GET method should be in methods map")
	}

	if !methods[constants.AllowedMethodPOST] {
		t.Error("POST method should be in methods map")
	}

	// Test webhook events
	events := map[constants.WebhookEventType]bool{
		constants.WebhookEventUserCreated:         true,
		constants.WebhookEventUserUpdated:         true,
		constants.WebhookEventUserDeleted:         true,
		constants.WebhookEventOrganizationCreated: true,
		constants.WebhookEventOrganizationUpdated: true,
		constants.WebhookEventOrganizationDeleted: true,
	}

	if !events[constants.WebhookEventUserCreated] {
		t.Error("User created event should be in events map")
	}

	if !events[constants.WebhookEventOrganizationCreated] {
		t.Error("Organization created event should be in events map")
	}
}

func TestConstants_UsageInSwitch(t *testing.T) {
	// Test that constants can be used in switch statements
	method := constants.AllowedMethodGET

	switch method {
	case constants.AllowedMethodGET:
		// This should match
	case constants.AllowedMethodPOST:
		t.Error("Should not match POST method")
	default:
		t.Error("Should match GET method")
	}

	eventType := constants.WebhookEventUserCreated

	switch eventType {
	case constants.WebhookEventUserCreated:
		// This should match
	case constants.WebhookEventUserUpdated:
		t.Error("Should not match user updated event")
	default:
		t.Error("Should match user created event")
	}
}

func TestConstants_TypeSafety(t *testing.T) {
	// Test that constants maintain their types
	var method constants.AllowedMethod = constants.AllowedMethodGET
	var eventType constants.WebhookEventType = constants.WebhookEventUserCreated

	// These should compile without issues
	if method == constants.AllowedMethodGET {
		// Type-safe comparison
	}

	if eventType == constants.WebhookEventUserCreated {
		// Type-safe comparison
	}

	// Test that we can't accidentally mix types
	// This would cause a compile error if uncommented:
	// if method == eventType { }
}

// Benchmark tests
func BenchmarkAllowedMethodComparison(b *testing.B) {
	method := constants.AllowedMethodGET
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = method == constants.AllowedMethodGET
	}
}

func BenchmarkWebhookEventTypeComparison(b *testing.B) {
	eventType := constants.WebhookEventUserCreated
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = eventType == constants.WebhookEventUserCreated
	}
}

func BenchmarkTimeoutComparison(b *testing.B) {
	timeout := constants.WriteTimeout
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = timeout == 15*time.Second
	}
}
