package postrepo

import (
	"context"
	"database/sql"
	"errors"
	"strings"
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

type PostRepository interface {
	PostsPaginated(ctx context.Context, skip int, limit int) ([]models.Post, error)
	Posts(context.Context) ([]models.Post, error)
	PostByID(context.Context, uuid.UUID) (models.Post, error)
	AddPosts(context.Context, ...models.Post) ([]models.Post, error)
	UpdatePost(context.Context, models.Post) (models.Post, error)
	DeletePost(context.Context, uuid.UUID) error
}

type Repo struct {
	db *bun.DB
}

func NewRepo(db *bun.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) PostsPaginated(ctx context.Context, skip, limit int) ([]models.Post, error) {
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
            posts p
        ORDER BY
            p.id
        LIMIT ?
        OFFSET ?;
    `

	_, err := r.db.NewRaw(query, limit, skip).Exec(ctx, &posts)
	if err != nil {
		return []models.Post{}, err
	}
	return posts, nil
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

func (r *Repo) AddPosts(ctx context.Context, posts ...models.Post) ([]models.Post, error) {
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
        VALUES 
    `
	// Create a slice to store the args
	args := make([]interface{}, 0)
	// Create a slice to store the placeholders
	placeholders := make([]string, 0)

	// Iterate over the posts
	for _, post := range posts {
		timestamp := time.Now()
		post.CreatedAt = timestamp
		post.UpdatedAt = timestamp
		post.CreatedAtUnix = timestamp.Unix()

		// Append the placeholders
		placeholders = append(placeholders, "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")

		// Append the values
		args = append(args,
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
			post.UpdatedAt,
		)
	}

	// Join the placeholders
	query += strings.Join(placeholders, ", ")
	query += " RETURNING *"

	if _, err := r.db.NewRaw(query, args...).Exec(ctx, &posts); err != nil {
		var pgdriverErr pgdriver.Error
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgUniqueViolation {
			return nil, ErrPostDuplicateID
		}
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgConstraintViolation {
			return nil, ErrPostParentTableRecordNotFound
		}
		return nil, err
	}

	return posts, nil
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
