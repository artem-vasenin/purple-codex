package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Addr, DatabaseURL, JWTSecret string
	AccessTTL, RefreshTTL        time.Duration
	CookieSecure                 bool
}

func Load() (Config, error) {
	c := Config{Addr: env("APP_ADDR", ":8080"), DatabaseURL: env("DATABASE_URL", "host=localhost user=postgres password=postgres dbname=app port=5433 sslmode=disable"), JWTSecret: os.Getenv("JWT_SECRET"), AccessTTL: duration("ACCESS_TTL", 15*time.Minute), RefreshTTL: duration("REFRESH_TTL", 30*24*time.Hour), CookieSecure: boolValue("COOKIE_SECURE", false)}
	if c.JWTSecret == "" {
		return Config{}, fmt.Errorf("JWT_SECRET is required")
	}
	return c, nil
}
func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
func duration(key string, fallback time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if parsed, err := time.ParseDuration(value); err == nil {
			return parsed
		}
	}
	return fallback
}
func boolValue(key string, fallback bool) bool {
	if value := os.Getenv(key); value != "" {
		parsed, err := strconv.ParseBool(value)
		if err == nil {
			return parsed
		}
	}
	return fallback
}
