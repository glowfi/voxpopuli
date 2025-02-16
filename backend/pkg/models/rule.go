package models

import "github.com/google/uuid"

type Rule struct {
	ID          uuid.UUID `json:"id"`
	VoxsphereID uuid.UUID `json:"voxsphere_id"`
	ShortName   string    `json:"short_name"`
	Description string    `json:"description"`
}
