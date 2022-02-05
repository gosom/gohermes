package controllers

import (
	"net/http"

	"github.com/gosom/gohermes/pkg/container"
	"github.com/gosom/gohermes/pkg/user"
	"github.com/gosom/gohermes/pkg/utils"

	"github.com/gosom/gohermes/examples/todo/models"
	"github.com/gosom/gohermes/examples/todo/payloads"
	"github.com/gosom/gohermes/examples/todo/services"
)

func CreateTodo(di *container.ServiceContainer) http.HandlerFunc {
	service := services.GetFromDi(di)
	return func(w http.ResponseWriter, r *http.Request) {
		var payload payloads.CreateTodoPayload
		if err := utils.Bind(r.Body, &payload, true); err != nil {
			ae := utils.NewBadRequestError(err.Error())
			utils.RenderJson(r, w, ae.StatusCode, ae)
			return
		}
		payload.User = user.GetAuthenticatedUser(r.Context()).(models.CustomUser)

		todo, err := service.Create(payload)
		if err != nil {
			ae := utils.NewInternalServerError(err.Error())
			utils.RenderJson(r, w, ae.StatusCode, ae)
			return
		}

		utils.RenderJson(r, w, http.StatusOK, todo)
	}
}

func GetTodoByID(di *container.ServiceContainer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user_id, err := utils.GetUintUrlParam(r, "user_id")
		if err != nil {
			ae := err.(*utils.ApiError)
			utils.RenderJson(r, w, ae.StatusCode, ae)
			return
		}
		id, err := utils.GetUintUrlParam(r, "todo_id")
		if err != nil {
			ae := err.(*utils.ApiError)
			utils.RenderJson(r, w, ae.StatusCode, ae)
			return
		}

		var todo models.Todo
		err = di.DB.Where("id = ? and user_profile_id = ?", id, user_id).Take(&todo).Error
		if notFound, ae := utils.IsErrNotFound(err, "todo", id); notFound {
			utils.RenderJson(r, w, ae.StatusCode, ae)
			return
		}
		if err != nil {
			ae := utils.NewInternalServerError(err.Error())
			utils.RenderJson(r, w, ae.StatusCode, ae)
			return
		}
		utils.RenderJson(r, w, http.StatusOK, todo)
	}
}

func GetTodos(di *container.ServiceContainer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user_id, err := utils.GetUintUrlParam(r, "user_id")
		if err != nil {
			ae := err.(*utils.ApiError)
			utils.RenderJson(r, w, ae.StatusCode, ae)
			return
		}
		var items []models.Todo
		err = di.DB.Where("user_profile_id = ?", user_id).Find(&items).Error
		if err != nil {
			ae := utils.NewInternalServerError(err.Error())
			utils.RenderJson(r, w, ae.StatusCode, ae)
			return
		}
		utils.RenderJson(r, w, http.StatusOK, items)
	}
}
