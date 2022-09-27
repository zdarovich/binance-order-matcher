package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	OrderMatcherService struct {
		URL string
	}

	Order struct {
		Type   int
		Symbol string
		Size   float64
		Price  float64
	}

	Database struct {
		URL string
	}
}

func Init() (*Config, error) {
	viper.SetConfigType("yaml")
	viper.SetConfigFile("./config/config.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	constant := &Config{}
	err = viper.Unmarshal(constant)
	if err != nil {
		return nil, err
	}
	return constant, nil
}
