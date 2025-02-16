package trophy

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
	pgUniqueViolation = "23505"
)

var (
	ErrTrophyNotFound           = errors.New("trophy not found")
	ErrTrophyDuplicateIDorTitle = errors.New("trophy duplicate id or title")
)

type Repository interface {
	Trophies() ([]models.Trophy, error)
	TrophyByID(context.Context, uuid.UUID) (models.Trophy, error)
	AddTrophy(context.Context, models.Trophy) (models.Trophy, error)
	UpdateTrophy(context.Context, models.Trophy) (models.Trophy, error)
	DeleteTrophy(context.Context, uuid.UUID) error
}

type Repo struct {
	db *bun.DB
}

func NewRepo(db *bun.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) Trophies(ctx context.Context) ([]models.Trophy, error) {
	var trophies []models.Trophy

	query := `
                SELECT
                    id,
                    title,
                    description,
                    image_link
                FROM
                    trophies;
            `

	_, err := r.db.NewRaw(query).Exec(ctx, &trophies)
	if err != nil {
		return []models.Trophy{}, err
	}
	return trophies, nil
}

func (r *Repo) TrophyByID(ctx context.Context, ID uuid.UUID) (models.Trophy, error) {
	var trophy models.Trophy

	query := `
                SELECT
                    id,
                    title,
                    description,
                    image_link
                FROM
                    trophies
                WHERE
                    id = ?;
            `
	_, err := r.db.NewRaw(query, ID).Exec(ctx, &trophy)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Trophy{}, ErrTrophyNotFound
		}
		return models.Trophy{}, err
	}
	return trophy, nil
}

func (r *Repo) AddTrophy(ctx context.Context, trophy models.Trophy) (models.Trophy, error) {
	query := `
                INSERT INTO
                    trophies (
                        id,
                        title,
                        description,
                        image_link
                    )
                VALUES (
                    ?,
                    ?,
                    ?,
                    ?
                )
            `

	if _, err := r.db.NewRaw(query, trophy.ID, trophy.Title, trophy.Description, trophy.ImageLink).Exec(ctx); err != nil {
		var pgdriverErr pgdriver.Error
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgUniqueViolation {
			return models.Trophy{}, ErrTrophyDuplicateIDorTitle
		}
		return models.Trophy{}, err
	}

	return trophy, nil
}

func (r *Repo) UpdateTrophy(ctx context.Context, trophy models.Trophy) (models.Trophy, error) {
	query := `
                UPDATE
                    trophies
                SET
                    title = ?,
                    description = ?,
                    image_link = ?
                WHERE
                    id = ?
            `

	res, err := r.db.NewRaw(query, trophy.Title, trophy.Description, trophy.ImageLink, trophy.ID).Exec(ctx)
	if err != nil {
		return models.Trophy{}, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return models.Trophy{}, err
	}
	if rowsAffected == 0 {
		return models.Trophy{}, ErrTrophyNotFound
	}
	return trophy, nil
}

func (r *Repo) DeleteTrophy(ctx context.Context, ID uuid.UUID) error {
	query := `
                DELETE FROM 
                    trophies
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
		return ErrTrophyNotFound
	}
	return nil
}
