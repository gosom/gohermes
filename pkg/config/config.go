package config

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	ServerAddress string `envconfig:"SERVER_ADDRESS" default:":8080"`
	// TODO add more http server configs

	PostgresHost     string `envconfig:"POSTGRES_HOST" default:"db"`
	PostgresPort     string `envconfig:"POSTGRES_PORT" default:"5432"`
	PostgresDb       string `envconfig:"POSTGRES_DB" default:"demo"`
	PostgresUser     string `envconfig:"POSTGRES_USER" default:"postgres"`
	PostgresPassword string `envconfig:"POSTGRES_PASSWORD" default:"postgres"`

	TokenIssuer          string        `envconfig:"TOKEN_ISSUER" default:"myapp"`
	TokenSecret          string        `envconfig:"TOKEN_SECRET" default:"something secret"`
	AccessTokenDuration  time.Duration `envconfig:"ACCESS_TOKEN_DURATION" default:"15m"`
	RefreshTokenDuration time.Duration `envconfig:"REFRRESH_TOKEN_DURATION" default:"30m"`

	Debug bool `split_words:"true" default:"true"`
}

func (o Config) DSN() string {
	dsn := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s",
		o.PostgresHost, o.PostgresPort, o.PostgresDb, o.PostgresUser, o.PostgresPassword)
	return dsn
}

func New(prefix string) (*Config, error) {
	var cfg Config
	err := envconfig.Process(prefix, &cfg)
	return &cfg, err
}
