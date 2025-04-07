package ram_storage

//
//import (
//	"errors"
//	"github.com/google/uuid"
//	"homework-dontpanicw/app/domain"
//	"homework-dontpanicw/app/repository"
//)
//
//type TaskRepository struct {
//	task map[uuid.UUID]*domain.Task
//}
//
//func NewTask() repository.Task {
//	return &TaskRepository{
//		task: make(map[uuid.UUID]*domain.Task),
//	}
//}
//
//func (rs *TaskRepository) PostTask(id uuid.UUID, task domain.Task) error {
//		if _, exists := rs.task[id]; exists {
//			return errors.New("task id already exists")
//		}
//		rs.task[id] = &domain.Task{
//			PhotoId:   id,
//			Parameter: task.Parameter,
//			Filter:    task.Filter,
//			Status:    task.Status,
//		}
//		return nil
//}
//
//func (rs *TaskRepository) GetTask(id uuid.UUID) (*domain.Task, error) {
//	task, exists := rs.task[id]
//	if !exists {
//		return nil, repository.TaskNotFound
//	}
//	return task, nil
//}
//
//func (rs *TaskRepository) GetAllTasks() map[uuid.UUID]*domain.Task {
//	return rs.task
//}
//
//func (rs *TaskRepository) UpdateTask(id uuid.UUID, task domain.Task) error {
//	if _, exists := rs.task[id]; !exists {
//		return errors.New("task not found")
//	}
//	rs.task[id] = &task
//	return nil
//}
