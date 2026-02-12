package utils

import "os"

func GetEnvOrDefault(key, fallback string) string {
	if result := os.Getenv(key); result != "" {
		return result
	}
	return fallback
}
