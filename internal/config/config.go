package config

import (
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	GitHub   GitHubConfig
	Session  SessionConfig
	Token    TokenConfig
}

type TokenConfig struct {
	Key []byte // 32 bytes for AES-256; nil disables encryption (dev only)
}

type ServerConfig struct {
	Port string
	Env  string
}

type DatabaseConfig struct {
	URL string
}

type RedisConfig struct {
	URL string
}

type GitHubConfig struct {
	ClientID     string
	ClientSecret string
	CallbackURL  string
}

type SessionConfig struct {
	Secret string
	MaxAge int
}

func Load() (*Config, error) {
	var tokenKey []byte
	if raw := getEnv("TOKEN_ENCRYPTION_KEY", ""); raw != "" {
		k, err := hex.DecodeString(raw)
		if err != nil {
			return nil, fmt.Errorf("TOKEN_ENCRYPTION_KEY: must be hex string: %w", err)
		}
		tokenKey = k
	}

	cfg := &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
			Env:  getEnv("ENV", "development"),
		},
		Database: DatabaseConfig{
			URL: getEnv("DATABASE_URL", "postgres://gitrpg:gitrpg@localhost:5433/gitrpg?sslmode=disable"),
		},
		Redis: RedisConfig{
			URL: getEnv("REDIS_URL", "redis://localhost:6380"),
		},
		GitHub: GitHubConfig{
			ClientID:     getEnv("GITHUB_CLIENT_ID", ""),
			ClientSecret: getEnv("GITHUB_CLIENT_SECRET", ""),
			CallbackURL:  getEnv("GITHUB_CALLBACK_URL", "http://localhost:8080/auth/github/callback"),
		},
		Session: SessionConfig{
			Secret: getEnv("SESSION_SECRET", ""),
			MaxAge: getEnvInt("SESSION_MAX_AGE", 86400*30),
		},
		Token: TokenConfig{Key: tokenKey},
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) validate() error {
	if c.Server.Env == "production" {
		if c.GitHub.ClientID == "" {
			return fmt.Errorf("GITHUB_CLIENT_ID required in production")
		}
		if c.GitHub.ClientSecret == "" {
			return fmt.Errorf("GITHUB_CLIENT_SECRET required in production")
		}
		if c.Session.Secret == "" {
			return fmt.Errorf("SESSION_SECRET required in production")
		}
		if len(c.Token.Key) == 0 {
			return fmt.Errorf("TOKEN_ENCRYPTION_KEY required in production")
		}
	}
	if len(c.Token.Key) > 0 && len(c.Token.Key) != 32 {
		return fmt.Errorf("TOKEN_ENCRYPTION_KEY must be 32 bytes (64 hex chars), got %d bytes", len(c.Token.Key))
	}
	return nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		return fallback
	}
	return i
}
