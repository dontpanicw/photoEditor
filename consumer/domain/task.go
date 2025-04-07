package domain

import "github.com/google/uuid"

type Task struct {
	PhotoId   uuid.UUID `json:"photo_id"`
	Filter    string    `json:"filter"`
	Parameter float64   `json:"parameter"`
	Status    string    `json:"status"`
}
