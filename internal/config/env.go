//revive:disable:package-comments
package config

import (
	"fmt"
	"os"
)

// RequireEnv returns the value of the given environment variable or an error
// if it is not set.
func RequireEnv(key string) (string, error) {
	val := os.Getenv(key)
	if val == "" {
		return "", fmt.Errorf("%s is required", key)
	}
	return val, nil
}

// StringEnv returns the value of the given environment variable, or the
// fallback if not set.
func StringEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}
