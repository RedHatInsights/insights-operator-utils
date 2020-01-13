package env

import (
	"os"
)

// GetEnv return value of environment variable if it exists, or fallback otherwise
func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
