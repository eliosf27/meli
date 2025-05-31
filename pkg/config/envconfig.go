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
		Port           string `envconfig:"PORT" default:"8000" required:"true"`
		Environment    string `envconfig:"ENVIRONMENT" default:"prod" required:"true"`
		BaseEndpoint   string `envconfig:"BASE_ENDPOINT" required:"true"`
		Redis          RedisSpecification
		Postgres       PostgresSpecification
		Elasticsearch  ElasticsearchSpecification
	}

	ElasticsearchSpecification struct {
		Addresses []string `envconfig:"ELASTICSEARCH_ADDRESSES" required:"true"` // Comma-separated list
		Username  string   `envconfig:"ELASTICSEARCH_USERNAME"`
		Password  string   `envconfig:"ELASTICSEARCH_PASSWORD"`
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
	DevelopEnvironment = "develop"
	TestingEnvironment = "testing"
)

func NewConfig() Config {
	Configs.Environment = os.Getenv(environmentKey)
	if Configs.Environment == DevelopEnvironment {
		_ = godotenv.Load(".env")
	}
	if Configs.Environment == TestingEnvironment {
		_ = godotenv.Load(".env.testing")
	}
	if err := envconfig.Process("", &Configs); err != nil {
		panic(err.Error())
	}

	return Configs
}
