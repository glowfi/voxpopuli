package topic

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
	ErrTopicNotFound          = errors.New("topic not found")
	ErrTopicDuplicateIDorName = errors.New("topic duplicate id or name")
)

type Repository interface {
	Topics() ([]models.Topic, error)
	TopicByID(context.Context, uuid.UUID) (models.Topic, error)
	AddTopic(context.Context, models.Topic) (models.Topic, error)
	UpdateTopic(context.Context, models.Topic) (models.Topic, error)
	DeleteTopic(context.Context, uuid.UUID) error
}

type Repo struct {
	db *bun.DB
}

func NewRepo(db *bun.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) Topics(ctx context.Context) ([]models.Topic, error) {
	var topics []models.Topic

	query := `
            SELECT
                id,
                name
            FROM
                topics;
            `

	_, err := r.db.NewRaw(query).Exec(ctx, &topics)
	if err != nil {
		return []models.Topic{}, err
	}
	return topics, nil
}

func (r *Repo) TopicByID(ctx context.Context, ID uuid.UUID) (models.Topic, error) {
	var topic models.Topic

	query := `
                SELECT
                    id,
                    name
                FROM
                    topics
                WHERE
                    id = ?;
            `
	_, err := r.db.NewRaw(query, ID).Exec(ctx, &topic)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Topic{}, ErrTopicNotFound
		}
		return models.Topic{}, err
	}
	return topic, nil
}

func (r *Repo) AddTopic(ctx context.Context, topic models.Topic) (models.Topic, error) {
	query := `
                INSERT INTO
                    topics (
                        id,
                        name
                    )
                VALUES (
                    ?,
                    ?
                )
                RETURNING *
            `

	if _, err := r.db.NewRaw(query, topic.ID, topic.Name).Exec(ctx, &topic); err != nil {
		var pgdriverErr pgdriver.Error
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgUniqueViolation {
			return models.Topic{}, ErrTopicDuplicateIDorName
		}
		return models.Topic{}, err
	}

	return topic, nil
}

func (r *Repo) UpdateTopic(ctx context.Context, topic models.Topic) (models.Topic, error) {
	query := `
                UPDATE
                    topics
                SET
                    name = ?
                WHERE
                    id = ?
                RETURNING *
            `

	res, err := r.db.NewRaw(query, topic.Name, topic.ID).Exec(ctx, &topic)
	if errors.Is(err, sql.ErrNoRows) {
		return models.Topic{}, ErrTopicNotFound
	}
	if err != nil {
		return models.Topic{}, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return models.Topic{}, err
	}
	if rowsAffected == 0 {
		return models.Topic{}, ErrTopicNotFound
	}
	return topic, nil
}

func (r *Repo) DeleteTopic(ctx context.Context, ID uuid.UUID) error {
	query := `
                DELETE FROM 
                    topics
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
		return ErrTopicNotFound
	}
	return nil
}
