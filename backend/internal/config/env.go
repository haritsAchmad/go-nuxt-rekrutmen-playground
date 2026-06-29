package config

import (
	"os"
	"strconv"
)

// GetString reads a value from environment variables.
// If the value is empty, it returns the provided fallback.
func GetString(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}

// GetInt reads an integer from environment variables.
// If the value is empty or invalid, it returns the provided fallback.
func GetInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsedValue, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}

	return parsedValue
}

// GetBool reads a boolean from environment variables.
// If the value is empty or invalid, it returns the provided fallback.
func GetBool(key string, fallback bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsedValue, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}

	return parsedValue
}
