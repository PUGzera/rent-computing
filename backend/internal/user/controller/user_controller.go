package user_controller

import (
	"context"
	"time"
)

type UserController interface {
	Login(ctx context.Context, username, password string) (*UserRepresentation, error)
	Register(ctx context.Context, username, email, password string) (*UserRepresentation, error)
	Profile(ctx context.Context, id string) (*UserRepresentation, error)
	ProfileFromUsername(ctx context.Context, username, callerId string) (*UserRepresentation, error)
	ProfileFromEmail(ctx context.Context, email, callerId string) (*UserRepresentation, error)
}

type UserRepresentation struct {
	Id        string    `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"createdAt"`
}
