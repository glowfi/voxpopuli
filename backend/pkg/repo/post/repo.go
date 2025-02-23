package postrepo

import (
	"context"
	"database/sql"
	"errors"
	"time"

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
	ErrPostNotFound                  = errors.New("post not found")
	ErrPostDuplicateID               = errors.New("post duplicate id")
	ErrPostParentTableRecordNotFound = errors.New("record does not exist in the parent table")
)

type Repository interface {
	Posts() ([]models.Post, error)
	PostByID(context.Context, uuid.UUID) (models.Post, error)
	AddPost(context.Context, models.Post) (models.Post, error)
	UpdatePost(context.Context, models.Post) (models.Post, error)
	DeletePost(context.Context, uuid.UUID) error
}

type Repo struct {
	db *bun.DB
}

func NewRepo(db *bun.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) Posts(ctx context.Context) ([]models.Post, error) {
	var posts []models.Post

	query := `
        SELECT
            p.id,
            p.author_id,
            p.voxsphere_id,
            p.title,
            p.text,
            p.text_html,
            p.ups,
            p.over18,
            p.spoiler,
            p.created_at,
            p.created_at_unix,
            p.updated_at
        FROM
            posts p;
    `

	_, err := r.db.NewRaw(query).Exec(ctx, &posts)
	if err != nil {
		return []models.Post{}, err
	}
	return posts, nil
}

func (r *Repo) PostByID(ctx context.Context, ID uuid.UUID) (models.Post, error) {
	var post models.Post

	query := `
        SELECT
            p.id,
            p.author_id,
            p.voxsphere_id,
            p.title,
            p.text,
            p.text_html,
            p.ups,
            p.over18,
            p.spoiler,
            p.created_at,
            p.created_at_unix,
            p.updated_at
        FROM
            posts p
        WHERE
            p.id = ?;
    `

	_, err := r.db.NewRaw(query, ID).Exec(ctx, &post)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Post{}, ErrPostNotFound
		}
		return models.Post{}, err
	}
	return post, nil
}

func (r *Repo) AddPost(ctx context.Context, post models.Post) (models.Post, error) {
	query := `
        INSERT INTO
            posts (
                id,
                author_id,
                voxsphere_id,
                title,
                text,
                text_html,
                ups,
                over18,
                spoiler,
                created_at,
                created_at_unix,
                updated_at
            )
        VALUES (
            ?,
            ?,
            ?,
            ?,
            ?,
            ?,
            ?,
            ?,
            ?,
            ?,
            ?,
            ?
        )
        RETURNING *
    `

	timestamp := time.Now()
	post.CreatedAt = timestamp
	post.UpdatedAt = timestamp
	post.CreatedAtUnix = timestamp.Unix()

	if _, err := r.db.NewRaw(query,
		post.ID,
		post.AuthorID,
		post.VoxsphereID,
		post.Title,
		post.Text,
		post.TextHtml,
		post.Ups,
		post.Over18,
		post.Spoiler,
		post.CreatedAt,
		post.CreatedAtUnix,
		post.UpdatedAt).Exec(ctx, &post); err != nil {
		var pgdriverErr pgdriver.Error
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgUniqueViolation {
			return models.Post{}, ErrPostDuplicateID
		}
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgConstraintViolation {
			return models.Post{}, ErrPostParentTableRecordNotFound
		}
		return models.Post{}, err
	}

	return post, nil
}

func (r *Repo) UpdatePost(ctx context.Context, post models.Post) (models.Post, error) {
	query := `
        UPDATE
            posts
        SET
            author_id = ?,
            voxsphere_id = ?,
            title = ?,
            text = ?,
            text_html = ?,
            ups = ?,
            over18 = ?,
            spoiler = ?,
            updated_at = ?
        WHERE
            id = ?
        RETURNING *
    `

	timestamp := time.Now()
	post.UpdatedAt = timestamp

	res, err := r.db.NewRaw(query,
		post.AuthorID,
		post.VoxsphereID,
		post.Title,
		post.Text,
		post.TextHtml,
		post.Ups,
		post.Over18,
		post.Spoiler,
		post.UpdatedAt,
		post.ID,
	).Exec(ctx, &post)
	if err != nil {
		var pgdriverErr pgdriver.Error
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgConstraintViolation {
			return models.Post{}, ErrPostParentTableRecordNotFound
		}
		if errors.Is(err, sql.ErrNoRows) {
			return models.Post{}, ErrPostNotFound
		}
		return models.Post{}, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return models.Post{}, err
	}
	if rowsAffected == 0 {
		return models.Post{}, ErrPostNotFound
	}
	return post, nil
}

func (r *Repo) DeletePost(ctx context.Context, ID uuid.UUID) error {
	query := `
        DELETE FROM 
            posts
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
		return ErrPostNotFound
	}
	return nil
}
