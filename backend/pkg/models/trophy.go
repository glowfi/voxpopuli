package models

import "github.com/google/uuid"

type Trophy struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	ImageLink   string    `json:"image_link"`
}
