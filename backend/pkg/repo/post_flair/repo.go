package postflair

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
	ErrPostFlairNotFound                  = errors.New("post flair not found")
	ErrPostFlairDuplicateID               = errors.New("post flair duplicate id")
	ErrPostFlairParentTableRecordNotFound = errors.New("record does not exist in the parent table")
)

type Repository interface {
	PostFlairs(ctx context.Context) ([]models.PostFlair, error)
	PostFlairByID(ctx context.Context, ID uuid.UUID) (models.PostFlair, error)
	AddPostFlairs(ctx context.Context, postFlair ...models.PostFlair) (models.PostFlair, error)
	UpdatePostFlair(ctx context.Context, postFlair models.PostFlair) (models.PostFlair, error)
	DeletePostFlair(ctx context.Context, ID uuid.UUID) error
}

type Repo struct {
	db *bun.DB
}

func NewRepo(db *bun.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) PostFlairs(ctx context.Context) ([]models.PostFlair, error) {
	var postFlairs []models.PostFlair

	query := `
                SELECT
                    id,
                    voxsphere_id,
                    full_text,
                    background_color
                FROM
                    post_flairs;
            `

	_, err := r.db.NewRaw(query).Exec(ctx, &postFlairs)
	if err != nil {
		return []models.PostFlair{}, err
	}
	return postFlairs, nil
}

func (r *Repo) PostFlairByID(ctx context.Context, ID uuid.UUID) (models.PostFlair, error) {
	var postFlair models.PostFlair

	query := `
                SELECT
                    id,
                    voxsphere_id,
                    full_text,
                    background_color
                FROM
                    post_flairs
                WHERE
                    id = ?;
            `
	_, err := r.db.NewRaw(query, ID).Exec(ctx, &postFlair)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.PostFlair{}, ErrPostFlairNotFound
		}
		return models.PostFlair{}, err
	}
	return postFlair, nil
}

func (r *Repo) AddPostFlairs(ctx context.Context, postFlairs ...models.PostFlair) ([]models.PostFlair, error) {
	query := `
        INSERT INTO
            post_flairs (
                id,
                voxsphere_id,
                full_text,
                background_color
            )
        VALUES 
    `
	args := make([]interface{}, 0)
	placeholders := make([]string, 0)

	for _, postFlair := range postFlairs {
		placeholders = append(placeholders, "(?, ?, ?, ?)")

		args = append(args,
			postFlair.ID,
			postFlair.VoxsphereID,
			postFlair.FullText,
			postFlair.BackgroundColor,
		)
	}

	query += strings.Join(placeholders, ", ")
	query += " RETURNING *"

	if _, err := r.db.NewRaw(
		query, args...).Exec(ctx, &postFlairs); err != nil {
		var pgdriverErr pgdriver.Error
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgUniqueViolation {
			return nil, ErrPostFlairDuplicateID
		}
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgConstraintViolation {
			return nil, ErrPostFlairParentTableRecordNotFound
		}
		return nil, err
	}

	return postFlairs, nil
}

func (r *Repo) UpdatePostFlair(ctx context.Context, postFlair models.PostFlair) (models.PostFlair, error) {
	query := `
                UPDATE
                    post_flairs
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
		postFlair.VoxsphereID,
		postFlair.FullText,
		postFlair.BackgroundColor,
		postFlair.ID,
	).Exec(ctx, &postFlair)
	if errors.Is(err, sql.ErrNoRows) {
		return models.PostFlair{}, ErrPostFlairNotFound
	}
	if err != nil {
		var pgdriverErr pgdriver.Error
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgConstraintViolation {
			return models.PostFlair{}, ErrPostFlairParentTableRecordNotFound
		}
		return models.PostFlair{}, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return models.PostFlair{}, err
	}
	if rowsAffected == 0 {
		return models.PostFlair{}, ErrPostFlairNotFound
	}
	return postFlair, nil
}

func (r *Repo) DeletePostFlair(ctx context.Context, ID uuid.UUID) error {
	query := `
                DELETE FROM 
                    post_flairs
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
		return ErrPostFlairNotFound
	}
	return nil
}
