package user_controller

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	user "rent-computing/internal/user/data"
	user_repo "rent-computing/internal/user/repo"

	"golang.org/x/crypto/bcrypt"
)

type UserControllerImpl struct {
	userRepo            user_repo.UserRepository
	hashPassword        func(password string) (string, error)
	compareHashPassword func(hash, password string) error
}

type Options struct {
	HashPassword        func(password string) (string, error)
	CompareHashPassword func(hash, password string) error
	UserRepo            user_repo.UserRepository
}

func New(options Options) (*UserControllerImpl, error) {
	userController := &UserControllerImpl{
		userRepo: options.UserRepo,
	}

	if options.HashPassword == nil {
		userController.hashPassword = func(password string) (string, error) {
			hash, err := bcrypt.GenerateFromPassword([]byte(password), 0)
			if err != nil {
				return "", err
			}
			return string(hash), nil
		}
	}

	if options.CompareHashPassword == nil {
		userController.compareHashPassword = func(hash, password string) error {
			return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
		}
	}

	return userController, nil
}

func (c *UserControllerImpl) Login(ctx context.Context, username, password string) (*UserRepresentation, error) {
	user, err := c.userRepo.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	err = c.compareHashPassword(user.Password, password)
	if err != nil {
		return nil, err
	}

	return &UserRepresentation{
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}, nil
}

func (c *UserControllerImpl) Register(ctx context.Context, username, email, password string) (*UserRepresentation, error) {
	address, err := mail.ParseAddress(email)
	if err != nil {
		return nil, err
	}

	if len(username) < 3 || len(username) > 10 {
		return nil, fmt.Errorf("%s not valid, username needs to be between length 3-10", username)
	}

	if len(password) < 8 {
		return nil, errors.New("password not valid, password needs to be longer than eight characters")
	}

	hashedPassword, err := c.hashPassword(password)
	if err != nil {
		return nil, err
	}

	userOpts := user.Options{
		Username: username,
		Email:    address.Address,
		Password: hashedPassword,
	}

	user, err := user.New(userOpts)
	if err != nil {
		return nil, err
	}

	err = c.userRepo.CreateUser(ctx, *user)
	if err != nil {
		return nil, err
	}

	return &UserRepresentation{
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}, nil

}

func (c *UserControllerImpl) ProfileFromUsername(ctx context.Context, username string) (*UserRepresentation, error) {
	user, err := c.userRepo.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	return &UserRepresentation{
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}, nil
}

func (c *UserControllerImpl) ProfileFromEmail(ctx context.Context, email string) (*UserRepresentation, error) {
	user, err := c.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return &UserRepresentation{
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}, nil
}
