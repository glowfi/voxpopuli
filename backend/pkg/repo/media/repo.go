package media

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
	pgUniqueViolation     = "23505"
	pgConstraintViolation = "23503"
)

var (
	ErrNotFound                  = errors.New("not found")
	ErrDuplicateID               = errors.New("duplicate id")
	ErrParentTableRecordNotFound = errors.New("record does not exist in the parent table")
)

type Repository interface {
	// post media
	PostMedias(context.Context) ([]models.PostMedia, error)
	PostMediaByID(context.Context, uuid.UUID) (models.PostMedia, error)
	AddPostMedia(context.Context, models.PostMedia) (models.PostMedia, error)
	UpdatePostMedia(context.Context, models.PostMedia) (models.PostMedia, error)
	DeletePostMedia(context.Context, uuid.UUID) error

	// image
	Images(context.Context) ([]models.Image, error)
	ImageByID(context.Context, uuid.UUID) (models.Image, error)
	AddImage(context.Context, models.Image) (models.Image, error)
	UpdateImage(context.Context, models.Image) (models.Image, error)
	DeleteImage(context.Context, uuid.UUID) error

	// image metadata
	ImageMetadatas(context.Context) ([]models.ImageMetadata, error)
	ImageMetadataByID(context.Context, uuid.UUID) (models.ImageMetadata, error)
	AddImageMetadata(context.Context, models.ImageMetadata) (models.ImageMetadata, error)
	UpdateImageMetadata(context.Context, models.ImageMetadata) (models.ImageMetadata, error)
	DeleteImageMetadata(context.Context, uuid.UUID) error

	// gif
	Gifs(context.Context) ([]models.Gif, error)
	GifByID(context.Context, uuid.UUID) (models.Gif, error)
	AddGif(context.Context, models.Gif) (models.Gif, error)
	UpdateGif(context.Context, models.Gif) (models.Gif, error)
	DeleteGif(context.Context, uuid.UUID) error

	// gif metadata
	GifMetadatas(context.Context) ([]models.GifMetadata, error)
	GifMetadataByID(context.Context, uuid.UUID) (models.GifMetadata, error)
	AddGifMetadata(context.Context, models.GifMetadata) (models.GifMetadata, error)
	UpdateGifMetadata(context.Context, models.GifMetadata) (models.GifMetadata, error)
	DeleteGifMetadata(context.Context, uuid.UUID) error

	// gallery
	Galleries(context.Context) ([]models.Gallery, error)
	GalleryByID(context.Context, uuid.UUID) (models.Gallery, error)
	AddGallery(context.Context, models.Gallery) (models.Gallery, error)
	UpdateGallery(context.Context, models.Gallery) (models.Gallery, error)
	DeleteGallery(context.Context, uuid.UUID) error

	// gallery metadata
	GalleryMetadatas(context.Context) ([]models.GalleryMetadata, error)
	GalleryMetadataByID(context.Context, uuid.UUID) (models.GalleryMetadata, error)
	AddGalleryMetadata(context.Context, models.GalleryMetadata) (models.GalleryMetadata, error)
	UpdateGalleryMetadata(context.Context, models.GalleryMetadata) (models.GalleryMetadata, error)
	DeleteGalleryMetadata(context.Context, uuid.UUID) error

	// video
	Videos(context.Context) ([]models.Video, error)
	VideoByID(context.Context, uuid.UUID) (models.Video, error)
	AddVideo(context.Context, models.Video) (models.Video, error)
	UpdateVideo(context.Context, models.Video) (models.Video, error)
	DeleteVideo(context.Context, uuid.UUID) error

	// link
	Links(context.Context) ([]models.Link, error)
	LinkByID(context.Context, uuid.UUID) (models.Link, error)
	AddLink(context.Context, models.Link) (models.Link, error)
	UpdateLink(context.Context, models.Link) (models.Link, error)
	DeleteLink(context.Context, uuid.UUID) error
}

type Repo struct {
	db *bun.DB
}

func NewRepo(db *bun.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) PostMedias(ctx context.Context) ([]models.PostMedia, error) {
	var postMedias []models.PostMedia

	query := `
        SELECT 
            id,
            post_id,
            media_type
        FROM 
            post_medias;
    `

	_, err := r.db.NewRaw(query).Exec(ctx, &postMedias)
	if err != nil {
		return []models.PostMedia{}, err
	}
	return postMedias, nil
}

func (r *Repo) PostMediaByID(ctx context.Context, ID uuid.UUID) (models.PostMedia, error) {
	var postMedia models.PostMedia

	query := `
        SELECT 
            id,
            post_id,
            media_type
        FROM 
            post_medias
        WHERE 
            id = ?;
    `

	_, err := r.db.NewRaw(query, ID).Exec(ctx, &postMedia)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.PostMedia{}, ErrNotFound
		}
		return models.PostMedia{}, err
	}
	return postMedia, nil
}

