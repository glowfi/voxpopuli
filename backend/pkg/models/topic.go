package models

import "github.com/google/uuid"

type Topic struct {
	ID   uuid.UUID `json:"id" bun:",pk"`
	Name string    `json:"name"`
}
