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
	ErrVoxsphereNotFound           = errors.New("voxsphere not found")
	ErrVoxsphereDuplicateIDorTitle = errors.New("voxsphere duplicate id or title")
)

type Repository interface {
	Voxspheres() ([]models.Voxsphere, error)
	VoxsphereByID(context.Context, uuid.UUID) (models.Voxsphere, error)
	AddVoxsphere(context.Context, models.Voxsphere) (models.Voxsphere, error)
	UpdateVoxsphere(context.Context, models.Voxsphere) (models.Voxsphere, error)
	DeleteVoxsphere(context.Context, uuid.UUID) error
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
	            v.id,
	            v.title,
	            v.topic_id,
	            json_build_object('id', t.id,'name', t.name) as topic,
	            v.public_description,
	            v.community_icon,
	            v.banner_background_image,
	            v.banner_background_color,
	            v.key_color,
	            v.primary_color,
	            v.over18,
	            v.spoilers_enabled,
	            v.created_at,
	            v.created_at_unix,
	            v.updated_at
	        FROM
	            voxspheres v
	        JOIN
	            topics t ON v.topic_id = t.id;
	    `

	_, err := r.db.NewRaw(query).Exec(ctx, &voxspheres)
	// err := r.db.NewSelect().Model(&voxspheres).Relation("Topic").Scan(ctx)
	if err != nil {
		return []models.Voxsphere{}, err
	}
	return voxspheres, nil
}

func (r *Repo) VoxsphereByID(ctx context.Context, ID uuid.UUID) (models.Voxsphere, error) {
	var voxsphere models.Voxsphere

	query := `
	        SELECT
	            v.id,
	            v.title,
	            v.topic_id,
	            json_build_object('id', t.id,'name', t.name) as topic,
	            v.public_description,
	            v.community_icon,
	            v.banner_background_image,
	            v.banner_background_color,
	            v.key_color,
	            v.primary_color,
	            v.over18,
	            v.spoilers_enabled,
	            v.created_at,
	            v.created_at_unix,
	            v.updated_at
	        FROM
	            voxspheres v
	        JOIN
	            topics t ON v.topic_id = t.id
	        WHERE
	            v.id = ?;
	    `
	_, err := r.db.NewRaw(query, ID).Exec(ctx, &voxsphere)
	// err := r.db.NewSelect().Model(&voxsphere).Relation("Topic").Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Voxsphere{}, ErrVoxsphereNotFound
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
	                topic_id,
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
	            ?,
	            ?,
	            ?
	        )
	    `

	timestamp := time.Now()
	voxsphere.CreatedAt = timestamp
	voxsphere.UpdatedAt = timestamp

	if _, err := r.db.NewRaw(query,
		voxsphere.ID,
		voxsphere.Title,
		voxsphere.TopicID,
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
		voxsphere.UpdatedAt).Exec(ctx); err != nil {
		// if _, err := r.db.NewInsert().Model(&voxsphere).Exec(context.Background()); err != nil {
		var pgdriverErr pgdriver.Error
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgUniqueViolation {
			return models.Voxsphere{}, ErrVoxsphereDuplicateIDorTitle
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
	            title = ?,
	            public_description = ?,
	            community_icon = ?,
	            banner_background_image = ?,
	            banner_background_color = ?,
	            key_color = ?,
	            primary_color = ?,
	            over18 = ?,
	            spoilers_enabled = ?,
	            updated_at = ?
	        WHERE
	            id = ?
	    `

	voxsphere.UpdatedAt = time.Now()

	res, err := r.db.NewRaw(query,
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
	).Exec(ctx)
	// res, err := r.db.NewUpdate().
	// 	Model(&voxsphere).
	// 	Column("id", "topic_id", "title", "public_description", "community_icon", "banner_background_image", "banner_background_color", "key_color", "primary_color", "over18", "spoilers_enabled", "created_at", "created_at_unix", "updated_at").
	// 	Where("id = ?", voxsphere.ID).
	// 	Exec(context.Background())
	if err != nil {
		return models.Voxsphere{}, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return models.Voxsphere{}, err
	}
	if rowsAffected == 0 {
		return models.Voxsphere{}, ErrVoxsphereNotFound
	}
	return voxsphere, nil
}

func (r *Repo) DeleteVoxsphere(ctx context.Context, ID uuid.UUID) error {
	query := `
        DELETE FROM 
            voxspheres
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
		return ErrVoxsphereNotFound
	}
	return nil
}
