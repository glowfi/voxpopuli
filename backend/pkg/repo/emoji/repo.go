package emoji

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
	ErrEmojiNotFound          = errors.New("emoji not found")
	ErrEmojiDuplicateIDorText = errors.New("emoji duplicate id or text")
)

type Repository interface {
	Emojis() ([]models.Emoji, error)
	EmojiByID(context.Context, uuid.UUID) (models.Emoji, error)
	AddEmoji(context.Context, models.Emoji) (models.Emoji, error)
	UpdateEmoji(context.Context, models.Emoji) (models.Emoji, error)
	DeleteEmoji(context.Context, uuid.UUID) error
}

type Repo struct {
	db *bun.DB
}

func NewRepo(db *bun.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) Emojis(ctx context.Context) ([]models.Emoji, error) {
	var emojis []models.Emoji

	query := `
                SELECT
                    id,
                    title
                FROM
                    emojis;
            `

	_, err := r.db.NewRaw(query).Exec(ctx, &emojis)
	if err != nil {
		return []models.Emoji{}, err
	}
	return emojis, nil
}

func (r *Repo) EmojiByID(ctx context.Context, ID uuid.UUID) (models.Emoji, error) {
	var emoji models.Emoji

	query := `
                SELECT
                    id,
                    title
                FROM
                    emojis
                WHERE
                    id = ?;
            `
	_, err := r.db.NewRaw(query, ID).Exec(ctx, &emoji)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Emoji{}, ErrEmojiNotFound
		}
		return models.Emoji{}, err
	}
	return emoji, nil
}

func (r *Repo) AddEmoji(ctx context.Context, emoji models.Emoji) (models.Emoji, error) {
	query := `
                INSERT INTO
                    emojis (
                        id,
                        title
                    )
                VALUES (
                    ?,
                    ?
                )
                RETURNING *
            `

	if _, err := r.db.NewRaw(query, emoji.ID, emoji.Title).Exec(ctx, &emoji); err != nil {
		var pgdriverErr pgdriver.Error
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgUniqueViolation {
			return models.Emoji{}, ErrEmojiDuplicateIDorText
		}
		return models.Emoji{}, err
	}

	return emoji, nil
}

func (r *Repo) UpdateEmoji(ctx context.Context, emoji models.Emoji) (models.Emoji, error) {
	query := `
                UPDATE
                    emojis
                SET
                    title = ?
                WHERE
                    id = ?
                RETURNING *
            `

	res, err := r.db.NewRaw(query, emoji.Title, emoji.ID).Exec(ctx, &emoji)
	if errors.Is(err, sql.ErrNoRows) {
		return models.Emoji{}, ErrEmojiNotFound
	}
	if err != nil {
		return models.Emoji{}, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return models.Emoji{}, err
	}
	if rowsAffected == 0 {
		return models.Emoji{}, ErrEmojiNotFound
	}
	return emoji, nil
}

func (r *Repo) DeleteEmoji(ctx context.Context, ID uuid.UUID) error {
	query := `
                DELETE FROM 
                    emojis
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
		return ErrEmojiNotFound
	}
	return nil
}
