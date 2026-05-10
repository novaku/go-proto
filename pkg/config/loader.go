package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// LoadConfig reads configuration from file and environment variables using Viper.
// Supports JSON, YAML, TOML, and other formats. Environment variables override file settings.
func LoadConfig(path string) (*Config, error) {
	v := viper.New()

	// Set file-related configurations
	dir, filename := splitPath(path)
	if dir != "" {
		v.AddConfigPath(dir)
	}
	v.SetConfigName(strings.TrimSuffix(filename, ".json"))
	v.SetConfigType("json")

	// set to use environment variables
	v.SetEnvPrefix("APP")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Read the config file
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Unmarshal into Config struct
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}

// LoadConfigWithViper provides low-level access to the Viper instance for advanced use cases.
// Returns both the Config struct and the Viper instance for direct access if needed.
func LoadConfigWithViper(path string) (*Config, *viper.Viper, error) {
	v := viper.New()

	// Set file-related configurations
	dir, filename := splitPath(path)
	if dir != "" {
		v.AddConfigPath(dir)
	}
	v.SetConfigName(strings.TrimSuffix(filename, ".json"))
	v.SetConfigType("json")

	// Set to use environment variables
	v.SetEnvPrefix("APP")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Read the config file
	if err := v.ReadInConfig(); err != nil {
		return nil, nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Unmarshal into Config struct
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, v, nil
}

// splitPath splits a file path into directory and filename
func splitPath(path string) (string, string) {
	if path == "" {
		return "", ""
	}
	lastSlash := strings.LastIndex(path, "/")
	if lastSlash == -1 {
		return "", path
	}
	return path[:lastSlash], path[lastSlash+1:]
}
