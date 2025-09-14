package uuid_test

import (
	"regexp"
	"strings"
	"testing"

	"{{.Module}}/internal/shared/uuid"
)

func TestGenerateNamespaceUUID(t *testing.T) {
	tests := []struct {
		name      string
		namespace string
	}{
		{
			name:      "user namespace",
			namespace: "user",
		},
		{
			name:      "session namespace",
			namespace: "session",
		},
		{
			name:      "order namespace",
			namespace: "order",
		},
		{
			name:      "empty namespace",
			namespace: "",
		},
		{
			name:      "namespace with special characters",
			namespace: "test-namespace_123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := uuid.GenerateNamespaceUUID(tt.namespace)

			// Check that result is not empty
			if result == "" {
				t.Error("Generated UUID should not be empty")
			}

			// Check that result starts with namespace
			if tt.namespace != "" && !strings.HasPrefix(result, tt.namespace+"_") {
				t.Errorf("Generated UUID should start with namespace '%s_', got: %s", tt.namespace, result)
			}

			// Check UUID format (should have dashes)
			uuidPattern := regexp.MustCompile(`^[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}$`)
			if tt.namespace == "" {
				if !uuidPattern.MatchString(result) {
					t.Errorf("Generated UUID should match UUID format, got: %s", result)
				}
			} else {
				// Extract UUID part after namespace
				parts := strings.Split(result, "_")
				if len(parts) < 2 {
					t.Errorf("Generated UUID should have namespace prefix, got: %s", result)
				} else {
					uuidPart := parts[1]
					if !uuidPattern.MatchString(uuidPart) {
						t.Errorf("UUID part should match UUID format, got: %s", uuidPart)
					}
				}
			}
		})
	}
}

func TestGenerateNamespaceUUID_Uniqueness(t *testing.T) {
	namespace := "test"
	generated := make(map[string]bool)
	count := 1000

	for i := 0; i < count; i++ {
		uuid := uuid.GenerateNamespaceUUID(namespace)
		if generated[uuid] {
			t.Errorf("Generated duplicate UUID: %s", uuid)
		}
		generated[uuid] = true
	}

	if len(generated) != count {
		t.Errorf("Expected %d unique UUIDs, got %d", count, len(generated))
	}
}

func TestGenerateShortUUID(t *testing.T) {
	result := uuid.GenerateShortUUID()

	// Check that result is not empty
	if result == "" {
		t.Error("Generated short UUID should not be empty")
	}

	// Check that result is 16 characters (8 bytes in hex)
	if len(result) != 16 {
		t.Errorf("Generated short UUID should be 16 characters, got: %d", len(result))
	}

	// Check that result contains only hexadecimal characters
	hexPattern := regexp.MustCompile(`^[a-f0-9]+$`)
	if !hexPattern.MatchString(result) {
		t.Errorf("Generated short UUID should contain only hexadecimal characters, got: %s", result)
	}
}

func TestGenerateShortUUID_Uniqueness(t *testing.T) {
	generated := make(map[string]bool)
	count := 1000

	for i := 0; i < count; i++ {
		uuid := uuid.GenerateShortUUID()
		if generated[uuid] {
			t.Errorf("Generated duplicate short UUID: %s", uuid)
		}
		generated[uuid] = true
	}

	if len(generated) != count {
		t.Errorf("Expected %d unique short UUIDs, got %d", count, len(generated))
	}
}

func TestGenerateNamespaceShortUUID(t *testing.T) {
	tests := []struct {
		name      string
		namespace string
	}{
		{
			name:      "user namespace",
			namespace: "user",
		},
		{
			name:      "session namespace",
			namespace: "session",
		},
		{
			name:      "order namespace",
			namespace: "order",
		},
		{
			name:      "empty namespace",
			namespace: "",
		},
		{
			name:      "namespace with special characters",
			namespace: "test-namespace_123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := uuid.GenerateNamespaceShortUUID(tt.namespace)

			// Check that result is not empty
			if result == "" {
				t.Error("Generated namespace short UUID should not be empty")
			}

			// Check that result starts with namespace
			if tt.namespace != "" && !strings.HasPrefix(result, tt.namespace+"_") {
				t.Errorf("Generated namespace short UUID should start with namespace '%s_', got: %s", tt.namespace, result)
			}

			// Check short UUID format (should be 16 hex characters)
			hexPattern := regexp.MustCompile(`^[a-f0-9]{16}$`)
			if tt.namespace == "" {
				if !hexPattern.MatchString(result) {
					t.Errorf("Generated namespace short UUID should match hex format, got: %s", result)
				}
			} else {
				// Extract UUID part after namespace
				parts := strings.Split(result, "_")
				if len(parts) < 2 {
					t.Errorf("Generated namespace short UUID should have namespace prefix, got: %s", result)
				} else {
					uuidPart := parts[1]
					if !hexPattern.MatchString(uuidPart) {
						t.Errorf("Short UUID part should match hex format, got: %s", uuidPart)
					}
				}
			}
		})
	}
}

