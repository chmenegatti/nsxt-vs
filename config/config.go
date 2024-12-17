package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Databases   map[string]DatabaseConfig `yaml:"databases"`
	NSXTServers map[string]NSXtConfig     `yaml:"nsxt_servers"`
	Token       string                    `yaml:"token"`
}

type DatabaseConfig struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	DBName   string `yaml:"dbname"`
}

type NSXtConfig struct {
	SessionID string `yaml:"session_id"`
	Auth      string `yaml:"auth"`
	URL       string `yaml:"url"`
}

func LoadConfig(filename string) (*Config, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cfg Config
	decoder := yaml.NewDecoder(f)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (c *Config) GetDatabaseConfig(serverName string) (DatabaseConfig, error) {
	cfg, exists := c.Databases[serverName]
	if !exists {
		return DatabaseConfig{}, fmt.Errorf("database server '%s' not found in configuration", serverName)
	}
	return cfg, nil
}

func (c *Config) GetNSXtConfig(serverName string) (NSXtConfig, error) {
	cfg, exists := c.NSXTServers[serverName]
	if !exists {
		return NSXtConfig{}, fmt.Errorf("NSX-T server '%s' not found in configuration", serverName)
	}
	return cfg, nil
}

func (c *Config) GetToken() string {
	return c.Token
}
