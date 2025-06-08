package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("❌ Warning: No .env file found. Using default environment variables.")
	}
}

func GetEnv(key string, defaultValue ...string) string {
	value, exists := os.LookupEnv(key)
	if !exists && len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return value
}
