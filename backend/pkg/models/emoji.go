package models

import "github.com/google/uuid"

type Emoji struct {
	ID    uuid.UUID `json:"id"`
	Title string    `json:"title"`
}
