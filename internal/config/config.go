package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/jieggii/groshi/internal/loggers"
	"github.com/jieggii/lookupcfg"
)

type Config struct {
	// server settings:
	Host         string `env:"GROSHI_HOST" $default:"0.0.0.0"`
	Port         int    `env:"GROSHI_PORT" $default:"8080"`
	JWTSecretKey []byte `env:"GROSHI_JWT_SECRET_KEY" $default:"secret-key"`

	// postgresql settings:
	PostgresHost     string `env:"GROSHI_POSTGRES_HOST" $default:"localhost"`
	PostgresPort     int    `env:"GROSHI_POSTGRES_PORT" $default:"5432"`
	PostgresUser     string `env:"GROSHI_POSTGRES_USER" $default:"jieggii"`
	PostgresPassword string `env:"GROSHI_POSTGRES_PASSWORD" $default:""`
	PostgresDatabase string `env:"GROSHI_POSTGRES_DATABASE" $default:"groshi"`
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
			"missing the following necessary environ variables: %v\n",
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
			"incorrect values of environmental variables: %v\n",
			strings.Join(incorrectTypeFieldsFmt, ", "),
		)
	}
	if !success {
		loggers.Error.Fatalln("exiting due to previous errors")
	}
	return config
}
