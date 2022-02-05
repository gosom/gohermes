package user

import (
	"errors"
	"fmt"

	"github.com/gosom/gohermes/pkg/container"
	"github.com/gosom/gohermes/pkg/utils"
	"gorm.io/gorm"
)

var ErrCannotCastToIUserSrv error = errors.New("cannot cast to IUserSrv")

type IUserSrv interface {
	Register(db *gorm.DB, email, password string) (User, error)
}

type Enforcer interface {
	AddRolesForUser(user string, roles []string, domain ...string) (bool, error)
	AddPolicies(rules [][]string) (bool, error)
}

func CastIUserSrv(iface interface{}) (IUserSrv, error) {
	ans, ok := iface.(IUserSrv)
	if !ok {
		return nil, ErrCannotCastToIUserSrv
	}
	return ans, nil
}

func CastEnforcer(i interface{}) (Enforcer, error) {
	ans, ok := i.(Enforcer)
	if !ok {
		return nil, fmt.Errorf("type is not Enforcer")
	}
	return ans, nil
}

type UserService struct {
	di *container.ServiceContainer
}

func NewUserService(di *container.ServiceContainer) *UserService {
	ans := UserService{
		di: di,
	}
	return &ans
}

// TODO proper errors
func (o *UserService) Register(db *gorm.DB, email, password string) (User, error) {
	u := User{}
	err := utils.ValidatePassword(password)
	if err != nil {
		return u, err
	}
	hash, err := utils.HashPassword(password)
	if err != nil {
		return u, err
	}

	u.Email = email
	u.Password = hash
	u.Active = true // TODO remove me

	err = utils.Validate.Struct(u)
	if err != nil {
		return u, err
	}
	exists, err := o.Exists(db, email)
	if err != nil {
		return u, err
	}
	if exists {
		return u, fmt.Errorf("user with email: %s exists", email)
	}

	if err = db.Create(&u).Error; err != nil {
		return u, err
	}

	ei, err := o.di.GetService("enforcer")
	if err != nil {
		return u, err
	}
	e, err := CastEnforcer(ei)
	if err != nil {
		return u, err
	}
	obj := fmt.Sprint(u.ID)
	rules := [][]string{
		{obj, fmt.Sprintf("/users/%d", u.ID), "*"},
		{obj, fmt.Sprintf("/users/%d/*", u.ID), "*"},
	}
	if _, err := e.AddPolicies(rules); err != nil {
		return u, err
	}
	_, err = e.AddRolesForUser(fmt.Sprint(u.ID), []string{RoleNameMember})
	if err != nil {
		return u, err
	}
	return u, err
}

func (o *UserService) Exists(db *gorm.DB, email string) (bool, error) {
	var exists bool
	err := db.Raw("select exists(select 1 from users where email = ?)", email).
		Scan(&exists).Error
	return exists, err
}
