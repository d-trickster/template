package config

import (
	"fmt"
	"os"
	"time"
)

type Config struct {
	Env     string
	LogPath string
}

func MustLoad() *Config {
	return &Config{
		Env:     getEnv("APP_ENV", "dev"),
		LogPath: getEnv("APP_LOG_PATH", ""),
	}
}

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	DBName   string
	Password string
}

func MustLoadPostgres() *PostgresConfig {
	return &PostgresConfig{
		Host:     getEnv("POSTGRES_HOST", "localhost"),
		Port:     getEnv("POSTGRES_PORT", "5432"),
		User:     getEnv("POSTGRES_USER", "postgres"),
		DBName:   getEnv("POSTGRES_DBNAME", "postgres"),
		Password: getEnv("POSTGRES_PASSWORD", "password"),
	}
}

func (c *PostgresConfig) ConnString() string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s",
		c.User, c.Password, c.Host, c.Port, c.DBName,
	)
}

type CacheConfig struct {
	TTL       time.Duration
	RedisAddr string
}

func MustLoadCache() *CacheConfig {
	ttl, err := time.ParseDuration(getEnv("REDIS_CACHE_TTL", "10m"))
	if err != nil {
		panic(err)
	}
	return &CacheConfig{
		TTL:       ttl,
		RedisAddr: getEnv("REDIS_CACHE_ADDR", "localhost:6379"),
	}
}

func getEnv(key string, defaultValue string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultValue
}
