package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	DatabaseURL      string `mapstructure:"DATABASE_URL"`
	HTTPServerPort   string `mapstructure:"HTTP_SERVER_PORT"`
	ArgonMemory      uint32 `mapstructure:"ARGON_MEMORY"`
	ArgonIterations  uint32 `mapstructure:"ARGON_ITERATIONS"`
	ArgonParallelism uint8  `mapstructure:"ARGON_PARALLELISM"`
	ArgonSaltLength  uint32 `mapstructure:"ARGON_SALT_LENGTH"`
	ArgonKeyLength   uint32 `mapstructure:"ARGON_KEY_LENGTH"`
}

func LoadConfigFromEnvFile() (config Config, err error) {
	viper.AddConfigPath(".")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)

	return
}
