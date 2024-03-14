package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// RedisConfig represents the Redis configuration
type RedisConfig struct {
	Params map[string]string
}

// NewRedisConfig creates a new instance of RedisConfig
func NewRedisConfig() *RedisConfig {
	newConfig := RedisConfig{
		Params: make(map[string]string),
	}
	if err := newConfig.ReadConfig("redis.conf"); err != nil {
		newConfig.setDefaultConfig()
	}

	// // Print configuration parameters
	// fmt.Println("Redis configuration parameters:")
	// for key, value := range newConfig.Params {
	// 	slog.Debug(fmt.Sprintf("%s: %s\n", key, value))
	// }
	return &newConfig
}

// ReadConfig reads the Redis configuration from a file
func (rc *RedisConfig) ReadConfig(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			// Configuration file doesn't exist, create a new one
			rc.setDefaultConfig()
			return rc.WriteConfig(filename)
		}
		return err
	}
	defer file.Close()

	var key, value string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) >= 1 {
			key = parts[0]
			value = strings.Join(parts[1:], " ")
		}
		rc.Params[key] = value
	}
	return nil
}

// WriteConfig writes the Redis configuration to a file
func (rc *RedisConfig) WriteConfig(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write configuration parameters to file
	for key, value := range rc.Params {
		_, err := fmt.Fprintf(file, "%s %s\n", key, value)
		if err != nil {
			return err
		}
	}

	return nil
}

// setDefaultConfig sets default Redis configuration parameters
func (rc *RedisConfig) setDefaultConfig() {
	// Add more default configuration parameters as needed
	rc.Params["port"] = "6379"
	rc.Params["bind"] = "127.0.0.1"
	rc.Params["timeout"] = "300"
	rc.Params["loglevel"] = "notice"
	rc.Params["databases"] = "16"
	rc.Params["maxclients"] = "10000"
	rc.Params["maxmemory"] = "100mb"
	rc.Params["maxmemory-policy"] = "volatile-lru"
	rc.Params["save"] = "3600 1 300 100 60 10000"
	rc.Params["appendonly"] = "no"
}

func (rc *RedisConfig) getParam(key string) (string, bool) {
	key = strings.ToLower(key)
	paramValue, found := rc.Params[key]
	if !found {
		return "Unknown", false
	}
	return paramValue, true
}

// Usage example
var redisConfig = NewRedisConfig()
