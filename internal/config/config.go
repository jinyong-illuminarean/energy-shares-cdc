package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	DebeziumURL     string                 `json:"debezium_url"`
	ConnectorConfig map[string]interface{} `json:"connector_config"`
}

func LoadConfig(env string, dbPassword string) (*Config, error) {
	filename := fmt.Sprintf("config_%s.json", env)
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	var cfg Config
	if err := json.NewDecoder(file).Decode(&cfg); err != nil {
		return nil, err
	}

	// 민감한 정보를 환경 변수에서 로드
	cfg.ConnectorConfig["database.password"] = dbPassword

	return &cfg, nil
}

func (c *Config) Compare(config map[string]interface{}) bool {
	for k, v := range config {
		if c.ConnectorConfig[k] != v {
			return false
		}
	}
	return true
}
