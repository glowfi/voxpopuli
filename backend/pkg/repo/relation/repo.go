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
            `

	if _, err := r.db.NewRaw(query, ut.UserID, ut.TrophyID).Exec(ctx); err != nil {
		var pgdriverErr pgdriver.Error
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgUniqueViolation {
			return models.UserTrophy{}, ErrDuplicateID
		}
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgConstraintViolation {
			return models.UserTrophy{}, ErrParentTableRecordNotFound
		}
		return models.UserTrophy{}, err
	}

	return models.UserTrophy{
		UserID:   ut.UserID,
		TrophyID: ut.TrophyID,
	}, nil
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
            `

	if _, err := r.db.NewRaw(query, vme.VoxsphereID, vme.UserID).Exec(ctx); err != nil {
		var pgdriverErr pgdriver.Error
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgUniqueViolation {
			return models.VoxsphereMember{}, ErrDuplicateID
		}
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgConstraintViolation {
			return models.VoxsphereMember{}, ErrParentTableRecordNotFound
		}
		return models.VoxsphereMember{}, err
	}

	return models.VoxsphereMember{
		VoxsphereID: vme.VoxsphereID,
		UserID:      vme.UserID,
	}, nil
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
            `

	if _, err := r.db.NewRaw(query, vmod.VoxsphereID, vmod.UserID).Exec(ctx); err != nil {
		var pgdriverErr pgdriver.Error
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgUniqueViolation {
			return models.VoxsphereModerator{}, ErrDuplicateID
		}
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgConstraintViolation {
			return models.VoxsphereModerator{}, ErrParentTableRecordNotFound
		}
		return models.VoxsphereModerator{}, err
	}

	return models.VoxsphereModerator{
		VoxsphereID: vmod.VoxsphereID,
		UserID:      vmod.UserID,
	}, nil
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

	return models.UserFlairEmoji{
		EmojiID:     ufe.EmojiID,
		UserFlairID: ufe.UserFlairID,
		OrderIndex:  ufe.OrderIndex,
	}, nil
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
            `

	if _, err := r.db.NewRaw(query, ufce.CustomEmojiID, ufce.UserFlairID, ufce.OrderIndex).Exec(ctx); err != nil {
		var pgdriverErr pgdriver.Error
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgUniqueViolation {
			return models.UserFlairCustomEmoji{}, ErrDuplicateID
		}
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgConstraintViolation {
			return models.UserFlairCustomEmoji{}, ErrParentTableRecordNotFound
		}
		return models.UserFlairCustomEmoji{}, err
	}

	return models.UserFlairCustomEmoji{
		CustomEmojiID: ufce.CustomEmojiID,
		UserFlairID:   ufce.UserFlairID,
		OrderIndex:    ufce.OrderIndex,
	}, nil
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
            `

	if _, err := r.db.NewRaw(query, ufd.UserFlairID, ufd.OrderIndex, ufd.Description).Exec(ctx); err != nil {
		var pgdriverErr pgdriver.Error
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgUniqueViolation {
			return models.UserFlairDescription{}, ErrDuplicateID
		}
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgConstraintViolation {
			return models.UserFlairDescription{}, ErrParentTableRecordNotFound
		}
		return models.UserFlairDescription{}, err
	}

	return models.UserFlairDescription{
		UserFlairID: ufd.UserFlairID,
		OrderIndex:  ufd.OrderIndex,
		Description: ufd.Description,
	}, nil
}
