package award

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
	pgUniqueViolation = "23505"
)

var (
	ErrAwardNotFound           = errors.New("award not found")
	ErrAwardDuplicateIDorTitle = errors.New("award duplicate id or title")
)

type AwardRepository interface {
	Awards(context.Context) ([]models.Award, error)
	AwardByID(context.Context, uuid.UUID) (models.Award, error)
	AddAwards(context.Context, ...models.Award) ([]models.Award, error)
	UpdateAward(context.Context, models.Award) (models.Award, error)
	DeleteAward(context.Context, uuid.UUID) error
}

type Repo struct {
	db *bun.DB
}

func NewRepo(db *bun.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) Awards(ctx context.Context) ([]models.Award, error) {
	var awards []models.Award

	query := `
                SELECT
                    id,
                    title,
                    image_link
                FROM
                    awards;
            `

	_, err := r.db.NewRaw(query).Exec(ctx, &awards)
	if err != nil {
		return []models.Award{}, err
	}
	return awards, nil
}

func (r *Repo) AwardByID(ctx context.Context, ID uuid.UUID) (models.Award, error) {
	var award models.Award

	query := `
                SELECT
                    id,
                    title,
                    image_link
                FROM
                    awards
                WHERE
                    id = ?;
            `
	_, err := r.db.NewRaw(query, ID).Exec(ctx, &award)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Award{}, ErrAwardNotFound
		}
		return models.Award{}, err
	}
	return award, nil
}

func (r *Repo) AddAwards(ctx context.Context, awards ...models.Award) ([]models.Award, error) {
	query := `
        INSERT INTO
            awards (
                id,
                title,
                image_link
            )
        VALUES 
    `

	args := make([]interface{}, 0)
	placeholders := make([]string, 0)
	for _, award := range awards {
		placeholders = append(placeholders, "(?, ?, ?)")
		args = append(args, award.ID, award.Title, award.ImageLink)
	}
	query += strings.Join(placeholders, ", ") + " RETURNING *"

	if _, err := r.db.NewRaw(query, args...).Exec(ctx, &awards); err != nil {
		var pgdriverErr pgdriver.Error
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgUniqueViolation {
			return nil, ErrAwardDuplicateIDorTitle
		}
		return nil, err
	}

	return awards, nil
}

func (r *Repo) UpdateAward(ctx context.Context, award models.Award) (models.Award, error) {
	query := `
                UPDATE
                    awards
                SET
                    title = ?,
                    image_link = ?
                WHERE
                    id = ?
                RETURNING *
            `

	res, err := r.db.NewRaw(query, award.Title, award.ImageLink, award.ID).Exec(ctx, &award)
	if errors.Is(err, sql.ErrNoRows) {
		return models.Award{}, ErrAwardNotFound
	}
	if err != nil {
		return models.Award{}, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return models.Award{}, err
	}
	if rowsAffected == 0 {
		return models.Award{}, ErrAwardNotFound
	}
	return award, nil
}

func (r *Repo) DeleteAward(ctx context.Context, ID uuid.UUID) error {
	query := `
                DELETE FROM 
                    awards
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
		return ErrAwardNotFound
	}
	return nil
}
