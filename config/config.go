package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Configuration interface {
	Load(filename string) error
	GetDatabaseConfig(serverName string) (DatabaseConfig, error)
	GetNSXtConfig(serverName string) (NSXtConfig, error)
}

type DatabaseConfig struct {
	User     string
	Password string
	Host     string
	Port     int
	DBName   string
}

type NSXtConfig struct {
	SessionID string
	Auth      string
	URL       string
}

type YAMLConfig struct {
	Databases   map[string]DatabaseConfig `yaml:"databases"`
	NSXTServers map[string]NSXtConfig     `yaml:"nsxt_servers"`
}

func (c *YAMLConfig) Load(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(c); err != nil {
		return fmt.Errorf("failed to decode YAML: %w", err)
	}
	return nil
}

func (c *YAMLConfig) GetDatabaseConfig(serverName string) (DatabaseConfig, error) {
	config, exists := c.Databases[serverName]
	if !exists {
		return DatabaseConfig{}, fmt.Errorf("database server '%s' not found", serverName)
	}
	return config, nil
}

func (c *YAMLConfig) GetNSXtConfig(serverName string) (NSXtConfig, error) {
	config, exists := c.NSXTServers[serverName]
	if !exists {
		return NSXtConfig{}, fmt.Errorf("NSX-T server '%s' not found", serverName)
	}
	return config, nil
}
