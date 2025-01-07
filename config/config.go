package config

import (
	"errors"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type DatabaseConfig struct {
	User      string `yaml:"user"`
	Password  string `yaml:"password"`
	Host      string `yaml:"host"`
	Port      int    `yaml:"port"`
	DBName    string `yaml:"dbname"`
	URL       string `yaml:"url"`
	SessionId string `yaml:"sessionid"`
	Auth      string `yaml:"auth"`
	Server    string `yaml:"server"`
}

type Config struct {
	Server map[string]DatabaseConfig `yaml:"Server"`
	Token  string                    `yaml:"token"`
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

	_, exists := config.Server[edge]
	if !exists {
		return nil, errors.New("server configuration not found for edge: " + edge)
	}

	return config, nil
}
