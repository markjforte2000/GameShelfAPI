package util

import (
	"log"
	"os"
)

func GetEnvironOrFail(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Failed to get environment key %v", key)
	}
	return value
}
