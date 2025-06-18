package user_repo

import (
	"context"
	user "rent-computing/internal/user/data"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user user.User) error
	ListUsers(ctx context.Context) ([]user.User, error)
	GetUser(ctx context.Context, id string) (*user.User, error)
	GetUserByUsername(ctx context.Context, username string) (*user.User, error)
	GetUserByEmail(ctx context.Context, email string) (*user.User, error)
	UpdateUser(ctx context.Context, user user.User) error
	DeleteUser(ctx context.Context, id string) error
	DeleteUserByUsername(ctx context.Context, username string) error
	DeleteUserByEmail(ctx context.Context, email string) error
}
