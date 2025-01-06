package config

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/valkey-io/valkey-go"
)

type ValkeyConfig struct {
	Host      string `json:"host"`
	Port      string `json:"port"`
	MaxConns  int    `json:"max_conns"`
	TimeoutMs int    `json:"timeout_ms"`
}

func LoadValkeyConfig(file string) (ValkeyConfig, error) {
	var config ValkeyConfig
	configFile, err := os.Open(file)
	if err != nil {
		return config, fmt.Errorf("failed to open config file: %w", err)
	}
	defer configFile.Close()

	if err := json.NewDecoder(configFile).Decode(&config); err != nil {
		return config, fmt.Errorf("failed to parse Valkey config: %w", err)
	}

	// Override with environment variables if set
	if host := os.Getenv("VALKEY_HOST"); host != "" {
		config.Host = host
	}
	if port := os.Getenv("VALKEY_PORT"); port != "" {
		config.Port = port
	}
	if maxConns := os.Getenv("VALKEY_MAX_CONNS"); maxConns != "" {
		if val, err := strconv.Atoi(maxConns); err == nil {
			config.MaxConns = val
		}
	}
	if timeoutMs := os.Getenv("VALKEY_TIMEOUT_MS"); timeoutMs != "" {
		if val, err := strconv.Atoi(timeoutMs); err == nil {
			config.TimeoutMs = val
		}
	}

	return config, nil
}

func ConnectToValkey(configPath string) (valkey.Client, error) {
	fmt.Println("Loading Valkey configuration...")
	valkeyConfig, err := LoadValkeyConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load Valkey config: %v", err)
	}
	fmt.Println("Valkey configuration loaded successfully!")

	// Setup Valkey client options
	options := valkey.ClientOption{
		InitAddress: []string{fmt.Sprintf("%s:%s", valkeyConfig.Host, valkeyConfig.Port)},
		// Additional options can be added here if required
	}

	fmt.Printf("Connecting to Valkey at %s:%s...\n", valkeyConfig.Host, valkeyConfig.Port)
	client, err := valkey.NewClient(options)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Valkey: %v", err)
	}

	// Optionally, perform a ping to validate the connection
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(valkeyConfig.TimeoutMs))
	defer cancel()

	if err := client.Do(ctx, client.B().Ping().Build()).Error(); err != nil {
		client.Close()
		return nil, fmt.Errorf("failed to ping Valkey: %v", err)
	}

	log.Println("Connected to Valkey successfully!")
	return client, nil
}

func DisconnectFromValkey(client valkey.Client) error {
	fmt.Println("Disconnecting from Valkey...")
	client.Close()
	log.Println("Disconnected from Valkey successfully!")
	return nil
}
