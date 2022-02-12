package utils

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
)

func IsErrNotFound(err error, resource string, id int) *ApiError {
	if err == nil {
		return nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		ae := NewResourceNotFoundError(resource, id)
		return &ae
	}
	return nil
}

type ApiError struct {
	StatusCode int    `json:"-"`
	Msg        string `json:"error,omitempty"`
}

func (o *ApiError) Error() string {
	return fmt.Sprintf("%d: %s", o.StatusCode, o.Msg)
}

func ApiErrorFromErr(err error) *ApiError {
	if err == nil {
		return nil
	}
	ae, ok := err.(*ApiError)
	if !ok {
		o := NewInternalServerError(err.Error())
		ae = &o
	}
	return ae
}

func NewResourceNotFoundError(resourceName string, id int) ApiError {
	ae := ApiError{
		StatusCode: http.StatusNotFound,
		Msg:        fmt.Sprintf("not found %s with id %d", resourceName, id),
	}
	return ae
}

func NewAuthenticationError(msg string) ApiError {
	if len(msg) == 0 {
		msg = http.StatusText(http.StatusUnauthorized)
	}
	return ApiError{StatusCode: http.StatusUnauthorized, Msg: msg}
}

func NewAuthorizationError(msg string) ApiError {
	if len(msg) == 0 {
		msg = http.StatusText(http.StatusForbidden)
	}
	return ApiError{StatusCode: http.StatusForbidden, Msg: msg}
}

func NewInternalServerError(msg string) ApiError {
	return ApiError{StatusCode: http.StatusInternalServerError, Msg: msg}
}

func NewBadRequestError(msg string) ApiError {
	return ApiError{StatusCode: http.StatusBadRequest, Msg: msg}
}
