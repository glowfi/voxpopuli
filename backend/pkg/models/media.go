package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type MediaType string

const (
	MediaTypeImage   MediaType = "image"
	MediaTypeGif     MediaType = "gif"
	MediaTypeVideo   MediaType = "video"
	MediaTypeGallery MediaType = "gallery"
	MediaTypeLink    MediaType = "link"
	MediaTypeMulti   MediaType = "multi"
)

type PostMedia struct {
	bun.BaseModel `bun:"table:post_medias"`
	ID            uuid.UUID `json:"id"`
	PostID        uuid.UUID `json:"post_id"`
	MediaType     MediaType `json:"media_type"`
}

type Image struct {
	ID            uuid.UUID       `json:"id"`
	MediaID       uuid.UUID       `json:"media_id"`
	ImageMetadata []ImageMetadata `json:"image_metadata" bun:",scanonly"`
}

type ImageMetadata struct {
	bun.BaseModel `bun:"table:image_metadatas"`
	ID            uuid.UUID `json:"id"`
	ImageID       uuid.UUID `json:"image_id"`
	Height        int32     `json:"height"`
	Width         int32     `json:"width"`
	Url           string    `json:"url"`
	CreatedAt     time.Time `json:"created_at"`
	CreatedAtUnix int64     `json:"created_at_unix"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Gif struct {
	ID          uuid.UUID     `json:"id"`
	MediaID     uuid.UUID     `json:"media_id"`
	GifMetadata []GifMetadata `json:"image_metadata" bun:",scanonly"`
}

type GifMetadata struct {
	bun.BaseModel `bun:"table:gif_metadatas"`
	ID            uuid.UUID `json:"id"`
	GifID         uuid.UUID `json:"gif_id"`
	Height        int32     `json:"height"`
	Width         int32     `json:"width"`
	Url           string    `json:"url"`
	CreatedAt     time.Time `json:"created_at"`
	CreatedAtUnix int64     `json:"created_at_unix"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Gallery struct {
	ID              uuid.UUID         `json:"id"`
	MediaID         uuid.UUID         `json:"media_id"`
	GalleryMetadata []GalleryMetadata `json:"image_metadata" bun:",scanonly"`
}

type GalleryMetadata struct {
	bun.BaseModel `bun:"table:gallery_metadatas"`
	ID            uuid.UUID `json:"id"`
	GalleryID     uuid.UUID `json:"gallery_id"`
	OrderIndex    int32     `json:"order_index"`
	Height        int32     `json:"height"`
	Width         int32     `json:"width"`
	Url           string    `json:"url"`
	CreatedAt     time.Time `json:"created_at"`
	CreatedAtUnix int64     `json:"created_at_unix"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Video struct {
	ID            uuid.UUID `json:"id"`
	MediaID       uuid.UUID `json:"media_id"`
	Url           string    `json:"url"`
	Height        int32     `json:"height"`
	Width         int32     `json:"width"`
	CreatedAt     time.Time `json:"created_at"`
	CreatedAtUnix int64     `json:"created_at_unix"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Link struct {
	ID            uuid.UUID `json:"id"`
	MediaID       uuid.UUID `json:"media_id"`
	Link          string    `json:"link"`
	CreatedAt     time.Time `json:"created_at"`
	CreatedAtUnix int64     `json:"created_at_unix"`
	UpdatedAt     time.Time `json:"updated_at"`
}
