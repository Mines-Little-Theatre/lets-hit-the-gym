package util

import (
	"log"
	"os"
)

// ReadEnv reads an environment variable, panicking if it is not set
func ReadEnv(key string) string {
	result, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalln(key, "environment variable must be set")
	}
	return result
}
