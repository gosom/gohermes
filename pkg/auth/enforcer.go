package auth

import (
	"database/sql"
	"io"
	"io/ioutil"

	sqladapter "github.com/Blank-Xu/sql-adapter"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
)

const defaultCasbinPolicy = `
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

func NewEnforcer(db *sql.DB, r io.Reader, a interface{}) (*casbin.Enforcer, error) {
	m, err := getCasbinPolicy(r)
	if err != nil {
		return nil, err
	}
	var adapter interface{}
	if a == nil {
		adapter, err = sqladapter.NewAdapter(db, "postgres", "")
		if err != nil {
			return nil, err
		}
	}
	enforcer, err := casbin.NewEnforcer(m, adapter)
	if err != nil {
		return nil, err
	}
	return enforcer, err
}

func getCasbinPolicy(r io.Reader) (model.Model, error) {
	var casbinPolicy string
	if r == nil {
		casbinPolicy = defaultCasbinPolicy
	} else {
		policyBytes, err := ioutil.ReadAll(r)
		if err != nil {
			return nil, err
		}
		casbinPolicy = string(policyBytes)
	}
	return model.NewModelFromString(casbinPolicy)
}
