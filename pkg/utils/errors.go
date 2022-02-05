package utils

import (
	"errors"
	"fmt"
	"net/http"

	"gorm.io/gorm"
)

func IsErrNotFound(err error, resource string, id uint) (bool, *ApiError) {
	if err == nil {
		return false, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		ae := NewResourceNotFoundError(resource, id)
		return true, &ae
	}
	return false, nil
}

type ApiError struct {
	StatusCode int    `json:"-"`
	Msg        string `json:"error,omitempty"`
}

func (o *ApiError) Error() string {
	return fmt.Sprintf("%d: %s", o.StatusCode, o.Msg)
}

func NewResourceNotFoundError(resourceName string, id uint) ApiError {
	ae := ApiError{
		StatusCode: http.StatusNotFound,
		Msg:        fmt.Sprintf("not found %s with id %d", resourceName, id),
	}
	return ae
}

func NewAuthenticationError(msg string) ApiError {
	return ApiError{StatusCode: http.StatusForbidden, Msg: msg}
}

func NewAuthorizationError(msg string) ApiError {
	return ApiError{StatusCode: http.StatusUnauthorized, Msg: msg}
}

func NewInternalServerError(msg string) ApiError {
	return ApiError{StatusCode: http.StatusInternalServerError, Msg: msg}
}

func NewBadRequestError(msg string) ApiError {
	return ApiError{StatusCode: http.StatusBadRequest, Msg: msg}
}
