package auth

import "context"

type IUser interface {
	GetID() int
	GetRoles() []string
}

type IUserSrv interface {
	GetUserFromApiKey(ctx context.Context, apiKey string) (IUser, error)
}

type AuthenticatedUser struct {
	ID    int
	Roles []string
}
