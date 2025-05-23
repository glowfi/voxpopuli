package models

import "github.com/google/uuid"

type Topic struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Category string    `json:"category"`
}
