package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

// AppConfig структура, содержащая всю конфигурацию приложения.
type AppConfig struct {
	Redis        RedisConfig        `yaml:"redis"`
	FloodControl FloodControlConfig `yaml:"floodControl"`
}

// RedisConfig структура для конфигурации Redis сервера.
type RedisConfig struct {
	Address  string `yaml:"address"`
	Password string `yaml:"password,omitempty"`
	DB       int    `yaml:"db"`
}

// FloodControlConfig структура для параметров контроля флуда.
type FloodControlConfig struct {
	RequestLimit  int `yaml:"requestLimit"`
	PeriodSeconds int `yaml:"periodSeconds"`
}

// LoadConfig функция для загрузки конфигурации.
func LoadConfig() (*AppConfig, error) {
	configPath := fetchConfigPath()

	// Проверяем, указан ли путь к файлу конфигурации
	if configPath == "" {
		return nil, fmt.Errorf("path to config file is not provided")
	}

	// Проверяем, существует ли файл конфигурации
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file does not exist: %s", configPath)
	}

	var cfg AppConfig
	// Читаем конфигурацию из файла
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		return nil, fmt.Errorf("failed to read config: %v", err)
	}

	return &cfg, nil
}

// fetchConfigPath функция для получения пути к файлу конфигурации.
func fetchConfigPath() string {
	var configPath string
	flag.StringVar(&configPath, "config", "", "path to the config file")
	flag.Parse()

	// Если флаг не указан, пробуем получить путь из переменной окружения
	if configPath == "" {
		configPath = os.Getenv("CONFIG_PATH")
	}

	return configPath
}
