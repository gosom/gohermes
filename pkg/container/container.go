package container

import (
	"errors"
	"sync"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"

	"github.com/gosom/gohermes/pkg/config"
)

var ErrServiceNotFound = errors.New("service not found")

type ServiceContainer struct {
	Cfg      *config.Config
	Logger   zerolog.Logger
	DB       *gorm.DB
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

func (o *ServiceContainer) AutoMigrate(dst ...interface{}) error {
	return o.DB.AutoMigrate(dst...)
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
	var gormConfig gorm.Config
	if ans.Cfg.Debug {
		level = zerolog.DebugLevel
		gormConfig.Logger = gormLogger.Default.LogMode(gormLogger.Info)
	} else {
		level = zerolog.InfoLevel
	}
	ans.Logger = log.Level(level)
	// if you set environment variable PG_HOST='' then do not initialize postgres connection
	if len(ans.Cfg.PgHost) > 0 {
		// TODO make gorm use zerolog
		ans.DB, err = gorm.Open(postgres.Open(ans.Cfg.DSN()), &gormConfig)
		if err != nil {
			return
		}
	}
	return
}

func GetLogger(key, value string) zerolog.Logger {
	return log.With().Str(key, value).Logger()
}
