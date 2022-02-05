package services

import (
	"github.com/gosom/gohermes/examples/todo/models"
	"github.com/gosom/gohermes/examples/todo/payloads"
	"github.com/gosom/gohermes/pkg/container"
)

type TodoService struct {
	di *container.ServiceContainer
}

func (o *TodoService) Create(p payloads.CreateTodoPayload) (models.Todo, error) {
	todo := models.Todo{
		UserID:  p.User.ID,
		Title:   p.Title,
		Content: p.Content,
	}
	return todo, o.di.DB.Create(&todo).Error
}

func NewTodoService(di *container.ServiceContainer) *TodoService {
	ans := TodoService{
		di: di,
	}
	return &ans
}

func GetFromDi(di *container.ServiceContainer) *TodoService {
	i, err := di.GetService("todo")
	if err != nil {
		panic(err)
	}
	srv := i.(*TodoService)
	return srv
}
