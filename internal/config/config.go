package config

import (
	"encoding/json"
	"os"
)

// read mssql connection string from config.json

type Config struct {
	Port int    `json:"port"`
	DSN  string `json:"dsn"`

	Limiter struct {
		Rps     float64 `json:"rps"`
		Burst   int     `json:"burst"`
		Enabled bool    `json:"enabled"`
	} `json:"limiter"`

	HTTPClient struct {
		RequestTimeout int `json:"request_timeout"`
	} `json:"http_client"`
	Log struct {
		Level  string `json:"level"`
		Format string `json:"format"`
	}
}

func ReadConfigFile(path string) (*Config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	config := &Config{}
	if err := json.Unmarshal(b, &config); err != nil {
		return nil, err
	}

	return config, nil
}
