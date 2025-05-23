package models

import "github.com/google/uuid"

type Award struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	ImageLink string    `json:"image_link"`
}
