package relation

import (
	"context"
	"errors"
	"strings"

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

type RelationRepository interface {
	UserTrophies(context.Context) ([]models.UserTrophy, error)
	LinkUserTrophies(context.Context, ...models.UserTrophy) ([]models.UserTrophy, error)

	VoxsphereMembers(context.Context) ([]models.VoxsphereMember, error)
	LinkVoxsphereMembers(context.Context, ...models.VoxsphereMember) ([]models.VoxsphereMember, error)

	VoxsphereModerators(context.Context) ([]models.VoxsphereModerator, error)
	LinkVoxsphereModerators(context.Context, ...models.VoxsphereModerator) ([]models.VoxsphereModerator, error)

	PostPostFlairs(context.Context) ([]models.PostPostFlair, error)
	LinkPostPostFlairs(context.Context, ...models.PostPostFlair) ([]models.PostPostFlair, error)

	UserUserFlair(context.Context) ([]models.UserUserFlair, error)
	LinkUserUserFlair(context.Context, ...models.UserUserFlair) ([]models.UserUserFlair, error)

	UserFlairEmojis(context.Context) ([]models.UserFlairEmoji, error)
	LinkUserFlairEmojis(context.Context, ...models.UserFlairEmoji) ([]models.UserFlairEmoji, error)

	UserFlairCustomEmojis(context.Context) ([]models.UserFlairCustomEmoji, error)
	LinkUserFlairCustomEmojis(context.Context, ...models.UserFlairCustomEmoji) ([]models.UserFlairCustomEmoji, error)

	UserFlairDescriptions(context.Context) ([]models.UserFlairDescription, error)
	LinkUserFlairDescriptions(context.Context, ...models.UserFlairDescription) ([]models.UserFlairDescription, error)

	PostFlairEmojis(context.Context) ([]models.PostFlairEmoji, error)
	LinkPostFlairEmojis(context.Context, ...models.PostFlairEmoji) ([]models.PostFlairEmoji, error)

	PostFlairCustomEmojis(context.Context) ([]models.PostFlairCustomEmoji, error)
	LinkPostFlairCustomEmojis(context.Context, ...models.PostFlairCustomEmoji) ([]models.PostFlairCustomEmoji, error)

	PostFlairDescriptions(context.Context) ([]models.PostFlairDescription, error)
	LinkPostFlairDescriptions(context.Context, ...models.PostFlairDescription) ([]models.PostFlairDescription, error)

	PostAwards(context.Context) ([]models.PostAward, error)
	LinkPostAwards(context.Context, ...models.PostAward) ([]models.PostAward, error)
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

func (r *Repo) LinkUserTrophies(ctx context.Context, uts ...models.UserTrophy) ([]models.UserTrophy, error) {
	query := `
        INSERT INTO user_trophies
             (
                user_id,
                trophy_id
            )
        VALUES 
    `
	args := make([]interface{}, 0)
	placeholders := make([]string, 0)

	for _, ut := range uts {
		placeholders = append(placeholders, "(?, ?)")

		args = append(args,
			ut.UserID,
			ut.TrophyID,
		)
	}

	query += strings.Join(placeholders, ", ")
	query += " RETURNING *"

	if _, err := r.db.NewRaw(query, args...).Exec(ctx, &uts); err != nil {
		var pgdriverErr pgdriver.Error
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgUniqueViolation {
			return nil, ErrDuplicateID
		}
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgConstraintViolation {
			return nil, ErrParentTableRecordNotFound
		}
		return nil, err
	}

	return uts, nil
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

func (r *Repo) LinkVoxsphereMembers(ctx context.Context, vmes ...models.VoxsphereMember) ([]models.VoxsphereMember, error) {
	query := `
        INSERT INTO voxsphere_members
             (
                voxsphere_id,
                user_id
            )
        VALUES 
    `
	args := make([]interface{}, 0)
	placeholders := make([]string, 0)

	for _, vme := range vmes {
		placeholders = append(placeholders, "(?, ?)")

		args = append(args,
			vme.VoxsphereID,
			vme.UserID,
		)
	}

	query += strings.Join(placeholders, ", ")
	query += " RETURNING *"

	if _, err := r.db.NewRaw(query, args...).Exec(ctx, &vmes); err != nil {
		var pgdriverErr pgdriver.Error
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgUniqueViolation {
			return nil, ErrDuplicateID
		}
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgConstraintViolation {
			return nil, ErrParentTableRecordNotFound
		}
		return nil, err
	}

	return vmes, nil
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

func (r *Repo) LinkVoxsphereModerators(ctx context.Context, vmods ...models.VoxsphereModerator) ([]models.VoxsphereModerator, error) {
	query := `
        INSERT INTO voxsphere_moderators
             (
                voxsphere_id,
                user_id
            )
        VALUES 
    `
	args := make([]interface{}, 0)
	placeholders := make([]string, 0)

	for _, vmod := range vmods {
		placeholders = append(placeholders, "(?, ?)")

		args = append(args,
			vmod.VoxsphereID,
			vmod.UserID,
		)
	}

	query += strings.Join(placeholders, ", ")
	query += " RETURNING *"

	if _, err := r.db.NewRaw(query, args...).Exec(ctx, &vmods); err != nil {
		var pgdriverErr pgdriver.Error
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgUniqueViolation {
			return nil, ErrDuplicateID
		}
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgConstraintViolation {
			return nil, ErrParentTableRecordNotFound
		}
		return nil, err
	}

	return vmods, nil
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

func (r *Repo) LinkUserFlairEmojis(ctx context.Context, ufes ...models.UserFlairEmoji) ([]models.UserFlairEmoji, error) {
	query := `
        INSERT INTO user_flair_emojis
             (
                emoji_id,
                user_flair_id,
                order_index
            )
        VALUES 
    `
	args := make([]interface{}, 0)
	placeholders := make([]string, 0)

	for _, ufe := range ufes {
		placeholders = append(placeholders, "(?, ?, ?)")

		args = append(args,
			ufe.EmojiID,
			ufe.UserFlairID,
			ufe.OrderIndex,
		)
	}

	query += strings.Join(placeholders, ", ")
	query += " RETURNING *"

	if _, err := r.db.NewRaw(query, args...).Exec(ctx, &ufes); err != nil {
		var pgdriverErr pgdriver.Error
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgUniqueViolation {
			return nil, ErrDuplicateID
		}
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgConstraintViolation {
			return nil, ErrParentTableRecordNotFound
		}
		return nil, err
	}

	return ufes, nil
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

func (r *Repo) LinkUserFlairCustomEmojis(ctx context.Context, ufces ...models.UserFlairCustomEmoji) ([]models.UserFlairCustomEmoji, error) {
	query := `
        INSERT INTO user_flair_custom_emojis
             (
                custom_emoji_id,
                user_flair_id,
                order_index
            )
        VALUES 
    `
	args := make([]interface{}, 0)
	placeholders := make([]string, 0)

	for _, ufce := range ufces {
		placeholders = append(placeholders, "(?, ?, ?)")

		args = append(args,
			ufce.CustomEmojiID,
			ufce.UserFlairID,
			ufce.OrderIndex,
		)
	}

	query += strings.Join(placeholders, ", ")
	query += " RETURNING *"

	if _, err := r.db.NewRaw(query, args...).Exec(ctx, &ufces); err != nil {
		var pgdriverErr pgdriver.Error
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgUniqueViolation {
			return nil, ErrDuplicateID
		}
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgConstraintViolation {
			return nil, ErrParentTableRecordNotFound
		}
		return nil, err
	}

	return ufces, nil
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

func (r *Repo) LinkUserFlairDescriptions(ctx context.Context, ufds ...models.UserFlairDescription) ([]models.UserFlairDescription, error) {
	query := `
        INSERT INTO user_flair_descriptions
             (
                user_flair_id,
                order_index,
                description
            )
        VALUES 
    `

	args := make([]interface{}, 0)
	placeholders := make([]string, 0)

	for _, ufd := range ufds {
		placeholders = append(placeholders, "(?, ?, ?)")

		args = append(args,
			ufd.UserFlairID,
			ufd.OrderIndex,
			ufd.Description,
		)
	}

	query += strings.Join(placeholders, ", ")
	query += " RETURNING *"

	if _, err := r.db.NewRaw(query, args...).Exec(ctx, &ufds); err != nil {
		var pgdriverErr pgdriver.Error
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgUniqueViolation {
			return nil, ErrDuplicateID
		}
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgConstraintViolation {
			return nil, ErrParentTableRecordNotFound
		}
		return nil, err
	}

	return ufds, nil
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

func (r *Repo) LinkPostFlairEmojis(ctx context.Context, pfes ...models.PostFlairEmoji) ([]models.PostFlairEmoji, error) {
	query := `
        INSERT INTO post_flair_emojis
             (
                emoji_id,
                post_flair_id,
                order_index
            )
        VALUES 
    `

	args := make([]interface{}, 0)
	placeholders := make([]string, 0)
	for _, pfe := range pfes {
		placeholders = append(placeholders, "(?, ?, ?)")
		args = append(args, pfe.EmojiID, pfe.PostFlairID, pfe.OrderIndex)
	}
	query += strings.Join(placeholders, ", ") + " RETURNING *"

	if _, err := r.db.NewRaw(query, args...).Exec(ctx, &pfes); err != nil {
		var pgdriverErr pgdriver.Error
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgUniqueViolation {
			return nil, ErrDuplicateID
		}
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgConstraintViolation {
			return nil, ErrParentTableRecordNotFound
		}
		return nil, err
	}

	return pfes, nil
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

func (r *Repo) LinkPostFlairCustomEmojis(ctx context.Context, pfces ...models.PostFlairCustomEmoji) ([]models.PostFlairCustomEmoji, error) {
	query := `
        INSERT INTO post_flair_custom_emojis
             (
                custom_emoji_id,
                post_flair_id,
                order_index
            )
        VALUES 
    `

	args := make([]interface{}, 0)
	placeholders := make([]string, 0)
	for _, pfce := range pfces {
		placeholders = append(placeholders, "(?, ?, ?)")
		args = append(args, pfce.CustomEmojiID, pfce.PostFlairID, pfce.OrderIndex)
	}
	query += strings.Join(placeholders, ", ") + " RETURNING *"

	if _, err := r.db.NewRaw(query, args...).Exec(ctx, &pfces); err != nil {
		var pgdriverErr pgdriver.Error
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgUniqueViolation {
			return nil, ErrDuplicateID
		}
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgConstraintViolation {
			return nil, ErrParentTableRecordNotFound
		}
		return nil, err
	}

	return pfces, nil
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

func (r *Repo) LinkPostFlairDescriptions(ctx context.Context, pfds ...models.PostFlairDescription) ([]models.PostFlairDescription, error) {
	query := `
        INSERT INTO post_flair_descriptions
             (
                post_flair_id,
                order_index,
                description
            )
        VALUES 
    `

	args := make([]interface{}, 0)
	placeholders := make([]string, 0)
	for _, pfd := range pfds {
		placeholders = append(placeholders, "(?, ?, ?)")
		args = append(args, pfd.PostFlairID, pfd.OrderIndex, pfd.Description)
	}
	query += strings.Join(placeholders, ", ") + " RETURNING *"

	if _, err := r.db.NewRaw(query, args...).Exec(ctx, &pfds); err != nil {
		var pgdriverErr pgdriver.Error
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgUniqueViolation {
			return nil, ErrDuplicateID
		}
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgConstraintViolation {
			return nil, ErrParentTableRecordNotFound
		}
		return nil, err
	}

	return pfds, nil
}

func (r *Repo) PostAwards(ctx context.Context) ([]models.PostAward, error) {
	var post_awards []models.PostAward

	query := `
                SELECT 
                    post_id,
                    award_id 
                FROM
                    post_awards
            `

	_, err := r.db.NewRaw(query).Exec(ctx, &post_awards)
	if err != nil {
		return []models.PostAward{}, err
	}
	return post_awards, nil
}

func (r *Repo) LinkPostAwards(ctx context.Context, pas ...models.PostAward) ([]models.PostAward, error) {
	query := `
        INSERT INTO post_awards
             (
                post_id,
                award_id
            )
        VALUES 
    `

	args := make([]interface{}, 0)
	placeholders := make([]string, 0)
	for _, pa := range pas {
		placeholders = append(placeholders, "(?, ?)")
		args = append(args, pa.PostID, pa.AwardID)
	}
	query += strings.Join(placeholders, ", ") + " RETURNING *"

	if _, err := r.db.NewRaw(query, args...).Exec(ctx, &pas); err != nil {
		var pgdriverErr pgdriver.Error
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgUniqueViolation {
			return nil, ErrDuplicateID
		}
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgConstraintViolation {
			return nil, ErrParentTableRecordNotFound
		}
		return nil, err
	}

	return pas, nil
}

func (r *Repo) PostPostFlairs(ctx context.Context) ([]models.PostPostFlair, error) {
	var postPostFlairs []models.PostPostFlair

	query := `
		SELECT 
			post_id,
			post_flair_id
		FROM
			post_post_flairs
	`

	_, err := r.db.NewRaw(query).Exec(ctx, &postPostFlairs)
	if err != nil {
		return []models.PostPostFlair{}, err
	}
	return postPostFlairs, nil
}

func (r *Repo) LinkPostPostFlairs(ctx context.Context, ppfs ...models.PostPostFlair) ([]models.PostPostFlair, error) {
	query := `
		INSERT INTO post_post_flairs
			(post_id, post_flair_id)
		VALUES 
	`

	args := make([]interface{}, 0)
	placeholders := make([]string, 0)
	for _, ppf := range ppfs {
		placeholders = append(placeholders, "(?, ?)")
		args = append(args, ppf.PostID, ppf.PostFlairID)
	}
	query += strings.Join(placeholders, ", ") + " RETURNING *"

	if _, err := r.db.NewRaw(query, args...).Exec(ctx, &ppfs); err != nil {
		var pgdriverErr pgdriver.Error
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgUniqueViolation {
			return nil, ErrDuplicateID
		}
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgConstraintViolation {
			return nil, ErrParentTableRecordNotFound
		}
		return nil, err
	}

	return ppfs, nil
}

func (r *Repo) UserUserFlairs(ctx context.Context) ([]models.UserUserFlair, error) {
	var userUserFlairs []models.UserUserFlair

	query := `
		SELECT 
			user_id,
			user_flair_id
		FROM
			user_user_flairs
	`

	_, err := r.db.NewRaw(query).Exec(ctx, &userUserFlairs)
	if err != nil {
		return []models.UserUserFlair{}, err
	}
	return userUserFlairs, nil
}

func (r *Repo) LinkUserUserFlairs(ctx context.Context, uufs ...models.UserUserFlair) ([]models.UserUserFlair, error) {
	query := `
		INSERT INTO user_user_flairs
			(user_id, user_flair_id)
		VALUES 
	`

	args := make([]interface{}, 0)
	placeholders := make([]string, 0)
	for _, uuf := range uufs {
		placeholders = append(placeholders, "(?, ?)")
		args = append(args, uuf.UserID, uuf.UserFlairID)
	}
	query += strings.Join(placeholders, ", ") + " RETURNING *"

	if _, err := r.db.NewRaw(query, args...).Exec(ctx, &uufs); err != nil {
		var pgdriverErr pgdriver.Error
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgUniqueViolation {
			return nil, ErrDuplicateID
		}
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgConstraintViolation {
			return nil, ErrParentTableRecordNotFound
		}
		return nil, err
	}

	return uufs, nil
}
