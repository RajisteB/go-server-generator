package environment_test

import (
	"os"
	"path/filepath"
	"testing"

	"{{.Module}}/internal/shared/environment"
)

func TestLoadEnvVarsFromEnv(t *testing.T) {
	// Create a temporary directory for test .env file
	tempDir := t.TempDir()
	envFile := filepath.Join(tempDir, ".env")

	// Create test .env file
	envContent := `TEST_VAR=test_value
TEST_NUMBER=123
TEST_BOOL=true
TEST_EMPTY=
TEST_QUOTED="quoted value"
TEST_SPACES=value with spaces
`
	err := os.WriteFile(envFile, []byte(envContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test .env file: %v", err)
	}

	// Change to temp directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(originalDir)

	err = os.Chdir(tempDir)
	if err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	// Clear any existing environment variables
	os.Unsetenv("TEST_VAR")
	os.Unsetenv("TEST_NUMBER")
	os.Unsetenv("TEST_BOOL")
	os.Unsetenv("TEST_EMPTY")
	os.Unsetenv("TEST_QUOTED")
	os.Unsetenv("TEST_SPACES")

	// Load environment variables
	environment.LoadEnvVarsFromEnv()

	// Test that environment variables were loaded
	tests := []struct {
		name     string
		envVar   string
		expected string
	}{
		{
			name:     "simple variable",
			envVar:   "TEST_VAR",
			expected: "test_value",
		},
		{
			name:     "numeric variable",
			envVar:   "TEST_NUMBER",
			expected: "123",
		},
		{
			name:     "boolean variable",
			envVar:   "TEST_BOOL",
			expected: "true",
		},
		{
			name:     "empty variable",
			envVar:   "TEST_EMPTY",
			expected: "",
		},
		{
			name:     "quoted variable",
			envVar:   "TEST_QUOTED",
			expected: "quoted value",
		},
		{
			name:     "variable with spaces",
			envVar:   "TEST_SPACES",
			expected: "value with spaces",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := os.Getenv(tt.envVar)
			if actual != tt.expected {
				t.Errorf("Expected %s=%s, got %s", tt.envVar, tt.expected, actual)
			}
		})
	}
}

func TestLoadEnvVarsFromEnv_NoEnvFile(t *testing.T) {
	// Create a temporary directory without .env file
	tempDir := t.TempDir()

	// Change to temp directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(originalDir)

	err = os.Chdir(tempDir)
	if err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	// Set a system environment variable
	os.Setenv("SYSTEM_VAR", "system_value")
	defer os.Unsetenv("SYSTEM_VAR")

	// Load environment variables (should not fail even without .env file)
	environment.LoadEnvVarsFromEnv()

	// Check that system environment variable is still available
	if os.Getenv("SYSTEM_VAR") != "system_value" {
		t.Error("System environment variable should still be available")
	}
}

