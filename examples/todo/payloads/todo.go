package payloads

type CreateTodoPayload struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}
