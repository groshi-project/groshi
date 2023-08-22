package config

import (
	"github.com/jieggii/groshi/internal/loggers"
	"github.com/jieggii/lookupcfg"
	"os"
	"strings"
)

// EnvVars TODO
type EnvVars struct {
	// application settings
	Host                string `env:"GROSHI_HOST"`
	Port                int    `env:"GROSHI_PORT"`
	JWTSecretKeyFile    string `env:"GROSHI_JWT_SECRET_KEY_FILE"`
	ExchangeRatesAPIKey string `env:"GROSHI_EXCHANGERATES_API_KEY"`

	// MongoDB settings
	MongoHost string `env:"GROSHI_MONGO_HOST"`
	MongoPort int    `env:"GROSHI_MONGO_PORT"`

	MongoUsernameFile string `env:"GROSHI_MONGO_USERNAME_FILE"`
	MongoPasswordFile string `env:"GROSHI_MONGO_PASSWORD_FILE"`
	MongoDatabaseFile string `env:"GROSHI_MONGO_DATABASE_FILE"`
}

// handleConfigPopulationError TODO
func catchConfigPopulationError(result *lookupcfg.ConfigPopulationResult) {
	die := false
	if len(result.MissingFields) != 0 {
		loggers.Error.Printf("missing fields environmental variables: %v", result.MissingFields)
		die = true
	}
	if len(result.IncorrectTypeFields) != 0 {
		loggers.Error.Printf("incorrect values of environmental variables: %v", result.IncorrectTypeFields)
		die = true
	}
	if die {
		loggers.Error.Fatal("exiting due to previous errors")
	}
}

// ReadEnvVars TODO
func ReadEnvVars() *EnvVars {
	config := EnvVars{}
	result := lookupcfg.PopulateConfig("env", os.LookupEnv, &config)
	catchConfigPopulationError(result)
	return &config
}

// ReadDockerSecret TODO
func ReadDockerSecret(filePath string) string {
	data, err := os.ReadFile(filePath)
	if err != nil {
		loggers.Error.Fatalf("error reading docker secret file: %v", err)
	}
	content := string(data)
	if len(content) == 0 {
		loggers.Error.Fatalf("empty docker secret file %v", filePath)
	}
	content = strings.TrimRight(content, "\n")
	return content
}
