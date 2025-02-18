package models

import (
	"time"

	"github.com/google/uuid"
)

type Voxsphere struct {
	ID                    uuid.UUID `json:"id"`
	TopicID               uuid.UUID `json:"topic_id"`
	Topic                 Topic     `json:"topic" bun:",scanonly"`
	Title                 string    `json:"title"`
	PublicDescription     *string   `json:"public_description"`
	CommunityIcon         *string   `json:"community_icon"`
	BannerBackgroundImage *string   `json:"banner_background_image"`
	BannerBackgroundColor *string   `json:"banner_background_color"`
	KeyColor              *string   `json:"key_color"`
	PrimaryColor          *string   `json:"primary_color"`
	Over18                bool      `json:"over18"`
	SpoilersEnabled       bool      `json:"spoilers_enabled"`
	CreatedAt             time.Time `json:"created_at"`
	CreatedAtUnix         int64     `json:"created_at_unix"`
	UpdatedAt             time.Time `json:"updated_at"`
}
