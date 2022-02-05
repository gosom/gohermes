package models

import (
	"github.com/gosom/gohermes/pkg/user"
	"gorm.io/gorm"
)

type CustomUser struct {
	ID uint `json:"id" gorm:"primarykey"`

	Email string `json:"email" gorm:"type:varchar(100);uniqueIndex" validate:"required,lte=100,email"`
}

func (o CustomUser) GetID() uint {
	return o.ID
}

func FetchDbUser(db *gorm.DB, identifier string) (user.IUser, error) {
	var u CustomUser
	if err := db.Where("id = ?", 2).Take(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}
