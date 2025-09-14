package uuid

import (
	"crypto/rand"
	"fmt"
)

// GenerateNamespaceUUID generates a UUID with a namespace prefix
func GenerateNamespaceUUID(namespace string) string {
	// Generate a random UUID
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}

	// Format as UUID
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])

	// Add namespace prefix
	return fmt.Sprintf("%s_%s", namespace, uuid)
}

// GenerateShortUUID generates a shorter UUID without dashes
func GenerateShortUUID() string {
	b := make([]byte, 8)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("%x", b)
}

// GenerateNamespaceShortUUID generates a short UUID with namespace prefix
func GenerateNamespaceShortUUID(namespace string) string {
	return fmt.Sprintf("%s_%s", namespace, GenerateShortUUID())
}
