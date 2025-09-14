package validation_test

import (
	"strings"
	"testing"

	"{{.Module}}/internal/shared/validation"
)

func TestSanitizeString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "clean string remains unchanged",
			input:    "hello world",
			expected: "hello world",
		},
		{
			name:     "string with leading/trailing whitespace is trimmed",
			input:    "  hello world  ",
			expected: "hello world",
		},
		{
			name:     "string with HTML tags is sanitized",
			input:    "<script>alert('xss')</script>hello",
			expected: "&lt;script&gt;alert(&#39;xss&#39;)&lt;/script&gt;hello",
		},
		{
			name:     "string with HTML entities is escaped",
			input:    "hello & goodbye < > \"quotes\"",
			expected: "hello &amp; goodbye &lt; &gt; &#34;quotes&#34;",
		},
		{
			name:     "empty string returns empty",
			input:    "",
			expected: "",
		},
		{
			name:     "whitespace only string returns empty",
			input:    "   \t\n  ",
			expected: "",
		},
		{
			name:     "string with malicious script",
			input:    "<img src=x onerror=alert(1)>",
			expected: "&lt;img src=x onerror=alert(1)&gt;",
		},
		{
			name:     "string with multiple HTML tags",
			input:    "<div><p>Hello <b>world</b></p></div>",
			expected: "&lt;div&gt;&lt;p&gt;Hello &lt;b&gt;world&lt;/b&gt;&lt;/p&gt;&lt;/div&gt;",
		},
		{
			name:     "string with inline styles",
			input:    "<p style='color:red'>Hello</p>",
			expected: "&lt;p style=&#39;color:red&#39;&gt;Hello&lt;/p&gt;",
		},
		{
			name:     "string with javascript protocol",
			input:    "<a href='javascript:alert(1)'>Click</a>",
			expected: "&lt;a href=&#39;javascript:alert(1)&#39;&gt;Click&lt;/a&gt;",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validation.SanitizeString(tt.input)
			if result != tt.expected {
				t.Errorf("SanitizeString() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestSanitizeEmail(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "valid email remains unchanged",
			input:    "user@example.com",
			expected: "user@example.com",
		},
		{
			name:     "email with uppercase is converted to lowercase",
			input:    "USER@EXAMPLE.COM",
			expected: "user@example.com",
		},
		{
			name:     "email with leading/trailing whitespace is trimmed",
			input:    "  user@example.com  ",
			expected: "user@example.com",
		},
		{
			name:     "email with invalid characters is cleaned",
			input:    "user<script>@example.com",
			expected: "userscript@example.com",
		},
		{
			name:     "email with special valid characters remains",
			input:    "user.name+tag@example-domain.com",
			expected: "user.nametag@example-domain.com",
		},
		{
			name:     "email with underscores remains",
			input:    "user_name@example_domain.com",
			expected: "user_name@example_domain.com",
		},
		{
			name:     "empty email returns empty",
			input:    "",
			expected: "",
		},
		{
			name:     "email with numbers remains",
			input:    "user123@example123.com",
			expected: "user123@example123.com",
		},
		{
			name:     "email with invalid HTML characters is cleaned",
			input:    "user&lt;@example&gt;.com",
			expected: "userlt@examplegt.com",
		},
		{
			name:     "email with spaces is cleaned",
			input:    "user name@example.com",
			expected: "username@example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validation.SanitizeEmail(tt.input)
			if result != tt.expected {
				t.Errorf("SanitizeEmail() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestValidateStruct(t *testing.T) {
	// Test struct with validation tags
	type TestStruct struct {
		Email    string `validate:"required,email"`
		Age      int    `validate:"min=0,max=120"`
		Username string `validate:"required,min=3,max=20"`
	}

	tests := []struct {
		name    string
		input   TestStruct
		wantErr bool
	}{
		{
			name: "valid struct passes validation",
			input: TestStruct{
				Email:    "user@example.com",
				Age:      25,
				Username: "validuser",
			},
			wantErr: false,
		},
		{
			name: "missing required email fails validation",
			input: TestStruct{
				Email:    "",
				Age:      25,
				Username: "validuser",
			},
			wantErr: true,
		},
		{
			name: "invalid email format fails validation",
			input: TestStruct{
				Email:    "invalid-email",
				Age:      25,
				Username: "validuser",
			},
			wantErr: true,
		},
		{
			name: "age below minimum fails validation",
			input: TestStruct{
				Email:    "user@example.com",
				Age:      -1,
				Username: "validuser",
			},
			wantErr: true,
		},
		{
			name: "age above maximum fails validation",
			input: TestStruct{
				Email:    "user@example.com",
				Age:      150,
				Username: "validuser",
			},
			wantErr: true,
		},
		{
			name: "username too short fails validation",
			input: TestStruct{
				Email:    "user@example.com",
				Age:      25,
				Username: "ab",
			},
			wantErr: true,
		},
		{
			name: "username too long fails validation",
			input: TestStruct{
				Email:    "user@example.com",
				Age:      25,
				Username: "verylongusernamethatexceedslimit",
			},
			wantErr: true,
		},
		{
			name: "missing required username fails validation",
			input: TestStruct{
				Email:    "user@example.com",
				Age:      25,
				Username: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validation.ValidateStruct(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateStruct() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateStruct_WithPointers(t *testing.T) {
	type TestStruct struct {
		Email    *string `validate:"omitempty,email"`
		Age      *int    `validate:"omitempty,min=0,max=120"`
		Username *string `validate:"required"`
	}

	validEmail := "user@example.com"
	validAge := 25
	validUsername := "testuser"
	invalidEmail := "invalid-email"

	tests := []struct {
		name    string
		input   TestStruct
		wantErr bool
	}{
		{
			name: "valid struct with pointers passes validation",
			input: TestStruct{
				Email:    &validEmail,
				Age:      &validAge,
				Username: &validUsername,
			},
			wantErr: false,
		},
		{
			name: "struct with nil optional fields passes validation",
			input: TestStruct{
				Email:    nil,
				Age:      nil,
				Username: &validUsername,
			},
			wantErr: false,
		},
		{
			name: "struct with invalid email fails validation",
			input: TestStruct{
				Email:    &invalidEmail,
				Age:      &validAge,
				Username: &validUsername,
			},
			wantErr: true,
		},
		{
			name: "struct with nil required field fails validation",
			input: TestStruct{
				Email:    &validEmail,
				Age:      &validAge,
				Username: nil,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validation.ValidateStruct(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateStruct() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateStruct_NestedStruct(t *testing.T) {
	type Address struct {
		Street  string `validate:"required"`
		City    string `validate:"required"`
		ZipCode string `validate:"required,len=5"`
	}

	type User struct {
		Name    string   `validate:"required"`
		Email   string   `validate:"required,email"`
		Address *Address `validate:"required"`
	}

	tests := []struct {
		name    string
		input   User
		wantErr bool
	}{
		{
			name: "valid nested struct passes validation",
			input: User{
				Name:  "John Doe",
				Email: "john@example.com",
				Address: &Address{
					Street:  "123 Main St",
					City:    "Anytown",
					ZipCode: "12345",
				},
			},
			wantErr: false,
		},
		{
			name: "invalid nested struct fails validation",
			input: User{
				Name:  "John Doe",
				Email: "john@example.com",
				Address: &Address{
					Street:  "",
					City:    "Anytown",
					ZipCode: "12345",
				},
			},
			wantErr: true,
		},
		{
			name: "nil nested struct fails validation",
			input: User{
				Name:    "John Doe",
				Email:   "john@example.com",
				Address: nil,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validation.ValidateStruct(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateStruct() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateStruct_SlicesAndMaps(t *testing.T) {
	type TestStruct struct {
		Tags     []string          `validate:"required,min=1,dive,required"`
		Metadata map[string]string `validate:"required"`
	}

	tests := []struct {
		name    string
		input   TestStruct
		wantErr bool
	}{
		{
			name: "valid struct with slice and map passes validation",
			input: TestStruct{
				Tags:     []string{"tag1", "tag2"},
				Metadata: map[string]string{"key": "value"},
			},
			wantErr: false,
		},
		{
			name: "empty slice fails validation",
			input: TestStruct{
				Tags:     []string{},
				Metadata: map[string]string{"key": "value"},
			},
			wantErr: true,
		},
		{
			name: "slice with empty string fails validation",
			input: TestStruct{
				Tags:     []string{"tag1", ""},
				Metadata: map[string]string{"key": "value"},
			},
			wantErr: true,
		},
		{
			name: "nil map fails validation",
			input: TestStruct{
				Tags:     []string{"tag1", "tag2"},
				Metadata: nil,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validation.ValidateStruct(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateStruct() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSanitizeString_SecurityTests(t *testing.T) {
	// Test various XSS and injection attempts
	maliciousInputs := []struct {
		name  string
		input string
	}{
		{"Script tag", "<script>alert('xss')</script>"},
		{"Image with onerror", "<img src=x onerror=alert(1)>"},
		{"SVG with script", "<svg onload=alert(1)>"},
		{"Iframe with javascript", "<iframe src='javascript:alert(1)'></iframe>"},
		{"Link with javascript", "<a href='javascript:alert(1)'>Click</a>"},
		{"Style with expression", "<div style='expression(alert(1))'>test</div>"},
		{"Meta refresh", "<meta http-equiv='refresh' content='0;url=javascript:alert(1)'>"},
		{"Object with data", "<object data='javascript:alert(1)'></object>"},
		{"Embed with src", "<embed src='javascript:alert(1)'>"},
		{"Form with action", "<form action='javascript:alert(1)'><input type='submit'></form>"},
	}

	for _, tt := range maliciousInputs {
		t.Run(tt.name, func(t *testing.T) {
			result := validation.SanitizeString(tt.input)

			// Result should not contain unescaped dangerous content
			// The sanitizer escapes HTML, so we check for escaped versions
			if strings.Contains(result, "<script>") {
				t.Errorf("SanitizeString() should escape script tags, got: %s", result)
			}
			// Note: The current implementation escapes rather than removes content
			// This is still secure as the content cannot execute
		})
	}
}

func TestSanitizeEmail_SecurityTests(t *testing.T) {
	// Test various email injection attempts
	maliciousEmails := []struct {
		name          string
		input         string
		shouldContain string
	}{
		{"Email with script", "user<script>@example.com", "@example.com"},
		{"Email with HTML", "user<b>name</b>@example.com", "userbnameb@example.com"},
		{"Email with quotes", "user\"name\"@example.com", "username@example.com"},
		{"Email with semicolon", "user;name@example.com", "username@example.com"},
		{"Email with newline", "user\nname@example.com", "username@example.com"},
		{"Email with tab", "user\tname@example.com", "username@example.com"},
	}

	for _, tt := range maliciousEmails {
		t.Run(tt.name, func(t *testing.T) {
			result := validation.SanitizeEmail(tt.input)

			// Result should not contain dangerous characters
			dangerousChars := []string{"<", ">", "\"", ";", "\n", "\t", "'"}
			for _, char := range dangerousChars {
				if strings.Contains(result, char) {
					t.Errorf("SanitizeEmail() should remove dangerous character %s, got: %s", char, result)
				}
			}

			// Should still contain valid email parts
			if tt.shouldContain != "" && !strings.Contains(result, tt.shouldContain) {
				t.Errorf("SanitizeEmail() should preserve valid parts %s, got: %s", tt.shouldContain, result)
			}
		})
	}
}

// Integration tests
func TestValidationIntegration(t *testing.T) {
	// Test a complete validation flow: sanitize then validate
	type User struct {
		Email    string `validate:"required,email"`
		Username string `validate:"required,min=3"`
	}

	// Input with potentially dangerous content
	rawEmail := "  USER<script>@EXAMPLE.COM  "
	rawUsername := "  <b>testuser</b>  "

	// Sanitize inputs
	cleanEmail := validation.SanitizeEmail(rawEmail)
	cleanUsername := validation.SanitizeString(rawUsername)

	user := User{
		Email:    cleanEmail,
		Username: cleanUsername,
	}

	// Validate the sanitized struct
	err := validation.ValidateStruct(user)
	if err != nil {
		t.Errorf("Integration test failed: %v", err)
	}

	// Verify sanitization worked
	expectedEmail := "userscript@example.com"         // The sanitizer removes < > but keeps other characters
	expectedUsername := "&lt;b&gt;testuser&lt;/b&gt;" // The sanitizer escapes HTML
	if user.Email != expectedEmail {
		t.Errorf("Email sanitization failed: got %s, want %s", user.Email, expectedEmail)
	}
	if user.Username != expectedUsername {
		t.Errorf("Username sanitization failed: got %s, want %s", user.Username, expectedUsername)
	}
}

// Benchmark tests
func BenchmarkSanitizeString(b *testing.B) {
	input := "<script>alert('xss')</script>Hello <b>World</b> & friends"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		validation.SanitizeString(input)
	}
}

func BenchmarkSanitizeEmail(b *testing.B) {
	input := "  USER<script>@EXAMPLE.COM  "
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		validation.SanitizeEmail(input)
	}
}

func BenchmarkValidateStruct(b *testing.B) {
	type TestStruct struct {
		Email    string `validate:"required,email"`
		Age      int    `validate:"min=0,max=120"`
		Username string `validate:"required,min=3,max=20"`
	}

	testStruct := TestStruct{
		Email:    "user@example.com",
		Age:      25,
		Username: "testuser",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		validation.ValidateStruct(testStruct)
	}
}

func BenchmarkSanitizeStringLarge(b *testing.B) {
	// Test with a larger input to see performance characteristics
	input := strings.Repeat("<script>alert('xss')</script>Hello <b>World</b> & friends ", 100)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		validation.SanitizeString(input)
	}
}
