package models

type Todo struct {
	ID      uint   `json:"id" gorm:"primarykey"`
	UserID  uint   `json:"user_id"`
	Title   string `json:"title"`
	Content string `json:"content"`

	User CustomUser `json:"-"`
}
