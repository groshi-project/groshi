package config

import (
	"github.com/jieggii/groshi/groshi/logger"
	"github.com/jieggii/lookupcfg"
	"os"
)

type Config struct {
	Host         string `env:"GROSHI_HOST"`
	Port         int    `env:"GROSHI_PORT"`
	JWTSecretKey []byte `env:"GROSHI_JWT_SECRET_KEY"`

	SuperuserUsername string `env:"GROSHI_SUPERUSER_USERNAME"`
	SuperuserPassword string `env:"GROSHI_SUPERUSER_PASSWORD"`

	PostgresHost     string `env:"GROSHI_POSTGRES_HOST"`
	PostgresPort     int    `env:"GROSHI_POSTGRES_PORT"`
	PostgresUser     string `env:"GROSHI_POSTGRES_USER"`
	PostgresPassword string `env:"GROSHI_POSTGRES_PASSWORD"`
	PostgresDatabase string `env:"GROSHI_POSTGRES_DATABASE"`
}

func ReadFromEnv() *Config {
	config := &Config{}
	result := lookupcfg.PopulateConfig("env", os.LookupEnv, config)
	ok := true

	if len(result.MissingFields) != 0 {
		ok = false
		logger.Fatal.Printf("Missing: %v\n", result.MissingFields)
	}
	if len(result.IncorrectTypeFields) != 0 {
		ok = false
		logger.Fatal.Printf("Incorrect type: %v\n", result.IncorrectTypeFields)
	}
	if !ok {
		logger.Fatal.Fatalln("Exiting according to the previous errors.")
	}
	return config
}
