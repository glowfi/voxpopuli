package models

import "github.com/google/uuid"

type UserTrophy struct {
	UserID   uuid.UUID `json:"user_id"`
	TrophyID uuid.UUID `json:"trophy_id"`
}
