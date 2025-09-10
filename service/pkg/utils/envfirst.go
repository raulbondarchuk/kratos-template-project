package utils

import "os"

// envFirst returns the first non-empty environment variable from the provided keys
func EnvFirst(keys ...string) string {
	for _, k := range keys {
		if v := os.Getenv(k); v != "" {
			return v
		}
	}
	return ""
}
