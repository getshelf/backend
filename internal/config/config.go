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
}

func (config DBConfig) URL() string {
	return (&url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(config.Username, config.Password),
		Host:   config.Host + ":" + strconv.Itoa(config.Port),
		Path:   config.Name,
		RawQuery: url.Values{
			"sslmode": []string{config.SSLMode},
		}.Encode(),
	}).String()
}

type HTTPConfig struct {
	Host string
	Port int
}

func (config HTTPConfig) Address() string {
	return config.Host + ":" + strconv.Itoa(config.Port)
}

func Load() (Config, error) {
	if err := godotenv.Load(); err != nil {
		// Environment variables are also supported in containers.
		// Missing .env is therefore not a configuration error.
	}

	return Config{
		DB: DBConfig{
			Host:     env("DB_HOST", "localhost"),
			Port:     envInt("DB_PORT", 5432),
			Username: env("DB_USERNAME", "getshelf"),
			Password: env("DB_PASSWORD", "getshelf"),
			Name:     env("DB_NAME", "getshelf"),
			SSLMode:  env("DB_SSLMODE", "disable"),
		},
		HTTP: HTTPConfig{
			Host: env("HTTP_HOST", "0.0.0.0"),
			Port: envInt("HTTP_PORT", 8080),
		},
		JWT: JWTConfig{
			Secret:     env("JWT_SECRET", "change-me-in-production"),
			AccessTTL:  time.Duration(envInt("JWT_ACCESS_TTL_MINUTES", 15)) * time.Minute,
			RefreshTTL: time.Duration(envInt("JWT_REFRESH_TTL_HOURS", 168)) * time.Hour,
		},
		Redis: RedisConfig{Host: env("REDIS_HOST", "localhost"), Port: envInt("REDIS_PORT", 6379)},
	}, nil
}
