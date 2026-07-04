package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Port    int      `yaml:"port"`
	Servers []string `yaml:"servers"`
}

func LoadConfig(filepath string) (*Config, error) {
	file, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(file, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