func (r *Repo) AddPostMedia(ctx context.Context, postMedia models.PostMedia) (models.PostMedia, error) {
	query := `
        INSERT INTO 
            post_medias (
                id,
                post_id,
                media_type
            )
        VALUES (
            ?,
            ?,
            ?
        )
        RETURNING *
    `

	if _, err := r.db.NewRaw(query,
		postMedia.ID,
		postMedia.PostID,
		postMedia.MediaType).Exec(ctx, &postMedia); err != nil {
		var pgdriverErr pgdriver.Error
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgUniqueViolation {
			return models.PostMedia{}, ErrDuplicateID
		}
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgConstraintViolation {
			return models.PostMedia{}, ErrParentTableRecordNotFound
		}
		return models.PostMedia{}, err
	}

	return postMedia, nil
}

func (r *Repo) UpdatePostMedia(ctx context.Context, postMedia models.PostMedia) (models.PostMedia, error) {
	query := `
        UPDATE 
            post_medias
        SET 
            post_id = ?,
            media_type = ?
        WHERE 
            id = ?
        RETURNING *
    `

	res, err := r.db.NewRaw(query,
		postMedia.PostID,
		postMedia.MediaType,
		postMedia.ID).Exec(ctx, &postMedia)
	if err != nil {
		var pgdriverErr pgdriver.Error
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgConstraintViolation {
			return models.PostMedia{}, ErrParentTableRecordNotFound
		}
		if errors.Is(err, sql.ErrNoRows) {
			return models.PostMedia{}, ErrNotFound
		}
		return models.PostMedia{}, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return models.PostMedia{}, err
	}
	if rowsAffected == 0 {
		return models.PostMedia{}, ErrNotFound
	}
	return postMedia, nil
}

