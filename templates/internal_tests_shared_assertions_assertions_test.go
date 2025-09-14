package assertions_test

import (
	"testing"

	"{{.Module}}/internal/shared/assertions"
)

func TestAssertNonEmptyString(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{
			name:    "non-empty string should not return error",
			value:   "test",
			wantErr: false,
		},
		{
			name:    "empty string should return error",
			value:   "",
			wantErr: true,
		},
		{
			name:    "whitespace string should return error",
			value:   "   ",
			wantErr: true,
		},
		{
			name:    "string with content should not return error",
			value:   "hello world",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := assertions.AssertNonEmptyString(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("AssertNonEmptyString() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAssertNonZeroInt(t *testing.T) {
	tests := []struct {
		name    string
		value   int
		wantErr bool
	}{
		{
			name:    "positive integer should not return error",
			value:   5,
			wantErr: false,
		},
		{
			name:    "negative integer should not return error",
			value:   -5,
			wantErr: false,
		},
		{
			name:    "zero should return error",
			value:   0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := assertions.AssertNonZeroInt(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("AssertNonZeroInt() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAssertNonZeroInt64(t *testing.T) {
	tests := []struct {
		name    string
		value   int64
		wantErr bool
	}{
		{
			name:    "positive int64 should not return error",
			value:   9223372036854775807,
			wantErr: false,
		},
		{
			name:    "negative int64 should not return error",
			value:   -9223372036854775808,
			wantErr: false,
		},
		{
			name:    "zero should return error",
			value:   0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := assertions.AssertNonZeroInt64(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("AssertNonZeroInt64() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAssertPositiveInt(t *testing.T) {
	tests := []struct {
		name    string
		value   int
		wantErr bool
	}{
		{
			name:    "positive integer should not return error",
			value:   5,
			wantErr: false,
		},
		{
			name:    "zero should return error",
			value:   0,
			wantErr: true,
		},
		{
			name:    "negative integer should return error",
			value:   -5,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := assertions.AssertPositiveInt(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("AssertPositiveInt() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAssertPositiveInt64(t *testing.T) {
	tests := []struct {
		name    string
		value   int64
		wantErr bool
	}{
		{
			name:    "positive int64 should not return error",
			value:   9223372036854775807,
			wantErr: false,
		},
		{
			name:    "zero should return error",
			value:   0,
			wantErr: true,
		},
		{
			name:    "negative int64 should return error",
			value:   -9223372036854775808,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := assertions.AssertPositiveInt64(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("AssertPositiveInt64() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAssertNonEmptyString_ErrorMessages(t *testing.T) {
	err := assertions.AssertNonEmptyString("")
	if err == nil {
		t.Error("Expected error for empty string")
	}
	if err.Error() != "string cannot be empty" {
		t.Errorf("Expected error message 'string cannot be empty', got '%s'", err.Error())
	}
}

func TestAssertNonZeroInt_ErrorMessages(t *testing.T) {
	err := assertions.AssertNonZeroInt(0)
	if err == nil {
		t.Error("Expected error for zero integer")
	}
	if err.Error() != "integer cannot be zero" {
		t.Errorf("Expected error message 'integer cannot be zero', got '%s'", err.Error())
	}
}

func TestAssertNonZeroInt64_ErrorMessages(t *testing.T) {
	err := assertions.AssertNonZeroInt64(0)
	if err == nil {
		t.Error("Expected error for zero int64")
	}
	if err.Error() != "int64 cannot be zero" {
		t.Errorf("Expected error message 'int64 cannot be zero', got '%s'", err.Error())
	}
}

func TestAssertPositiveInt_ErrorMessages(t *testing.T) {
	err := assertions.AssertPositiveInt(-1)
	if err == nil {
		t.Error("Expected error for negative integer")
	}
	if err.Error() != "integer must be positive" {
		t.Errorf("Expected error message 'integer must be positive', got '%s'", err.Error())
	}
}

func TestAssertPositiveInt64_ErrorMessages(t *testing.T) {
	err := assertions.AssertPositiveInt64(-1)
	if err == nil {
		t.Error("Expected error for negative int64")
	}
	if err.Error() != "int64 must be positive" {
		t.Errorf("Expected error message 'int64 must be positive', got '%s'", err.Error())
	}
}

// Benchmark tests
func BenchmarkAssertNonEmptyString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		assertions.AssertNonEmptyString("test string")
	}
}

func BenchmarkAssertNonZeroInt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		assertions.AssertNonZeroInt(42)
	}
}

func BenchmarkAssertNonZeroInt64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		assertions.AssertNonZeroInt64(42)
	}
}

func BenchmarkAssertPositiveInt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		assertions.AssertPositiveInt(42)
	}
}

func BenchmarkAssertPositiveInt64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		assertions.AssertPositiveInt64(42)
	}
}
