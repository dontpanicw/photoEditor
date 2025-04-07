package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"homework-dontpanicw/app/domain"
	"homework-dontpanicw/app/repository"
)

type PostgresStorageUser struct {
	db *sql.DB
}

func (psu *PostgresStorageUser) GetDb() *sql.DB {
	return psu.db
}

func NewUserPostgresStorage(connStr string) (*PostgresStorageUser, error) {
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

	return &PostgresStorageUser{db: db}, nil
}

func (psu *PostgresStorageUser) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	query := `SELECT * FROM users WHERE username = $1;`
	var user domain.User
	err := psu.db.QueryRowContext(ctx, query, username).Scan(&user.Id, &user.Username, &user.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.TaskNotFound
		}
	}
	return &user, nil
}

func (psu *PostgresStorageUser) RegisterNewUser(ctx context.Context, username string, password string) error {
	user, _ := psu.GetUserByUsername(ctx, username)
	if user != nil {
		return errors.New("user already exists")
	}
	id := uuid.New()
	query := `INSERT INTO users (user_id, username, password)VALUES ($1, $2, $3)`
	_, err := psu.db.ExecContext(ctx, query, id, username, password)
	if err != nil {
		return fmt.Errorf("failed to insert task: %w", err)
	}
	return nil
}
