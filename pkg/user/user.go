package user

import (
	"time"

	"gorm.io/gorm"
)

type IUser interface {
	GetID() uint
}

const (
	RoleNameAnonymous = "anonymous"
	RoleNameAdmin     = "admin"
	RoleNameMember    = "member"
)

var afterUserCreateHooks []func(tx *gorm.DB, u *User) error

func RegisterUserAfterCreateHook(f func(tx *gorm.DB, u *User) error) {
	afterUserCreateHooks = append(afterUserCreateHooks, f)
}

type User struct {
	ID uint `json:"id" gorm:"primarykey"`

	Email           string `json:"email" gorm:"type:varchar(100);uniqueIndex" validate:"required,lte=100,email"`
	Password        []byte `json:"-"`
	Active          bool   `json:"active"`
	ActivationToken string `json:"-"`

	ApiKeys []ApiKey `json:"-"`

	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (u *User) AfterCreate(tx *gorm.DB) (err error) {
	for i := range afterUserCreateHooks {
		if err := afterUserCreateHooks[i](tx, u); err != nil {
			return err
		}
	}
	return nil
}

type ApiKey struct {
	ID        uint           `gorm:"primarykey"`
	UserID    uint           `gorm:"uniqueIndex:unq_idx_user_id_api_key"`
	ApiKey    string         `gorm:"varchar(64);uniqueIndex:unq_idx_user_id_api_key"`
	CreatedAt time.Time      `json:"createdAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

/*
func AutoMigrateModels() []interface{} {
	ans := []interface{}{
		User{},
		ApiKey{},
	}
	return ans
}
*/
