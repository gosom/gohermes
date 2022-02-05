package utils

import (
	"golang.org/x/crypto/bcrypt"
)

const hashCost = bcrypt.DefaultCost

func HashPassword(password string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), hashCost)
	return hash, err
}
