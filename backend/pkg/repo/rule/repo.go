package rule

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
	ErrRuleNotFound          = errors.New("rule not found")
	ErrRuleDuplicateIDorName = errors.New("rule duplicate id or name")
)

type RuleRepository interface {
	Rules() ([]models.Rule, error)
	RuleByID(context.Context, uuid.UUID) (models.Rule, error)
	RulesByVoxsphereID(context.Context, uuid.UUID) ([]models.Rule, error)
	AddRules(context.Context, ...models.Rule) ([]models.Rule, error)
	UpdateRule(context.Context, models.Rule) (models.Rule, error)
	DeleteRule(context.Context, uuid.UUID) error
}

type RuleRepo struct {
	db *bun.DB
}

func NewRepo(db *bun.DB) *RuleRepo {
	return &RuleRepo{db: db}
}

func (r *RuleRepo) Rules(ctx context.Context) ([]models.Rule, error) {
	var rules []models.Rule

	query := `
                SELECT
                    id,
                    voxsphere_id,
                    short_name,
                    description
                FROM
                    rules;
            `

	_, err := r.db.NewRaw(query).Exec(ctx, &rules)
	if err != nil {
		return []models.Rule{}, err
	}
	return rules, nil
}

func (r *RuleRepo) RuleByID(ctx context.Context, ID uuid.UUID) (models.Rule, error) {
	var rule models.Rule

	query := `
                SELECT
                    id,
                    voxsphere_id,
                    short_name,
                    description
                FROM
                    rules
                WHERE
                    id = ?;
            `
	_, err := r.db.NewRaw(query, ID).Exec(ctx, &rule)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Rule{}, ErrRuleNotFound
		}
		return models.Rule{}, err
	}
	return rule, nil
}

func (r *RuleRepo) RulesByVoxsphereID(ctx context.Context, voxsphereID uuid.UUID) ([]models.Rule, error) {
	var rules []models.Rule

	query := `
                SELECT
                    id,
                    voxsphere_id,
                    short_name,
                    description
                FROM
                    rules
                WHERE
                    voxsphere_id = ?;
            `
	_, err := r.db.NewRaw(query, voxsphereID).Exec(ctx, &rules)
	if err != nil {
		return []models.Rule{}, err
	}
	return rules, nil
}

func (r *RuleRepo) AddRules(ctx context.Context, rules ...models.Rule) ([]models.Rule, error) {
	query := `
        INSERT INTO
            rules (
                id,
                voxsphere_id,
                short_name,
                description
            )
        VALUES 
    `
	args := make([]interface{}, 0)
	placeholders := make([]string, 0)

	for _, rule := range rules {
		placeholders = append(placeholders, "(?, ?, ?, ?)")

		args = append(args,
			rule.ID,
			rule.VoxsphereID,
			rule.ShortName,
			rule.Description,
		)
	}

	query += strings.Join(placeholders, ", ")
	query += " RETURNING *"

	if _, err := r.db.NewRaw(query, args...).Exec(ctx, &rules); err != nil {
		var pgdriverErr pgdriver.Error
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgUniqueViolation {
			return nil, ErrRuleDuplicateIDorName
		}
		return nil, err
	}

	return rules, nil
}

func (r *RuleRepo) UpdateRule(ctx context.Context, rule models.Rule) (models.Rule, error) {
	query := `
                UPDATE
                    rules
                SET
                    voxsphere_id = ?,
                    short_name = ?,
                    description = ?
                WHERE
                    id = ?
                RETURNING *
            `

	res, err := r.db.NewRaw(query,
		rule.VoxsphereID,
		rule.ShortName,
		rule.Description,
		rule.ID,
	).Exec(ctx, &rule)
	if errors.Is(err, sql.ErrNoRows) {
		return models.Rule{}, ErrRuleNotFound
	}
	if err != nil {
		return models.Rule{}, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return models.Rule{}, err
	}
	if rowsAffected == 0 {
		return models.Rule{}, ErrRuleNotFound
	}
	return rule, nil
}

func (r *RuleRepo) DeleteRule(ctx context.Context, ID uuid.UUID) error {
	query := `
        DELETE FROM 
            rules
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
		return ErrRuleNotFound
	}
	return nil
}
