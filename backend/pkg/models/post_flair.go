package models

import "github.com/google/uuid"

type PostFlair struct {
	ID              uuid.UUID `json:"id"`
	PostID          uuid.UUID `json:"post_id"`
	VoxsphereID     uuid.UUID `json:"voxsphere_id"`
	FullText        string    `json:"full_text"`
	BackgroundColor string    `json:"background_color"`
}
