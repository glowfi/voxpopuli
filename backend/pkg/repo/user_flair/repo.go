package userflair

import (
	"context"
	"database/sql"
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
	ErrUserFlairNotFound                  = errors.New("user flair not found")
	ErrUserFlairDuplicateID               = errors.New("user flair duplicate id")
	ErrUserFlairParentTableRecordNotFound = errors.New("record does not exist in the parent table")
)

type Repository interface {
	UserFlairs() ([]models.UserFlair, error)
	UserFlairByID(context.Context, uuid.UUID) (models.UserFlair, error)
	AddUserFlair(context.Context, models.UserFlair) (models.UserFlair, error)
	UpdateUserFlair(context.Context, models.UserFlair) (models.UserFlair, error)
	DeleteUserFlair(context.Context, uuid.UUID) error
}

type Repo struct {
	db *bun.DB
}

func NewRepo(db *bun.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) UserFlairs(ctx context.Context) ([]models.UserFlair, error) {
	var userFlairs []models.UserFlair

	query := `
                SELECT
                    id,
                    user_id,
                    voxsphere_id,
                    full_text,
                    background_color
                FROM
                    user_flairs;
            `

	_, err := r.db.NewRaw(query).Exec(ctx, &userFlairs)
	if err != nil {
		return []models.UserFlair{}, err
	}
	return userFlairs, nil
}

func (r *Repo) UserFlairByID(ctx context.Context, ID uuid.UUID) (models.UserFlair, error) {
	var userFlair models.UserFlair

	query := `
                SELECT
                    id,
                    user_id,
                    voxsphere_id,
                    full_text,
                    background_color
                FROM
                    user_flairs
                WHERE
                    id = ?;
            `
	_, err := r.db.NewRaw(query, ID).Exec(ctx, &userFlair)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.UserFlair{}, ErrUserFlairNotFound
		}
		return models.UserFlair{}, err
	}
	return userFlair, nil
}

func (r *Repo) AddUserFlair(ctx context.Context, userFlair models.UserFlair) (models.UserFlair, error) {
	query := `
                INSERT INTO
                    user_flairs (
                        id,
                        user_id,
                        voxsphere_id,
                        full_text,
                        background_color
                    )
                VALUES (
                    ?,
                    ?,
                    ?,
                    ?,
                    ?
                )
                RETURNING *
            `

	if _, err := r.db.NewRaw(
		query,
		userFlair.ID,
		userFlair.UserID,
		userFlair.VoxsphereID,
		userFlair.FullText,
		userFlair.BackgroundColor,
	).Exec(ctx, &userFlair); err != nil {
		var pgdriverErr pgdriver.Error
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgUniqueViolation {
			return models.UserFlair{}, ErrUserFlairDuplicateID
		}
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgConstraintViolation {
			return models.UserFlair{}, ErrUserFlairParentTableRecordNotFound
		}
		return models.UserFlair{}, err
	}

	return userFlair, nil
}

func (r *Repo) UpdateUserFlair(ctx context.Context, userFlair models.UserFlair) (models.UserFlair, error) {
	query := `
                UPDATE
                    user_flairs
                SET
                    user_id = ?,
                    voxsphere_id = ?,
                    full_text = ?,
                    background_color = ?
                WHERE
                    id = ?
                RETURNING *
            `

	res, err := r.db.NewRaw(
		query,
		userFlair.UserID,
		userFlair.VoxsphereID,
		userFlair.FullText,
		userFlair.BackgroundColor,
		userFlair.ID,
	).Exec(ctx, &userFlair)
	if errors.Is(err, sql.ErrNoRows) {
		return models.UserFlair{}, ErrUserFlairNotFound
	}
	if err != nil {
		var pgdriverErr pgdriver.Error
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgConstraintViolation {
			return models.UserFlair{}, ErrUserFlairParentTableRecordNotFound
		}
		return models.UserFlair{}, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return models.UserFlair{}, err
	}
	if rowsAffected == 0 {
		return models.UserFlair{}, ErrUserFlairNotFound
	}
	return userFlair, nil
}

func (r *Repo) DeleteUserFlair(ctx context.Context, ID uuid.UUID) error {
	query := `
                DELETE FROM 
                    user_flairs
                WHERE 
                    id = ?
            `
	res, err := r.db.NewRaw(query, ID).Exec(ctx)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrUserFlairNotFound
	}
	return nil
}
