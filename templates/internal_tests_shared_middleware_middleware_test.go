package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"{{.Module}}/internal/shared/middleware"
)

func TestNewMiddleware(t *testing.T) {
	clerkClient := nil // Mock clerk client
	clerkSecret := "test-secret"

	m := middleware.NewMiddleware(clerkClient, clerkSecret)

	if m.ClerkClient != clerkClient {
		t.Error("ClerkClient should be set correctly")
	}

	if m.ClerkSecret != clerkSecret {
		t.Error("ClerkSecret should be set correctly")
	}

	if m.RateLimiter == nil {
		t.Error("RateLimiter should be initialized")
	}
}

func TestLoggerMiddleware(t *testing.T) {
	m := middleware.NewMiddleware(nil, "test-secret")

	// Create a test handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test response"))
	})

	// Wrap with logger middleware
	loggerHandler := m.LoggerMiddleware(handler)

	// Create test request
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("User-Agent", "test-agent")
	req.Header.Set("Content-Type", "application/json")
	req.RemoteAddr = "127.0.0.1:12345"

	// Create response recorder
	w := httptest.NewRecorder()

	// Execute request
	loggerHandler.ServeHTTP(w, req)

	// Check that the response was handled correctly
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	if w.Body.String() != "test response" {
		t.Errorf("Expected response body 'test response', got '%s'", w.Body.String())
	}
}

func TestRateLimiterMiddleware(t *testing.T) {
	m := middleware.NewMiddleware(nil, "test-secret")

	// Create a test handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	})

	// Wrap with rate limiter middleware
	rateLimitHandler := m.RateLimiterMiddleware(handler)

	// Test multiple requests
	req := httptest.NewRequest("GET", "/test", nil)

	// First few requests should succeed
	for i := 0; i < 5; i++ {
		w := httptest.NewRecorder()
		rateLimitHandler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Request %d should succeed, got status %d", i+1, w.Code)
		}
	}
}

func TestClerkAuthMiddleware(t *testing.T) {
	m := middleware.NewMiddleware(nil, "test-secret")

	// Create a test handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check that user context was set
		userID := r.Context().Value("user_id")
		sessionID := r.Context().Value("session_id")

		if userID == nil {
			http.Error(w, "user_id not set", http.StatusInternalServerError)
			return
		}

		if sessionID == nil {
			http.Error(w, "session_id not set", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("authenticated"))
	})

	// Wrap with auth middleware
	authHandler := m.ClerkAuthMiddleware(handler)

	tests := []struct {
		name           string
		authHeader     string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "valid bearer token",
			authHeader:     "Bearer valid-token",
			expectedStatus: http.StatusOK,
			expectedBody:   "authenticated",
		},
		{
			name:           "missing authorization header",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "Authorization header required\n",
		},
		{
			name:           "invalid authorization format",
			authHeader:     "InvalidFormat token",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "Invalid authorization header format\n",
		},
		{
			name:           "missing token",
			authHeader:     "Bearer",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "Invalid authorization header format\n",
		},
		{
			name:           "extra parts in header",
			authHeader:     "Bearer token extra",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "Invalid authorization header format\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/test", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			w := httptest.NewRecorder()
			authHandler.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatus, w.Code)
			}

			if w.Body.String() != tt.expectedBody {
				t.Errorf("Expected body '%s', got '%s'", tt.expectedBody, w.Body.String())
			}
		})
	}
}

func TestClerkWebhookMiddleware(t *testing.T) {
	m := middleware.NewMiddleware(nil, "test-secret")

	// Create a test handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("webhook processed"))
	})

	// Wrap with webhook middleware
	webhookHandler := m.ClerkWebhookMiddleware(handler)

	tests := []struct {
		name           string
		signature      string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "valid signature",
			signature:      "valid-signature",
			expectedStatus: http.StatusOK,
			expectedBody:   "webhook processed",
		},
		{
			name:           "missing signature",
			signature:      "",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "Missing svix-signature header\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/webhook", nil)
			if tt.signature != "" {
				req.Header.Set("svix-signature", tt.signature)
			}

			w := httptest.NewRecorder()
			webhookHandler.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatus, w.Code)
			}

			if w.Body.String() != tt.expectedBody {
				t.Errorf("Expected body '%s', got '%s'", tt.expectedBody, w.Body.String())
			}
		})
	}
}

func TestSecurityHeadersMiddleware(t *testing.T) {
	m := middleware.NewMiddleware(nil, "test-secret")

	config := middleware.SecurityConfig{
		CSPPolicy:          "default-src 'self'",
		HSTSMaxAge:         31536000,
		FrameOptions:       "DENY",
		ContentTypeOptions: true,
		ReferrerPolicy:     "strict-origin-when-cross-origin",
		PermissionsPolicy:  "geolocation=(), microphone=()",
	}

	// Create a test handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("secure response"))
	})

	// Wrap with security headers middleware
	securityHandler := m.SecurityHeadersMiddleware(config)

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	securityHandler.ServeHTTP(w, req)

	// Check security headers
	expectedHeaders := map[string]string{
		"Content-Security-Policy":   "default-src 'self'",
		"Strict-Transport-Security": "max-age=31536000",
		"X-Frame-Options":           "DENY",
		"X-Content-Type-Options":    "nosniff",
		"Referrer-Policy":           "strict-origin-when-cross-origin",
		"Permissions-Policy":        "geolocation=(), microphone=()",
	}

	for header, expectedValue := range expectedHeaders {
		actualValue := w.Header().Get(header)
		if actualValue != expectedValue {
			t.Errorf("Expected header %s: %s, got: %s", header, expectedValue, actualValue)
		}
	}
}

