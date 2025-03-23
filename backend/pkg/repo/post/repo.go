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
	PostsPaginated(ctx context.Context, skip, limit int) ([]models.PostPaginated, error)
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

func (r *Repo) PostsPaginated(ctx context.Context, skip, limit int) ([]models.PostPaginated, error) {
	var posts []models.PostPaginated

	query := `
        WITH
          ps AS (
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
            LIMIT
              ?
            OFFSET
              ?
          )
        SELECT
          ps.*,
          m.media_type as media_type,
          CASE
            WHEN m.media_type = 'image' THEN (
              SELECT
                JSON_AGG(imageMetadata)
              FROM
                (
                  SELECT
                    JSON_BUILD_OBJECT(
                      'id',
                      imageMetadata.id,
                      'image_id',
                      imageMetadata.image_id,
                      'height',
                      imageMetadata.height,
                      'width',
                      imageMetadata.width,
                      'url',
                      imageMetadata.url,
                      'created_at',
                      TO_CHAR(
                        imageMetadata.created_at AT TIME ZONE 'UTC',
                        'YYYY-MM-DD"T"HH24:MI:SS"Z"'
                      ),
                      'created_at_unix',
                      imageMetadata.created_at_unix,
                      'updated_at',
                      TO_CHAR(
                        imageMetadata.updated_at AT TIME ZONE 'UTC',
                        'YYYY-MM-DD"T"HH24:MI:SS"Z"'
                      )
                    ) AS imageMetadata
                  FROM
                    image_metadatas imageMetadata
                  WHERE
                    imageMetadata.image_id IN (
                      SELECT
                        i.id
                      FROM
                        images i
                      WHERE
                        i.media_id = m.id
                    )
                  ORDER BY
                    imageMetadata.height
                )
            )
            WHEN m.media_type = 'gif' THEN (
              SELECT
                JSON_AGG(gifMetadata)
              FROM
                (
                  SELECT
                    JSON_BUILD_OBJECT(
                      'id',
                      gifMetadata.id,
                      'gif_id',
                      gifMetadata.gif_id,
                      'height',
                      gifMetadata.height,
                      'width',
                      gifMetadata.width,
                      'url',
                      gifMetadata.url,
                      'created_at',
                      TO_CHAR(
                        gifMetadata.created_at AT TIME ZONE 'UTC',
                        'YYYY-MM-DD"T"HH24:MI:SS"Z"'
                      ),
                      'created_at_unix',
                      gifMetadata.created_at_unix,
                      'updated_at',
                      TO_CHAR(
                        gifMetadata.updated_at AT TIME ZONE 'UTC',
                        'YYYY-MM-DD"T"HH24:MI:SS"Z"'
                      )
                    ) AS gifMetadata
                  FROM
                    gif_metadatas gifMetadata
                  WHERE
                    gifMetadata.gif_id IN (
                      SELECT
                        g.id
                      FROM
                        gifs g
                      WHERE
                        g.media_id = m.id
                    )
                  ORDER BY
                    gifMetadata.height
                )
            )
            WHEN m.media_type = 'gallery' THEN (
              SELECT
                JSON_AGG(galleryMetadata)
              FROM
                (
                  SELECT
                    JSON_BUILD_OBJECT(
                      'id',
                      galleryMetadata.id,
                      'gallery_id',
                      galleryMetadata.gallery_id,
                      'order_index',
                      galleryMetadata.order_index,
                      'height',
                      galleryMetadata.height,
                      'width',
                      galleryMetadata.width,
                      'url',
                      galleryMetadata.url,
                      'created_at',
                      TO_CHAR(
                        galleryMetadata.created_at AT TIME ZONE 'UTC',
                        'YYYY-MM-DD"T"HH24:MI:SS"Z"'
                      ),
                      'created_at_unix',
                      galleryMetadata.created_at_unix,
                      'updated_at',
                      TO_CHAR(
                        galleryMetadata.updated_at AT TIME ZONE 'UTC',
                        'YYYY-MM-DD"T"HH24:MI:SS"Z"'
                      )
                    ) AS galleryMetadata
                  FROM
                    gallery_metadatas galleryMetadata
                  WHERE
                    galleryMetadata.gallery_id IN (
                      select
                        id
                      from
                        galleries
                      where
                        galleries.media_id = m.id
                    )
                  ORDER BY
                    galleryMetadata.order_index
                )
            )
            WHEN m.media_type = 'video' THEN (
              SELECT
                JSON_AGG(
                  JSON_BUILD_OBJECT(
                    'id',
                    videos.id,
                    'media_id',
                    videos.media_id,
                    'url',
                    videos.url,
                    'height',
                    videos.height,
                    'width',
                    videos.width,
                    'created_at',
                    TO_CHAR(
                      videos.created_at AT TIME ZONE 'UTC',
                      'YYYY-MM-DD"T"HH24:MI:SS"Z"'
                    ),
                    'created_at_unix',
                    videos.created_at_unix,
                    'updated_at',
                    TO_CHAR(
                      videos.updated_at AT TIME ZONE 'UTC',
                      'YYYY-MM-DD"T"HH24:MI:SS"Z"'
                    )
                  )
                )
              FROM
                videos
              WHERE
                videos.media_id = m.id
            )
            WHEN m.media_type = 'link' THEN (
              SELECT
                JSON_AGG(
                  JSON_BUILD_OBJECT(
                    'id',
                    links.id,
                    'media_id',
                    links.media_id,
                    'link',
                    links.link,
                    'image',
                    (
                      SELECT
                        JSON_AGG(
                          JSON_BUILD_OBJECT(
                            'id',
                            imageMetadata.id,
                            'image_id',
                            imageMetadata.image_id,
                            'height',
                            imageMetadata.height,
                            'width',
                            imageMetadata.width,
                            'url',
                            imageMetadata.url,
                            'created_at',
                            TO_CHAR(
                              imageMetadata.created_at AT TIME ZONE 'UTC',
                              'YYYY-MM-DD"T"HH24:MI:SS"Z"'
                            ),
                            'created_at_unix',
                            imageMetadata.created_at_unix,
                            'updated_at',
                            TO_CHAR(
                              imageMetadata.updated_at AT TIME ZONE 'UTC',
                              'YYYY-MM-DD"T"HH24:MI:SS"Z"'
                            )
                          )
                          ORDER BY
                            imageMetadata.height
                        )
                      FROM
                        image_metadatas imageMetadata
                      WHERE
                        imageMetadata.image_id IN (
                          SELECT
                            i.id
                          FROM
                            images i
                          WHERE
                            i.media_id = m.id
                        )
                    ),
                    'created_at',
                    TO_CHAR(
                      links.created_at AT TIME ZONE 'UTC',
                      'YYYY-MM-DD"T"HH24:MI:SS"Z"'
                    ),
                    'created_at_unix',
                    links.created_at_unix,
                    'updated_at',
                    TO_CHAR(
                      links.updated_at AT TIME ZONE 'UTC',
                      'YYYY-MM-DD"T"HH24:MI:SS"Z"'
                    )
                  )
                )
              FROM
                links
              WHERE
                links.media_id = m.id
            )
            ELSE NULL
          END AS medias
        FROM
          ps
          JOIN post_medias m ON ps.id = m.post_id;
    `

	_, err := r.db.NewRaw(query, limit, skip).Exec(ctx, &posts)
	if err != nil {
		return []models.PostPaginated{}, err
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
