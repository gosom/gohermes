package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	ServerAddress string `split_words:"true" default:":8080"`
	// TODO add more http server configs

	PgHost   string `split_words:"true" default:"db"`
	PgPort   string `split_words:"true" default:"5432"`
	PgDb     string `split_words:"true" default:"demo"`
	PgUser   string `split_words:"true" default:"postgres"`
	PgPasswd string `split_words:"true" default:"postgres"`

	Debug bool `split_words:"true" default:"true"`
}

func (o Config) DSN() string {
	dsn := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s",
		o.PgHost, o.PgPort, o.PgDb, o.PgUser, o.PgPasswd)
	return dsn
}

func New(prefix string) (*Config, error) {
	var cfg Config
	err := envconfig.Process(prefix, &cfg)
	return &cfg, err
}
