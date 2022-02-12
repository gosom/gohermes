package auth

import (
	"net/http"
	"strconv"

	"github.com/casbin/casbin/v2"
	"github.com/gosom/gohermes/pkg/container"
	"github.com/gosom/gohermes/pkg/utils"
)

func Authorization(di *container.ServiceContainer) func(next http.Handler) http.Handler {
	e := EnforcerFromDi(di)
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			u := GetAuthenticatedUser(r.Context())
			canAccess, err := canAccess(r, u, e)
			if err != nil {
				ae := utils.NewInternalServerError(err.Error())
				utils.RenderJson(r, w, ae.StatusCode, ae)
				return
			}
			if !canAccess {
				ae := utils.NewAuthorizationError("forbidden to access this resource")
				utils.RenderJson(r, w, ae.StatusCode, ae)
				return
			}
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

func canAccess(r *http.Request, u AuthenticatedUser, e *casbin.Enforcer) (bool, error) {
	if err := e.LoadPolicy(); err != nil { // probably better to use filtered policy
		return false, err
	}
	canAccessWithUID, err := e.Enforce(strconv.Itoa(u.ID), r.URL.Path, r.Method)
	if err == nil && canAccessWithUID {
		return true, nil
	}
	if err != nil {
		return false, err
	}
	requests := make([][]interface{}, len(u.Roles), len(u.Roles))
	for i := range u.Roles {
		requests[i] = []interface{}{u.Roles[i], r.URL.Path, r.Method}
	}
	boolArr, err := e.BatchEnforce(requests)
	if err != nil {
		return false, err
	}
	for _, ok := range boolArr {
		if ok {
			return true, nil
		}
	}
	return false, nil
}
