package repository

import "errors"

var (
	TaskNotFound = errors.New("Task not found")
	BadRequest   = errors.New("Bad request")
	UserNotFound = errors.New("User not found")
	Unauthorized = errors.New("Incorrect username/password")
)
