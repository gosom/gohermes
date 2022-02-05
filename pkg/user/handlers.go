package user

import (
	"context"
	"fmt"
	"net/http"

	"github.com/casbin/casbin/v2"
	"github.com/gosom/gohermes/pkg/container"
	"github.com/gosom/gohermes/pkg/utils"
	"gorm.io/gorm"
)

type RegisterUserPayload struct {
	Email    string `json:"email" validate:"required,lte=100,email"`
	Password string `json:"password" validate:"required"`
}

func RegisterUserHandler(di *container.ServiceContainer) http.HandlerFunc {
	isrv, err := di.GetService("users")
	if err != nil {
		panic(err)
	}
	srv, err := CastIUserSrv(isrv)
	if err != nil {
		panic(err)
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var payload RegisterUserPayload
		if err := utils.Bind(r.Body, &payload, true); err != nil {
			apiErr := utils.NewBadRequestError(err.Error())
			utils.RenderJson(r, w, http.StatusBadRequest, apiErr)
			return
		}
		if err := utils.ValidatePassword(payload.Password); err != nil {
			apiErr := utils.NewBadRequestError(err.Error())
			utils.RenderJson(r, w, http.StatusBadRequest, apiErr)
			return
		}
		u, err := srv.Register(di.DB, payload.Email, payload.Password)
		// TODO errors should be correct
		if err != nil {
			apiErr := utils.NewInternalServerError(err.Error())
			utils.RenderJson(r, w, http.StatusBadRequest, apiErr)
			return
		}
		utils.RenderJson(r, w, http.StatusCreated, u)
	}
}

func GetHandler(di *container.ServiceContainer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := GetAuthenticatedUser(r.Context())
		utils.RenderJson(r, w, http.StatusOK, u)
	}
}

func Authentication(di *container.ServiceContainer, userFunc func(db *gorm.DB, s string) (IUser, error)) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			//logger := di.Logger.With().Str("request-id", utils.GetReqID(r.Context())).Logger()
			u, err := userFunc(di.DB, "g+1@gkomninos.com")
			if err != nil {
				panic(err)
			}
			ctx := context.WithValue(r.Context(), utils.Authenticated, u)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func Authorizer(di *container.ServiceContainer) func(next http.Handler) http.Handler {
	enfIface, err := di.GetService("enforcer")
	if err != nil {
		panic(err)
	}
	e := enfIface.(*casbin.Enforcer)
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			u := GetAuthenticatedUser(r.Context())
			e.LoadPolicy()
			canAccess, err := e.Enforce(fmt.Sprint(u.GetID()), r.URL.Path, r.Method)
			if err != nil {
				ae := utils.NewInternalServerError(err.Error())
				utils.RenderJson(r, w, ae.StatusCode, ae)
				return
			}
			if !canAccess {
				ae := utils.NewAuthenticationError("forbidden to access this resource")
				utils.RenderJson(r, w, ae.StatusCode, ae)
				return
			}
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
