package http_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	httpHelpers "{{.Module}}/internal/shared/http"
)

func TestRespondWithJSON(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		data       interface{}
		expected   string
	}{
		{
			name:       "simple object",
			statusCode: http.StatusOK,
			data:       map[string]string{"message": "success"},
			expected:   `{"message":"success"}`,
		},
		{
			name:       "array response",
			statusCode: http.StatusOK,
			data:       []string{"item1", "item2", "item3"},
			expected:   `["item1","item2","item3"]`,
		},
		{
			name:       "struct response",
			statusCode: http.StatusCreated,
			data: struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
			}{
				ID:   1,
				Name: "test",
			},
			expected: `{"id":1,"name":"test"}`,
		},
		{
			name:       "nil data",
			statusCode: http.StatusNoContent,
			data:       nil,
			expected:   `null`,
		},
		{
			name:       "empty object",
			statusCode: http.StatusOK,
			data:       map[string]interface{}{},
			expected:   `{}`,
		},
		{
			name:       "nested object",
			statusCode: http.StatusOK,
			data: map[string]interface{}{
				"user": map[string]interface{}{
					"id":   1,
					"name": "John",
					"tags": []string{"admin", "user"},
				},
			},
			expected: `{"user":{"id":1,"name":"John","tags":["admin","user"]}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			err := httpHelpers.RespondWithJSON(w, tt.statusCode, tt.data)
			if err != nil {
				t.Fatalf("RespondWithJSON returned error: %v", err)
			}

			// Check status code
			if w.Code != tt.statusCode {
				t.Errorf("Expected status code %d, got %d", tt.statusCode, w.Code)
			}

			// Check content type
			contentType := w.Header().Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("Expected Content-Type 'application/json', got '%s'", contentType)
			}

			// Check response body
			var actual, expected interface{}
			if err := json.Unmarshal(w.Body.Bytes(), &actual); err != nil {
				t.Fatalf("Failed to unmarshal actual response: %v", err)
			}
			if err := json.Unmarshal([]byte(tt.expected), &expected); err != nil {
				t.Fatalf("Failed to unmarshal expected response: %v", err)
			}

			// Compare the unmarshaled objects
			actualJSON, _ := json.Marshal(actual)
			expectedJSON, _ := json.Marshal(expected)
			if string(actualJSON) != string(expectedJSON) {
				t.Errorf("Expected response body %s, got %s", tt.expected, string(actualJSON))
			}
		})
	}
}

func TestRespondWithJSON_ErrorHandling(t *testing.T) {
	// Test with data that cannot be marshaled
	w := httptest.NewRecorder()

	// Create a channel (cannot be marshaled to JSON)
	unmarshallableData := make(chan int)

	err := httpHelpers.RespondWithJSON(w, http.StatusOK, unmarshallableData)
	if err == nil {
		t.Error("Expected error when marshaling unmarshallable data")
	}

	// Response should still have status code and content type set
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type 'application/json', got '%s'", contentType)
	}
}

func TestRespondWithError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected string
	}{
		{
			name:     "simple error",
			err:      errors.New("something went wrong"),
			expected: `{"error":"something went wrong"}`,
		},
		{
			name:     "empty error message",
			err:      errors.New(""),
			expected: `{"error":""}`,
		},
		{
			name:     "nil error",
			err:      nil,
			expected: `{"error":""}`,
		},
		{
			name:     "error with special characters",
			err:      errors.New("error with \"quotes\" and 'apostrophes'"),
			expected: `{"error":"error with \"quotes\" and 'apostrophes'"}`,
		},
		{
			name:     "error with newlines",
			err:      errors.New("error\nwith\nnewlines"),
			expected: `{"error":"error\nwith\nnewlines"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			err := httpHelpers.RespondWithError(w, tt.err)
			if err != nil {
				t.Fatalf("RespondWithError returned error: %v", err)
			}

			// Check status code (should default to 500)
			if w.Code != http.StatusInternalServerError {
				t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, w.Code)
			}

			// Check content type
			contentType := w.Header().Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("Expected Content-Type 'application/json', got '%s'", contentType)
			}

			// Check response body
			var actual, expected interface{}
			if err := json.Unmarshal(w.Body.Bytes(), &actual); err != nil {
				t.Fatalf("Failed to unmarshal actual response: %v", err)
			}
			if err := json.Unmarshal([]byte(tt.expected), &expected); err != nil {
				t.Fatalf("Failed to unmarshal expected response: %v", err)
			}

			// Compare the unmarshaled objects
			actualJSON, _ := json.Marshal(actual)
			expectedJSON, _ := json.Marshal(expected)
			if string(actualJSON) != string(expectedJSON) {
				t.Errorf("Expected response body %s, got %s", tt.expected, string(actualJSON))
			}
		})
	}
}

