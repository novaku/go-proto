# Viper Configuration Guide

Your application now uses **Viper** for configuration management, which provides powerful features for loading config from files and environment variables.

## Features

- ✅ Load configuration from JSON files
- ✅ Override settings via environment variables
- ✅ Support for multiple formats (JSON, YAML, TOML, HCL)
- ✅ Nested configuration with easy access
- ✅ Type-safe configuration through structs

## Configuration Priority

Environment variables **override** file settings:
1. Environment variables (highest priority)
2. Config file values (lowest priority)

## Environment Variables

Environment variables use the **`APP_`** prefix and support nested keys with underscores.

### Format
```
APP_<SECTION>_<KEY>=value
```

### Examples

```bash
# Server settings
APP_SERVER_PORT=8000
APP_SERVER_HTTPPORT=9000

# Database settings
APP_DATABASE_HOST=prod.db.example.com
APP_DATABASE_PORT=3306
APP_DATABASE_USER=produser
APP_DATABASE_PASSWORD=securepass
APP_DATABASE_DBNAME=proddb

# Redis settings
APP_REDIS_HOST=redis.example.com
APP_REDIS_PORT=6379
APP_REDIS_PASSWORD=redispass
APP_REDIS_DB=1
APP_REDIS_ENABLED=true

# JWT settings
APP_JWT_SECRETKEY=your-secret-key
APP_JWT_TOKENDURATION=24
APP_JWT_ISSUER=guestbook-service

# Environment selection
APP_ENV=production
APP_ENVIRONMENT=production
```

## Running with Environment Variables

### Local development
```bash
./server
```

### Production with environment overrides
```bash
APP_DATABASE_HOST=prod.db.com \
APP_DATABASE_PASSWORD=securepass \
APP_JWT_SECRETKEY=prod-secret-key \
APP_ENV=production \
./server
```

### Using .env file (if implemented with godotenv)
Create a `.env` file:
```env
APP_DATABASE_HOST=mydb.example.com
APP_DATABASE_PASSWORD=secretpass
APP_JWT_SECRETKEY=my-secret
```

Then load before running:
```bash
source .env
./server
```

## Config Files

- **Development**: `pkg/config/config.local.json`
- **Production**: `pkg/config/config.prod.json`
- **Selection**: Based on `APP_ENV` environment variable

## Code Usage

### Basic Configuration Loading
```go
cfg, err := config.LoadConfig("pkg/config/config.local.json")
if err != nil {
    log.Fatal(err)
}

// Access config values
port := cfg.Server.Port
dbHost := cfg.Database.Host
```

### Advanced: Direct Viper Access
```go
cfg, v, err := config.LoadConfigWithViper("pkg/config/config.local.json")
if err != nil {
    log.Fatal(err)
}

// Use Viper directly for advanced features
secretValue := v.GetString("jwt.secret_key")
portValue := v.GetInt("server.port")
```

## Configuration Struct Tags

The configuration uses JSON tags for file mapping:

```go
type Config struct {
    Server      ServerConfig   `json:"server"`
    Database    DatabaseConfig `json:"database"`
    Redis       RedisConfig    `json:"redis"`
    JWT         JWTConfig      `json:"jwt"`
    Environment string         `json:"environment"`
}
```

## Migration from Old System

The refactoring maintains backward compatibility:
- Old function `LoadFromReader()` is deprecated
- Use `LoadConfig(path)` instead
- All existing code continues to work without changes
- Environment variable support added automatically

## Benefits

1. **No code changes needed** - Existing configuration files work as-is
2. **Environment variable support** - Perfect for containers and cloud deployments
3. **Type safety** - Configuration strongly typed with Go structs
4. **Flexible** - Can switch to YAML/TOML/HCL any time
5. **Industry standard** - Viper is widely used in Go applications

## Troubleshooting

### Config not loading
- Verify file path is correct
- Check file permissions
- Ensure JSON is valid

### Environment variables not working
- Must use **`APP_`** prefix
- Use underscores for nested keys
- Example: `APP_DATABASE_HOST` not `APP_DATABASE.HOST`

### Case sensitivity
- Environment variable keys are converted to lowercase
- Use underscores in env vars: `APP_JWT_SECRETKEY` → `jwt.secretkey`
