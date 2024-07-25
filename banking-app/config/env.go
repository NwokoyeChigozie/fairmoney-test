package config

type Configuration struct {
	ServerPort        int
	MongoDb           MongoDb
	ThirdPartyBaseUrl string
}

type MongoDb struct {
	ConnectionString string
	DbName           string
}

type EnvModel struct {
	SERVER_PORT               int    `mapstructure:"SERVER_PORT"`
	MONGODB_CONNECTION_STRING string `mapstructure:"MONGODB_CONNECTION_STRING"`
	MONGODB_DATABASE_NAME     string `mapstructure:"MONGODB_DATABASE_NAME"`
	THIRD_PARTY_BASE_URL      string `mapstructure:"THIRD_PARTY_BASE_URL"`
}

func (env *EnvModel) UpdateConfiguration() *Configuration {
	return &Configuration{
		ServerPort:        env.SERVER_PORT,
		ThirdPartyBaseUrl: env.THIRD_PARTY_BASE_URL,
		MongoDb: MongoDb{
			ConnectionString: env.MONGODB_CONNECTION_STRING,
			DbName:           env.MONGODB_DATABASE_NAME,
		},
	}
}
