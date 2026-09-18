package config

import (
	"net/url"
	"os"
	"strconv"
)

func lookupEnv(key string) string {
	return os.Getenv(key)
}

func env(key, fallback string) string {
	if value := lookupEnv(key); value != "" {
		return value
	}
	return fallback
}

func envInt(key string, fallback int) int {
	value, err := strconv.Atoi(env(key, strconv.Itoa(fallback)))
	if err != nil {
		return fallback
	}
	return value
}

func envBool(key string, fallback bool) bool {
	value, err := strconv.ParseBool(env(key, strconv.FormatBool(fallback)))
	if err != nil {
		return fallback
	}
	return value
}

func envURL(key, fallback string) *url.URL {
	value, err := url.Parse(env(key, fallback))
	if err != nil {
		return nil
	}
	return value
}
