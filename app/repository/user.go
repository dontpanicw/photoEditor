package repository

import (
	"context"
	"homework-dontpanicw/app/domain"
)

type User interface {
	RegisterNewUser(ctx context.Context, username string, password string) error
	GetUserByUsername(ctx context.Context, username string) (*domain.User, error)
}
