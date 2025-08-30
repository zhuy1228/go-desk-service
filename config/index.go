package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type AppConfig struct {
	Port         int            `yaml:"port"`
	Database     DatabaseConfig `yaml:"database"`
	WorkerID     int64          `yaml:"workerID"`
	DatacenterID int64          `yaml:"datacenterID"`
}

type DatabaseConfig struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Database string `yaml:"database"`
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
