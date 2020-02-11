package db

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	DatabaseDialect string
	DatabaseURI string
}

func InitConfig() (*Config, error) {
	config := &Config{
		DatabaseDialect: viper.GetString("DatabaseDialect"),
		DatabaseURI: viper.GetString("DatabaseURI"),
	}

	if config.DatabaseDialect == "" {
		config.DatabaseDialect = "sqlite3"
	}

	if config.DatabaseURI == "" {
		return nil, fmt.Errorf("DatabaseURI must be set")
	}

	return config, nil
}
