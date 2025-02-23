package user

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
	ErrUserNotFound          = errors.New("user not found")
	ErrUserDuplicateIDorName = errors.New("user duplicate id or name")
)

type Repository interface {
	Users() ([]models.User, error)
	UserByID(context.Context, uuid.UUID) (models.User, error)
	AddUser(context.Context, models.User) (models.User, error)
	UpdateUser(context.Context, models.User) (models.User, error)
	DeleteUser(context.Context, uuid.UUID) error
}

type Repo struct {
	db *bun.DB
}

func NewRepo(db *bun.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) Users(ctx context.Context) ([]models.User, error) {
	var users []models.User

	query := `
                SELECT
                    id,
                    name,
                    public_description,
                    avatar_img,
                    banner_img,
                    iconcolor,
                    keycolor,
                    primarycolor,
                    over18,
                    suspended,
                    created_at,
                    created_at_unix,
                    updated_at
                FROM
                    users;
            `

	_, err := r.db.NewRaw(query).Exec(ctx, &users)
	if err != nil {
		return []models.User{}, err
	}
	return users, nil
}

func (r *Repo) UserByID(ctx context.Context, ID uuid.UUID) (models.User, error) {
	var user models.User

	query := `
                SELECT
                    id,
                    name,
                    public_description,
                    avatar_img,
                    banner_img,
                    iconcolor,
                    keycolor,
                    primarycolor,
                    over18,
                    suspended,
                    created_at,
                    created_at_unix,
                    updated_at
                FROM
                    users
                WHERE
                    id = ?;
            `
	_, err := r.db.NewRaw(query, ID).Exec(ctx, &user)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, ErrUserNotFound
		}
		return models.User{}, err
	}
	return user, nil
}

func (r *Repo) AddUser(ctx context.Context, user models.User) (models.User, error) {
	query := `
                INSERT INTO
                    users (
                        id,
                        name,
                        public_description,
                        avatar_img,
                        banner_img,
                        iconcolor,
                        keycolor,
                        primarycolor,
                        over18,
                        suspended,
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
                    ?
                )
                RETURNING *
            `

	timestamp := time.Now()
	user.CreatedAt = timestamp
	user.UpdatedAt = timestamp
	user.CreatedAtUnix = timestamp.Unix()

	if _, err := r.db.NewRaw(
		query,
		user.ID,
		user.Name,
		user.PublicDescription,
		user.AvatarImg,
		user.BannerImg,
		user.Iconcolor,
		user.Keycolor,
		user.Primarycolor,
		user.Over18,
		user.Suspended,
		user.CreatedAt,
		user.CreatedAtUnix,
		user.UpdatedAt,
	).Exec(ctx, &user); err != nil {
		var pgdriverErr pgdriver.Error
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgUniqueViolation {
			return models.User{}, ErrUserDuplicateIDorName
		}
		return models.User{}, err
	}

	return user, nil
}

func (r *Repo) UpdateUser(ctx context.Context, user models.User) (models.User, error) {
	query := `
                UPDATE
                    users
                SET
                    name = ?,
                    public_description = ?,
                    avatar_img = ?,
                    banner_img = ?,
                    iconcolor = ?,
                    keycolor = ?,
                    primarycolor = ?,
                    over18 = ?,
                    suspended = ?,
                    updated_at = ?
                WHERE
                    id = ?
                RETURNING *
            `

	user.UpdatedAt = time.Now()

	res, err := r.db.NewRaw(
		query,
		user.Name,
		user.PublicDescription,
		user.AvatarImg,
		user.BannerImg,
		user.Iconcolor,
		user.Keycolor,
		user.Primarycolor,
		user.Over18,
		user.Suspended,
		time.Now(),
		user.ID,
	).Exec(ctx, &user)
	if errors.Is(err, sql.ErrNoRows) {
		return models.User{}, ErrUserNotFound
	}
	if err != nil {
		return models.User{}, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return models.User{}, err
	}
	if rowsAffected == 0 {
		return models.User{}, ErrUserNotFound
	}
	return user, nil
}

func (r *Repo) DeleteUser(ctx context.Context, ID uuid.UUID) error {
	query := `
                DELETE FROM 
                    users
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
		return ErrUserNotFound
	}
	return nil
}