func TestSecurityHeadersMiddleware_EmptyConfig(t *testing.T) {
	m := middleware.NewMiddleware(nil, "test-secret")

	config := middleware.SecurityConfig{}

	// Create a test handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("response"))
	})

	// Wrap with security headers middleware
	securityHandler := m.SecurityHeadersMiddleware(config)

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	securityHandler.ServeHTTP(w, req)

	// No security headers should be set
	securityHeaders := []string{
		"Content-Security-Policy",
		"Strict-Transport-Security",
		"X-Frame-Options",
		"X-Content-Type-Options",
		"Referrer-Policy",
		"Permissions-Policy",
	}

	for _, header := range securityHeaders {
		if w.Header().Get(header) != "" {
			t.Errorf("Header %s should not be set with empty config", header)
		}
	}
}

func TestRequestSizeLimitMiddleware(t *testing.T) {
	m := middleware.NewMiddleware(nil, "test-secret")

	config := middleware.RequestLimitsConfig{
		MaxRequestSize: 1024, // 1KB
	}

	// Create a test handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	})

	// Wrap with request size limit middleware
	sizeLimitHandler := m.RequestSizeLimitMiddleware(config)

	tests := []struct {
		name           string
		body           string
		expectedStatus int
	}{
		{
			name:           "small request",
			body:           "small data",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "large request",
			body:           strings.Repeat("a", 2048), // 2KB
			expectedStatus: http.StatusRequestEntityTooLarge,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/test", strings.NewReader(tt.body))
			req.ContentLength = int64(len(tt.body))

			w := httptest.NewRecorder()
			sizeLimitHandler.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestRequestTimeoutMiddleware(t *testing.T) {
	m := middleware.NewMiddleware(nil, "test-secret")

	config := middleware.RequestLimitsConfig{
		ReadTimeout: 1, // 1 second
	}

	// Create a test handler that takes longer than timeout
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate long-running operation
		time.Sleep(2 * time.Second)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	})

	// Wrap with timeout middleware
	timeoutHandler := m.RequestTimeoutMiddleware(config)

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	start := time.Now()
	timeoutHandler.ServeHTTP(w, req)
	duration := time.Since(start)

	// Should timeout and return context deadline exceeded
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// Should complete quickly due to timeout
	if duration > 2*time.Second {
		t.Errorf("Request should timeout quickly, took %v", duration)
	}
}

func TestSafeJSONDecoder(t *testing.T) {
	type TestStruct struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	tests := []struct {
		name      string
		json      string
		maxSize   int64
		expectErr bool
	}{
		{
			name:      "valid JSON",
			json:      `{"name": "test", "value": 123}`,
			maxSize:   1024,
			expectErr: false,
		},
		{
			name:      "invalid JSON",
			json:      `{"name": "test", "value": 123`, // missing closing brace
			maxSize:   1024,
			expectErr: true,
		},
		{
			name:      "unknown fields",
			json:      `{"name": "test", "value": 123, "unknown": "field"}`,
			maxSize:   1024,
			expectErr: true,
		},
		{
			name:      "too large",
			json:      `{"name": "` + strings.Repeat("a", 2048) + `", "value": 123}`,
			maxSize:   1024,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/test", strings.NewReader(tt.json))
			req.Header.Set("Content-Type", "application/json")

			var result TestStruct
			err := middleware.SafeJSONDecoder(req, &result, tt.maxSize)

			if tt.expectErr && err == nil {
				t.Error("Expected error but got none")
			}

			if !tt.expectErr && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}

			if !tt.expectErr {
				if result.Name != "test" {
					t.Errorf("Expected name 'test', got '%s'", result.Name)
				}
				if result.Value != 123 {
					t.Errorf("Expected value 123, got %d", result.Value)
				}
			}
		})
	}
}

func TestResponseWriter(t *testing.T) {
	// Test the custom responseWriter
	w := httptest.NewRecorder()
	rw := &middleware.ResponseWriter{
		ResponseWriter: w,
		StatusCode:     http.StatusOK,
	}

	// Test WriteHeader
	rw.WriteHeader(http.StatusNotFound)
	if rw.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, rw.StatusCode)
	}

	// Test Write
	data := []byte("test data")
	n, err := rw.Write(data)
	if err != nil {
		t.Errorf("Write returned error: %v", err)
	}
	if n != len(data) {
		t.Errorf("Expected to write %d bytes, wrote %d", len(data), n)
	}

	// Check that data was written to underlying ResponseWriter
	if w.Body.String() != "test data" {
		t.Errorf("Expected body 'test data', got '%s'", w.Body.String())
	}
}

// Benchmark tests
func BenchmarkLoggerMiddleware(b *testing.B) {
	m := middleware.NewMiddleware(nil, "test-secret")
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	loggerHandler := m.LoggerMiddleware(handler)
	req := httptest.NewRequest("GET", "/test", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		loggerHandler.ServeHTTP(w, req)
	}
}

func BenchmarkRateLimiterMiddleware(b *testing.B) {
	m := middleware.NewMiddleware(nil, "test-secret")
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	rateLimitHandler := m.RateLimiterMiddleware(handler)
	req := httptest.NewRequest("GET", "/test", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		rateLimitHandler.ServeHTTP(w, req)
	}
}

func BenchmarkSecurityHeadersMiddleware(b *testing.B) {
	m := middleware.NewMiddleware(nil, "test-secret")
	config := middleware.SecurityConfig{
		CSPPolicy:          "default-src 'self'",
		HSTSMaxAge:         31536000,
		FrameOptions:       "DENY",
		ContentTypeOptions: true,
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	securityHandler := m.SecurityHeadersMiddleware(config)
	req := httptest.NewRequest("GET", "/test", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		securityHandler.ServeHTTP(w, req)
	}
}
