package utils

import (
	"github.com/go-playground/validator/v10"
	passwordvalidator "github.com/wagslane/go-password-validator"
)

const minEntropyBits = 60

var Validate *validator.Validate

func init() {
	Validate = validator.New()
}

func ValidatePassword(password string) error {
	return passwordvalidator.Validate(password, minEntropyBits)
}
