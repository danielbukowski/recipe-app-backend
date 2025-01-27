package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL      string `env:"DATABASE_URL,notEmpty"`
	HTTPServerPort   string `env:"HTTP_SERVER_PORT,notEmpty"`
	ArgonMemory      uint32 `env:"ARGON_MEMORY,notEmpty"`
	ArgonIterations  uint32 `env:"ARGON_ITERATIONS,notEmpty"`
	ArgonParallelism uint8  `env:"ARGON_PARALLELISM,notEmpty"`
	ArgonSaltLength  uint32 `env:"ARGON_SALT_LENGTH,notEmpty"`
	ArgonKeyLength   uint32 `env:"ARGON_KEY_LENGTH,notEmpty"`
	AppEnv           string `env:"APP_ENV,notEmpty"`
	MemcachedServer  string `env:"MEMCACHE_SERVER,notEmpty"`
}

func LoadEnvironmentVariablesToConfig() (cfg Config, err error) {
	// if .env file does not exist then just ignore it.
	_ = godotenv.Load()

	// Read environment variables from OS and validate them.
	err = env.Parse(&cfg)

	return
}