func TestLoadEnvVarsFromEnv_OverrideSystemEnv(t *testing.T) {
	// Create a temporary directory for test .env file
	tempDir := t.TempDir()
	envFile := filepath.Join(tempDir, ".env")

	// Create test .env file with variable that overrides system env
	envContent := `OVERRIDE_VAR=env_file_value
`
	err := os.WriteFile(envFile, []byte(envContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test .env file: %v", err)
	}

	// Change to temp directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(originalDir)

	err = os.Chdir(tempDir)
	if err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	// Set system environment variable
	os.Setenv("OVERRIDE_VAR", "system_value")
	defer os.Unsetenv("OVERRIDE_VAR")

	// Load environment variables
	environment.LoadEnvVarsFromEnv()

	// Check that .env file value overrides system environment variable
	actual := os.Getenv("OVERRIDE_VAR")
	expected := "env_file_value"
	if actual != expected {
		t.Errorf("Expected OVERRIDE_VAR=%s, got %s", expected, actual)
	}
}

func TestLoadEnvVarsFromEnv_InvalidEnvFile(t *testing.T) {
	// Create a temporary directory for test .env file
	tempDir := t.TempDir()
	envFile := filepath.Join(tempDir, ".env")

	// Create test .env file with invalid content
	envContent := `INVALID_LINE_WITHOUT_EQUALS
TEST_VAR=valid_value
ANOTHER_INVALID_LINE
`
	err := os.WriteFile(envFile, []byte(envContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test .env file: %v", err)
	}

	// Change to temp directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(originalDir)

	err = os.Chdir(tempDir)
	if err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	// Load environment variables (should not fail with invalid lines)
	environment.LoadEnvVarsFromEnv()

	// Check that valid variables are still loaded
	actual := os.Getenv("TEST_VAR")
	expected := "valid_value"
	if actual != expected {
		t.Errorf("Expected TEST_VAR=%s, got %s", expected, actual)
	}
}

func TestLoadEnvVarsFromEnv_EmptyEnvFile(t *testing.T) {
	// Create a temporary directory for test .env file
	tempDir := t.TempDir()
	envFile := filepath.Join(tempDir, ".env")

	// Create empty .env file
	err := os.WriteFile(envFile, []byte(""), 0644)
	if err != nil {
		t.Fatalf("Failed to create test .env file: %v", err)
	}

	// Change to temp directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(originalDir)

	err = os.Chdir(tempDir)
	if err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	// Load environment variables (should not fail with empty file)
	environment.LoadEnvVarsFromEnv()

	// This test mainly ensures the function doesn't panic with empty file
	// No specific assertions needed as the function should handle empty files gracefully
}

func TestLoadEnvVarsFromEnv_CommentsAndWhitespace(t *testing.T) {
	// Create a temporary directory for test .env file
	tempDir := t.TempDir()
	envFile := filepath.Join(tempDir, ".env")

	// Create test .env file with comments and whitespace
	envContent := `# This is a comment
TEST_VAR=test_value
# Another comment
TEST_NUMBER=123

# Empty line above
TEST_BOOL=true
`
	err := os.WriteFile(envFile, []byte(envContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test .env file: %v", err)
	}

	// Change to temp directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(originalDir)

	err = os.Chdir(tempDir)
	if err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	// Load environment variables
	environment.LoadEnvVarsFromEnv()

	// Test that environment variables were loaded correctly
	tests := []struct {
		name     string
		envVar   string
		expected string
	}{
		{
			name:     "variable after comment",
			envVar:   "TEST_VAR",
			expected: "test_value",
		},
		{
			name:     "variable between comments",
			envVar:   "TEST_NUMBER",
			expected: "123",
		},
		{
			name:     "variable after empty line",
			envVar:   "TEST_BOOL",
			expected: "true",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := os.Getenv(tt.envVar)
			if actual != tt.expected {
				t.Errorf("Expected %s=%s, got %s", tt.envVar, tt.expected, actual)
			}
		})
	}
}

func TestLoadEnvVarsFromEnv_SpecialCharacters(t *testing.T) {
	// Create a temporary directory for test .env file
	tempDir := t.TempDir()
	envFile := filepath.Join(tempDir, ".env")

	// Create test .env file with special characters
	envContent := `TEST_SPECIAL=value with special chars: !@#$%^&*()
TEST_URL=https://example.com/path?param=value&other=123
TEST_JSON={"key": "value", "number": 123}
TEST_MULTILINE=line1\nline2\nline3
`
	err := os.WriteFile(envFile, []byte(envContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test .env file: %v", err)
	}

	// Change to temp directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(originalDir)

	err = os.Chdir(tempDir)
	if err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	// Load environment variables
	environment.LoadEnvVarsFromEnv()

	// Test that environment variables with special characters were loaded
	tests := []struct {
		name     string
		envVar   string
		expected string
	}{
		{
			name:     "special characters",
			envVar:   "TEST_SPECIAL",
			expected: "value with special chars: !@#$%^&*()",
		},
		{
			name:     "URL",
			envVar:   "TEST_URL",
			expected: "https://example.com/path?param=value&other=123",
		},
		{
			name:     "JSON",
			envVar:   "TEST_JSON",
			expected: `{"key": "value", "number": 123}`,
		},
		{
			name:     "multiline",
			envVar:   "TEST_MULTILINE",
			expected: "line1\\nline2\\nline3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := os.Getenv(tt.envVar)
			if actual != tt.expected {
				t.Errorf("Expected %s=%s, got %s", tt.envVar, tt.expected, actual)
			}
		})
	}
}

// Benchmark test
func BenchmarkLoadEnvVarsFromEnv(b *testing.B) {
	// Create a temporary directory for test .env file
	tempDir := b.TempDir()
	envFile := filepath.Join(tempDir, ".env")

	// Create test .env file
	envContent := `TEST_VAR1=value1
TEST_VAR2=value2
TEST_VAR3=value3
TEST_VAR4=value4
TEST_VAR5=value5
`
	err := os.WriteFile(envFile, []byte(envContent), 0644)
	if err != nil {
		b.Fatalf("Failed to create test .env file: %v", err)
	}

	// Change to temp directory
	originalDir, err := os.Getwd()
	if err != nil {
		b.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(originalDir)

	err = os.Chdir(tempDir)
	if err != nil {
		b.Fatalf("Failed to change to temp directory: %v", err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		environment.LoadEnvVarsFromEnv()
	}
}
