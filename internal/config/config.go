package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/jieggii/groshi/internal/loggers"
	"github.com/jieggii/lookupcfg"
)

type Config struct {
	// server settings
	Host         string `env:"GROSHI_HOST"`
	Port         int    `env:"GROSHI_PORT"`
	JWTSecretKey []byte `env:"GROSHI_JWT_SECRET_KEY"`

	// postgresql settings
	PostgresHost     string `env:"GROSHI_POSTGRES_HOST"`
	PostgresPort     int    `env:"GROSHI_POSTGRES_PORT"`
	PostgresUser     string `env:"GROSHI_POSTGRES_USER"`
	PostgresPassword string `env:"GROSHI_POSTGRES_PASSWORD"`
	PostgresDatabase string `env:"GROSHI_POSTGRES_DATABASE"`
}

func ReadFromEnv() *Config {
	config := &Config{}
	result := lookupcfg.PopulateConfig("env", os.LookupEnv, config)
	success := true

	if len(result.MissingFields) != 0 {
		success = false

		var envVarNames []string
		for _, field := range result.MissingFields {
			envVarNames = append(envVarNames, field.SourceName)
		}
		loggers.Error.Printf(
			"missing the following necessary environmental variables: %v\n",
			strings.Join(envVarNames, ", "),
		)
	}
	if len(result.IncorrectTypeFields) != 0 {
		success = false

		var incorrectTypeFieldsFmt []string
		for _, field := range result.IncorrectTypeFields {
			incorrectTypeFieldsFmt = append(
				incorrectTypeFieldsFmt,
				fmt.Sprintf(
					"%v (got `%v`, but expected value type: %v)",
					field.SourceName,
					field.RawValue,
					field.ExpectedValueType.String(),
				),
			)
		}
		loggers.Error.Printf(
			"incorrect values of the following environmental variables: %v\n",
			strings.Join(incorrectTypeFieldsFmt, ", "),
		)
	}
	if !success {
		loggers.Error.Fatalln("exiting due to previous errors")
	}
	return config
}