func TestGenerateNamespaceShortUUID_Uniqueness(t *testing.T) {
	namespace := "test"
	generated := make(map[string]bool)
	count := 1000

	for i := 0; i < count; i++ {
		uuid := uuid.GenerateNamespaceShortUUID(namespace)
		if generated[uuid] {
			t.Errorf("Generated duplicate namespace short UUID: %s", uuid)
		}
		generated[uuid] = true
	}

	if len(generated) != count {
		t.Errorf("Expected %d unique namespace short UUIDs, got %d", count, len(generated))
	}
}

func TestGenerateNamespaceUUID_Consistency(t *testing.T) {
	namespace := "consistency_test"

	// Generate multiple UUIDs and check they all follow the same pattern
	for i := 0; i < 100; i++ {
		result := uuid.GenerateNamespaceUUID(namespace)

		// Should always start with namespace
		if !strings.HasPrefix(result, namespace+"_") {
			t.Errorf("UUID should always start with namespace '%s_', got: %s", namespace, result)
		}

		// Should always have the same length
		expectedLength := len(namespace) + 1 + 36 // namespace + "_" + UUID
		if len(result) != expectedLength {
			t.Errorf("UUID should always have length %d, got: %d", expectedLength, len(result))
		}
	}
}

func TestGenerateShortUUID_Consistency(t *testing.T) {
	// Generate multiple short UUIDs and check they all follow the same pattern
	for i := 0; i < 100; i++ {
		result := uuid.GenerateShortUUID()

		// Should always be 16 characters
		if len(result) != 16 {
			t.Errorf("Short UUID should always be 16 characters, got: %d", len(result))
		}

		// Should always be hexadecimal
		hexPattern := regexp.MustCompile(`^[a-f0-9]{16}$`)
		if !hexPattern.MatchString(result) {
			t.Errorf("Short UUID should always be hexadecimal, got: %s", result)
		}
	}
}

func TestGenerateNamespaceShortUUID_Consistency(t *testing.T) {
	namespace := "consistency_test"

	// Generate multiple namespace short UUIDs and check they all follow the same pattern
	for i := 0; i < 100; i++ {
		result := uuid.GenerateNamespaceShortUUID(namespace)

		// Should always start with namespace
		if !strings.HasPrefix(result, namespace+"_") {
			t.Errorf("Namespace short UUID should always start with namespace '%s_', got: %s", namespace, result)
		}

		// Should always have the same length
		expectedLength := len(namespace) + 1 + 16 // namespace + "_" + short UUID
		if len(result) != expectedLength {
			t.Errorf("Namespace short UUID should always have length %d, got: %d", expectedLength, len(result))
		}
	}
}

// Benchmark tests
func BenchmarkGenerateNamespaceUUID(b *testing.B) {
	namespace := "benchmark"
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		uuid.GenerateNamespaceUUID(namespace)
	}
}

func BenchmarkGenerateShortUUID(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		uuid.GenerateShortUUID()
	}
}

func BenchmarkGenerateNamespaceShortUUID(b *testing.B) {
	namespace := "benchmark"
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		uuid.GenerateNamespaceShortUUID(namespace)
	}
}

// Test edge cases
func TestGenerateNamespaceUUID_EdgeCases(t *testing.T) {
	// Test with very long namespace
	longNamespace := strings.Repeat("a", 1000)
	result := uuid.GenerateNamespaceUUID(longNamespace)

	if !strings.HasPrefix(result, longNamespace+"_") {
		t.Errorf("Should handle long namespace, got: %s", result)
	}

	// Test with namespace containing underscores
	namespaceWithUnderscores := "test_namespace_with_underscores"
	result = uuid.GenerateNamespaceUUID(namespaceWithUnderscores)

	if !strings.HasPrefix(result, namespaceWithUnderscores+"_") {
		t.Errorf("Should handle namespace with underscores, got: %s", result)
	}
}

func TestGenerateShortUUID_EdgeCases(t *testing.T) {
	// Generate multiple short UUIDs to ensure randomness
	results := make([]string, 1000)
	for i := 0; i < 1000; i++ {
		results[i] = uuid.GenerateShortUUID()
	}

	// Check that we don't get all the same value
	first := results[0]
	allSame := true
	for _, result := range results {
		if result != first {
			allSame = false
			break
		}
	}

	if allSame {
		t.Error("Generated short UUIDs should not all be the same")
	}
}
