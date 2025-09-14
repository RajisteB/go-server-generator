package validation

import (
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
	return strings.TrimSpace(sanitizer.Sanitize(input))
}

// SanitizeEmail sanitizes an email input
func SanitizeEmail(input string) string {
	sanitized := SanitizeString(input)
	return strings.ToLower(sanitized)
}

// SanitizeHTML sanitizes HTML content
func SanitizeHTML(input string) string {
	return sanitizer.Sanitize(input)
}
