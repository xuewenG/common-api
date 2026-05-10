package config

import (
	"log"
	"os"

	"github.com/samber/do/v2"
	"gopkg.in/yaml.v3"
)

type Config struct {
	LogLevel   string     `yaml:"log_level"`
	Port       string     `yaml:"port"`
	CorsOrigin string     `yaml:"cors_origin"`
	Bark       BarkConfig `yaml:"bark"`
	Holidays   []Holiday  `yaml:"holidays"`
}

type Holiday struct {
	Name string `yaml:"name"`
	Date string `yaml:"date"`
}

type BarkConfig struct {
	ServerURL string `yaml:"server_url"`
	DeviceKey string `yaml:"device_key"`
	Username  string `yaml:"username"`
	Password  string `yaml:"password"`
}

func NewConfig(i do.Injector) (*Config, error) {
	configBytes, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("Read config failed, %v\n", err)
	}

	var config Config
	err = yaml.Unmarshal(configBytes, &config)
	if err != nil {
		log.Fatalf("Decode config failed: %v\n", err)
	}

	if config.Port == "" {
		config.Port = "8080"
	}

	return &config, nil
}
