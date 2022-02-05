package user

import (
	"context"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"

	"github.com/gosom/gohermes/pkg/utils"
)

func GetAuthenticatedUser(ctx context.Context) IUser {
	return ctx.Value(utils.Authenticated).(IUser)
}

var text = `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && keyMatch(r.obj,p.obj) && keyMatch(r.act,p.act)
`

func NewEnforcer(db *gorm.DB) (*casbin.Enforcer, error) {
	m, err := model.NewModelFromString(text)
	if err != nil {
		return nil, err
	}
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		return nil, err
	}
	enforcer, err := casbin.NewEnforcer(m, adapter)
	return enforcer, err
}
