package relation

import (
	"context"
	"errors"

	"github.com/glowfi/voxpopuli/backend/pkg/models"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/driver/pgdriver"
)

const (
	pgUniqueViolation     = "23505"
	pgConstraintViolation = "23503"
)

var (
	ErrParentTableRecordNotFound = errors.New("record does not exist in the parent table")
	ErrDuplicateID               = errors.New("duplicate id")
)

type Repository interface {
	UserTrophies(ctx context.Context) ([]models.UserTrophy, error)
	LinkUserTrophy(ctx context.Context, ut models.UserTrophy) (models.UserTrophy, error)

	VoxsphereMembers(ctx context.Context) ([]models.VoxsphereMember, error)
	LinkVoxsphereMember(ctx context.Context, vme models.VoxsphereMember) (models.VoxsphereMember, error)

	VoxsphereModerators(ctx context.Context) ([]models.VoxsphereModerator, error)
	LinkVoxsphereModerator(ctx context.Context, vmod models.VoxsphereModerator) (models.VoxsphereModerator, error)

	UserFlairEmojis(ctx context.Context) ([]models.UserFlairEmoji, error)
	LinkUserFlairEmoji(ctx context.Context, ufe models.UserFlairEmoji) (models.UserFlairEmoji, error)

	UserFlairCustomEmojis(ctx context.Context) ([]models.UserFlairCustomEmoji, error)
	LinkUserFlairCustomEmoji(ctx context.Context, ufce models.UserFlairCustomEmoji) (models.UserFlairCustomEmoji, error)

	UserFlairDescriptions(ctx context.Context) ([]models.UserFlairDescription, error)
	LinkUserFlairDescription(ctx context.Context, ufd models.UserFlairDescription) (models.UserFlairDescription, error)

	PostFlairEmojis(ctx context.Context) ([]models.PostFlairEmoji, error)
	LinkPostFlairEmoji(ctx context.Context, pfe models.PostFlairEmoji) (models.PostFlairEmoji, error)

	PostFlairCustomEmojis(ctx context.Context) ([]models.PostFlairCustomEmoji, error)
	LinkPostFlairCustomEmoji(ctx context.Context, pfce models.PostFlairCustomEmoji) (models.PostFlairCustomEmoji, error)

	PostFlairDescriptions(ctx context.Context) ([]models.PostFlairDescription, error)
	LinkPostFlairDescription(ctx context.Context, pfd models.PostFlairDescription) (models.PostFlairDescription, error)
}

type Repo struct {
	db *bun.DB
}

func NewRepo(db *bun.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) UserTrophies(ctx context.Context) ([]models.UserTrophy, error) {
	var user_trophies []models.UserTrophy

	query := `
                SELECT 
                    user_id,
                    trophy_id 
                FROM
                    user_trophies
            `

	_, err := r.db.NewRaw(query).Exec(ctx, &user_trophies)
	if err != nil {
		return []models.UserTrophy{}, err
	}
	return user_trophies, nil
}

func (r *Repo) LinkUserTrophy(ctx context.Context, ut models.UserTrophy) (models.UserTrophy, error) {
	query := `
                INSERT INTO user_trophies
                     (
                        user_id,
                        trophy_id
                    )
                VALUES (
                    ?,
                    ?
                )
                RETURNING *
            `

	if _, err := r.db.NewRaw(query, ut.UserID, ut.TrophyID).Exec(ctx, &ut); err != nil {
		var pgdriverErr pgdriver.Error
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgUniqueViolation {
			return models.UserTrophy{}, ErrDuplicateID
		}
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgConstraintViolation {
			return models.UserTrophy{}, ErrParentTableRecordNotFound
		}
		return models.UserTrophy{}, err
	}

	return ut, nil
}

func (r *Repo) VoxsphereMembers(ctx context.Context) ([]models.VoxsphereMember, error) {
	var voxsphereMembers []models.VoxsphereMember

	query := `
                SELECT 
                    voxsphere_id,
                    user_id 
                FROM
                    voxsphere_members
            `

	_, err := r.db.NewRaw(query).Exec(ctx, &voxsphereMembers)
	if err != nil {
		return []models.VoxsphereMember{}, err
	}
	return voxsphereMembers, nil
}

func (r *Repo) LinkVoxsphereMember(ctx context.Context, vme models.VoxsphereMember) (models.VoxsphereMember, error) {
	query := `
                INSERT INTO voxsphere_members
                     (
                        voxsphere_id,
                        user_id
                    )
                VALUES (
                    ?,
                    ?
                )
                RETURNING *
            `

	if _, err := r.db.NewRaw(query, vme.VoxsphereID, vme.UserID).Exec(ctx, &vme); err != nil {
		var pgdriverErr pgdriver.Error
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgUniqueViolation {
			return models.VoxsphereMember{}, ErrDuplicateID
		}
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgConstraintViolation {
			return models.VoxsphereMember{}, ErrParentTableRecordNotFound
		}
		return models.VoxsphereMember{}, err
	}

	return vme, nil
}

