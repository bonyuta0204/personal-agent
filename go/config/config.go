package config

import (
	"fmt"
	"os"
)

// Config holds all configuration for the application
type Config struct {
	// Database configuration
	Database DatabaseConfig
	Memory   MemoryConfig
}

// MemoryConfig holds all memory related configuration
type MemoryConfig struct {
	// GitHub repository URL
	Repo string
}

// DatabaseConfig holds all database related configuration
type DatabaseConfig struct {
	User     string
	Password string
	Name     string
	Host     string
	Port     string
}

// LoadConfig loads the configuration from environment variables
func LoadConfig() (*Config, error) {
	config := &Config{
		Database: DatabaseConfig{
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Name:     os.Getenv("DB_NAME"),
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
		},
		Memory: MemoryConfig{
			Repo: os.Getenv("MEMORY_REPO"),
		},
	}

	// Set default values
	if config.Database.Port == "" {
		config.Database.Port = "5432" // Default PostgreSQL port
	}

	// Validate required fields
	if err := validateConfig(config); err != nil {
		return nil, err
	}

	return config, nil
}

// validateConfig validates that all required configuration is present
func validateConfig(config *Config) error {
	// Validate Database configuration
	if config.Database.User == "" {
		return fmt.Errorf("DB_USER is required")
	}
	if config.Database.Password == "" {
		return fmt.Errorf("DB_PASSWORD is required")
	}
	if config.Database.Name == "" {
		return fmt.Errorf("DB_NAME is required")
	}
	if config.Database.Host == "" {
		return fmt.Errorf("DB_HOST is required")
	}
	if config.Database.Port == "" {
		return fmt.Errorf("DB_PORT is required")
	}

	// Validate Memory configuration
	if config.Memory.Repo == "" {
		return fmt.Errorf("MEMORY_REPO is required")
	}

	return nil
}

// GetDSN returns the database connection string
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		c.User, c.Password, c.Host, c.Port, c.Name)
}
