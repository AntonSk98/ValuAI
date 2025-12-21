package common

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// DatabaseConfig holds the minimal database configuration
type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Database string `yaml:"database"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	SSLMode  string `yaml:"sslmode"`
}

// DatabaseYAML represents the YAML structure
type DatabaseYAML struct {
	Database DatabaseConfig `yaml:"database"`
}

// CreateConnection creates a PostgreSQL connection from YAML config file
func CreateConnection(configPath string) (*sql.DB, error) {
	var yamlConfig DatabaseYAML
	if err := LoadProperties(configPath, &yamlConfig); err != nil {
		return nil, fmt.Errorf("failed to load database config: %w", err)
	}

	// Build connection string
	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		yamlConfig.Database.Username,
		yamlConfig.Database.Password,
		yamlConfig.Database.Host,
		yamlConfig.Database.Port,
		yamlConfig.Database.Database,
		yamlConfig.Database.SSLMode,
	)

	// Open database connection
	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}
