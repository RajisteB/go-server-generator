package assertions

import (
	"errors"
	"strings"
)

// AssertNonEmptyString checks if a string is not empty
func AssertNonEmptyString(value string) error {
	if strings.TrimSpace(value) == "" {
		return errors.New("string cannot be empty")
	}
	return nil
}

// AssertNonZeroInt checks if an integer is not zero
func AssertNonZeroInt(value int) error {
	if value == 0 {
		return errors.New("integer cannot be zero")
	}
	return nil
}

// AssertNonZeroInt64 checks if an int64 is not zero
func AssertNonZeroInt64(value int64) error {
	if value == 0 {
		return errors.New("int64 cannot be zero")
	}
	return nil
}

// AssertPositiveInt checks if an integer is positive
func AssertPositiveInt(value int) error {
	if value <= 0 {
		return errors.New("integer must be positive")
	}
	return nil
}

// AssertPositiveInt64 checks if an int64 is positive
func AssertPositiveInt64(value int64) error {
	if value <= 0 {
		return errors.New("int64 must be positive")
	}
	return nil
}