func TestHandlerFunc_ServeHTTP(t *testing.T) {
	tests := []struct {
		name           string
		handlerFunc    httpHelpers.HandlerFunc
		expectedStatus int
		expectError    bool
	}{
		{
			name: "successful handler",
			handlerFunc: httpHelpers.HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
				return httpHelpers.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "success"})
			}),
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name: "handler that returns error",
			handlerFunc: httpHelpers.HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
				return errors.New("handler error")
			}),
			expectedStatus: http.StatusInternalServerError,
			expectError:    true,
		},
		{
			name: "handler that returns nil error",
			handlerFunc: httpHelpers.HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
				w.WriteHeader(http.StatusNoContent)
				return nil
			}),
			expectedStatus: http.StatusNoContent,
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/test", nil)

			handler := tt.handlerFunc
			handler.ServeHTTP(w, r)

			// Check status code
			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatus, w.Code)
			}

			// Check content type for error responses
			if tt.expectError {
				contentType := w.Header().Get("Content-Type")
				if contentType != "application/json" {
					t.Errorf("Expected Content-Type 'application/json' for error response, got '%s'", contentType)
				}

				// Check that error response is valid JSON
				var errorResponse map[string]interface{}
				if err := json.Unmarshal(w.Body.Bytes(), &errorResponse); err != nil {
					t.Errorf("Error response should be valid JSON: %v", err)
				}

				// Check that error response has error field
				if _, exists := errorResponse["error"]; !exists {
					t.Error("Error response should have 'error' field")
				}
			}
		})
	}
}

func TestHandlerFunc_ImplementsHandler(t *testing.T) {
	// Test that HandlerFunc implements http.Handler interface
	var handler http.Handler = httpHelpers.HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
		return nil
	})

	if handler == nil {
		t.Error("HandlerFunc should implement http.Handler interface")
	}

	// Test that we can call ServeHTTP
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/test", nil)
	handler.ServeHTTP(w, r)

	// Should not panic and should complete successfully
}

func TestRespondWithJSON_HeadersPreserved(t *testing.T) {
	w := httptest.NewRecorder()

	// Set some custom headers before calling RespondWithJSON
	w.Header().Set("X-Custom-Header", "custom-value")
	w.Header().Set("Cache-Control", "no-cache")

	data := map[string]string{"message": "test"}
	err := httpHelpers.RespondWithJSON(w, http.StatusOK, data)
	if err != nil {
		t.Fatalf("RespondWithJSON returned error: %v", err)
	}

	// Check that custom headers are preserved
	if w.Header().Get("X-Custom-Header") != "custom-value" {
		t.Error("Custom header should be preserved")
	}

	if w.Header().Get("Cache-Control") != "no-cache" {
		t.Error("Cache-Control header should be preserved")
	}

	// Check that Content-Type is set
	if w.Header().Get("Content-Type") != "application/json" {
		t.Error("Content-Type should be set to application/json")
	}
}

func TestRespondWithError_HeadersPreserved(t *testing.T) {
	w := httptest.NewRecorder()

	// Set some custom headers before calling RespondWithError
	w.Header().Set("X-Custom-Header", "custom-value")
	w.Header().Set("Cache-Control", "no-cache")

	err := httpHelpers.RespondWithError(w, errors.New("test error"))
	if err != nil {
		t.Fatalf("RespondWithError returned error: %v", err)
	}

	// Check that custom headers are preserved
	if w.Header().Get("X-Custom-Header") != "custom-value" {
		t.Error("Custom header should be preserved")
	}

	if w.Header().Get("Cache-Control") != "no-cache" {
		t.Error("Cache-Control header should be preserved")
	}

	// Check that Content-Type is set
	if w.Header().Get("Content-Type") != "application/json" {
		t.Error("Content-Type should be set to application/json")
	}
}

func TestRespondWithJSON_LargeData(t *testing.T) {
	// Test with large data structure
	w := httptest.NewRecorder()

	// Create a large slice
	largeData := make([]map[string]interface{}, 1000)
	for i := 0; i < 1000; i++ {
		largeData[i] = map[string]interface{}{
			"id":    i,
			"name":  "item" + string(rune(i)),
			"value": i * 2,
		}
	}

	err := httpHelpers.RespondWithJSON(w, http.StatusOK, largeData)
	if err != nil {
		t.Fatalf("RespondWithJSON returned error: %v", err)
	}

	// Check that response was successful
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// Check that response is valid JSON
	var response []map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Errorf("Large data response should be valid JSON: %v", err)
	}

	// Check that we got the expected number of items
	if len(response) != 1000 {
		t.Errorf("Expected 1000 items, got %d", len(response))
	}
}

// Benchmark tests
func BenchmarkRespondWithJSON(b *testing.B) {
	w := httptest.NewRecorder()
	data := map[string]interface{}{
		"id":    1,
		"name":  "test",
		"value": 123.45,
		"items": []string{"item1", "item2", "item3"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		httpHelpers.RespondWithJSON(w, http.StatusOK, data)
		w.Body.Reset()
		w.Code = 0
		w.HeaderMap = make(http.Header)
	}
}

func BenchmarkRespondWithError(b *testing.B) {
	w := httptest.NewRecorder()
	err := errors.New("test error")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		httpHelpers.RespondWithError(w, err)
		w.Body.Reset()
		w.Code = 0
		w.HeaderMap = make(http.Header)
	}
}

func BenchmarkHandlerFunc_ServeHTTP(b *testing.B) {
	handler := httpHelpers.HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
		return httpHelpers.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "success"})
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/test", nil)
		handler.ServeHTTP(w, r)
	}
}
