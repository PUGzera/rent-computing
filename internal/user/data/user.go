package user

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id        string    `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"createdAt"`
}

type Options struct {
	Email    string
	Username string
	Password string
}

func New(options Options) (*User, error) {
	return &User{
		Id:        uuid.NewString(),
		Email:     options.Email,
		Username:  options.Username,
		Password:  options.Password,
		CreatedAt: time.Now(),
	}, nil
}
