package utils

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

var config *Config

type Config struct {
	Port string `envconfig:"SERVER_PORT" default:"3000"`

	MySQL struct {
		Masters      []string `envconfig:"DB_MASTER_HOSTS"`
		Workers      []string `envconfig:"DB_WORKERS_HOSTS"`
		Port         string   `envconfig:"DB_PORT" default:"3306"`
		DBName       string   `envconfig:"DB_NAME"`
		User         string   `envconfig:"DB_USER"`
		Pass         string   `envconfig:"DB_PASS"`
		MaxIdleConns int      `envconfig:"DB_MAX_IDLE_CONNECTIONS" default:"2"`
		MaxOpenConns int      `envconfig:"DB_MAX_OPEN_CONNECTIONS" default:"5"`
	}
}

func init() {
	config = &Config{}

	// Pass env from system to Config struct
	err := envconfig.Process("", config)
	if err != nil {
		panic(fmt.Sprintf("Failed to decode config env: %v", err))
	}
}

// GetConfig .
func GetConfig() *Config {
	return config
}
