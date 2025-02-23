package models

import "github.com/google/uuid"

type UserTrophy struct {
	UserID   uuid.UUID `json:"user_id"`
	TrophyID uuid.UUID `json:"trophy_id"`
}

type VoxsphereMember struct {
	VoxsphereID uuid.UUID `json:"voxsphere_id"`
	UserID      uuid.UUID `json:"user_id"`
}

type VoxsphereModerator struct {
	VoxsphereID uuid.UUID `json:"voxsphere_id"`
	UserID      uuid.UUID `json:"user_id"`
}

type UserFlairEmoji struct {
	EmojiID     uuid.UUID `json:"emoji_id"`
	UserFlairID uuid.UUID `json:"user_flair_id"`
	OrderIndex  int32     `json:"order_index"`
}

type UserFlairCustomEmoji struct {
	CustomEmojiID uuid.UUID `json:"custom_emoji_id"`
	UserFlairID   uuid.UUID `json:"user_flair_id"`
	OrderIndex    int32     `json:"order_index"`
}

type UserFlairDescription struct {
	UserFlairID uuid.UUID `json:"user_flair_id"`
	OrderIndex  int32     `json:"order_index"`
	Description string    `json:"description"`
}

type PostFlairEmoji struct {
	EmojiID     uuid.UUID `json:"emoji_id"`
	PostFlairID uuid.UUID `json:"post_flair_id"`
	OrderIndex  int32     `json:"order_index"`
}

type PostFlairCustomEmoji struct {
	CustomEmojiID uuid.UUID `json:"custom_emoji_id"`
	PostFlairID   uuid.UUID `json:"post_flair_id"`
	OrderIndex    int32     `json:"order_index"`
}

type PostFlairDescription struct {
	PostFlairID uuid.UUID `json:"post_flair_id"`
	OrderIndex  int32     `json:"order_index"`
	Description string    `json:"description"`
}
