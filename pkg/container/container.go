package container

import (
	"database/sql"
	"errors"
	"sync"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/volatiletech/sqlboiler/v4/boil"

	"github.com/gosom/gohermes/pkg/config"
)

var ErrServiceNotFound = errors.New("service not found")

type ServiceContainer struct {
	Cfg      *config.Config
	Logger   zerolog.Logger
	DB       *sql.DB
	Lock     *sync.RWMutex
	Services map[string]interface{}
}

func (o *ServiceContainer) GetService(name string) (interface{}, error) {
	o.Lock.RLock()
	defer o.Lock.RUnlock()
	s, ok := o.Services[name]
	if !ok {
		return nil, ErrServiceNotFound
	}
	return s, nil
}

func (o *ServiceContainer) RegisterService(name string, s interface{}) {
	o.Lock.Lock()
	defer o.Lock.Unlock()
	o.Services[name] = s
}

func NewDefault() (ans *ServiceContainer, err error) {
	ans = &ServiceContainer{
		Lock:     &sync.RWMutex{},
		Services: make(map[string]interface{}),
	}
	ans.Cfg, err = config.New("")
	if err != nil {
		return
	}
	var level zerolog.Level
	if ans.Cfg.Debug {
		level = zerolog.DebugLevel
		boil.DebugMode = true
	} else {
		level = zerolog.InfoLevel
	}
	ans.Logger = log.Level(level)
	ans.DB, err = sql.Open("pgx", ans.Cfg.DSN())
	if err != nil {
		return
	}
	err = ans.DB.Ping()
	return
}

func GetLogger(key, value string) zerolog.Logger {
	return log.With().Str(key, value).Logger()
}