func (r *Repo) VoxsphereModerators(ctx context.Context) ([]models.VoxsphereModerator, error) {
	var voxsphereModerators []models.VoxsphereModerator

	query := `
                SELECT 
                    voxsphere_id,
                    user_id
                FROM
                    voxsphere_moderators
            `

	_, err := r.db.NewRaw(query).Exec(ctx, &voxsphereModerators)
	if err != nil {
		return []models.VoxsphereModerator{}, err
	}
	return voxsphereModerators, nil
}

func (r *Repo) LinkVoxsphereModerator(ctx context.Context, vmod models.VoxsphereModerator) (models.VoxsphereModerator, error) {
	query := `
                INSERT INTO voxsphere_moderators
                     (
                        voxsphere_id,
                        user_id
                    )
                VALUES (
                    ?,
                    ?
                )
                RETURNING *
            `

	if _, err := r.db.NewRaw(query, vmod.VoxsphereID, vmod.UserID).Exec(ctx, &vmod); err != nil {
		var pgdriverErr pgdriver.Error
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgUniqueViolation {
			return models.VoxsphereModerator{}, ErrDuplicateID
		}
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgConstraintViolation {
			return models.VoxsphereModerator{}, ErrParentTableRecordNotFound
		}
		return models.VoxsphereModerator{}, err
	}

	return vmod, nil
}

func (r *Repo) UserFlairEmojis(ctx context.Context) ([]models.UserFlairEmoji, error) {
	var UserFlairEmojis []models.UserFlairEmoji

	query := `
                SELECT 
                    emoji_id,
                    user_flair_id,
                    order_index
                FROM
                    user_flair_emojis
            `

	_, err := r.db.NewRaw(query).Exec(ctx, &UserFlairEmojis)
	if err != nil {
		return []models.UserFlairEmoji{}, err
	}
	return UserFlairEmojis, nil
}

func (r *Repo) LinkUserFlairEmoji(ctx context.Context, ufe models.UserFlairEmoji) (models.UserFlairEmoji, error) {
	query := `
                INSERT INTO user_flair_emojis
                     (
                        emoji_id,
                        user_flair_id,
                        order_index
                    )
                VALUES (
                    ?,
                    ?,
                    ?
                )
                RETURNING *
            `

	if _, err := r.db.NewRaw(query, ufe.EmojiID, ufe.UserFlairID, ufe.OrderIndex).Exec(ctx); err != nil {
		var pgdriverErr pgdriver.Error
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgUniqueViolation {
			return models.UserFlairEmoji{}, ErrDuplicateID
		}
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgConstraintViolation {
			return models.UserFlairEmoji{}, ErrParentTableRecordNotFound
		}
		return models.UserFlairEmoji{}, err
	}

	return ufe, nil
}

func (r *Repo) UserFlairCustomEmojis(ctx context.Context) ([]models.UserFlairCustomEmoji, error) {
	var UserFlairCustomEmojis []models.UserFlairCustomEmoji

	query := `
                SELECT 
                    custom_emoji_id,
                    user_flair_id,
                    order_index
                FROM
                    user_flair_custom_emojis
            `

	_, err := r.db.NewRaw(query).Exec(ctx, &UserFlairCustomEmojis)
	if err != nil {
		return []models.UserFlairCustomEmoji{}, err
	}
	return UserFlairCustomEmojis, nil
}

func (r *Repo) LinkUserFlairCustomEmoji(ctx context.Context, ufce models.UserFlairCustomEmoji) (models.UserFlairCustomEmoji, error) {
	query := `
                INSERT INTO user_flair_custom_emojis
                     (
                        custom_emoji_id,
                        user_flair_id,
                        order_index
                    )
                VALUES (
                    ?,
                    ?,
                    ?
                )
                RETURNING *
            `

	if _, err := r.db.NewRaw(query, ufce.CustomEmojiID, ufce.UserFlairID, ufce.OrderIndex).Exec(ctx, &ufce); err != nil {
		var pgdriverErr pgdriver.Error
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgUniqueViolation {
			return models.UserFlairCustomEmoji{}, ErrDuplicateID
		}
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgConstraintViolation {
			return models.UserFlairCustomEmoji{}, ErrParentTableRecordNotFound
		}
		return models.UserFlairCustomEmoji{}, err
	}

	return ufce, nil
}

func (r *Repo) UserFlairDescriptions(ctx context.Context) ([]models.UserFlairDescription, error) {
	var UserFlairDescriptions []models.UserFlairDescription

	query := `
                SELECT
                    user_flair_id,
                    order_index,
                    description
                FROM
                    user_flair_descriptions
            `

	_, err := r.db.NewRaw(query).Exec(ctx, &UserFlairDescriptions)
	if err != nil {
		return []models.UserFlairDescription{}, err
	}
	return UserFlairDescriptions, nil
}

