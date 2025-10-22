package util

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// MustGetEnv loads .env and returns a required variable.
func MustGetEnv(key string) string {
	_ = godotenv.Load()
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("Missing environment variable: %s", key)
	}
	return val
}
