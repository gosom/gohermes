package scheduler

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	ServerAddress string `envconfig:"SCHEDULER_SERVER_ADDRESS" default:":50051"`

	PostgresHost     string `envconfig:"SCHEDULER_POSTGRES_HOST" default:"localhost"`
	PostgresPort     string `envconfig:"SCHEDULER_POSTGRES_PORT" default:"5432"`
	PostgresDb       string `envconfig:"SCHEDULER_POSTGRES_DB" default:"scheduler"`
	PostgresUser     string `envconfig:"SCHEDULER_POSTGRES_USER" default:"postgres"`
	PostgresPassword string `envconfig:"SCHEDULER_POSTGRES_PASSWORD" default:"postgres"`

	Secret string `envconfig:"SCHEDULER_SECRET" default:"secret"`
	ApiKey string `envconfig:"SCHEDULER_API_KEY" default:"apikey"`

	Debug     bool `envconfig:"SCHEDULER_DEBUG" default:"true"`
	Executors int  `envconfig:"SCHEDULER_EXECUTORS" default:"10"`
}

func (o Config) DSN() string {
	dsn := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s",
		o.PostgresHost, o.PostgresPort, o.PostgresDb, o.PostgresUser, o.PostgresPassword)
	return dsn
}

func NewConfig(prefix string) (*Config, error) {
	var cfg Config
	err := envconfig.Process(prefix, &cfg)
	return &cfg, err
}