func (r *Repo) DeletePostMedia(ctx context.Context, ID uuid.UUID) error {
	query := `
        DELETE FROM 
            post_medias
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
		return ErrNotFound
	}
	return nil
}

func (r *Repo) Images(ctx context.Context) ([]models.Image, error) {
	var images []models.Image

	query := `
            SELECT 
                i.id,
                i.media_id,
                (
                    SELECT 
                        JSON_AGG(imageMetadata)
                    FROM 
                        (
                            SELECT 
                                JSON_BUILD_OBJECT(
                                    'id', imageMetadata.id,
                                    'image_id', imageMetadata.image_id,
                                    'height', imageMetadata.height,
                                    'width', imageMetadata.width,
                                    'url', imageMetadata.url,
                                    'created_at', TO_CHAR(imageMetadata.created_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"'),
                                    'created_at_unix', imageMetadata.created_at_unix,
                                    'updated_at', TO_CHAR(imageMetadata.updated_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"')
                                ) AS imageMetadata
                            FROM 
                                image_metadatas imageMetadata
                            WHERE 
                                imageMetadata.image_id = i.id
                            ORDER BY 
                                imageMetadata.height
                        ) AS orderedMetadata
                ) AS image_metadata
            FROM 
                images i;
`

	_, err := r.db.NewRaw(query).Exec(ctx, &images)
	if err != nil {
		return []models.Image{}, err
	}
	return images, nil
}

func (r *Repo) ImageByID(ctx context.Context, ID uuid.UUID) (models.Image, error) {
	var image models.Image

	query := `
            SELECT 
                i.id,
                i.media_id,
                (
                    SELECT 
                        JSON_AGG(imageMetadata)
                    FROM 
                        (
                            SELECT 
                                JSON_BUILD_OBJECT(
                                    'id', imageMetadata.id,
                                    'image_id', imageMetadata.image_id,
                                    'height', imageMetadata.height,
                                    'width', imageMetadata.width,
                                    'url', imageMetadata.url,
                                    'created_at', TO_CHAR(imageMetadata.created_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"'),
                                    'created_at_unix', imageMetadata.created_at_unix,
                                    'updated_at', TO_CHAR(imageMetadata.updated_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"')
                                ) AS imageMetadata
                            FROM 
                                image_metadatas imageMetadata
                            WHERE 
                                imageMetadata.image_id = i.id
                            ORDER BY 
                                imageMetadata.height
                        ) AS orderedMetadata
                ) AS image_metadata
            FROM 
                images i
            WHERE
                i.id = ?;
    `

	if _, err := r.db.NewRaw(query, ID).Exec(ctx, &image); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Image{}, ErrNotFound
		}
		return models.Image{}, err
	}
	return image, nil
}

func (r *Repo) AddImage(ctx context.Context, image models.Image) (models.Image, error) {
	query := `
        INSERT INTO 
            images (
                id,
                media_id
            )
        VALUES (
            ?,
            ?
        )
        RETURNING *
    `

	if _, err := r.db.NewRaw(query,
		image.ID,
		image.MediaID).Exec(ctx, &image); err != nil {
		var pgdriverErr pgdriver.Error
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgUniqueViolation {
			return models.Image{}, ErrDuplicateID
		}
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgConstraintViolation {
			return models.Image{}, ErrParentTableRecordNotFound
		}
		return models.Image{}, err
	}

	return image, nil
}

func (r *Repo) UpdateImage(ctx context.Context, image models.Image) (models.Image, error) {
	query := `
        UPDATE 
            images
        SET 
            media_id = ?
        WHERE 
            id = ?
        RETURNING *
    `

	res, err := r.db.NewRaw(query,
		image.MediaID,
		image.ID).Exec(ctx, &image)
	if err != nil {
		var pgdriverErr pgdriver.Error
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgConstraintViolation {
			return models.Image{}, ErrParentTableRecordNotFound
		}
		if errors.Is(err, sql.ErrNoRows) {
			return models.Image{}, ErrNotFound
		}
		return models.Image{}, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return models.Image{}, err
	}
	if rowsAffected == 0 {
		return models.Image{}, ErrNotFound
	}
	return image, nil
}

func (r *Repo) DeleteImage(ctx context.Context, ID uuid.UUID) error {
	query := `
        DELETE FROM 
            images
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
		return ErrNotFound
	}
	return nil
}

func (r *Repo) ImageMetadatas(ctx context.Context) ([]models.ImageMetadata, error) {
	var imageMetadatas []models.ImageMetadata

	query := `
        SELECT
            id,
            image_id,
            height,
            width,
            url,
            created_at,
            created_at_unix,
            updated_at
        FROM
            image_metadatas;
    `

	_, err := r.db.NewRaw(query).Exec(ctx, &imageMetadatas)
	if err != nil {
		return []models.ImageMetadata{}, err
	}
	return imageMetadatas, nil
}

func (r *Repo) ImageMetadataByID(ctx context.Context, ID uuid.UUID) (models.ImageMetadata, error) {
	var imageMetadata models.ImageMetadata

	query := `
        SELECT
            id,
            image_id,
            height,
            width,
            url,
            created_at,
            created_at_unix,
            updated_at
        FROM
            image_metadatas
        WHERE
            id = ?;
    `

	_, err := r.db.NewRaw(query, ID).Exec(ctx, &imageMetadata)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.ImageMetadata{}, ErrNotFound
		}
		return models.ImageMetadata{}, err
	}
	return imageMetadata, nil
}

func (r *Repo) AddImageMetadata(ctx context.Context, imageMetadata models.ImageMetadata) (models.ImageMetadata, error) {
	query := `
        INSERT INTO
            image_metadatas (
                id,
                image_id,
                height,
                width,
                url,
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
            ?
        )
        RETURNING *
    `

	timestamp := time.Now()
	imageMetadata.CreatedAt = timestamp
	imageMetadata.CreatedAtUnix = timestamp.Unix()
	imageMetadata.UpdatedAt = timestamp

	if _, err := r.db.NewRaw(query,
		imageMetadata.ID,
		imageMetadata.ImageID,
		imageMetadata.Height,
		imageMetadata.Width,
		imageMetadata.Url,
		imageMetadata.CreatedAt,
		imageMetadata.CreatedAtUnix,
		imageMetadata.UpdatedAt).Exec(ctx, &imageMetadata); err != nil {
		var pgdriverErr pgdriver.Error
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgUniqueViolation {
			return models.ImageMetadata{}, ErrDuplicateID
		}
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgConstraintViolation {
			return models.ImageMetadata{}, ErrParentTableRecordNotFound
		}
		return models.ImageMetadata{}, err
	}
	return imageMetadata, nil
}

func (r *Repo) UpdateImageMetadata(ctx context.Context, imageMetadata models.ImageMetadata) (models.ImageMetadata, error) {
	query := `
        UPDATE
            image_metadatas
        SET
            image_id = ?,
            height = ?,
            width = ?,
            url = ?,
            updated_at = ?
        WHERE
            id = ?
        RETURNING *
    `

	imageMetadata.UpdatedAt = time.Now()

	res, err := r.db.NewRaw(query,
		imageMetadata.ImageID,
		imageMetadata.Height,
		imageMetadata.Width,
		imageMetadata.Url,
		imageMetadata.UpdatedAt,
		imageMetadata.ID,
	).Exec(ctx, &imageMetadata)
	if err != nil {
		var pgdriverErr pgdriver.Error
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgConstraintViolation {
			return models.ImageMetadata{}, ErrParentTableRecordNotFound
		}
		if errors.Is(err, sql.ErrNoRows) {
			return models.ImageMetadata{}, ErrNotFound
		}
		return models.ImageMetadata{}, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return models.ImageMetadata{}, err
	}
	if rowsAffected == 0 {
		return models.ImageMetadata{}, ErrNotFound
	}
	return imageMetadata, nil
}

func (r *Repo) DeleteImageMetadata(ctx context.Context, ID uuid.UUID) error {
	query := `
        DELETE FROM 
            image_metadatas
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
		return ErrNotFound
	}
	return nil
}

func (r *Repo) Gifs(ctx context.Context) ([]models.Gif, error) {
	var gifs []models.Gif

	query := `
    SELECT 
        g.id,
        g.media_id,
        (
            SELECT 
                JSON_AGG(gifMetadata)
            FROM 
                (
                    SELECT 
                        JSON_BUILD_OBJECT(
                            'id', gifMetadata.id,
                            'gif_id', gifMetadata.gif_id,
                            'height', gifMetadata.height,
                            'width', gifMetadata.width,
                            'url', gifMetadata.url,
                            'created_at', TO_CHAR(gifMetadata.created_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"'),
                            'created_at_unix', gifMetadata.created_at_unix,
                            'updated_at', TO_CHAR(gifMetadata.updated_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"')
                        ) AS gifMetadata
                    FROM 
                        gif_metadatas gifMetadata
                    WHERE 
                        gifMetadata.gif_id = g.id
                    ORDER BY 
                        gifMetadata.height
                ) AS orderedMetadata
        ) AS gif_metadata
    FROM 
        gifs g;
`

	_, err := r.db.NewRaw(query).Exec(ctx, &gifs)
	if err != nil {
		return []models.Gif{}, err
	}
	return gifs, nil
}

func (r *Repo) GifByID(ctx context.Context, ID uuid.UUID) (models.Gif, error) {
	var gif models.Gif

	query := `
        SELECT 
            g.id,
            g.media_id,
            (
                SELECT 
                    JSON_AGG(gifMetadata)
                FROM 
                    (
                        SELECT 
                            JSON_BUILD_OBJECT(
                                'id', gifMetadata.id,
                                'gif_id', gifMetadata.gif_id,
                                'height', gifMetadata.height,
                                'width', gifMetadata.width,
                                'url', gifMetadata.url,
                                'created_at', TO_CHAR(gifMetadata.created_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"'),
                                'created_at_unix', gifMetadata.created_at_unix,
                                'updated_at', TO_CHAR(gifMetadata.updated_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"')
                            ) AS gifMetadata
                        FROM 
                            gif_metadatas gifMetadata
                        WHERE 
                            gifMetadata.gif_id = g.id
                        ORDER BY 
                            gifMetadata.height
                    ) AS orderedMetadata
            ) AS gif_metadata
        FROM 
            gifs g
        WHERE 
            id = ?;
    `

	_, err := r.db.NewRaw(query, ID).Exec(ctx, &gif)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Gif{}, ErrNotFound
		}
		return models.Gif{}, err
	}
	return gif, nil
}

func (r *Repo) AddGif(ctx context.Context, gif models.Gif) (models.Gif, error) {
	query := `
        INSERT INTO 
            gifs (
                id,
                media_id
            )
        VALUES (
            ?,
            ?
        )
        RETURNING *
    `

	if _, err := r.db.NewRaw(query,
		gif.ID,
		gif.MediaID).Exec(ctx); err != nil {
		var pgdriverErr pgdriver.Error
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgUniqueViolation {
			return models.Gif{}, ErrDuplicateID
		}
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgConstraintViolation {
			return models.Gif{}, ErrParentTableRecordNotFound
		}
		return models.Gif{}, err
	}

	return gif, nil
}

func (r *Repo) UpdateGif(ctx context.Context, gif models.Gif) (models.Gif, error) {
	query := `
        UPDATE 
            gifs
        SET 
            media_id = ?
        WHERE 
            id = ?
        RETURNING *
    `

	res, err := r.db.NewRaw(query,
		gif.MediaID,
		gif.ID).Exec(ctx, &gif)
	if err != nil {
		var pgdriverErr pgdriver.Error
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgConstraintViolation {
			return models.Gif{}, ErrParentTableRecordNotFound
		}
		if errors.Is(err, sql.ErrNoRows) {
			return models.Gif{}, ErrNotFound
		}
		return models.Gif{}, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return models.Gif{}, err
	}
	if rowsAffected == 0 {
		return models.Gif{}, ErrNotFound
	}
	return gif, nil
}

func (r *Repo) DeleteGif(ctx context.Context, ID uuid.UUID) error {
	query := `
        DELETE FROM 
            gifs
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
		return ErrNotFound
	}
	return nil
}

func (r *Repo) GifMetadatas(ctx context.Context) ([]models.GifMetadata, error) {
	var gifMetadatas []models.GifMetadata

	query := `
        SELECT
            id,
            gif_id,
            height,
            width,
            url,
            created_at,
            created_at_unix,
            updated_at
        FROM
            gif_metadatas;
    `

	_, err := r.db.NewRaw(query).Exec(ctx, &gifMetadatas)
	if err != nil {
		return []models.GifMetadata{}, err
	}
	return gifMetadatas, nil
}

func (r *Repo) GifMetadataByID(ctx context.Context, ID uuid.UUID) (models.GifMetadata, error) {
	var gifMetadata models.GifMetadata

	query := `
        SELECT
            id,
            gif_id,
            height,
            width,
            url,
            created_at,
            created_at_unix,
            updated_at
        FROM
            gif_metadatas
        WHERE
            id = ?;
    `

	_, err := r.db.NewRaw(query, ID).Exec(ctx, &gifMetadata)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.GifMetadata{}, ErrNotFound
		}
		return models.GifMetadata{}, err
	}
	return gifMetadata, nil
}

func (r *Repo) AddGifMetadata(ctx context.Context, gifMetadata models.GifMetadata) (models.GifMetadata, error) {
	query := `
        INSERT INTO
            gif_metadatas (
                id,
                gif_id,
                height,
                width,
                url,
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
            ?
        )
        RETURNING *
    `

	timestamp := time.Now()
	gifMetadata.CreatedAt = timestamp
	gifMetadata.CreatedAtUnix = timestamp.Unix()
	gifMetadata.UpdatedAt = timestamp

	if _, err := r.db.NewRaw(query,
		gifMetadata.ID,
		gifMetadata.GifID,
		gifMetadata.Height,
		gifMetadata.Width,
		gifMetadata.Url,
		gifMetadata.CreatedAt,
		gifMetadata.CreatedAtUnix,
		gifMetadata.UpdatedAt).Exec(ctx); err != nil {
		var pgdriverErr pgdriver.Error
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgUniqueViolation {
			return models.GifMetadata{}, ErrDuplicateID
		}
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgConstraintViolation {
			return models.GifMetadata{}, ErrParentTableRecordNotFound
		}
		return models.GifMetadata{}, err
	}
	return gifMetadata, nil
}

func (r *Repo) UpdateGifMetadata(ctx context.Context, gifMetadata models.GifMetadata) (models.GifMetadata, error) {
	query := `
        UPDATE
            gif_metadatas
        SET
            gif_id = ?,
            height = ?,
            width = ?,
            url = ?,
            updated_at = ?
        WHERE
            id = ?
        RETURNING *
    `

	gifMetadata.UpdatedAt = time.Now()

	res, err := r.db.NewRaw(query,
		gifMetadata.GifID,
		gifMetadata.Height,
		gifMetadata.Width,
		gifMetadata.Url,
		gifMetadata.UpdatedAt,
		gifMetadata.ID,
	).Exec(ctx, &gifMetadata)
	if err != nil {
		var pgdriverErr pgdriver.Error
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgConstraintViolation {
			return models.GifMetadata{}, ErrParentTableRecordNotFound
		}
		if errors.Is(err, sql.ErrNoRows) {
			return models.GifMetadata{}, ErrNotFound
		}
		return models.GifMetadata{}, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return models.GifMetadata{}, err
	}
	if rowsAffected == 0 {
		return models.GifMetadata{}, ErrNotFound
	}
	return gifMetadata, nil
}

func (r *Repo) DeleteGifMetadata(ctx context.Context, ID uuid.UUID) error {
	query := `
        DELETE FROM 
            gif_metadatas
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
		return ErrNotFound
	}
	return nil
}

func (r *Repo) Galleries(ctx context.Context) ([]models.Gallery, error) {
	var galleries []models.Gallery

	query := `
        SELECT 
            g.id,
            g.media_id,
            (
                SELECT 
                    JSON_AGG(galleryMetadata)
                FROM 
                    (
                        SELECT 
                            JSON_BUILD_OBJECT(
                                'id', galleryMetadata.id,
                                'gallery_id', galleryMetadata.gallery_id,
                                'order_index', galleryMetadata.order_index,
                                'height', galleryMetadata.height,
                                'width', galleryMetadata.width,
                                'url', galleryMetadata.url,
                                'created_at', TO_CHAR(galleryMetadata.created_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"'),
                                'created_at_unix', galleryMetadata.created_at_unix,
                                'updated_at', TO_CHAR(galleryMetadata.updated_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"')
                            ) AS galleryMetadata
                        FROM 
                            gallery_metadatas galleryMetadata
                        WHERE 
                            galleryMetadata.gallery_id = g.id
                        ORDER BY 
                            galleryMetadata.order_index
                    ) AS orderedMetadata
            ) AS gallery_metadata
        FROM 
            galleries g;
    `

	_, err := r.db.NewRaw(query).Exec(ctx, &galleries)
	if err != nil {
		return []models.Gallery{}, err
	}
	return galleries, nil
}

func (r *Repo) GalleryByID(ctx context.Context, ID uuid.UUID) (models.Gallery, error) {
	var gallery models.Gallery

	query := `
        SELECT 
            g.id,
            g.media_id,
            (
                SELECT 
                    JSON_AGG(galleryMetadata)
                FROM 
                    (
                        SELECT 
                            JSON_BUILD_OBJECT(
                                'id', galleryMetadata.id,
                                'gallery_id', galleryMetadata.gallery_id,
                                'order_index', galleryMetadata.order_index,
                                'height', galleryMetadata.height,
                                'width', galleryMetadata.width,
                                'url', galleryMetadata.url,
                                'created_at', TO_CHAR(galleryMetadata.created_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"'),
                                'created_at_unix', galleryMetadata.created_at_unix,
                                'updated_at', TO_CHAR(galleryMetadata.updated_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"')
                            ) AS galleryMetadata
                        FROM 
                            gallery_metadatas galleryMetadata
                        WHERE 
                            galleryMetadata.gallery_id = g.id
                        ORDER BY 
                            galleryMetadata.order_index
                    ) AS orderedMetadata
            ) AS gallery_metadata
        FROM 
            galleries g
        WHERE 
            id = ?;
    `

	_, err := r.db.NewRaw(query, ID).Exec(ctx, &gallery)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Gallery{}, ErrNotFound
		}
		return models.Gallery{}, err
	}
	return gallery, nil
}

func (r *Repo) AddGallery(ctx context.Context, gallery models.Gallery) (models.Gallery, error) {
	query := `
        INSERT INTO 
            galleries (
                id,
                media_id
            )
        VALUES (
            ?,
            ?
        )
        RETURNING *
    `

	if _, err := r.db.NewRaw(query,
		gallery.ID,
		gallery.MediaID).Exec(ctx, &gallery); err != nil {
		var pgdriverErr pgdriver.Error
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgUniqueViolation {
			return models.Gallery{}, ErrDuplicateID
		}
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgConstraintViolation {
			return models.Gallery{}, ErrParentTableRecordNotFound
		}
		return models.Gallery{}, err
	}

	return gallery, nil
}

func (r *Repo) UpdateGallery(ctx context.Context, gallery models.Gallery) (models.Gallery, error) {
	query := `
        UPDATE 
            galleries
        SET 
            media_id = ?
        WHERE 
            id = ?
        RETURNING *
    `

	res, err := r.db.NewRaw(query,
		gallery.MediaID,
		gallery.ID).Exec(ctx, &gallery)
	if err != nil {
		var pgdriverErr pgdriver.Error
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgConstraintViolation {
			return models.Gallery{}, ErrParentTableRecordNotFound
		}
		if errors.Is(err, sql.ErrNoRows) {
			return models.Gallery{}, ErrNotFound
		}
		return models.Gallery{}, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return models.Gallery{}, err
	}
	if rowsAffected == 0 {
		return models.Gallery{}, ErrNotFound
	}
	return gallery, nil
}

func (r *Repo) DeleteGallery(ctx context.Context, ID uuid.UUID) error {
	query := `
        DELETE FROM 
            galleries
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
		return ErrNotFound
	}
	return nil
}

func (r *Repo) GalleryMetadatas(ctx context.Context) ([]models.GalleryMetadata, error) {
	var galleryMetadatas []models.GalleryMetadata

	query := `
        SELECT 
            id,
            gallery_id,
            order_index,
            height,
            width,
            url,
            created_at,
            created_at_unix,
            updated_at
        FROM 
            gallery_metadatas;
    `

	_, err := r.db.NewRaw(query).Exec(ctx, &galleryMetadatas)
	if err != nil {
		return []models.GalleryMetadata{}, err
	}
	return galleryMetadatas, nil
}

func (r *Repo) GalleryMetadataByID(ctx context.Context, ID uuid.UUID) (models.GalleryMetadata, error) {
	var galleryMetadata models.GalleryMetadata

	query := `
        SELECT 
            id,
            gallery_id,
            order_index,
            height,
            width,
            url,
            created_at,
            created_at_unix,
            updated_at
        FROM 
            gallery_metadatas
        WHERE 
            id = ?;
    `

	_, err := r.db.NewRaw(query, ID).Exec(ctx, &galleryMetadata)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.GalleryMetadata{}, ErrNotFound
		}
		return models.GalleryMetadata{}, err
	}
	return galleryMetadata, nil
}

func (r *Repo) AddGalleryMetadata(ctx context.Context, galleryMetadata models.GalleryMetadata) (models.GalleryMetadata, error) {
	query := `
        INSERT INTO 
            gallery_metadatas (
                id,
                gallery_id,
                order_index,
                height,
                width,
                url,
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
            ?
        )
        RETURNING *
    `

	timestamp := time.Now()
	galleryMetadata.CreatedAt = timestamp
	galleryMetadata.CreatedAtUnix = timestamp.Unix()
	galleryMetadata.UpdatedAt = timestamp

	if _, err := r.db.NewRaw(query,
		galleryMetadata.ID,
		galleryMetadata.GalleryID,
		galleryMetadata.OrderIndex,
		galleryMetadata.Height,
		galleryMetadata.Width,
		galleryMetadata.Url,
		galleryMetadata.CreatedAt,
		galleryMetadata.CreatedAtUnix,
		galleryMetadata.UpdatedAt).Exec(ctx, &galleryMetadata); err != nil {
		var pgdriverErr pgdriver.Error
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgUniqueViolation {
			return models.GalleryMetadata{}, ErrDuplicateID
		}
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgConstraintViolation {
			return models.GalleryMetadata{}, ErrParentTableRecordNotFound
		}
		return models.GalleryMetadata{}, err
	}

	return galleryMetadata, nil
}

func (r *Repo) UpdateGalleryMetadata(ctx context.Context, galleryMetadata models.GalleryMetadata) (models.GalleryMetadata, error) {
	query := `
        UPDATE 
            gallery_metadatas
        SET 
            gallery_id = ?,
            order_index = ?,
            height = ?,
            width = ?,
            url = ?,
            updated_at = ?
        WHERE 
            id = ?
        RETURNING *
    `

	galleryMetadata.UpdatedAt = time.Now()

	res, err := r.db.NewRaw(query,
		galleryMetadata.GalleryID,
		galleryMetadata.OrderIndex,
		galleryMetadata.Height,
		galleryMetadata.Width,
		galleryMetadata.Url,
		galleryMetadata.UpdatedAt,
		galleryMetadata.ID).Exec(ctx, &galleryMetadata)
	if err != nil {
		var pgdriverErr pgdriver.Error
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgConstraintViolation {
			return models.GalleryMetadata{}, ErrParentTableRecordNotFound
		}
		if errors.Is(err, sql.ErrNoRows) {
			return models.GalleryMetadata{}, ErrNotFound
		}
		return models.GalleryMetadata{}, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return models.GalleryMetadata{}, err
	}
	if rowsAffected == 0 {
		return models.GalleryMetadata{}, ErrNotFound
	}
	return galleryMetadata, nil
}

func (r *Repo) DeleteGalleryMetadata(ctx context.Context, ID uuid.UUID) error {
	query := `
        DELETE FROM 
            gallery_metadatas
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
		return ErrNotFound
	}
	return nil
}

func (r *Repo) Videos(ctx context.Context) ([]models.Video, error) {
	var videos []models.Video

	query := `
        SELECT
            id,
            media_id,
            url,
            height,
            width,
            created_at,
            created_at_unix,
            updated_at
        FROM
            videos;
    `

	_, err := r.db.NewRaw(query).Exec(ctx, &videos)
	if err != nil {
		return []models.Video{}, err
	}
	return videos, nil
}

func (r *Repo) VideoByID(ctx context.Context, ID uuid.UUID) (models.Video, error) {
	var video models.Video

	query := `
        SELECT
            id,
            media_id,
            url,
            height,
            width,
            created_at,
            created_at_unix,
            updated_at
        FROM
            videos
        WHERE
            id = ?;
    `

	_, err := r.db.NewRaw(query, ID).Exec(ctx, &video)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Video{}, ErrNotFound
		}
		return models.Video{}, err
	}
	return video, nil
}

func (r *Repo) AddVideo(ctx context.Context, video models.Video) (models.Video, error) {
	query := `
        INSERT INTO
            videos (
                id,
                media_id,
                url,
                height,
                width,
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
            ?
        )
        RETURNING *
    `

	timestamp := time.Now()
	video.CreatedAt = timestamp
	video.CreatedAtUnix = timestamp.Unix()
	video.UpdatedAt = timestamp

	if _, err := r.db.NewRaw(query,
		video.ID,
		video.MediaID,
		video.Url,
		video.Height,
		video.Width,
		video.CreatedAt,
		video.CreatedAtUnix,
		video.UpdatedAt).Exec(ctx, &video); err != nil {
		var pgdriverErr pgdriver.Error
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgUniqueViolation {
			return models.Video{}, ErrDuplicateID
		}
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgConstraintViolation {
			return models.Video{}, ErrParentTableRecordNotFound
		}
		return models.Video{}, err
	}
	return video, nil
}

func (r *Repo) UpdateVideo(ctx context.Context, video models.Video) (models.Video, error) {
	query := `
        UPDATE
            videos
        SET
            media_id = ?,
            url = ?,
            height = ?,
            width = ?,
            updated_at = ?
        WHERE
            id = ?
        RETURNING *
    `

	video.UpdatedAt = time.Now()

	res, err := r.db.NewRaw(query,
		video.MediaID,
		video.Url,
		video.Height,
		video.Width,
		video.UpdatedAt,
		video.ID,
	).Exec(ctx, &video)
	if err != nil {
		var pgdriverErr pgdriver.Error
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgConstraintViolation {
			return models.Video{}, ErrParentTableRecordNotFound
		}
		if errors.Is(err, sql.ErrNoRows) {
			return models.Video{}, ErrNotFound
		}
		return models.Video{}, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return models.Video{}, err
	}
	if rowsAffected == 0 {
		return models.Video{}, ErrNotFound
	}
	return video, nil
}

func (r *Repo) DeleteVideo(ctx context.Context, ID uuid.UUID) error {
	query := `
        DELETE FROM 
            videos
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
		return ErrNotFound
	}
	return nil
}

func (r *Repo) Links(ctx context.Context) ([]models.Link, error) {
	var links []models.Link

	query := `
        SELECT
            id,
            media_id,
            link,
            created_at,
            created_at_unix,
            updated_at
        FROM
            links;
    `

	_, err := r.db.NewRaw(query).Exec(ctx, &links)
	if err != nil {
		return []models.Link{}, err
	}
	return links, nil
}

func (r *Repo) LinkByID(ctx context.Context, ID uuid.UUID) (models.Link, error) {
	var link models.Link

	query := `
        SELECT
            id,
            media_id,
            link,
            created_at,
            created_at_unix,
            updated_at
        FROM
            links
        WHERE
            id = ?;
    `

	_, err := r.db.NewRaw(query, ID).Exec(ctx, &link)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Link{}, ErrNotFound
		}
		return models.Link{}, err
	}
	return link, nil
}

