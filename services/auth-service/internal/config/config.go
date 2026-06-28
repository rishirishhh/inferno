package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	App AppConfig
	DB  DBConfig
	JWT JWTConfig
}

type AppConfig struct {
	Env  string
	Port string
}

type DBConfig struct {
	URL string
}

type JWTConfig struct {
	Secret               string
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
	BcryptCost           int
}

func Load() (*Config, error) {
	cfg := &Config{
		App: AppConfig{
			Env:  getEnv("APP_ENV", "development"),
			Port: getEnv("PORT", "3000"),
		},
		DB: DBConfig{
			URL: os.Getenv("DATABASE_URL"),
		},
		JWT: JWTConfig{
			Secret: os.Getenv("JWT_SECRET"),
		},
	}

	if cfg.DB.URL == "" {
		return nil, errors.New("Database_URL is required")
	}

	if cfg.JWT.Secret == "" {
		return nil, errors.New("JWT_SECRET is required")
	}

	var err error

	cfg.JWT.AccessTokenDuration, err = time.ParseDuration(getEnv("JWT_ACCESS_TOKEN_DURATION", "15m"))
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_ACCESS_TOKEN_DURATION: %w", err)
	}

	cfg.JWT.BcryptCost, err = strconv.Atoi(getEnv("BCRYPT_COST", "12"))
	if err != nil {
		return nil, fmt.Errorf("invalid BCRYPT_COST: %w", err)
	}

	if cfg.JWT.BcryptCost < 10 || cfg.JWT.BcryptCost > 15 {
		return nil, fmt.Errorf("bcrypt cost must be between 10 and 15")
	}

	return cfg, nil

}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}
