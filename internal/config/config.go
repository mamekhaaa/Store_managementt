package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	HTTPAddr          string
	DatabaseURL       string
	DBMaxOpenConns    int
	DBMaxIdleConns    int
	DBConnMaxLifetime time.Duration
	ShutdownTimeout   time.Duration
}

func Load() Config {
	return Config{
		HTTPAddr:          envString("HTTP_ADDR", ":8080"),
		DatabaseURL:       envString("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"),
		DBMaxOpenConns:    envInt("DB_MAX_OPEN_CONNS", 10),
		DBMaxIdleConns:    envInt("DB_MAX_IDLE_CONNS", 5),
		DBConnMaxLifetime: time.Duration(envInt("DB_CONN_MAX_LIFETIME_SECONDS", 300)) * time.Second,
		ShutdownTimeout:   time.Duration(envInt("SHUTDOWN_TIMEOUT_SECONDS", 10)) * time.Second,
	}
}

func envString(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func envInt(key string, fallback int) int {
	v, err := strconv.Atoi(os.Getenv(key))
	if err != nil {
		return fallback
	}
	return v
}