func (r *Repo) AddLink(ctx context.Context, link models.Link) (models.Link, error) {
	query := `
        INSERT INTO
            links (
                id,
                media_id,
                link,
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
            ?
        )
        RETURNING *
    `

	timestamp := time.Now()
	link.CreatedAt = timestamp
	link.CreatedAtUnix = timestamp.Unix()
	link.UpdatedAt = timestamp

	if _, err := r.db.NewRaw(query,
		link.ID,
		link.MediaID,
		link.Link,
		link.CreatedAt,
		link.CreatedAtUnix,
		link.UpdatedAt).Exec(ctx, &link); err != nil {
		var pgdriverErr pgdriver.Error
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgUniqueViolation {
			return models.Link{}, ErrDuplicateID
		}
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgConstraintViolation {
			return models.Link{}, ErrParentTableRecordNotFound
		}
		return models.Link{}, err
	}
	return link, nil
}

func (r *Repo) UpdateLink(ctx context.Context, link models.Link) (models.Link, error) {
	query := `
        UPDATE
            links
        SET
            media_id = ?,
            link = ?,
            updated_at = ?
        WHERE
            id = ?
        RETURNING *
    `

	link.UpdatedAt = time.Now()

	res, err := r.db.NewRaw(query,
		link.MediaID,
		link.Link,
		link.UpdatedAt,
		link.ID,
	).Exec(ctx, &link)
	if err != nil {
		var pgdriverErr pgdriver.Error
		if errors.As(err, &pgdriverErr) && pgdriverErr.Field('C') == pgConstraintViolation {
			return models.Link{}, ErrParentTableRecordNotFound
		}
		if errors.Is(err, sql.ErrNoRows) {
			return models.Link{}, ErrNotFound
		}
		return models.Link{}, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return models.Link{}, err
	}
	if rowsAffected == 0 {
		return models.Link{}, ErrNotFound
	}
	return link, nil
}

func (r *Repo) DeleteLink(ctx context.Context, ID uuid.UUID) error {
	query := `
        DELETE FROM 
            links
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
		return ErrNotFound
	}
	return nil
}
