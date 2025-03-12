package userflair

import (
	"context"
	"database/sql"
	"errors"
	"strings"

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

type UserFlairRepository interface {
	UserFlairs(context.Context) ([]models.UserFlair, error)
	UserFlairByID(context.Context, uuid.UUID) (models.UserFlair, error)
	AddUserFlairs(context.Context, ...models.UserFlair) ([]models.UserFlair, error)
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

func (r *Repo) AddUserFlairs(ctx context.Context, userFlairs ...models.UserFlair) ([]models.UserFlair, error) {
	query := `
        INSERT INTO
            user_flairs (
                id,
                voxsphere_id,
                full_text,
                background_color
            )
        VALUES 
    `
	args := make([]interface{}, 0)
	placeholders := make([]string, 0)

	for _, userFlair := range userFlairs {
		placeholders = append(placeholders, "(?, ?, ?, ?)")

		args = append(args,
			userFlair.ID,
			userFlair.VoxsphereID,
			userFlair.FullText,
			userFlair.BackgroundColor,
		)
	}

	query += strings.Join(placeholders, ", ")
	query += " RETURNING *"

	if _, err := r.db.NewRaw(
		query, args...).Exec(ctx, &userFlairs); err != nil {
		var pgdriverErr pgdriver.Error
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgUniqueViolation {
			return nil, ErrUserFlairDuplicateID
		}
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgConstraintViolation {
			return nil, ErrUserFlairParentTableRecordNotFound
		}
		return nil, err
	}

	return userFlairs, nil
}

func (r *Repo) UpdateUserFlair(ctx context.Context, userFlair models.UserFlair) (models.UserFlair, error) {
	query := `
                UPDATE
                    user_flairs
                SET
                    voxsphere_id = ?,
                    full_text = ?,
                    background_color = ?
                WHERE
                    id = ?
                RETURNING *
            `

	res, err := r.db.NewRaw(
		query,
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
