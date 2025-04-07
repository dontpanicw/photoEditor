package repository

import (
	"homework-dontpanicw/app/domain"
	"homework-dontpanicw/app/usecases"
)

type RabbitMQSender interface {
	SendTask(task domain.Task) (string, error)
	ListenForResponses(task usecases.Task)
}
