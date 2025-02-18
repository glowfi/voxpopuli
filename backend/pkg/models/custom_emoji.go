package models

import "github.com/google/uuid"

type CustomEmoji struct {
	ID          uuid.UUID `json:"id"`
	VoxsphereID uuid.UUID `json:"voxsphere_id"`
	Url         string    `json:"url"`
	Title       string    `json:"title"`
}
