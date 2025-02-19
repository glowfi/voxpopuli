package models

import "github.com/google/uuid"

type UserFlair struct {
	ID              uuid.UUID `json:"id"`
	UserID          uuid.UUID `json:"user_id"`
	VoxsphereID     uuid.UUID `json:"voxsphere_id"`
	FullText        string    `json:"full_text"`
	BackgroundColor string    `json:"background_color"`
}
