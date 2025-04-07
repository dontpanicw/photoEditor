package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"homework-dontpanicw/app/domain"
	"homework-dontpanicw/app/repository"
	"homework-dontpanicw/app/usecases"
)

type TaskService struct {
	repo   repository.Task
	sender repository.RabbitMQSender
}

func NewTask(repo repository.Task, sender repository.RabbitMQSender) usecases.Task {
	return &TaskService{
		repo:   repo,
		sender: sender,
	}
}

func (rs *TaskService) PostTask(ctx context.Context, id uuid.UUID, object domain.Task) error {
	_, err := rs.sender.SendTask(object)
	if err != nil {
		return err
	}
	return rs.repo.PostTask(ctx, id, object)
}

func (rs *TaskService) GetTask(ctx context.Context, id uuid.UUID) (*domain.Task, error) {
	return rs.repo.GetTask(ctx, id)
}

func (rs *TaskService) GetAllTasks(ctx context.Context) ([]*domain.Task, error) {
	return rs.repo.GetAllTasks(ctx)
}

func (rs *TaskService) DoingTask(ctx context.Context, id uuid.UUID) error {
	task, exists := rs.repo.GetTask(ctx, id)
	if exists != nil {
		return errors.New("task not found")
	}
	task.Status = "ready"
	err := rs.repo.UpdateTask(ctx, id, *task)
	if err != nil {
		return err
	}
	return nil
}
