package main

import (
	"fmt"
	"os"

	"github.com/casbin/casbin/v2"
	"github.com/go-chi/chi"

	"github.com/gosom/gohermes/examples/todo/controllers"
	"github.com/gosom/gohermes/examples/todo/models"
	"github.com/gosom/gohermes/examples/todo/services"
	"github.com/gosom/gohermes/pkg/container"
	"github.com/gosom/gohermes/pkg/server"
	"github.com/gosom/gohermes/pkg/user"
)

func main() {
	di, err := container.NewDefault()
	if err != nil {
		panic(err)
	}

	if err := migrations(di); err != nil {
		panic(err)
	}

	if err := registerServices(di); err != nil {
		panic(err)
	}

	srv, err := server.New(di)
	if err != nil {
		panic(err)
	}

	dbHooks(di)

	if err := setUpRoutes(di, srv); err != nil {
		panic(err)
	}

	if err := srv.Run(); err != nil {
		fmt.Println(err.Error())
		os.Exit(2)
	} else {
		os.Exit(0)
	}
}

func migrations(di *container.ServiceContainer) error {
	var migrate []interface{}
	migrate = append(migrate, models.AutoMigrateModels()...)
	if err := di.AutoMigrate(migrate...); err != nil {
		return err
	}
	return nil
}

func dbHooks(di *container.ServiceContainer) {
	/*
		user.RegisterUserAfterCreateHook(func(tx *gorm.DB, u *user.User) error {
			p := models.UserProfile{
				UserID: u.ID,
			}
			return tx.Create(&p).Error
		})
	*/
}

func registerServices(di *container.ServiceContainer) error {
	enforcer, err := user.NewEnforcer(di.DB)
	if err != nil {
		return err
	}
	di.RegisterService("enforcer", enforcer)

	users := user.NewUserService(di)
	di.RegisterService("users", users)

	todos := services.NewTodoService(di)
	di.RegisterService("todo", todos)
	return nil
}

func authPolicy(e *casbin.Enforcer) (err error) {
	rules := [][]string{}
	if len(rules) > 0 {
		_, err = e.AddPolicies(rules)
	}
	return err
}

func setUpRoutes(di *container.ServiceContainer, srv *server.Server) error {
	router := srv.GetRouter()

	router.Get("/health", server.HealthHandler(di))

	// public endpoints
	router.Post("/register", user.RegisterUserHandler(di))

	router.Route("/users", func(r chi.Router) {
		r.Use(user.Authentication(di, models.FetchDbUser))
		r.Use(user.Authorizer(di))

		r.Get(`/{id:\d+}`, user.GetHandler(di))

		r.Post(`/{user_id:\d+}/todo`, controllers.CreateTodo(di))
		r.Get(`/{user_id:\d+}/todo/{todo_id:\d+}`, controllers.GetTodoByID(di))
		r.Get(`/{user_id:\d+}/todo`, controllers.GetTodos(di))
	})
	return nil
}
