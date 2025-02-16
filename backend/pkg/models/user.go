package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID                uuid.UUID `json:"id"`
	Name              string    `json:"name"`
	PublicDescription *string   `json:"public_description"`
	AvatarImg         *string   `json:"avatar_img"`
	BannerImg         *string   `json:"banner_img"`
	Iconcolor         *string   `json:"iconcolor"`
	Keycolor          *string   `json:"keycolor"`
	Primarycolor      *string   `json:"primarycolor"`
	Over18            bool      `json:"over18"`
	Suspended         bool      `json:"suspended"`
	CreatedAt         time.Time `json:"created_at"`
	CreatedAtUnix     int64     `json:"created_at_unix"`
	UpdatedAt         time.Time `json:"updated_at"`
}
