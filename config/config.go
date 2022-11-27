package config

import (
	"os"
	"strconv"
	"time"
)

const (
	serverAddressDefault      = ":8080"
	serverWriteTimeoutDefault = time.Second * 15
	serverReadTimeoutDefault  = time.Second * 10
	serverIdleTimeoutDefault  = time.Second * 60
	cacheCapacityDefault      = 500
)

type Config struct {
	Server ServerConfig
	Cache  CacheConfig
}

type ServerConfig struct {
	Address      string
	WriteTimeout time.Duration
	ReadTimeout  time.Duration
	IdleTimeout  time.Duration
}

type CacheConfig struct {
	Capacity int
}

func New() *Config {
	return &Config{
		Server: ServerConfig{
			Address:      getEnv("SERVER_ADDRESS", serverAddressDefault),
			WriteTimeout: getEnvAsDuration("SERVER_WRITE_TIMEOUT", serverWriteTimeoutDefault),
			ReadTimeout:  getEnvAsDuration("SERVER_READ_TIMEOUT", serverReadTimeoutDefault),
			IdleTimeout:  getEnvAsDuration("SERVER_IDLE_TIMEOUT", serverIdleTimeoutDefault),
		},
		Cache: CacheConfig{
			Capacity: getEnvAsInt("CACHE_CAPACITY", cacheCapacityDefault),
		},
	}
}

func getEnv(key string, defaultValue string) string {
	if val, exist := os.LookupEnv(key); exist {
		return val
	}

	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value, err := strconv.Atoi(getEnv(key, "")); err == nil {
		return value
	}

	return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	if value, err := strconv.Atoi(getEnv(key, "")); err == nil {
		return time.Second * time.Duration(value)
	}

	return defaultValue
}
