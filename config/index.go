package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type AppConfig struct {
	Port       int              `yaml:"port"`
	Database   DatabaseConfig   `yaml:"database"`
	BitBrowser BitBrowserConfig `yaml:"bitBrowser"`
	AdsPower   AdsPowerConfig   `yaml:"adsPower"`
}

type DatabaseConfig struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Database string `yaml:"database"`
}

type BitBrowserConfig struct {
	ApiUrl   string `yaml:"apiUrl"`
	ApiToken string `yaml:"apiToken"`
}

type AdsPowerConfig struct {
	ApiUrl   string `yaml:"apiUrl"`
	ApiToken string `yaml:"apiToken"`
}

func LoadConfig() (*AppConfig, error) {
	// 读取文件内容
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		return nil, err
	}

	// 解析 YAML
	var cfg AppConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
