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
	ErrParentTableRecordNotFound = errors.New("record does not exist in the parent table")
	ErrDuplicateID               = errors.New("duplicate id")
)

type Repository interface {
	UserTrophies(ctx context.Context) ([]models.UserTrophy, error)
	LinkUserTrophy(ctx context.Context, userID uuid.UUID, trophyID uuid.UUID) (models.UserTrophy, error)

	VoxsphereMemebers(ctx context.Context) ([]models.VoxsphereMember, error)
	LinkVoxsphereMember(ctx context.Context, voxsphereID uuid.UUID, userID uuid.UUID) (models.VoxsphereMember, error)

	VoxsphereModerators(ctx context.Context) ([]models.VoxsphereModerator, error)
	LinkVoxsphereModerator(ctx context.Context, voxsphereID uuid.UUID, userID uuid.UUID) (models.VoxsphereModerator, error)
}

type Repo struct {
	db *bun.DB
}

func NewRepo(db *bun.DB) *Repo {
	return &Repo{db: db}
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
			return models.UserTrophy{}, ErrDuplicateID
		}
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgConstraintViolation {
			return models.UserTrophy{}, ErrParentTableRecordNotFound
		}
		return models.UserTrophy{}, err
	}

	return models.UserTrophy{
		UserID:   userID,
		TrophyID: trophyID,
	}, nil
}

func (r *Repo) VoxsphereMemebers(ctx context.Context) ([]models.VoxsphereMember, error) {
	var voxsphere_member []models.VoxsphereMember

	query := `
                SELECT 
                    voxsphere_id,
                    user_id 
                FROM
                    voxsphere_members
            `

	_, err := r.db.NewRaw(query).Exec(ctx, &voxsphere_member)
	if err != nil {
		return []models.VoxsphereMember{}, err
	}
	return voxsphere_member, nil
}

func (r *Repo) LinkVoxsphereMember(ctx context.Context, voxsphereID uuid.UUID, userID uuid.UUID) (models.VoxsphereMember, error) {
	query := `
                INSERT INTO voxsphere_members
                     (
                        voxsphere_id,
                        user_id
                    )
                VALUES (
                    ?,
                    ?
                )
            `

	if _, err := r.db.NewRaw(query, voxsphereID, userID).Exec(ctx); err != nil {
		var pgdriverErr pgdriver.Error
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgUniqueViolation {
			return models.VoxsphereMember{}, ErrDuplicateID
		}
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgConstraintViolation {
			return models.VoxsphereMember{}, ErrParentTableRecordNotFound
		}
		return models.VoxsphereMember{}, err
	}

	return models.VoxsphereMember{
		VoxsphereID: voxsphereID,
		UserID:      userID,
	}, nil
}
