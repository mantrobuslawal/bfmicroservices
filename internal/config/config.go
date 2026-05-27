package config

import (
	"fmt"
	"os"
)

// Config contains runtime configuration for Catalogue Service.
type Config struct {
	Environment string
	LogLevel    string
	GRPCPort    string
	Database    DatabaseConfig
}

// DatabaseConfig contains MySQL connection configuration.
type DatabaseConfig struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
}

// Load reads configuration from environment variables.
func Load() (Config, error) {
	cfg := Config{
		Environment: getEnv("ENVIRONMENT", "local"),
		LogLevel:    getEnv("LOG_LEVEL", "info"),
		GRPCPort:    getEnv("CATALOG_SERVICE_GRPC_PORT", "50051"),
		Database: DatabaseConfig{
			Host:     getEnv("MYSQL_HOST", "localhost"),
			Port:     getEnv("MYSQL_PORT", "3306"),
			Name:     getEnv("MYSQL_DATABASE", "bfstore_catalog"),
			User:     getEnv("MYSQL_USER", "bfstore_catalog_user"),
			Password: getEnv("MYSQL_PASSWORD", "bfstore_catalog_password"),
		},
	}

	if cfg.Database.Password == "" {
		return Config{}, fmt.Errorf("MYSQL_PASSWORD must be set")
	}

	return cfg, nil
}

// DSN returns a MySQL data source name suitable for database/sql.
func (d DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true&multiStatements=false&charset=utf8mb4,utf8",
		d.User,
		d.Password,
		d.Host,
		d.Port,
		d.Name,
	)
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
