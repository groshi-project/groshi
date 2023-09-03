package config

import (
	"github.com/groshi-project/groshi/internal/loggers"
	"github.com/jieggii/lookupcfg"
	"os"
	"strings"
)

// EnvVars represents environmental variables that are used to configure groshi.
type EnvVars struct {
	// groshi settings:
	Host    string `env:"GROSHI_HOST"`    // host to be listened to
	Port    int    `env:"GROSHI_PORT"`    // port to be listened to
	Swagger bool   `env:"GROSHI_SWAGGER"` // toggle Swagger API documentation at `/docs/index.html` route
	Debug   bool   `env:"GROSHI_DEBUG"`   // toggle debug mode (influences on gin mode)

	JWTSecretKeyFile    string `env:"GROSHI_JWT_SECRET_KEY_FILE"`   // file containing secret key for generating JWTs
	ExchangeRatesAPIKey string `env:"GROSHI_EXCHANGERATES_API_KEY"` // file containing API key for https://exchangeratesapi.io/

	// MongoDB settings:
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
