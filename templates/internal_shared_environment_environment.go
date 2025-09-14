package environment

import (
	"log"

	"github.com/joho/godotenv"
)

// LoadEnvVarsFromEnv loads environment variables from .env file
func LoadEnvVarsFromEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}
}
