package payloads

import "github.com/gosom/gohermes/examples/todo/models"

type CreateTodoPayload struct {
	Title   string `json:"title"`
	Content string `json:"content"`

	User models.CustomUser `json:"-" validate:"-"`
}
