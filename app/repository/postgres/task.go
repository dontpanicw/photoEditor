package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"homework-dontpanicw/app/domain"
	"homework-dontpanicw/app/repository"
	"log"
)

type PostgresStorageTask struct {
	db *sql.DB
}

func (ps *PostgresStorageTask) GetDb() *sql.DB {
	return ps.db
}

func NewTaskPostgresStorage(connStr string) (*PostgresStorageTask, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			db.Close()
		}
	}()

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &PostgresStorageTask{db: db}, nil
}

func (ps *PostgresStorageTask) PostTask(ctx context.Context, id uuid.UUID, task domain.Task) error {
	query := `INSERT INTO tasks (photo_id, parameter, filter, status)VALUES ($1, $2, $3, $4) ON CONFLICT (photo_id) DO NOTHING`
	cmdTag, err := ps.db.ExecContext(ctx, query, id, task.Parameter, task.Filter, task.Status)
	if err != nil {
		return fmt.Errorf("failed to insert task: %w", err)
	}
	rowsAffected, err := cmdTag.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return errors.New("task id already exists")
	}
	return nil
}

func (ps *PostgresStorageTask) GetTask(ctx context.Context, id uuid.UUID) (*domain.Task, error) {
	query := `SELECT photo_id, parameter, filter, status FROM tasks WHERE photo_id = $1`
	var task domain.Task
	err := ps.db.QueryRowContext(ctx, query, id).Scan(&task.PhotoId, &task.Parameter, &task.Filter, &task.Status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.TaskNotFound
		}
		return nil, fmt.Errorf("failed to get task: %w", err)
	}
	return &task, nil
}

func (ps *PostgresStorageTask) GetAllTasks(ctx context.Context) ([]*domain.Task, error) {
	query := `SELECT photo_id, parameter, filter, status FROM tasks`
	rows, err := ps.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query all tasks: %w", err)
	}
	defer rows.Close()

	var tasks []*domain.Task

	for rows.Next() {
		var task domain.Task
		err = rows.Scan(&task.PhotoId, &task.Parameter, &task.Filter, &task.Status)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}
		tasks = append(tasks, &task)
	}
	log.Printf("Retrieved tasks: %+v", tasks)
	return tasks, nil
}

func (ps *PostgresStorageTask) UpdateTask(ctx context.Context, id uuid.UUID, task domain.Task) error {
	query := `
		UPDATE tasks
		SET parameter = $1, filter = $2, status = $3
		WHERE photo_id = $4
		RETURNING photo_id
	`
	var photoID uuid.UUID
	err := ps.db.QueryRowContext(ctx, query, task.Parameter, task.Filter, task.Status, id).Scan(&photoID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("task not found")
		}
		return fmt.Errorf("failed to update task: %w", err)
	}
	return nil
}
