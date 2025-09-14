package validation

import (
	"html"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/microcosm-cc/bluemonday"
)

var validate *validator.Validate
var sanitizer *bluemonday.Policy

func init() {
	validate = validator.New()
	sanitizer = bluemonday.UGCPolicy()
}

// ValidateStruct validates a struct using struct tags
func ValidateStruct(s interface{}) error {
	return validate.Struct(s)
}

// SanitizeString sanitizes a string input
func SanitizeString(input string) string {
	// First escape HTML entities
	escaped := html.EscapeString(input)
	// Then trim whitespace
	return strings.TrimSpace(escaped)
}

// SanitizeEmail sanitizes an email input
func SanitizeEmail(input string) string {
	// Remove dangerous characters first
	dangerousChars := regexp.MustCompile(`[<>'"&;]`)
	cleaned := dangerousChars.ReplaceAllString(input, "")

	// Remove whitespace
	cleaned = strings.ReplaceAll(cleaned, " ", "")
	cleaned = strings.ReplaceAll(cleaned, "\t", "")
	cleaned = strings.ReplaceAll(cleaned, "\n", "")

	// Convert to lowercase
	return strings.ToLower(strings.TrimSpace(cleaned))
}

// SanitizeHTML sanitizes HTML content
func SanitizeHTML(input string) string {
	return sanitizer.Sanitize(input)
}
