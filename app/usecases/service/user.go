package service

import (
	"context"
	"errors"
	"homework-dontpanicw/app/domain"
	"homework-dontpanicw/app/repository"
	"homework-dontpanicw/app/usecases"
	"homework-dontpanicw/app/usecases/auth"
)

type UserService struct {
	userRepo    repository.User
	sessionRepo repository.Session
}

func NewUser(userRepo repository.User, sessionRepo repository.Session) usecases.User {
	return &UserService{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
	}
}

func (us *UserService) RegisterNewUser(ctx context.Context, username string, password string) error {
	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		return errors.New("Error hashing password")
	}
	return us.userRepo.RegisterNewUser(ctx, username, hashedPassword)
}

func (us *UserService) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	return us.userRepo.GetUserByUsername(ctx, username)
}

func (us *UserService) CreateNewSession(id int64) (int64, error) {
	return us.sessionRepo.CreateNewSession(id)
}

func (us *UserService) GetAllSessions() map[string]string {
	return us.sessionRepo.GetAllSessions()
}
