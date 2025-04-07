package usecases

import (
	"context"
	"github.com/google/uuid"
	"homework-dontpanicw/app/domain"
)

type Task interface {
	GetTask(ctx context.Context, id uuid.UUID) (*domain.Task, error)
	PostTask(ctx context.Context, id uuid.UUID, object domain.Task) error
	GetAllTasks(ctx context.Context) ([]*domain.Task, error)
	DoingTask(ctx context.Context, id uuid.UUID) error
}
