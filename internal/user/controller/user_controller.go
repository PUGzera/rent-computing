package user_controller

import (
	"context"
	"time"
)

type UserController interface {
	Login(ctx context.Context, username, password string) (*UserRepresentation, error)
	Register(ctx context.Context, username, email, password string) (*UserRepresentation, error)
	ProfileFromUsername(ctx context.Context, username string) (*UserRepresentation, error)
	ProfileFromEmail(ctx context.Context, email string) (*UserRepresentation, error)
}

type UserRepresentation struct {
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"createdAt"`
}
