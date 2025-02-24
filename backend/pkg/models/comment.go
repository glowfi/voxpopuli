package models

import (
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	ID              uuid.UUID `json:"id"`
	AuthorID        uuid.UUID `json:"author_id"`
	ParentCommentID uuid.UUID `json:"parent_comment_id"`
	PostID          uuid.UUID `json:"post_id"`
	Body            string    `json:"body"`
	BodyHtml        string    `json:"body_html"`
	Ups             int32     `json:"ups"`
	Score           int32     `json:"score"`
	CreatedAt       time.Time `json:"created_at"`
	CreatedAtUnix   int64     `json:"created_at_unix"`
	UpdatedAt       time.Time `json:"updated_at"`
}
