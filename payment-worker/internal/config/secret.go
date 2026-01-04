package config

import (
	"os"
	"strconv"
	"time"
)

func GetSecret(key string) string {
	if os.Getenv(key) == "" {
		panic("db connection string not set")
	}
	return os.Getenv(key)
}

func GetEnvInt32(key string, def int32) int32 {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.ParseInt(v, 10, 32); err == nil {
			return int32(i)
		}
	}
	return def
}

func GetEnvDuration(key string, def time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return def
}
