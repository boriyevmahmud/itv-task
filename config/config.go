package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cast"
)

// Config holds all environment variables
type Config struct {
	PostgresHost     string
	PostgresPort     int
	PostgresUser     string
	PostgresPassword string
	PostgresDatabase string
	JWTAccessSecret  string
	JWTRefreshSecret string

	ServiceName string
}

func Load() Config {
	// Load environment variables from .env file
	if err := godotenv.Load(".env"); err != nil {
		fmt.Println("Error loading .env file:", err)
	}

	return Config{
		PostgresHost:     cast.ToString(getOrDefault("POSTGRES_HOST", "localhost")),
		PostgresPort:     cast.ToInt(getOrDefault("POSTGRES_PORT", 5432)),
		PostgresUser:     cast.ToString(getOrDefault("POSTGRES_USER", "postgres")),
		PostgresPassword: cast.ToString(getOrDefault("POSTGRES_PASSWORD", "password")),
		PostgresDatabase: cast.ToString(getOrDefault("POSTGRES_DATABASE", "movies_db")),
		JWTAccessSecret:  cast.ToString(getOrDefault("JWT_ACCESS_SECRET", "access_secret")),
		JWTRefreshSecret: cast.ToString(getOrDefault("JWT_REFRESH_SECRET", "refresh_secret")),

		ServiceName: cast.ToString(getOrDefault("SERVICE_NAME", "movies_service")),
	}
}

// getOrDefault returns environment variable value or a default
func getOrDefault(key string, defaultValue interface{}) interface{} {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
