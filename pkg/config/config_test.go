package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		wantErr  bool
		validate func(*testing.T, *Config)
	}{
		{
			name: "valid config",
			content: `{
				"server": {
					"port": 50051,
					"httpPort": 8080,
					"name": "test-server"
				},
				"database": {
					"host": "localhost",
					"port": 3306,
					"user": "root",
					"password": "password",
					"dbname": "testdb"
				},
				"redis": {
					"host": "localhost",
					"port": 6379,
					"password": "",
					"db": 0,
					"enabled": true
				},
				"environment": "test"
			}`,
			wantErr: false,
			validate: func(t *testing.T, cfg *Config) {
				if cfg.Server.Port != 50051 {
					t.Errorf("Server.Port = %d; want 50051", cfg.Server.Port)
				}
				if cfg.Server.HttpPort != 8080 {
					t.Errorf("Server.HttpPort = %d; want 8080", cfg.Server.HttpPort)
				}
				if cfg.Server.Name != "test-server" {
					t.Errorf("Server.Name = %s; want test-server", cfg.Server.Name)
				}
				if cfg.Database.Host != "localhost" {
					t.Errorf("Database.Host = %s; want localhost", cfg.Database.Host)
				}
				if cfg.Database.Port != 3306 {
					t.Errorf("Database.Port = %d; want 3306", cfg.Database.Port)
				}
				if cfg.Redis.Enabled != true {
					t.Error("Redis.Enabled = false; want true")
				}
				if cfg.Environment != "test" {
					t.Errorf("Environment = %s; want test", cfg.Environment)
				}
			},
		},
		{
			name: "minimal config",
			content: `{
				"server": {"port": 9000, "httpPort": 9001, "name": "minimal"},
				"database": {"host": "db", "port": 5432, "user": "user", "password": "pass", "dbname": "db"},
				"redis": {"host": "redis", "port": 6379, "password": "", "db": 0, "enabled": false},
				"environment": "dev"
			}`,
			wantErr: false,
			validate: func(t *testing.T, cfg *Config) {
				if cfg.Server.Port != 9000 {
					t.Errorf("Server.Port = %d; want 9000", cfg.Server.Port)
				}
				if cfg.Redis.Enabled != false {
					t.Error("Redis.Enabled = true; want false")
				}
			},
		},
		{
			name:    "invalid json",
			content: `{"server": invalid}`,
			wantErr: true,
		},
		{
			name:    "empty json",
			content: `{}`,
			wantErr: false,
			validate: func(t *testing.T, cfg *Config) {
				if cfg == nil {
					t.Error("Config is nil")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary config file
			tmpDir := t.TempDir()
			tmpFile := filepath.Join(tmpDir, "config.json")

			err := os.WriteFile(tmpFile, []byte(tt.content), 0644)
			if err != nil {
				t.Fatalf("Failed to create temp config file: %v", err)
			}

			// Load the config
			cfg, err := LoadConfig(tmpFile)

			if (err != nil) != tt.wantErr {
				t.Errorf("LoadConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.validate != nil {
				tt.validate(t, cfg)
			}
		})
	}
}

func TestLoadFromReaderDeprecated(t *testing.T) {
	// Note: LoadFromReader is deprecated and no longer available with Viper.
	// Use LoadConfig() instead, which now supports environment variable overrides.
	t.Skip("LoadFromReader is deprecated - use LoadConfig instead")
}

func TestLoadConfig_FileNotFound(t *testing.T) {
	_, err := LoadConfig("/nonexistent/path/config.json")
	if err == nil {
		t.Error("Expected error for non-existent file, got nil")
	}
}

func TestLoadConfig_InvalidPath(t *testing.T) {
	_, err := LoadConfig("")
	if err == nil {
		t.Error("Expected error for empty path, got nil")
	}
}

func TestServerConfig(t *testing.T) {
	cfg := ServerConfig{
		Port:     50051,
		HttpPort: 8080,
		Name:     "test-server",
	}

	if cfg.Port != 50051 {
		t.Errorf("Port = %d; want 50051", cfg.Port)
	}
	if cfg.HttpPort != 8080 {
		t.Errorf("HttpPort = %d; want 8080", cfg.HttpPort)
	}
	if cfg.Name != "test-server" {
		t.Errorf("Name = %s; want test-server", cfg.Name)
	}
}

func TestDatabaseConfig(t *testing.T) {
	cfg := DatabaseConfig{
		Host:     "localhost",
		Port:     3306,
		User:     "root",
		Password: "password",
		DBName:   "testdb",
	}

	if cfg.Host != "localhost" {
		t.Errorf("Host = %s; want localhost", cfg.Host)
	}
	if cfg.Port != 3306 {
		t.Errorf("Port = %d; want 3306", cfg.Port)
	}
	if cfg.User != "root" {
		t.Errorf("User = %s; want root", cfg.User)
	}
	if cfg.Password != "password" {
		t.Errorf("Password = %s; want password", cfg.Password)
	}
	if cfg.DBName != "testdb" {
		t.Errorf("DBName = %s; want testdb", cfg.DBName)
	}
}

func TestRedisConfig(t *testing.T) {
	cfg := RedisConfig{
		Host:     "localhost",
		Port:     6379,
		Password: "redis-pass",
		DB:       1,
		Enabled:  true,
	}

	if cfg.Host != "localhost" {
		t.Errorf("Host = %s; want localhost", cfg.Host)
	}
	if cfg.Port != 6379 {
		t.Errorf("Port = %d; want 6379", cfg.Port)
	}
	if cfg.Password != "redis-pass" {
		t.Errorf("Password = %s; want redis-pass", cfg.Password)
	}
	if cfg.DB != 1 {
		t.Errorf("DB = %d; want 1", cfg.DB)
	}
	if !cfg.Enabled {
		t.Error("Enabled = false; want true")
	}
}

func TestConfig_AllFields(t *testing.T) {
	cfg := Config{
		Server: ServerConfig{
			Port:     50051,
			HttpPort: 8080,
			Name:     "test",
		},
		Database: DatabaseConfig{
			Host:     "localhost",
			Port:     3306,
			User:     "root",
			Password: "pass",
			DBName:   "test",
		},
		Redis: RedisConfig{
			Host:     "localhost",
			Port:     6379,
			Password: "",
			DB:       0,
			Enabled:  true,
		},
		Environment: "test",
	}

	if cfg.Server.Port != 50051 {
		t.Error("Server config not properly set")
	}
	if cfg.Database.Host != "localhost" {
		t.Error("Database config not properly set")
	}
	if cfg.Redis.Port != 6379 {
		t.Error("Redis config not properly set")
	}
	if cfg.Environment != "test" {
		t.Error("Environment not properly set")
	}
}

func TestLoadConfig_WithDifferentEnvironments(t *testing.T) {
	tests := []struct {
		name    string
		env     string
		content string
	}{
		{
			name: "production environment",
			env:  "production",
			content: `{
				"server": {"port": 50051, "httpPort": 8080, "name": "prod"},
				"database": {"host": "prod-db", "port": 3306, "user": "user", "password": "pass", "dbname": "proddb"},
				"redis": {"host": "prod-redis", "port": 6379, "password": "pass", "db": 0, "enabled": true},
				"environment": "production"
			}`,
		},
		{
			name: "development environment",
			env:  "development",
			content: `{
				"server": {"port": 50051, "httpPort": 8080, "name": "dev"},
				"database": {"host": "localhost", "port": 3306, "user": "root", "password": "root", "dbname": "devdb"},
				"redis": {"host": "localhost", "port": 6379, "password": "", "db": 0, "enabled": false},
				"environment": "development"
			}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			tmpFile := filepath.Join(tmpDir, "config.json")

			err := os.WriteFile(tmpFile, []byte(tt.content), 0644)
			if err != nil {
				t.Fatalf("Failed to create temp config file: %v", err)
			}

			cfg, err := LoadConfig(tmpFile)
			if err != nil {
				t.Fatalf("LoadConfig() error = %v", err)
			}

			if cfg.Environment != tt.env {
				t.Errorf("Environment = %s; want %s", cfg.Environment, tt.env)
			}
		})
	}
}
