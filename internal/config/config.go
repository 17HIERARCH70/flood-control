package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

// AppConfig structure containing the entire application configuration.
type AppConfig struct {
	Redis        RedisConfig        `yaml:"redis"`
	FloodControl FloodControlConfig `yaml:"floodControl"`
}

// RedisConfig structure for configuring a Redis server.
type RedisConfig struct {
	Address  string `yaml:"address"`
	Password string `yaml:"password,omitempty"`
	DB       int    `yaml:"db"`
}

// FloodControlConfig structure for flood control parameters.
type FloodControlConfig struct {
	RequestLimit  int `yaml:"requestLimit"`
	PeriodSeconds int `yaml:"periodSeconds"`
}

// LoadConfig function to load the configuration.
func LoadConfig() (*AppConfig, error) {
	configPath := fetchConfigPath()

	if configPath == "" {
		return nil, fmt.Errorf("path to config file is not provided")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file does not exist: %s", configPath)
	}

	var cfg AppConfig

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		return nil, fmt.Errorf("failed to read config: %v", err)
	}

	return &cfg, nil
}

// fetchConfigPath function to get the path to the configuration file.
func fetchConfigPath() string {
	var configPath string
	flag.StringVar(&configPath, "config", "", "path to the config file")
	flag.Parse()

	if configPath == "" {
		configPath = os.Getenv("CONFIG_PATH")
	}

	return configPath
}
