package api

import "github.com/spf13/viper"

type Config struct {
	Port       int
	ProxyCount int
}

func InitConfig() (*Config, error) {
	config := &Config{
		Port:       viper.GetInt("Port"),
		ProxyCount: 0,
	}

	if config.Port == 0 {
		config.Port = 9092
	}

	return config, nil
}
