package config

import (
	"net/url"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DB    DBConfig
	HTTP  HTTPConfig
	JWT   JWTConfig
	Redis RedisConfig
	SMTP  SMTPConfig
}

type JWTConfig struct {
	Secret     string
	AccessTTL  time.Duration
	RefreshTTL time.Duration
}

type RedisConfig struct {
	Host string
	Port int
}

func (config RedisConfig) Address() string {
	return config.Host + ":" + strconv.Itoa(config.Port)
}

type DBConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	Name     string
	SSLMode  string
	DATABASE_URL *url.URL
}

func (config DBConfig) URL() string {
	return (config.DATABASE_URL).String()
}

type HTTPConfig struct {
	Host string
	Port int
}

func (config HTTPConfig) Address() string {
	return config.Host + ":" + strconv.Itoa(config.Port)
}

type SMTPConfig struct {
	Password string
	Host     string
	Port     int
	User     string
	UseSSL   bool
	DevMode  bool
}

func Load() (Config, error) {
	if err := godotenv.Load(); err != nil {
		// Environment variables are also supported in containers.
		// Missing .env is therefore not a configuration error.
	}

	return Config{
		DB: DBConfig{
			DATABASE_URL: envURL("DATABASE_URL", "postgres://getshelf:getshelf@localhost:5432/getshelf?sslmode=disable"),
		},
		HTTP: HTTPConfig{
			Host: env("HTTP_HOST", "0.0.0.0"),
			Port: envInt("HTTP_PORT", 8080),
		},
		JWT: JWTConfig{
			Secret:     env("JWT_SECRET", "change-me-in-production"),
			AccessTTL:  time.Duration(envInt("JWT_ACCESS_TTL_MINUTES", 15)) * time.Minute,
			RefreshTTL: time.Duration(envInt("JWT_REFRESH_TTL_HOURS", 24)) * time.Hour,
		},
		Redis: RedisConfig{
			Host: env("REDIS_HOST", "localhost"),
			Port: envInt("REDIS_PORT", 6379),
		},
		SMTP: SMTPConfig{
			Password: env("SMTP_PASSWORD", "change-me-in-production"),
			Host:     env("SMTP_HOST", "change-me-in-production"),
			Port:     envInt("SMTP_PORT", 465),
			User:     env("SMTP_USER", "change-me-in-production"),
			UseSSL:   envBool("SMTP_USESSL", true),
			DevMode:  envBool("SMTP_DEV_MODE", true),
		},
	}, nil
}
