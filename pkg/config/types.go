package config

// Config is the root application configuration document.
type Config struct {
	Server      ServerConfig   `json:"server"`
	Database    DatabaseConfig `json:"database"`
	Redis       RedisConfig    `json:"redis"`
	JWT         JWTConfig      `json:"jwt"`
	Environment string         `json:"environment"`
}

// ServerConfig holds gRPC and HTTP listen settings.
type ServerConfig struct {
	Port     int    `json:"port"`
	HttpPort int    `json:"httpPort"`
	Name     string `json:"name"`
}

// DatabaseConfig holds MySQL connection parameters.
type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`
}

// RedisConfig holds optional cache backend settings.
type RedisConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Password string `json:"password"`
	DB       int    `json:"db"`
	Enabled  bool   `json:"enabled"`
}

// JWTConfig mirrors pkg/auth.JWTConfig JSON tags for file loading.
type JWTConfig struct {
	SecretKey     string `json:"secret_key"`
	TokenDuration int    `json:"token_duration"` // hours
	Issuer        string `json:"issuer"`
}
