package topic

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
	ErrTopicNotFound          = errors.New("topic not found")
	ErrTopicDuplicateIDorName = errors.New("topic duplicate id or name")
)

type TopicRepository interface {
	Topics(context.Context) ([]models.Topic, error)
	TopicByID(context.Context, uuid.UUID) (models.Topic, error)
	AddTopics(context.Context, ...models.Topic) ([]models.Topic, error)
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
                name,
                category
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
                    name,
                    category
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

func (r *Repo) AddTopics(ctx context.Context, topics ...models.Topic) ([]models.Topic, error) {
	query := `
        INSERT INTO
            topics (
                id,
                name,
                category
            )
        VALUES 
    `
	args := make([]interface{}, 0)
	placeholders := make([]string, 0)

	for _, topic := range topics {
		placeholders = append(placeholders, "(?, ?, ?)")

		args = append(args,
			topic.ID,
			topic.Name,
			topic.Category,
		)
	}

	query += strings.Join(placeholders, ", ")
	query += " RETURNING *"

	if _, err := r.db.NewRaw(query, args...).Exec(ctx, &topics); err != nil {
		var pgdriverErr pgdriver.Error
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgUniqueViolation {
			return nil, ErrTopicDuplicateIDorName
		}
		return nil, err
	}

	return topics, nil
}

func (r *Repo) UpdateTopic(ctx context.Context, topic models.Topic) (models.Topic, error) {
	query := `
                UPDATE
                    topics
                SET
                    name = ?,
                    category = ?
                WHERE
                    id = ?
                RETURNING *
            `

	res, err := r.db.NewRaw(query, topic.Name, topic.Category, topic.ID).Exec(ctx, &topic)
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
