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
