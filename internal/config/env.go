package config

import (
	"e-commerce/internal/logger"
	"os"
)

func GetEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		logger.Log.Error("Missing required env: %s" + key)
	}
	return val
}
