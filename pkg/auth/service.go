package auth

import (
	"github.com/casbin/casbin/v2"
	"github.com/gosom/gohermes/pkg/container"
)

func IUserSrvFromDi(di *container.ServiceContainer) IUserSrv {
	iface, err := di.GetService("users")
	if err != nil {
		panic(err)
	}
	srv := iface.(IUserSrv)
	return srv
}

func EnforcerFromDi(di *container.ServiceContainer) *casbin.Enforcer {
	i, err := di.GetService("enforcer")
	if err != nil {
		panic(err)
	}
	return i.(*casbin.Enforcer)
}
