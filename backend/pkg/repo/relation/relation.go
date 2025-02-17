package relation

import (
	"context"
	"errors"

	"github.com/glowfi/voxpopuli/backend/pkg/models"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/driver/pgdriver"
)

const (
	pgUniqueViolation     = "23505"
	pgConstraintViolation = "23503"
)

var (
	ErrUserTrophyParentTableRecordNotFound = errors.New("user or trophy does not exist in the parent table")
	ErrDuplicateTrophyIDorDuplicateUserID  = errors.New("duplicate user id or trophy id")
)

type Repository interface {
	LinkUserTrophy(ctx context.Context, userID uuid.UUID, trophyID uuid.UUID) (models.UserTrophy, error)
	UserTrophies(ctx context.Context) ([]models.UserTrophy, error)
}

type Repo struct {
	db *bun.DB
}

func NewRepo(db *bun.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) LinkUserTrophy(ctx context.Context, userID uuid.UUID, trophyID uuid.UUID) (models.UserTrophy, error) {
	query := `
                INSERT INTO user_trophies
                     (
                        user_id,
                        trophy_id
                    )
                VALUES (
                    ?,
                    ?
                )
            `

	if _, err := r.db.NewRaw(query, userID, trophyID).Exec(ctx); err != nil {
		var pgdriverErr pgdriver.Error
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgUniqueViolation {
			return models.UserTrophy{}, ErrDuplicateTrophyIDorDuplicateUserID
		}
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgConstraintViolation {
			return models.UserTrophy{}, ErrUserTrophyParentTableRecordNotFound
		}
		return models.UserTrophy{}, err
	}

	return models.UserTrophy{
		UserID:   userID,
		TrophyID: trophyID,
	}, nil
}

func (r *Repo) UserTrophies(ctx context.Context) ([]models.UserTrophy, error) {
	var user_trophies []models.UserTrophy

	query := `
                SELECT 
                    user_id,
                    trophy_id 
                FROM
                    user_trophies
            `

	_, err := r.db.NewRaw(query).Exec(ctx, &user_trophies)
	if err != nil {
		return []models.UserTrophy{}, err
	}
	return user_trophies, nil
}
