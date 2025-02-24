package comments

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
	ErrCommentNotFound                  = errors.New("comment not found")
	ErrCommentDuplicateID               = errors.New("comment duplicate id")
	ErrCommentParentTableRecordNotFound = errors.New("record does not exist in the parent table")
)

type Repository interface {
	Comments() ([]models.Comment, error)
	CommentByID(context.Context, uuid.UUID) (models.Comment, error)
	AddComment(context.Context, models.Comment) (models.Comment, error)
	UpdateComment(context.Context, models.Comment) (models.Comment, error)
	DeleteComment(context.Context, uuid.UUID) error
}

type Repo struct {
	db *bun.DB
}

func NewRepo(db *bun.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) Comments(ctx context.Context) ([]models.Comment, error) {
	var comments []models.Comment

	query := `
        SELECT
            c.id,
            c.author_id,
            c.parent_comment_id,
            c.post_id,
            c.body,
            c.body_html,
            c.ups,
            c.score,
            c.created_at,
            c.created_at_unix,
            c.updated_at
        FROM
            comments c;
    `

	_, err := r.db.NewRaw(query).Exec(ctx, &comments)
	if err != nil {
		return []models.Comment{}, err
	}
	return comments, nil
}

func (r *Repo) CommentByID(ctx context.Context, ID uuid.UUID) (models.Comment, error) {
	var comment models.Comment

	query := `
        SELECT
            c.id,
            c.author_id,
            c.parent_comment_id,
            c.post_id,
            c.body,
            c.body_html,
            c.ups,
            c.score,
            c.created_at,
            c.created_at_unix,
            c.updated_at
        FROM
            comments c
        WHERE
            c.id = ?;
    `

	_, err := r.db.NewRaw(query, ID).Exec(ctx, &comment)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Comment{}, ErrCommentNotFound
		}
		return models.Comment{}, err
	}
	return comment, nil
}

func (r *Repo) AddComment(ctx context.Context, comment models.Comment) (models.Comment, error) {
	query := `
        INSERT INTO
            comments (
                id,
                author_id,
                parent_comment_id,
                post_id,
                body,
                body_html,
                ups,
                score,
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
            ?
        )
        RETURNING *
    `

	timestamp := time.Now()
	comment.CreatedAt = timestamp
	comment.UpdatedAt = timestamp
	comment.CreatedAtUnix = timestamp.Unix()

	if _, err := r.db.NewRaw(query,
		comment.ID,
		comment.AuthorID,
		comment.ParentCommentID,
		comment.PostID,
		comment.Body,
		comment.BodyHtml,
		comment.Ups,
		comment.Score,
		comment.CreatedAt,
		comment.CreatedAtUnix,
		comment.UpdatedAt).Exec(ctx, &comment); err != nil {
		var pgdriverErr pgdriver.Error
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgUniqueViolation {
			return models.Comment{}, ErrCommentDuplicateID
		}
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgConstraintViolation {
			return models.Comment{}, ErrCommentParentTableRecordNotFound
		}
		return models.Comment{}, err
	}

	return comment, nil
}

func (r *Repo) UpdateComment(ctx context.Context, comment models.Comment) (models.Comment, error) {
	query := `
        UPDATE
            comments
        SET
            author_id = ?,
            parent_comment_id = ?,
            post_id = ?,
            body = ?,
            body_html = ?,
            ups = ?,
            score = ?,
            updated_at = ?
        WHERE
            id = ?
        RETURNING *
    `

	comment.UpdatedAt = time.Now()

	res, err := r.db.NewRaw(query,
		comment.AuthorID,
		comment.ParentCommentID,
		comment.PostID,
		comment.Body,
		comment.BodyHtml,
		comment.Ups,
		comment.Score,
		comment.UpdatedAt,
		comment.ID,
	).Exec(ctx, &comment)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Comment{}, ErrCommentNotFound
		}
		var pgdriverErr pgdriver.Error
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgConstraintViolation {
			return models.Comment{}, ErrCommentParentTableRecordNotFound
		}
		return models.Comment{}, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return models.Comment{}, err
	}
	if rowsAffected == 0 {
		return models.Comment{}, ErrCommentNotFound
	}
	return comment, nil
}

func (r *Repo) DeleteComment(ctx context.Context, ID uuid.UUID) error {
	query := `
        DELETE FROM 
            comments
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
		return ErrCommentNotFound
	}
	return nil
}
