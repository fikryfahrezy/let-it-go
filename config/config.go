package config

import (
	"os"
	"strconv"

	"github.com/fikryfahrezy/let-it-go/pkg/database"
	"github.com/fikryfahrezy/let-it-go/pkg/logger"
)

type Config struct {
	Server   ServerConfig
	Database database.Config
	Logger   logger.Config
}

type ServerConfig struct {
	Host string
	Port int
}

func Load() Config {
	return Config{
		Server: ServerConfig{
			Host: getEnv("SERVER_HOST", "localhost"),
			Port: getEnvAsInt("SERVER_PORT", 8080),
		},
		Database: database.Config{
			DSN: getEnv("DB_DSN", ""),
		},
		Logger: logger.Config{
			Level:  logger.ParseLevel(getEnv("LOG_LEVEL", "info")),
			Format: logger.ParseFormat(getEnv("LOG_FORMAT", "text")),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}