package customemoji

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
	ErrCustomEmojiNotFound                  = errors.New("custom emoji not found")
	ErrCustomEmojiDuplicateID               = errors.New("custom emoji duplicate id")
	ErrCustomEmojiParentTableRecordNotFound = errors.New("record does not exist in the parent table")
)

type Repository interface {
	CustomEmojis() ([]models.CustomEmoji, error)
	CustomEmojiByID(context.Context, uuid.UUID) (models.CustomEmoji, error)
	AddCustomEmoji(context.Context, models.CustomEmoji) (models.CustomEmoji, error)
	UpdateCustomEmoji(context.Context, models.CustomEmoji) (models.CustomEmoji, error)
	DeleteCustomEmoji(context.Context, uuid.UUID) error
}

type Repo struct {
	db *bun.DB
}

func NewRepo(db *bun.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) CustomEmojis(ctx context.Context) ([]models.CustomEmoji, error) {
	var customEmojis []models.CustomEmoji

	query := `
                SELECT
                    id,
                    voxsphere_id,
                    url,
                    title
                FROM
                    custom_emojis;
            `

	_, err := r.db.NewRaw(query).Exec(ctx, &customEmojis)
	if err != nil {
		return []models.CustomEmoji{}, err
	}
	return customEmojis, nil
}

func (r *Repo) CustomEmojiByID(ctx context.Context, ID uuid.UUID) (models.CustomEmoji, error) {
	var customEmoji models.CustomEmoji

	query := `
                SELECT
                    id,
                    voxsphere_id,
                    url,
                    title
                FROM
                    custom_emojis
                WHERE
                    id = ?;
            `
	_, err := r.db.NewRaw(query, ID).Exec(ctx, &customEmoji)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.CustomEmoji{}, ErrCustomEmojiNotFound
		}
		return models.CustomEmoji{}, err
	}
	return customEmoji, nil
}

func (r *Repo) AddCustomEmoji(ctx context.Context, customEmoji models.CustomEmoji) (models.CustomEmoji, error) {
	query := `
                INSERT INTO
                    custom_emojis (
                        id,
                        voxsphere_id,
                        url,
                        title
                    )
                VALUES (
                    ?,
                    ?,
                    ?,
                    ?
                )
                RETURNING *
            `

	if _, err := r.db.NewRaw(query, customEmoji.ID, customEmoji.VoxsphereID, customEmoji.Url, customEmoji.Title).Exec(ctx, &customEmoji); err != nil {
		var pgdriverErr pgdriver.Error
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgUniqueViolation {
			return models.CustomEmoji{}, ErrCustomEmojiDuplicateID
		}
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgConstraintViolation {
			return models.CustomEmoji{}, ErrCustomEmojiParentTableRecordNotFound
		}
		return models.CustomEmoji{}, err
	}

	return customEmoji, nil
}

func (r *Repo) UpdateCustomEmoji(ctx context.Context, customEmoji models.CustomEmoji) (models.CustomEmoji, error) {
	query := `
                UPDATE
                    custom_emojis
                SET
                    voxsphere_id = ?,
                    url = ?,
                    title = ?
                WHERE
                    id = ?
                RETURNING *
            `

	res, err := r.db.NewRaw(query, customEmoji.VoxsphereID, customEmoji.Url, customEmoji.Title, customEmoji.ID).Exec(ctx)
	if err != nil {
		var pgdriverErr pgdriver.Error
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgConstraintViolation {
			return models.CustomEmoji{}, ErrCustomEmojiParentTableRecordNotFound
		}
		if errors.Is(err, sql.ErrNoRows) {
			return models.CustomEmoji{}, ErrCustomEmojiNotFound
		}
		return models.CustomEmoji{}, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return models.CustomEmoji{}, err
	}
	if rowsAffected == 0 {
		return models.CustomEmoji{}, ErrCustomEmojiNotFound
	}
	return customEmoji, nil
}

func (r *Repo) DeleteCustomEmoji(ctx context.Context, ID uuid.UUID) error {
	query := `
                DELETE FROM 
                    custom_emojis
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
		return ErrCustomEmojiNotFound
	}
	return nil
}
