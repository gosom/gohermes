package services

import (
	"github.com/gosom/gohermes/pkg/auth"
	"github.com/gosom/gohermes/pkg/container"

	"github.com/gosom/gohermes/examples/todo/user"
)

func RegisterServices(di *container.ServiceContainer) error {
	enforcer, err := auth.NewEnforcer(di.DB, nil, nil)
	if err != nil {
		return err
	}
	di.RegisterService("enforcer", enforcer)

	users := user.NewUserService(di)
	di.RegisterService("users", users)

	return nil
}
