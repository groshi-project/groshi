package config

type GroshiConfig struct { // todo: parse config source according to tags, escape from code replication
	Host string `env:"GROSHI_HOST"`
	Port int    `env:"GROSHI_PORT"`

	MongoHost   string `env:"GROSHI_MONGO_HOST"`
	MongoPort   int    `env:"GROSHI_MONGO_PORT"`
	MongoDBName string `env:"GROSHI_MONGO_DB_NAME"`

	SuperuserPassword string `env:"GROSHI_SUPERUSER_PASSWORD"`
}

func ReadFromEnv() (error, *GroshiConfig) {
	//var missing []string
	return nil, nil
}
