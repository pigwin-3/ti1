package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Database struct {
		Host     string `json:"host"`
		Port     string `json:"port"`
		User     string `json:"user"`
		Password string `json:"password"`
		DBName   string `json:"dbname"`
		SSLMode  string `json:"sslmode"`
	} `json:"database"`
	Valkey struct {
		Host      string `json:"host"`
		Port      string `json:"port"`
		MaxConns  int    `json:"max_conns"`
		TimeoutMs int    `json:"timeout_ms"`
	} `json:"valkey"`
	Temp string `json:"temp"`
}

func LoadConfig(file string) (Config, error) {
	var config Config
	configFile, err := os.Open(file)
	if err != nil {
		return config, fmt.Errorf("failed to open config file: %w", err)
	}
	defer configFile.Close()

	jsonParser := json.NewDecoder(configFile)
	if err := jsonParser.Decode(&config); err != nil {
		return config, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Override with environment variables if they are set
	if host := os.Getenv("DB_HOST"); host != "" {
		config.Database.Host = host
	}
	if port := os.Getenv("DB_PORT"); port != "" {
		config.Database.Port = port
	}
	if user := os.Getenv("DB_USER"); user != "" {
		config.Database.User = user
	}
	if password := os.Getenv("DB_PASSWORD"); password != "" {
		config.Database.Password = password
	}
	if dbname := os.Getenv("DB_NAME"); dbname != "" {
		config.Database.DBName = dbname
	}
	if sslmode := os.Getenv("DB_SSLMODE"); sslmode != "" {
		config.Database.SSLMode = sslmode
	}
	if temp := os.Getenv("TEMP"); temp != "" {
		config.Temp = temp
	}

	// Override Valkey settings with environment variables
	if valkeyHost := os.Getenv("VALKEY_HOST"); valkeyHost != "" {
		config.Valkey.Host = valkeyHost
	}
	if valkeyPort := os.Getenv("VALKEY_PORT"); valkeyPort != "" {
		config.Valkey.Port = valkeyPort
	}
	if maxConns := os.Getenv("VALKEY_MAX_CONNS"); maxConns != "" {
		if val, err := strconv.Atoi(maxConns); err == nil {
			config.Valkey.MaxConns = val
		}
	}
	if timeoutMs := os.Getenv("VALKEY_TIMEOUT_MS"); timeoutMs != "" {
		if val, err := strconv.Atoi(timeoutMs); err == nil {
			config.Valkey.TimeoutMs = val
		}
	}

	return config, nil
}
