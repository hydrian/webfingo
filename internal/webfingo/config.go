package webfingo

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config represents the application configuration
type Config struct {
	DB       DBConfig       `json:"db"`
	Keycloak KeycloakConfig `json:"keycloak"`
}

// DBConfig holds database connection details
type DBConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

// KeycloakConfig holds Keycloak connection details
type KeycloakConfig struct {
	KeycloakHost string `json:"keycloak-host"`
}

// Load reads and parses the configuration file
func LoadConfig(filePath string) (*Config, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("error parsing config file: %w", err)
	}

	return &config, nil
}

// GetDBConfig returns the database configuration
func (c *Config) GetDBConfig() DBConfig {
	return c.DB
}

// GetKeycloakConfig returns the Keycloak configuration
func (c *Config) GetKeycloakConfig() KeycloakConfig {
	return c.Keycloak
}

// GetKeycloakHost returns the configured Keycloak host
func (c *Config) GetKeycloakHost() string {
	return c.Keycloak.KeycloakHost
}

// GetDBConnectionString returns a formatted PostgreSQL connection string
func (c *Config) GetDBConnectionString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		c.DB.User,
		c.DB.Password,
		c.DB.Host,
		c.DB.Port,
		c.DB.Name,
	)
}