func (r *Repo) LinkUserFlairDescription(ctx context.Context, ufd models.UserFlairDescription) (models.UserFlairDescription, error) {
	query := `
                INSERT INTO user_flair_descriptions
                     (
                        user_flair_id,
                        order_index,
                        description
                    )
                VALUES (
                    ?,
                    ?,
                    ?
                )
                RETURNING *
            `

	if _, err := r.db.NewRaw(query, ufd.UserFlairID, ufd.OrderIndex, ufd.Description).Exec(ctx, &ufd); err != nil {
		var pgdriverErr pgdriver.Error
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgUniqueViolation {
			return models.UserFlairDescription{}, ErrDuplicateID
		}
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgConstraintViolation {
			return models.UserFlairDescription{}, ErrParentTableRecordNotFound
		}
		return models.UserFlairDescription{}, err
	}

	return ufd, nil
}

func (r *Repo) PostFlairEmojis(ctx context.Context) ([]models.PostFlairEmoji, error) {
	var PostFlairEmojis []models.PostFlairEmoji

	query := `
                SELECT 
                    emoji_id,
                    post_flair_id,
                    order_index
                FROM
                    post_flair_emojis
            `

	_, err := r.db.NewRaw(query).Exec(ctx, &PostFlairEmojis)
	if err != nil {
		return []models.PostFlairEmoji{}, err
	}
	return PostFlairEmojis, nil
}

func (r *Repo) LinkPostFlairEmoji(ctx context.Context, pfe models.PostFlairEmoji) (models.PostFlairEmoji, error) {
	query := `
                INSERT INTO post_flair_emojis
                     (
                        emoji_id,
                        post_flair_id,
                        order_index
                    )
                VALUES (
                    ?,
                    ?,
                    ?
                )
                RETURNING *
            `

	if _, err := r.db.NewRaw(query, pfe.EmojiID, pfe.PostFlairID, pfe.OrderIndex).Exec(ctx); err != nil {
		var pgdriverErr pgdriver.Error
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgUniqueViolation {
			return models.PostFlairEmoji{}, ErrDuplicateID
		}
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgConstraintViolation {
			return models.PostFlairEmoji{}, ErrParentTableRecordNotFound
		}
		return models.PostFlairEmoji{}, err
	}

	return pfe, nil
}

func (r *Repo) PostFlairCustomEmojis(ctx context.Context) ([]models.PostFlairCustomEmoji, error) {
	var PostFlairCustomEmojis []models.PostFlairCustomEmoji

	query := `
                SELECT 
                    custom_emoji_id,
                    post_flair_id,
                    order_index
                FROM
                    post_flair_custom_emojis
            `

	_, err := r.db.NewRaw(query).Exec(ctx, &PostFlairCustomEmojis)
	if err != nil {
		return []models.PostFlairCustomEmoji{}, err
	}
	return PostFlairCustomEmojis, nil
}

func (r *Repo) LinkPostFlairCustomEmoji(ctx context.Context, pfce models.PostFlairCustomEmoji) (models.PostFlairCustomEmoji, error) {
	query := `
                INSERT INTO post_flair_custom_emojis
                     (
                        custom_emoji_id,
                        post_flair_id,
                        order_index
                    )
                VALUES (
                    ?,
                    ?,
                    ?
                )
                RETURNING *
            `

	if _, err := r.db.NewRaw(query, pfce.CustomEmojiID, pfce.PostFlairID, pfce.OrderIndex).Exec(ctx, &pfce); err != nil {
		var pgdriverErr pgdriver.Error
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgUniqueViolation {
			return models.PostFlairCustomEmoji{}, ErrDuplicateID
		}
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgConstraintViolation {
			return models.PostFlairCustomEmoji{}, ErrParentTableRecordNotFound
		}
		return models.PostFlairCustomEmoji{}, err
	}

	return pfce, nil
}

func (r *Repo) PostFlairDescriptions(ctx context.Context) ([]models.PostFlairDescription, error) {
	var PostFlairDescriptions []models.PostFlairDescription

	query := `
                SELECT
                    post_flair_id,
                    order_index,
                    description
                FROM
                    post_flair_descriptions
            `

	_, err := r.db.NewRaw(query).Exec(ctx, &PostFlairDescriptions)
	if err != nil {
		return []models.PostFlairDescription{}, err
	}
	return PostFlairDescriptions, nil
}

func (r *Repo) LinkPostFlairDescription(ctx context.Context, pfd models.PostFlairDescription) (models.PostFlairDescription, error) {
	query := `
                INSERT INTO post_flair_descriptions
                     (
                        post_flair_id,
                        order_index,
                        description
                    )
                VALUES (
                    ?,
                    ?,
                    ?
                )
                RETURNING *
            `

	if _, err := r.db.NewRaw(query, pfd.PostFlairID, pfd.OrderIndex, pfd.Description).Exec(ctx, &pfd); err != nil {
		var pgdriverErr pgdriver.Error
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgUniqueViolation {
			return models.PostFlairDescription{}, ErrDuplicateID
		}
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgConstraintViolation {
			return models.PostFlairDescription{}, ErrParentTableRecordNotFound
		}
		return models.PostFlairDescription{}, err
	}

	return pfd, nil
}
