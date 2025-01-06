package config

import (
	"errors"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type DatabaseConfig struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	DBName   string `yaml:"dbname"`
}

type NSXTConfig struct {
	URL       string `yaml:"url"`
	Auth      string `yaml:"auth"`
	SessionId string `yaml:"session_id"`
}

type Config struct {
	Databases   map[string]DatabaseConfig `yaml:"databases"`
	NSXTServers map[string]NSXTConfig     `yaml:"nsxt_servers"`
	Token       string                    `yaml:"token"`
}

func LoadConfig(edge string, logger *zap.Logger) (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		logger.Error("Failed to read config file", zap.Error(err))
		return nil, err
	}

	var config *Config
	if err := viper.Unmarshal(&config); err != nil {
		logger.Error("Failed to unmarshal config", zap.Error(err))
		return nil, err
	}

	_, exists := config.Databases[edge]
	if !exists {
		return nil, errors.New("database configuration not found for edge: " + edge)
	}

	_, exists = config.NSXTServers[edge]
	if !exists {
		return nil, errors.New("nsxt configuration not found for edge: " + edge)
	}

	return config, nil
}
