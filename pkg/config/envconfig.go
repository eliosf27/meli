package settings

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"os"
)

type (
	Config struct {
		ProjectName    string `default:"meli"`
		ProjectVersion string `default:"0.0.1"`
		Environment    string `envconfig:"ENVIRONMENT" default:"prod"`
		BaseEndpoint   string `envconfig:"BASE_ENDPOINT" required:"true"`
		Redis          RedisSpecification
		Postgres       PostgresSpecification
	}

	RedisSpecification struct {
		RedisHost string `envconfig:"REDIS_HOST" required:"true"`
		RedisPort string `envconfig:"REDIS_PORT" required:"true"`
	}

	PostgresSpecification struct {
		ItemsConnection string `envconfig:"POSTGRES_ITEMS_CONNECTION" required:"true"`
	}
)

var (
	Configs Config
)

const (
	environmentKey     = "ENVIRONMENT"
	DevelopEnvironment = "dev"
)

func NewConfig() Config {
	Configs.Environment = os.Getenv(environmentKey)
	if Configs.Environment == DevelopEnvironment {
		_ = godotenv.Load(".env")
	}
	if err := envconfig.Process("", &Configs); err != nil {
		panic(err.Error())
	}

	return Configs
}
