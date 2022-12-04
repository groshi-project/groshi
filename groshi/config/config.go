package config

import (
	"github.com/jieggii/groshi/groshi/logger"
	"os"
	"reflect"
	"strconv"
)

type Config struct { // todo: parse config source according to tags, escape from code replication
	Host string `env:"GROSHI_HOST"`
	Port int    `env:"GROSHI_PORT"`

	MongoHost   string `env:"GROSHI_MONGO_HOST"`
	MongoPort   int    `env:"GROSHI_MONGO_PORT"`
	MongoDBName string `env:"GROSHI_MONGO_DB_NAME"`

	SuperuserPassword string `env:"GROSHI_SUPERUSER_PASSWORD"`
}

func ReadFromEnv() *Config {
	config := Config{
		Host: "0.0.0.0",
		Port: 8080,

		MongoHost:   "localhost",
		MongoPort:   27017,
		MongoDBName: "groshi",

		SuperuserPassword: "password123",
	} // todo
	var missingEnvVars []string
	var mustBeIntVars []string

	configType := reflect.TypeOf(config)
	for i := 0; i < configType.NumField(); i++ {
		field := configType.Field(i)
		envVarName, envVarNameFound := field.Tag.Lookup("env")
		if !envVarNameFound {
			logger.Fatal.Panicf("Config field %v does not have env tag.", field.Name)
		}
		value, envVarFound := os.LookupEnv(envVarName)
		if !envVarFound {
			missingEnvVars = append(missingEnvVars, envVarName)
			continue
		}
		fieldValueObj := reflect.ValueOf(&config).Elem().Field(i)
		fieldType := fieldValueObj.Type().Name()
		if fieldType == "int" {
			n, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				mustBeIntVars = append(mustBeIntVars, envVarName)
			} else {
				fieldValueObj.SetInt(n)
			}
		} else if fieldType == "string" {
			fieldValueObj.SetString(value)
		} else {
			logger.Fatal.Panicf("Unimplemented Config struct type %v.", fieldType)
		}
	}

	//if len(missingEnvVars)+len(mustBeIntVars) != 0 {
	//	fmt.Printf("Missing: %v.\n", missingEnvVars)
	//	fmt.Printf("Must be int: %v.\n", mustBeIntVars)
	//	panic("Env config error")
	//} // todo
	return &config
}
