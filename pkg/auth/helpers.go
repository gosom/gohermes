package auth

import (
	"context"

	"github.com/gosom/gohermes/pkg/utils"
)

func GetAuthenticatedUser(ctx context.Context) AuthenticatedUser {
	return ctx.Value(utils.Authenticated).(AuthenticatedUser)
}

func IsAuthenticatedUser(ctx context.Context) bool {
	v := ctx.Value(utils.Authenticated)
	return v != nil
}
