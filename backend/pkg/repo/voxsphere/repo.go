package voxsphere

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
	pgUniqueViolation = "23505"
)

var (
	errVoxsphereNotFound    = errors.New("voxsphere not found")
	errVoxsphereDuplicateID = errors.New("voxsphere duplicate id")
)

type Repository interface {
	Voxspheres() ([]models.Voxsphere, error)
	VoxsphereByID(context.Context, uuid.UUID) (models.Voxsphere, error)
	AddVoxsphere(context.Context, models.Voxsphere) (models.Voxsphere, error)
	UpdateVoxsphere(context.Context, models.Voxsphere) (models.Voxsphere, error)
	DeleteVoxsphere(context.Context, uuid.UUID) (models.Voxsphere, error)
}

type Repo struct {
	db *bun.DB
}

func NewRepo(db *bun.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) Voxspheres(ctx context.Context) ([]models.Voxsphere, error) {
	var voxspheres []models.Voxsphere

	query := `
        SELECT 
            id,
            title,
            public_description,
            community_icon,
            banner_background_image,
            banner_background_color,
            key_color,
            primary_color,
            over18,
            spoilers_enabled,
            created_at,
            created_at_unix,
            updated_at
        FROM 
            voxspheres;
    `

	err := r.db.NewRaw(query).NewSelect().Model(&voxspheres).Relation("TopicExpanded").Scan(ctx)
	if err != nil {
		return []models.Voxsphere{}, err
	}
	return voxspheres, nil
}

func (r *Repo) VoxsphereByID(ctx context.Context, ID uuid.UUID) (models.Voxsphere, error) {
	var voxsphere models.Voxsphere

	query := `
        SELECT 
            id,
            title,
            public_description,
            community_icon,
            banner_background_image,
            banner_background_color,
            key_color,
            primary_color,
            over18,
            spoilers_enabled,
            created_at,
            created_at_unix,
            updated_at
        FROM 
            voxspheres
        WHERE 
            id = ?
    `
	err := r.db.NewRaw(query, ID).Scan(ctx, &voxsphere)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Voxsphere{}, errVoxsphereNotFound
		}
		return models.Voxsphere{}, err
	}
	return voxsphere, nil
}

func (r *Repo) AddVoxsphere(ctx context.Context, voxsphere models.Voxsphere) (models.Voxsphere, error) {
	query := `
        INSERT INTO 
            voxspheres (
                id,
                title,
                public_description,
                community_icon,
                banner_background_image,
                banner_background_color,
                key_color,
                primary_color,
                over18,
                spoilers_enabled,
                created_at,
                created_at_unix,
                updated_at
            )
        VALUES (
            $1,
            $2,
            $3,
            $4,
            $5,
            $6,
            $7,
            $8,
            $9,
            $10,
            $11,
            $12,
            $13
        )
    `

	timestamp := time.Now()
	voxsphere.CreatedAt = timestamp
	voxsphere.UpdatedAt = timestamp

	if _, err := r.db.ExecContext(ctx, query,
		voxsphere.ID,
		voxsphere.Title,
		voxsphere.PublicDescription,
		voxsphere.CommunityIcon,
		voxsphere.BannerBackgroundImage,
		voxsphere.BannerBackgroundColor,
		voxsphere.KeyColor,
		voxsphere.PrimaryColor,
		voxsphere.Over18,
		voxsphere.SpoilersEnabled,
		voxsphere.CreatedAt,
		voxsphere.CreatedAtUnix,
		voxsphere.UpdatedAt); err != nil {
		var pgdriverErr pgdriver.Error
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgUniqueViolation {
			return models.Voxsphere{}, errVoxsphereDuplicateID
		}
		return models.Voxsphere{}, err
	}

	return voxsphere, nil
}

func (r *Repo) UpdateVoxsphere(ctx context.Context, voxsphere models.Voxsphere) (models.Voxsphere, error) {
	query := `
        UPDATE 
            voxspheres
        SET 
            title = $1,
            public_description = $2,
            community_icon = $3,
            banner_background_image = $4,
            banner_background_color = $5,
            key_color = $6,
            primary_color = $7,
            over18 = $8,
            spoilers_enabled = $9,
            updated_at = $10
        WHERE 
            id = $11
    `

	voxsphere.UpdatedAt = time.Now()

	res, err := r.db.ExecContext(ctx, query,
		voxsphere.Title,
		voxsphere.PublicDescription,
		voxsphere.CommunityIcon,
		voxsphere.BannerBackgroundImage,
		voxsphere.BannerBackgroundColor,
		voxsphere.KeyColor,
		voxsphere.PrimaryColor,
		voxsphere.Over18,
		voxsphere.SpoilersEnabled,
		voxsphere.UpdatedAt,
		voxsphere.ID,
	)
	if err != nil {
		return models.Voxsphere{}, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return models.Voxsphere{}, err
	}
	if rowsAffected == 0 {
		return models.Voxsphere{}, errVoxsphereNotFound
	}
	return voxsphere, nil
}

func (r *Repo) DeleteVoxsphere(ctx context.Context, ID uuid.UUID) (models.Voxsphere, error) {
	query := `
        DELETE FROM 
            voxspheres
        WHERE 
            id = $1
    `
	res, err := r.db.ExecContext(ctx, query, ID)
	if err != nil {
		return models.Voxsphere{}, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return models.Voxsphere{}, err
	}
	if rowsAffected == 0 {
		return models.Voxsphere{}, errVoxsphereNotFound
	}
	return models.Voxsphere{}, nil
}
