package app

import (
	"os"
	"strconv"

	"github.com/vinzryyy/iot-backend/database"
)

type Config struct {
	AppPort     string
	AppEnv      string
	DB          database.Config
	JWTSecret   string
	JWTExpHours int
}

func LoadConfig() Config {
	expHours, _ := strconv.Atoi(getEnv("JWT_EXPIRES_HOURS", "24"))
	return Config{
		AppPort: getEnv("APP_PORT", "8080"),
		AppEnv:  getEnv("APP_ENV", "development"),
		DB: database.Config{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "iot_db"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		JWTSecret:   getEnv("JWT_SECRET", "change-me-in-production"),
		JWTExpHours: expHours,
	}
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}
