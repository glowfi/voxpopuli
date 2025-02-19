package models

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	ID            uuid.UUID `json:"id"`
	AuthorID      uuid.UUID `json:"author_id"`
	VoxsphereID   uuid.UUID `json:"voxsphere_id"`
	Title         string    `json:"title"`
	Text          string    `json:"text"`
	TextHtml      string    `json:"text_html"`
	Ups           int32     `json:"ups"`
	Over18        bool      `json:"over18"`
	Spoiler       bool      `json:"spoiler"`
	CreatedAt     time.Time `json:"created_at"`
	CreatedAtUnix int64     `json:"created_at_unix"`
	UpdatedAt     time.Time `json:"updated_at"`
}
