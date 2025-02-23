package media_test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"slices"
	"testing"
	"time"

	"github.com/glowfi/voxpopuli/backend/pkg/models"
	mediarepo "github.com/glowfi/voxpopuli/backend/pkg/repo/media"
	postrepo "github.com/glowfi/voxpopuli/backend/pkg/repo/post"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dbfixture"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

func connectPostgres(user, password, address, dbName string) *bun.DB {
	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", user, password, address, dbName)
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(sqldb, pgdialect.New())
	return db
}

func setupPostgres(t *testing.T, fixtureFiles ...string) *bun.DB {
	db := connectPostgres("postgres", "postgres", "127.0.0.1:5432", "voxpopuli")

	if err := db.Ping(); err != nil {
		t.Fatal("db error:", err)
	}
	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Log("db close error:", err)
		}
	})

	// add query logging hook
	db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))

	db.RegisterModel((*models.Topic)(nil))
	db.RegisterModel((*models.Voxsphere)(nil))
	db.RegisterModel((*models.User)(nil))
	db.RegisterModel((*models.Post)(nil))
	db.RegisterModel((*models.PostMedia)(nil))
	db.RegisterModel((*models.Image)(nil))
	db.RegisterModel((*models.ImageMetadata)(nil))
	db.RegisterModel((*models.Gif)(nil))
	db.RegisterModel((*models.GifMetadata)(nil))
	db.RegisterModel((*models.Gallery)(nil))
	db.RegisterModel((*models.GalleryMetadata)(nil))
	db.RegisterModel((*models.Video)(nil))
	db.RegisterModel((*models.Link)(nil))

	// drop all rows of the topics,voxspheres table
	if _, err := db.NewTruncateTable().Cascade().Model((*models.Topic)(nil)).Exec(context.Background()); err != nil {
		t.Fatal("truncate table failed:", err)
	}
	if _, err := db.NewTruncateTable().Cascade().Model((*models.Voxsphere)(nil)).Exec(context.Background()); err != nil {
		t.Fatal("truncate table failed:", err)
	}
	if _, err := db.NewTruncateTable().Cascade().Model((*models.User)(nil)).Exec(context.Background()); err != nil {
		t.Fatal("truncate table failed:", err)
	}
	if _, err := db.NewTruncateTable().Cascade().Model((*models.Post)(nil)).Exec(context.Background()); err != nil {
		t.Fatal("truncate table failed:", err)
	}
	if _, err := db.NewTruncateTable().Cascade().Model((*models.PostMedia)(nil)).Exec(context.Background()); err != nil {
		t.Fatal("truncate table failed:", err)
	}
	if _, err := db.NewTruncateTable().Cascade().Model((*models.Image)(nil)).Exec(context.Background()); err != nil {
		t.Fatal("truncate table failed:", err)
	}
	if _, err := db.NewTruncateTable().Cascade().Model((*models.ImageMetadata)(nil)).Exec(context.Background()); err != nil {
		t.Fatal("truncate table failed:", err)
	}
	if _, err := db.NewTruncateTable().Cascade().Model((*models.Gif)(nil)).Exec(context.Background()); err != nil {
		t.Fatal("truncate table failed:", err)
	}
	if _, err := db.NewTruncateTable().Cascade().Model((*models.GifMetadata)(nil)).Exec(context.Background()); err != nil {
		t.Fatal("truncate table failed:", err)
	}
	if _, err := db.NewTruncateTable().Cascade().Model((*models.Gallery)(nil)).Exec(context.Background()); err != nil {
		t.Fatal("truncate table failed:", err)
	}
	if _, err := db.NewTruncateTable().Cascade().Model((*models.GalleryMetadata)(nil)).Exec(context.Background()); err != nil {
		t.Fatal("truncate table failed:", err)
	}
	if _, err := db.NewTruncateTable().Cascade().Model((*models.Video)(nil)).Exec(context.Background()); err != nil {
		t.Fatal("truncate table failed:", err)
	}
	if _, err := db.NewTruncateTable().Cascade().Model((*models.Link)(nil)).Exec(context.Background()); err != nil {
		t.Fatal("truncate table failed:", err)
	}

	// load fixture
	fixture := dbfixture.New(db)
	if err := fixture.Load(context.Background(), os.DirFS("testdata"), fixtureFiles...); err != nil {
		t.Fatal("failed to load fixtures", err)
	}

	return db
}

func assertTimeWithinRange(t *testing.T, time, start, end time.Time) {
	t.Helper()

	assert.NotZero(
		t,
		time,
		"expect time to be a non-zeroed out value",
	)

	// Ensure all times are in UTC
	time = time.UTC()
	start = start.UTC()
	end = end.UTC()

	assert.WithinRange(t, time, start, end)
}

func assertPostMedias(t *testing.T, wantPostMedias, gotPostMedias []models.PostMedia) {
	t.Helper()

	if len(wantPostMedias) != len(gotPostMedias) {
		t.Fatal("length of wantPostMedias and gotPostMedias do not match")
	}

	for _, postMedia := range wantPostMedias {
		idx := slices.IndexFunc(gotPostMedias, func(v models.PostMedia) bool {
			return v.ID == postMedia.ID
		})

		if idx == -1 {
			t.Fatal(fmt.Sprintf("post media of ID %v is not present in gotPostMedias", postMedia.ID))
			return
		}
		assert.Equal(t, postMedia, gotPostMedias[idx], "expect post media to match")
	}
}

func assertImageMetadataWithoutTimestamp(t *testing.T, wantImageMetadata, gotImageMetadata models.ImageMetadata) {
	assert.Equal(t, wantImageMetadata.ID, gotImageMetadata.ID, "expected id to match")
	assert.Equal(t, wantImageMetadata.ImageID, gotImageMetadata.ImageID, "expected image id to match")
	assert.Equal(t, wantImageMetadata.Height, gotImageMetadata.Height, "expected height to match")
	assert.Equal(t, wantImageMetadata.Width, gotImageMetadata.Width, "expected width to match")
	assert.Equal(t, wantImageMetadata.Url, gotImageMetadata.Url, "expected url to match")
}

func assertImageMetadatasWithoutTimestamp(t *testing.T, wantImageMetadatas, gotImageMetadatas []models.ImageMetadata) {
	t.Helper()

	if len(wantImageMetadatas) != len(gotImageMetadatas) {
		t.Fatal("length of wantImageMetadatas and gotImageMetadatas do not match")
	}

	for _, imageMetadata := range wantImageMetadatas {
		idx := slices.IndexFunc(gotImageMetadatas, func(im models.ImageMetadata) bool {
			return im.ID == imageMetadata.ID
		})

		if idx == -1 {
			t.Fatal(fmt.Sprintf("image metadata %v of ID %v is not present in gotImageMetadatas", imageMetadata.Url, imageMetadata.ID))
			return
		}
		assertImageMetadataWithoutTimestamp(t, imageMetadata, gotImageMetadatas[idx])
	}
}

func assertImageMetadatasWithTimestamp(t *testing.T, wantImageMetadatas, gotImageMetadatas []models.ImageMetadata) {
	t.Helper()

	for _, imageMetadata := range wantImageMetadatas {
		idx := slices.IndexFunc(gotImageMetadatas, func(im models.ImageMetadata) bool {
			return im.ID == imageMetadata.ID
		})

		if idx == -1 {
			t.Fatal(fmt.Sprintf("image metadata %v of ID %v is not present in gotImageMetadatas", imageMetadata.Url, imageMetadata.ID))
			return
		}
		assert.Equal(t, imageMetadata, gotImageMetadatas[idx], "expected image metadata to match")
	}
}

func assertImage(t *testing.T, wantImage, gotImage models.Image) {
	assert.Equal(t, wantImage.ID, gotImage.ID, "expected id to match")
	assert.Equal(t, wantImage.MediaID, gotImage.MediaID, "expected media id to match")
	assertImageMetadatasWithTimestamp(t, wantImage.ImageMetadata, gotImage.ImageMetadata)
}

func assertImages(t *testing.T, wantImages, gotImages []models.Image) {
	t.Helper()

	if len(wantImages) != len(gotImages) {
		t.Fatal("length of wantImages and gotImages do not match")
	}

	for _, image := range wantImages {
		idx := slices.IndexFunc(gotImages, func(im models.Image) bool {
			return im.ID == image.ID
		})

		if idx == -1 {
			t.Fatal(fmt.Sprintf("image %v of ID %v is not present in gotImages", image.MediaID, image.ID))
			return
		}
		assertImage(t, image, gotImages[idx])
	}
}

func assertLinkWithoutTimestamp(t *testing.T, wantLink, gotLink models.Link) {
	assert.Equal(t, wantLink.ID, gotLink.ID, "expect link id to match")
	assert.Equal(t, wantLink.MediaID, gotLink.MediaID, "expect link media id to match")
	assert.Equal(t, wantLink.Link, gotLink.Link, "expect link url to match")
}

func assertLinksWithoutTimestamp(t *testing.T, wantLinks, gotLinks []models.Link) {
	t.Helper()

	if len(wantLinks) != len(gotLinks) {
		t.Fatal("length of wantLinks and gotLinks do not match")
	}

	for _, link := range wantLinks {
		idx := slices.IndexFunc(gotLinks, func(im models.Link) bool {
			return im.ID == link.ID
		})

		if idx == -1 {
			t.Fatal(fmt.Sprintf("link %v of url %v is not present in gotLinks", link.ID, link.Link))
			return
		}
		assertLinkWithoutTimestamp(t, link, gotLinks[idx])
	}
}

func TestRepo_PostMedias(t *testing.T) {
	tests := []struct {
		name           string
		fixtureFiles   []string
		wantPostMedias []models.PostMedia
		wantErr        error
	}{
		{
			name:         "post medias :POS",
			fixtureFiles: []string{"topics.yml", "voxspheres.yml", "users.yml", "posts.yml", "post_medias.yml"},
			wantPostMedias: []models.PostMedia{
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					PostID:    uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					MediaType: models.MediaTypeImage,
				},
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					PostID:    uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					MediaType: models.MediaTypeGif,
				},
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					PostID:    uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					MediaType: models.MediaTypeGallery,
				},
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000004"),
					PostID:    uuid.MustParse("00000000-0000-0000-0000-000000000004"),
					MediaType: models.MediaTypeVideo,
				},
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000005"),
					PostID:    uuid.MustParse("00000000-0000-0000-0000-000000000005"),
					MediaType: models.MediaTypeLink,
				},
			},
			wantErr: nil,
		},
		{
			name:           "no post medias :POS",
			fixtureFiles:   []string{},
			wantPostMedias: nil,
			wantErr:        nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupPostgres(t, tt.fixtureFiles...)
			pgrepo := mediarepo.NewRepo(db)

			gotPostMedias, gotErr := pgrepo.PostMedias(context.Background())

			assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")
			assert.Equal(t, tt.wantPostMedias, gotPostMedias, "expect post medias to match")
		})
	}
}

func TestRepo_PostMediaByID(t *testing.T) {
	type args struct {
		ID uuid.UUID
	}
	tests := []struct {
		name          string
		fixtureFiles  []string
		args          args
		wantPostMedia models.PostMedia
		wantErr       error
	}{
		{
			name:         "post media not found :NEG",
			fixtureFiles: []string{"topics.yml", "voxspheres.yml", "users.yml", "posts.yml", "post_medias.yml"},
			args: args{
				ID: uuid.MustParse("00000000-0000-0000-0000-000000000009"),
			},
			wantPostMedia: models.PostMedia{},
			wantErr:       mediarepo.ErrNotFound,
		},
		{
			name:         "post media found :POS",
			fixtureFiles: []string{"topics.yml", "voxspheres.yml", "users.yml", "posts.yml", "post_medias.yml"},
			args: args{
				ID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			},
			wantPostMedia: models.PostMedia{
				ID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				PostID:    uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				MediaType: models.MediaTypeImage,
			},

			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupPostgres(t, tt.fixtureFiles...)
			pgrepo := mediarepo.NewRepo(db)

			gotPostMedia, gotErr := pgrepo.PostMediaByID(context.Background(), tt.args.ID)

			assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")
			assert.Equal(t, tt.wantPostMedia, gotPostMedia, "expect post media to match")
		})
	}
}

func TestRepo_AddPostMedia(t *testing.T) {
	type args struct {
		postMedia models.PostMedia
	}
	tests := []struct {
		name           string
		fixtureFiles   []string
		args           args
		wantPostMedia  models.PostMedia
		wantPostMedias []models.PostMedia
		wantErr        error
	}{
		{
			name:         "duplicate media id :NEG",
			fixtureFiles: []string{"topics.yml", "voxspheres.yml", "users.yml", "posts.yml", "post_medias.yml"},
			args: args{
				postMedia: models.PostMedia{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					PostID:    uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					MediaType: models.MediaTypeImage,
				},
			},
			wantPostMedia: models.PostMedia{},
			wantPostMedias: []models.PostMedia{
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					PostID:    uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					MediaType: models.MediaTypeImage,
				},
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					PostID:    uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					MediaType: models.MediaTypeGif,
				},
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					PostID:    uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					MediaType: models.MediaTypeGallery,
				},
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000004"),
					PostID:    uuid.MustParse("00000000-0000-0000-0000-000000000004"),
					MediaType: models.MediaTypeVideo,
				},
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000005"),
					PostID:    uuid.MustParse("00000000-0000-0000-0000-000000000005"),
					MediaType: models.MediaTypeLink,
				},
			},
			wantErr: mediarepo.ErrDuplicateID,
		},
		{
			name:         "post not present in parent table :NEG",
			fixtureFiles: []string{"topics.yml", "voxspheres.yml", "users.yml", "posts.yml", "post_medias.yml"},
			args: args{
				postMedia: models.PostMedia{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000006"),
					PostID:    uuid.MustParse("00000000-0000-0000-0000-000000000009"),
					MediaType: models.MediaTypeImage,
				},
			},
			wantPostMedia: models.PostMedia{},
			wantPostMedias: []models.PostMedia{
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					PostID:    uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					MediaType: models.MediaTypeImage,
				},
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					PostID:    uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					MediaType: models.MediaTypeGif,
				},
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					PostID:    uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					MediaType: models.MediaTypeGallery,
				},
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000004"),
					PostID:    uuid.MustParse("00000000-0000-0000-0000-000000000004"),
					MediaType: models.MediaTypeVideo,
				},
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000005"),
					PostID:    uuid.MustParse("00000000-0000-0000-0000-000000000005"),
					MediaType: models.MediaTypeLink,
				},
			},
			wantErr: mediarepo.ErrParentTableRecordNotFound,
		},
		{
			name:         "add post media :POS",
			fixtureFiles: []string{"topics.yml", "voxspheres.yml", "users.yml", "posts.yml", "post_medias.yml"},
			args: args{
				postMedia: models.PostMedia{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000006"),
					PostID:    uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					MediaType: models.MediaTypeImage,
				},
			},
			wantPostMedia: models.PostMedia{
				ID:        uuid.MustParse("00000000-0000-0000-0000-000000000006"),
				PostID:    uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				MediaType: models.MediaTypeImage,
			},
			wantPostMedias: []models.PostMedia{
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					PostID:    uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					MediaType: models.MediaTypeImage,
				},
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					PostID:    uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					MediaType: models.MediaTypeGif,
				},
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					PostID:    uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					MediaType: models.MediaTypeGallery,
				},
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000004"),
					PostID:    uuid.MustParse("00000000-0000-0000-0000-000000000004"),
					MediaType: models.MediaTypeVideo,
				},
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000005"),
					PostID:    uuid.MustParse("00000000-0000-0000-0000-000000000005"),
					MediaType: models.MediaTypeLink,
				},
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000006"),
					PostID:    uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					MediaType: models.MediaTypeImage,
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupPostgres(t, tt.fixtureFiles...)
			pgrepo := mediarepo.NewRepo(db)

			gotPostMedia, gotErr := pgrepo.AddPostMedia(context.Background(), tt.args.postMedia)

			assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")
			assert.Equal(t, tt.wantPostMedia, gotPostMedia, "expect postMedia to match")

			gotPostMedias, err := pgrepo.PostMedias(context.Background())

			assert.NoError(t, err, "expect no error while getting emojis")
			assertPostMedias(t, tt.wantPostMedias, gotPostMedias)
		})
	}
}

func TestRepo_UpdatePostMedia(t *testing.T) {
	type args struct {
		postMedia models.PostMedia
	}
	tests := []struct {
		name           string
		fixtureFiles   []string
		args           args
		wantPostMedia  models.PostMedia
		wantPostMedias []models.PostMedia
		wantErr        error
	}{
		{
			name:         "post media not found :NEG",
			fixtureFiles: []string{"topics.yml", "voxspheres.yml", "users.yml", "posts.yml", "post_medias.yml"},
			args: args{
				postMedia: models.PostMedia{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000009"),
					PostID:    uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					MediaType: models.MediaTypeImage,
				},
			},
			wantPostMedia: models.PostMedia{},
			wantPostMedias: []models.PostMedia{
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					PostID:    uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					MediaType: models.MediaTypeImage,
				},
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					PostID:    uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					MediaType: models.MediaTypeGif,
				},
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					PostID:    uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					MediaType: models.MediaTypeGallery,
				},
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000004"),
					PostID:    uuid.MustParse("00000000-0000-0000-0000-000000000004"),
					MediaType: models.MediaTypeVideo,
				},
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000005"),
					PostID:    uuid.MustParse("00000000-0000-0000-0000-000000000005"),
					MediaType: models.MediaTypeLink,
				},
			},
			wantErr: mediarepo.ErrNotFound,
		},
		{
			name:         "post not present in parent table :NEG",
			fixtureFiles: []string{"topics.yml", "voxspheres.yml", "users.yml", "posts.yml", "post_medias.yml"},
			args: args{
				postMedia: models.PostMedia{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					PostID:    uuid.MustParse("00000000-0000-0000-0000-000000000009"),
					MediaType: models.MediaTypeImage,
				},
			},
			wantPostMedia: models.PostMedia{},
			wantPostMedias: []models.PostMedia{
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					PostID:    uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					MediaType: models.MediaTypeImage,
				},
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					PostID:    uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					MediaType: models.MediaTypeGif,
				},
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					PostID:    uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					MediaType: models.MediaTypeGallery,
				},
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000004"),
					PostID:    uuid.MustParse("00000000-0000-0000-0000-000000000004"),
					MediaType: models.MediaTypeVideo,
				},
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000005"),
					PostID:    uuid.MustParse("00000000-0000-0000-0000-000000000005"),
					MediaType: models.MediaTypeLink,
				},
			},
			wantErr: mediarepo.ErrParentTableRecordNotFound,
		},
		{
			name:         "update post media id :NEG",
			fixtureFiles: []string{"topics.yml", "voxspheres.yml", "users.yml", "posts.yml", "post_medias.yml"},
			args: args{
				postMedia: models.PostMedia{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					PostID:    uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					MediaType: models.MediaTypeImage,
				},
			},
			wantPostMedia: models.PostMedia{
				ID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				PostID:    uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				MediaType: models.MediaTypeImage,
			},
			wantPostMedias: []models.PostMedia{
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					PostID:    uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					MediaType: models.MediaTypeImage,
				},
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					PostID:    uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					MediaType: models.MediaTypeGif,
				},
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					PostID:    uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					MediaType: models.MediaTypeGallery,
				},
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000004"),
					PostID:    uuid.MustParse("00000000-0000-0000-0000-000000000004"),
					MediaType: models.MediaTypeVideo,
				},
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000005"),
					PostID:    uuid.MustParse("00000000-0000-0000-0000-000000000005"),
					MediaType: models.MediaTypeLink,
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupPostgres(t, tt.fixtureFiles...)
			pgrepo := mediarepo.NewRepo(db)

			gotPostMedia, gotErr := pgrepo.UpdatePostMedia(context.Background(), tt.args.postMedia)

			assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")
			assert.Equal(t, tt.wantPostMedia, gotPostMedia, "expect post media to match")

			gotPostMedias, err := pgrepo.PostMedias(context.Background())

			assert.NoError(t, err, "expect no error while getting post medias")
			assertPostMedias(t, tt.wantPostMedias, gotPostMedias)
		})
	}
}

func TestRepo_DeletePostMedia(t *testing.T) {
	type args struct {
		ID uuid.UUID
	}
	tests := []struct {
		name           string
		fixtureFiles   []string
		args           args
		wantPostMedias []models.PostMedia
		wantErr        error
	}{
		{
			name:         "post media not found :NEG",
			fixtureFiles: []string{"topics.yml", "voxspheres.yml", "users.yml", "posts.yml", "post_medias.yml"},
			args: args{
				ID: uuid.MustParse("00000000-0000-0000-0000-000000000009"),
			},
			wantPostMedias: []models.PostMedia{
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					PostID:    uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					MediaType: models.MediaTypeImage,
				},
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					PostID:    uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					MediaType: models.MediaTypeGif,
				},
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					PostID:    uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					MediaType: models.MediaTypeGallery,
				},
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000004"),
					PostID:    uuid.MustParse("00000000-0000-0000-0000-000000000004"),
					MediaType: models.MediaTypeVideo,
				},
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000005"),
					PostID:    uuid.MustParse("00000000-0000-0000-0000-000000000005"),
					MediaType: models.MediaTypeLink,
				},
			},
			wantErr: mediarepo.ErrNotFound,
		},
		{
			name:         "post media deleted :POS",
			fixtureFiles: []string{"topics.yml", "voxspheres.yml", "users.yml", "posts.yml", "post_medias.yml"},
			args: args{
				ID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			},
			wantPostMedias: []models.PostMedia{
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					PostID:    uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					MediaType: models.MediaTypeGif,
				},
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					PostID:    uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					MediaType: models.MediaTypeGallery,
				},
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000004"),
					PostID:    uuid.MustParse("00000000-0000-0000-0000-000000000004"),
					MediaType: models.MediaTypeVideo,
				},
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000005"),
					PostID:    uuid.MustParse("00000000-0000-0000-0000-000000000005"),
					MediaType: models.MediaTypeLink,
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupPostgres(t, tt.fixtureFiles...)
			pgrepo := mediarepo.NewRepo(db)

			gotErr := pgrepo.DeletePostMedia(context.Background(), tt.args.ID)

			assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")

			gotPostMedias, err := pgrepo.PostMedias(context.Background())

			assert.NoError(t, err, "expect no error while getting post medias")
			assertPostMedias(t, tt.wantPostMedias, gotPostMedias)
		})
	}
}

func TestRepo_PostMediaForeignKeyCascade(t *testing.T) {
	t.Run("on deleting post from the parent table ,no child references should exist in post_medias table", func(t *testing.T) {
		db := setupPostgres(t, "topics.yml", "voxspheres.yml", "users.yml", "posts.yml", "post_medias.yml")
		mediaPgrepo := mediarepo.NewRepo(db)
		postPgrepo := postrepo.NewRepo(db)

		err := postPgrepo.DeletePost(context.Background(), uuid.MustParse("00000000-0000-0000-0000-000000000001"))

		assert.NoError(t, err, "expect no error while deleting post")

		gotPostMedias, err := mediaPgrepo.PostMedias(context.Background())

		assert.NoError(t, err, "expect no error while getting post medias")
		wantPostMedias := []models.PostMedia{
			{
				ID:        uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				PostID:    uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				MediaType: models.MediaTypeGif,
			},
			{
				ID:        uuid.MustParse("00000000-0000-0000-0000-000000000003"),
				PostID:    uuid.MustParse("00000000-0000-0000-0000-000000000003"),
				MediaType: models.MediaTypeGallery,
			},
			{
				ID:        uuid.MustParse("00000000-0000-0000-0000-000000000004"),
				PostID:    uuid.MustParse("00000000-0000-0000-0000-000000000004"),
				MediaType: models.MediaTypeVideo,
			},
			{
				ID:        uuid.MustParse("00000000-0000-0000-0000-000000000005"),
				PostID:    uuid.MustParse("00000000-0000-0000-0000-000000000005"),
				MediaType: models.MediaTypeLink,
			},
		}
		assertPostMedias(t, wantPostMedias, gotPostMedias)
	})
}

func TestRepo_Images(t *testing.T) {
	tests := []struct {
		name         string
		fixtureFiles []string
		wantImages   []models.Image
		wantErr      error
	}{
		{
			name: "images :POS",
			fixtureFiles: []string{
				"topics.yml",
				"voxspheres.yml",
				"users.yml",
				"posts.yml",
				"post_medias.yml",
				"images.yml",
				"image_metadatas.yml",
			},
			wantImages: []models.Image{
				{
					ID:      uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					MediaID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					ImageMetadata: []models.ImageMetadata{
						{
							ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
							ImageID:       uuid.MustParse("00000000-0000-0000-0000-000000000001"),
							Height:        1080,
							Width:         1920,
							Url:           "https://example.com/image1.jpg",
							CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
							CreatedAtUnix: 1725091100,
							UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
						},
						{
							ID:            uuid.MustParse("00000000-0000-0000-0000-000000000002"),
							ImageID:       uuid.MustParse("00000000-0000-0000-0000-000000000001"),
							Height:        720,
							Width:         1280,
							Url:           "https://example.com/image2.jpg",
							CreatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
							CreatedAtUnix: 1725091101,
							UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
						},
					},
				},
			},
			wantErr: nil,
		},
		{
			name:         "no images :POS",
			fixtureFiles: []string{},
			wantImages:   nil,
			wantErr:      nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupPostgres(t, tt.fixtureFiles...)
			pgrepo := mediarepo.NewRepo(db)

			gotImages, gotErr := pgrepo.Images(context.Background())

			assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")
			assertImages(t, tt.wantImages, gotImages)
		})
	}
}

func TestRepo_ImageByID(t *testing.T) {
	type args struct {
		ID uuid.UUID
	}
	tests := []struct {
		name         string
		fixtureFiles []string
		args         args
		wantImage    models.Image
		wantErr      error
	}{
		{
			name: "image not found :NEG",
			fixtureFiles: []string{
				"topics.yml",
				"voxspheres.yml",
				"users.yml",
				"posts.yml",
				"post_medias.yml",
				"images.yml",
				"image_metadatas.yml",
			},
			args: args{
				ID: uuid.MustParse("00000000-0000-0000-0000-000000000009"),
			},
			wantImage: models.Image{},
			wantErr:   mediarepo.ErrNotFound,
		},
		{
			name: "image found :POS",
			fixtureFiles: []string{
				"topics.yml",
				"voxspheres.yml",
				"users.yml",
				"posts.yml",
				"post_medias.yml",
				"images.yml",
				"image_metadatas.yml",
			},
			args: args{
				ID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			},
			wantImage: models.Image{
				ID:      uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				MediaID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				ImageMetadata: []models.ImageMetadata{
					{
						ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
						ImageID:       uuid.MustParse("00000000-0000-0000-0000-000000000001"),
						Height:        1080,
						Width:         1920,
						Url:           "https://example.com/image1.jpg",
						CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
						CreatedAtUnix: 1725091100,
						UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					},
					{
						ID:            uuid.MustParse("00000000-0000-0000-0000-000000000002"),
						ImageID:       uuid.MustParse("00000000-0000-0000-0000-000000000001"),
						Height:        720,
						Width:         1280,
						Url:           "https://example.com/image2.jpg",
						CreatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
						CreatedAtUnix: 1725091101,
						UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
					},
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupPostgres(t, tt.fixtureFiles...)
			pgrepo := mediarepo.NewRepo(db)

			gotImage, gotErr := pgrepo.ImageByID(context.Background(), tt.args.ID)

			assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")
			assertImage(t, tt.wantImage, gotImage)
		})
	}
}

func TestRepo_AddImage(t *testing.T) {
	type args struct {
		image models.Image
	}
	tests := []struct {
		name         string
		fixtureFiles []string
		args         args
		wantImage    models.Image
		wantImages   []models.Image
		wantErr      error
	}{
		{
			name: "duplicate image id :NEG",
			fixtureFiles: []string{
				"topics.yml",
				"voxspheres.yml",
				"users.yml",
				"posts.yml",
				"post_medias.yml",
				"images.yml",
				"image_metadatas.yml",
			},
			args: args{
				image: models.Image{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					MediaID:       uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					ImageMetadata: []models.ImageMetadata{},
				},
			},
			wantImage: models.Image{},
			wantImages: []models.Image{
				{
					ID:      uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					MediaID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					ImageMetadata: []models.ImageMetadata{
						{
							ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
							ImageID:       uuid.MustParse("00000000-0000-0000-0000-000000000001"),
							Height:        1080,
							Width:         1920,
							Url:           "https://example.com/image1.jpg",
							CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
							CreatedAtUnix: 1725091100,
							UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
						},
						{
							ID:            uuid.MustParse("00000000-0000-0000-0000-000000000002"),
							ImageID:       uuid.MustParse("00000000-0000-0000-0000-000000000001"),
							Height:        720,
							Width:         1280,
							Url:           "https://example.com/image2.jpg",
							CreatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
							CreatedAtUnix: 1725091101,
							UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
						},
					},
				},
			},
			wantErr: mediarepo.ErrDuplicateID,
		},
		{
			name: "media not found in parent table :NEG",
			fixtureFiles: []string{
				"topics.yml",
				"voxspheres.yml",
				"users.yml",
				"posts.yml",
				"post_medias.yml",
				"images.yml",
				"image_metadatas.yml",
			},
			args: args{
				image: models.Image{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					MediaID:       uuid.MustParse("00000000-0000-0000-0000-000000000009"),
					ImageMetadata: []models.ImageMetadata{},
				},
			},
			wantImage: models.Image{},
			wantImages: []models.Image{
				{
					ID:      uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					MediaID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					ImageMetadata: []models.ImageMetadata{
						{
							ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
							ImageID:       uuid.MustParse("00000000-0000-0000-0000-000000000001"),
							Height:        1080,
							Width:         1920,
							Url:           "https://example.com/image1.jpg",
							CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
							CreatedAtUnix: 1725091100,
							UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
						},
						{
							ID:            uuid.MustParse("00000000-0000-0000-0000-000000000002"),
							ImageID:       uuid.MustParse("00000000-0000-0000-0000-000000000001"),
							Height:        720,
							Width:         1280,
							Url:           "https://example.com/image2.jpg",
							CreatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
							CreatedAtUnix: 1725091101,
							UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
						},
					},
				},
			},
			wantErr: mediarepo.ErrParentTableRecordNotFound,
		},
		{
			name: "add image :POS",
			fixtureFiles: []string{
				"topics.yml",
				"voxspheres.yml",
				"users.yml",
				"posts.yml",
				"post_medias.yml",
				"images.yml",
				"image_metadatas.yml",
			},
			args: args{
				image: models.Image{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					MediaID:       uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					ImageMetadata: []models.ImageMetadata{},
				},
			},
			wantImage: models.Image{
				ID:            uuid.MustParse("00000000-0000-0000-0000-000000000003"),
				MediaID:       uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				ImageMetadata: []models.ImageMetadata{},
			},
			wantImages: []models.Image{
				{
					ID:      uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					MediaID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					ImageMetadata: []models.ImageMetadata{
						{
							ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
							ImageID:       uuid.MustParse("00000000-0000-0000-0000-000000000001"),
							Height:        1080,
							Width:         1920,
							Url:           "https://example.com/image1.jpg",
							CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
							CreatedAtUnix: 1725091100,
							UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
						},
						{
							ID:            uuid.MustParse("00000000-0000-0000-0000-000000000002"),
							ImageID:       uuid.MustParse("00000000-0000-0000-0000-000000000001"),
							Height:        720,
							Width:         1280,
							Url:           "https://example.com/image2.jpg",
							CreatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
							CreatedAtUnix: 1725091101,
							UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
						},
					},
				},
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					MediaID:       uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					ImageMetadata: []models.ImageMetadata{},
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		db := setupPostgres(t, tt.fixtureFiles...)
		pgrepo := mediarepo.NewRepo(db)

		gotImage, gotErr := pgrepo.AddImage(context.Background(), tt.args.image)

		assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")
		assert.Equal(t, tt.wantImage, gotImage, "expect image to match")

		gotImages, err := pgrepo.Images(context.Background())

		assert.NoError(t, err, "expect no error while getting images")
		assertImages(t, tt.wantImages, gotImages)
	}
}

func TestRepo_UpdateImage(t *testing.T) {
	type args struct {
		image models.Image
	}
	tests := []struct {
		name         string
		fixtureFiles []string
		args         args
		wantImage    models.Image
		wantImages   []models.Image
		wantErr      error
	}{
		{
			name: "image not found :NEG",
			fixtureFiles: []string{
				"topics.yml",
				"voxspheres.yml",
				"users.yml",
				"posts.yml",
				"post_medias.yml",
				"images.yml",
				"image_metadatas.yml",
			},
			args: args{
				image: models.Image{
					ID:      uuid.MustParse("00000000-0000-0000-0000-000000000009"),
					MediaID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				},
			},
			wantImage: models.Image{},
			wantImages: []models.Image{
				{
					ID:      uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					MediaID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					ImageMetadata: []models.ImageMetadata{
						{
							ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
							ImageID:       uuid.MustParse("00000000-0000-0000-0000-000000000001"),
							Height:        1080,
							Width:         1920,
							Url:           "https://example.com/image1.jpg",
							CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
							CreatedAtUnix: 1725091100,
							UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
						},
						{
							ID:            uuid.MustParse("00000000-0000-0000-0000-000000000002"),
							ImageID:       uuid.MustParse("00000000-0000-0000-0000-000000000001"),
							Height:        720,
							Width:         1280,
							Url:           "https://example.com/image2.jpg",
							CreatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
							CreatedAtUnix: 1725091101,
							UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
						},
					},
				},
			},
			wantErr: mediarepo.ErrNotFound,
		},
		{
			name: "media is not present in parent table :NEG",
			fixtureFiles: []string{
				"topics.yml",
				"voxspheres.yml",
				"users.yml",
				"posts.yml",
				"post_medias.yml",
				"images.yml",
				"image_metadatas.yml",
			},
			args: args{
				image: models.Image{
					ID:      uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					MediaID: uuid.MustParse("00000000-0000-0000-0000-000000000009"),
				},
			},
			wantImage: models.Image{},
			wantImages: []models.Image{
				{
					ID:      uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					MediaID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					ImageMetadata: []models.ImageMetadata{
						{
							ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
							ImageID:       uuid.MustParse("00000000-0000-0000-0000-000000000001"),
							Height:        1080,
							Width:         1920,
							Url:           "https://example.com/image1.jpg",
							CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
							CreatedAtUnix: 1725091100,
							UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
						},
						{
							ID:            uuid.MustParse("00000000-0000-0000-0000-000000000002"),
							ImageID:       uuid.MustParse("00000000-0000-0000-0000-000000000001"),
							Height:        720,
							Width:         1280,
							Url:           "https://example.com/image2.jpg",
							CreatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
							CreatedAtUnix: 1725091101,
							UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
						},
					},
				},
			},
			wantErr: mediarepo.ErrParentTableRecordNotFound,
		},
		{
			name: "update image :POS",
			fixtureFiles: []string{
				"topics.yml",
				"voxspheres.yml",
				"users.yml",
				"posts.yml",
				"post_medias.yml",
				"images.yml",
				"image_metadatas.yml",
			},
			args: args{
				image: models.Image{
					ID:      uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					MediaID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				},
			},
			wantImage: models.Image{
				ID:      uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				MediaID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
			},
			wantImages: []models.Image{
				{
					ID:      uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					MediaID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					ImageMetadata: []models.ImageMetadata{
						{
							ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
							ImageID:       uuid.MustParse("00000000-0000-0000-0000-000000000001"),
							Height:        1080,
							Width:         1920,
							Url:           "https://example.com/image1.jpg",
							CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
							CreatedAtUnix: 1725091100,
							UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
						},
						{
							ID:            uuid.MustParse("00000000-0000-0000-0000-000000000002"),
							ImageID:       uuid.MustParse("00000000-0000-0000-0000-000000000001"),
							Height:        720,
							Width:         1280,
							Url:           "https://example.com/image2.jpg",
							CreatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
							CreatedAtUnix: 1725091101,
							UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
						},
					},
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupPostgres(t, tt.fixtureFiles...)
			pgrepo := mediarepo.NewRepo(db)

			gotImage, gotErr := pgrepo.UpdateImage(context.Background(), tt.args.image)

			assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")
			assert.Equal(t, tt.wantImage, gotImage, "expect image to match")

			gotImages, err := pgrepo.Images(context.Background())

			assert.NoError(t, err, "expect no error while getting images")
			assertImages(t, tt.wantImages, gotImages)
		})
	}
}

func TestRepo_DeleteImage(t *testing.T) {
	type args struct {
		ID uuid.UUID
	}
	tests := []struct {
		name         string
		fixtureFiles []string
		args         args
		wantImages   []models.Image
		wantErr      error
	}{
		{
			name: "image not found :NEG",
			fixtureFiles: []string{
				"topics.yml",
				"voxspheres.yml",
				"users.yml",
				"posts.yml",
				"post_medias.yml",
				"images.yml",
				"image_metadatas.yml",
			},
			args: args{
				ID: uuid.MustParse("00000000-0000-0000-0000-000000000009"),
			},
			wantImages: []models.Image{
				{
					ID:      uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					MediaID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					ImageMetadata: []models.ImageMetadata{
						{
							ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
							ImageID:       uuid.MustParse("00000000-0000-0000-0000-000000000001"),
							Height:        1080,
							Width:         1920,
							Url:           "https://example.com/image1.jpg",
							CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
							CreatedAtUnix: 1725091100,
							UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
						},
						{
							ID:            uuid.MustParse("00000000-0000-0000-0000-000000000002"),
							ImageID:       uuid.MustParse("00000000-0000-0000-0000-000000000001"),
							Height:        720,
							Width:         1280,
							Url:           "https://example.com/image2.jpg",
							CreatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
							CreatedAtUnix: 1725091101,
							UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
						},
					},
				},
			},
			wantErr: mediarepo.ErrNotFound,
		},
		{
			name: "image deleted :POS",
			fixtureFiles: []string{
				"topics.yml",
				"voxspheres.yml",
				"users.yml",
				"posts.yml",
				"post_medias.yml",
				"images.yml",
				"image_metadatas.yml",
			},
			args: args{
				ID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			},
			wantImages: nil,
			wantErr:    nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupPostgres(t, tt.fixtureFiles...)
			pgrepo := mediarepo.NewRepo(db)

			gotErr := pgrepo.DeleteImage(context.Background(), tt.args.ID)

			assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")

			gotImages, err := pgrepo.Images(context.Background())

			assert.NoError(t, err, "expect no error while getting images")
			assertImages(t, tt.wantImages, gotImages)
		})
	}
}

func TestRepo_ImageForeignKeyCascade(t *testing.T) {
	t.Run("on deleting media from parent table , no child references should exist in images table", func(t *testing.T) {
		db := setupPostgres(
			t,
			"topics.yml",
			"voxspheres.yml",
			"users.yml",
			"posts.yml",
			"post_medias.yml",
			"images.yml",
			"image_metadatas.yml",
		)
		pgrepo := mediarepo.NewRepo(db)

		err := pgrepo.DeletePostMedia(context.Background(), uuid.MustParse("00000000-0000-0000-0000-000000000001"))

		assert.NoError(t, err, "expect no error while deleting post media")

		gotImages, err := pgrepo.Images(context.Background())

		assert.NoError(t, err, "expect no error while getting images")
		assertImages(t, nil, gotImages)
	})
}

func TestRepo_ImageMetadatas(t *testing.T) {
	tests := []struct {
		name               string
		fixtureFiles       []string
		wantImageMetadatas []models.ImageMetadata
		wantErr            error
	}{
		{
			name: "image metadatas :POS",
			fixtureFiles: []string{
				"topics.yml",
				"voxspheres.yml",
				"users.yml",
				"posts.yml",
				"post_medias.yml",
				"images.yml",
				"image_metadatas.yml",
			},
			wantImageMetadatas: []models.ImageMetadata{
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					ImageID:       uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Height:        1080,
					Width:         1920,
					Url:           "https://example.com/image1.jpg",
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix: 1725091100,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					ImageID:       uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Height:        720,
					Width:         1280,
					Url:           "https://example.com/image2.jpg",
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
					CreatedAtUnix: 1725091101,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
				},
			},
			wantErr: nil,
		},
		{
			name:               "no image metadatas :POS",
			fixtureFiles:       []string{},
			wantImageMetadatas: nil,
			wantErr:            nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupPostgres(t, tt.fixtureFiles...)
			pgrepo := mediarepo.NewRepo(db)

			gotImageMetadatas, gotErr := pgrepo.ImageMetadatas(context.Background())

			assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")
			assertImageMetadatasWithTimestamp(t, tt.wantImageMetadatas, gotImageMetadatas)
		})
	}
}

func TestRepo_ImageMetadataByID(t *testing.T) {
	type args struct {
		ID uuid.UUID
	}
	tests := []struct {
		name              string
		fixtureFiles      []string
		args              args
		wantImageMetadata models.ImageMetadata
		wantErr           error
	}{
		{
			name: "image metadata not found :NEG",
			fixtureFiles: []string{
				"topics.yml",
				"voxspheres.yml",
				"users.yml",
				"posts.yml",
				"post_medias.yml",
				"images.yml",
				"image_metadatas.yml",
			},
			args: args{
				ID: uuid.MustParse("00000000-0000-0000-0000-000000000009"),
			},
			wantImageMetadata: models.ImageMetadata{},
			wantErr:           mediarepo.ErrNotFound,
		},
		{
			name: "image found :POS",
			fixtureFiles: []string{
				"topics.yml",
				"voxspheres.yml",
				"users.yml",
				"posts.yml",
				"post_medias.yml",
				"images.yml",
				"image_metadatas.yml",
			},
			args: args{
				ID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			},
			wantImageMetadata: models.ImageMetadata{
				ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				ImageID:       uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				Height:        1080,
				Width:         1920,
				Url:           "https://example.com/image1.jpg",
				CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				CreatedAtUnix: 1725091100,
				UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupPostgres(t, tt.fixtureFiles...)
			pgrepo := mediarepo.NewRepo(db)

			gotImageMetadata, gotErr := pgrepo.ImageMetadataByID(context.Background(), tt.args.ID)

			assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")
			assert.Equal(t, tt.wantImageMetadata, gotImageMetadata, "expect image metadata to match")
		})
	}
}

func TestRepo_AddImageMetadata(t *testing.T) {
	type args struct {
		imageMetadata models.ImageMetadata
	}
	tests := []struct {
		name               string
		fixtureFiles       []string
		args               args
		wantImageMetadata  models.ImageMetadata
		wantImageMetadatas []models.ImageMetadata
		wantErr            error
	}{
		{
			name: "duplicte image metadata :NEG",
			fixtureFiles: []string{
				"topics.yml",
				"voxspheres.yml",
				"users.yml",
				"posts.yml",
				"post_medias.yml",
				"images.yml",
				"image_metadatas.yml",
			},
			args: args{
				imageMetadata: models.ImageMetadata{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					ImageID:       uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Height:        2160,
					Width:         3840,
					Url:           "https://example.com/image3.jpg",
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 30, 0, time.UTC),
					CreatedAtUnix: 1725091300,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 30, 0, time.UTC),
				},
			},
			wantImageMetadata: models.ImageMetadata{},
			wantImageMetadatas: []models.ImageMetadata{
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					ImageID:       uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Height:        1080,
					Width:         1920,
					Url:           "https://example.com/image1.jpg",
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix: 1725091100,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					ImageID:       uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Height:        720,
					Width:         1280,
					Url:           "https://example.com/image2.jpg",
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
					CreatedAtUnix: 1725091101,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
				},
			},
			wantErr: mediarepo.ErrDuplicateID,
		},
		{
			name: "image not present in parent table :NEG",
			fixtureFiles: []string{
				"topics.yml",
				"voxspheres.yml",
				"users.yml",
				"posts.yml",
				"post_medias.yml",
				"images.yml",
				"image_metadatas.yml",
			},
			args: args{
				imageMetadata: models.ImageMetadata{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					ImageID:       uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					Height:        2160,
					Width:         3840,
					Url:           "https://example.com/image3.jpg",
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 30, 0, time.UTC),
					CreatedAtUnix: 1725091300,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 30, 0, time.UTC),
				},
			},
			wantImageMetadata: models.ImageMetadata{},
			wantImageMetadatas: []models.ImageMetadata{
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					ImageID:       uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Height:        1080,
					Width:         1920,
					Url:           "https://example.com/image1.jpg",
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix: 1725091100,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					ImageID:       uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Height:        720,
					Width:         1280,
					Url:           "https://example.com/image2.jpg",
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
					CreatedAtUnix: 1725091101,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
				},
			},
			wantErr: mediarepo.ErrParentTableRecordNotFound,
		},
		{
			name: "add image :POS",
			fixtureFiles: []string{
				"topics.yml",
				"voxspheres.yml",
				"users.yml",
				"posts.yml",
				"post_medias.yml",
				"images.yml",
				"image_metadatas.yml",
			},
			args: args{
				imageMetadata: models.ImageMetadata{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					ImageID:       uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Height:        2160,
					Width:         3840,
					Url:           "https://example.com/image3.jpg",
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 30, 0, time.UTC),
					CreatedAtUnix: 1725091300,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 30, 0, time.UTC),
				},
			},
			wantImageMetadata: models.ImageMetadata{
				ID:            uuid.MustParse("00000000-0000-0000-0000-000000000003"),
				ImageID:       uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				Height:        2160,
				Width:         3840,
				Url:           "https://example.com/image3.jpg",
				CreatedAt:     time.Date(2024, 10, 10, 10, 10, 30, 0, time.UTC),
				CreatedAtUnix: 1725091300,
				UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 30, 0, time.UTC),
			},
			wantImageMetadatas: []models.ImageMetadata{
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					ImageID:       uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Height:        1080,
					Width:         1920,
					Url:           "https://example.com/image1.jpg",
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix: 1725091100,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					ImageID:       uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Height:        720,
					Width:         1280,
					Url:           "https://example.com/image2.jpg",
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
					CreatedAtUnix: 1725091101,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
				},
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					ImageID:       uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Height:        2160,
					Width:         3840,
					Url:           "https://example.com/image3.jpg",
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 30, 0, time.UTC),
					CreatedAtUnix: 1725091300,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 30, 0, time.UTC),
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupPostgres(t, tt.fixtureFiles...)
			pgrepo := mediarepo.NewRepo(db)

			startTime := time.Now()
			gotImageMetadata, gotErr := pgrepo.AddImageMetadata(context.Background(), tt.args.imageMetadata)
			endTime := time.Now()

			assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")
			assert.Equal(
				t,
				gotImageMetadata.UpdatedAt,
				gotImageMetadata.CreatedAt,
				"expect CreatedAt and UpdatedAt to be same",
			)
			if tt.wantErr == nil {
				assertTimeWithinRange(t, gotImageMetadata.CreatedAt, startTime, endTime)
				assertTimeWithinRange(t, gotImageMetadata.UpdatedAt, startTime, endTime)
			}

			gotImageMetadatas, err := pgrepo.ImageMetadatas(context.Background())

			assert.NoError(t, err, "expect no error while getting image metadatas")
			assertImageMetadatasWithoutTimestamp(t, tt.wantImageMetadatas, gotImageMetadatas)
		})
	}
}

func TestRepo_UpdateImageMetadata(t *testing.T) {
	type args struct {
		imageMetadata models.ImageMetadata
	}
	tests := []struct {
		name               string
		fixtureFiles       []string
		args               args
		wantImageMetadata  models.ImageMetadata
		wantImageMetadatas []models.ImageMetadata
		wantErr            error
	}{
		{
			name: "image metadata not found :NEG",
			fixtureFiles: []string{
				"topics.yml",
				"voxspheres.yml",
				"users.yml",
				"posts.yml",
				"post_medias.yml",
				"images.yml",
				"image_metadatas.yml",
			},
			args: args{
				imageMetadata: models.ImageMetadata{
					ID:      uuid.MustParse("00000000-0000-0000-0000-000000000009"),
					ImageID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Height:  360,
					Width:   640,
					Url:     "https://example.com/image2-updated.jpg",
				},
			},
			wantImageMetadata: models.ImageMetadata{},
			wantImageMetadatas: []models.ImageMetadata{
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					ImageID:       uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Height:        1080,
					Width:         1920,
					Url:           "https://example.com/image1.jpg",
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix: 1725091100,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					ImageID:       uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Height:        720,
					Width:         1280,
					Url:           "https://example.com/image2.jpg",
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
					CreatedAtUnix: 1725091101,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
				},
			},
			wantErr: mediarepo.ErrNotFound,
		},
		{
			name: "image not found in parent table :NEG",
			fixtureFiles: []string{
				"topics.yml",
				"voxspheres.yml",
				"users.yml",
				"posts.yml",
				"post_medias.yml",
				"images.yml",
				"image_metadatas.yml",
			},
			args: args{
				imageMetadata: models.ImageMetadata{
					ID:      uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					ImageID: uuid.MustParse("00000000-0000-0000-0000-000000000009"),
					Height:  360,
					Width:   640,
					Url:     "https://example.com/image2-updated.jpg",
				},
			},
			wantImageMetadata: models.ImageMetadata{},
			wantImageMetadatas: []models.ImageMetadata{
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					ImageID:       uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Height:        1080,
					Width:         1920,
					Url:           "https://example.com/image1.jpg",
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix: 1725091100,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					ImageID:       uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Height:        720,
					Width:         1280,
					Url:           "https://example.com/image2.jpg",
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
					CreatedAtUnix: 1725091101,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
				},
			},
			wantErr: mediarepo.ErrParentTableRecordNotFound,
		},
		{
			name: "update image metadata :POS",
			fixtureFiles: []string{
				"topics.yml",
				"voxspheres.yml",
				"users.yml",
				"posts.yml",
				"post_medias.yml",
				"images.yml",
				"image_metadatas.yml",
			},
			args: args{
				imageMetadata: models.ImageMetadata{
					ID:      uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					ImageID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Height:  360,
					Width:   640,
					Url:     "https://example.com/image2-updated.jpg",
				},
			},
			wantImageMetadata: models.ImageMetadata{
				ID:      uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				ImageID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				Height:  360,
				Width:   640,
				Url:     "https://example.com/image2-updated.jpg",
			},
			wantImageMetadatas: []models.ImageMetadata{
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					ImageID:       uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Height:        720,
					Width:         1280,
					Url:           "https://example.com/image2.jpg",
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
					CreatedAtUnix: 1725091101,
				},
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					ImageID:       uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Height:        360,
					Width:         640,
					Url:           "https://example.com/image2-updated.jpg",
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix: 1725091100,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupPostgres(t, tt.fixtureFiles...)
			pgrepo := mediarepo.NewRepo(db)

			startTime := time.Now()
			gotImageMetadata, gotErr := pgrepo.UpdateImageMetadata(context.Background(), tt.args.imageMetadata)
			endTime := time.Now()

			assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")
			if tt.wantErr == nil {
				assertTimeWithinRange(t, gotImageMetadata.UpdatedAt, startTime, endTime)
			}

			gotImageMetadatas, err := pgrepo.ImageMetadatas(context.Background())

			assert.NoError(t, err, "expect no error while getting image metdatas")
			assertImageMetadatasWithoutTimestamp(t, tt.wantImageMetadatas, gotImageMetadatas)
		})
	}
}

func TestRepo_DeleteImageMetadata(t *testing.T) {
	type args struct {
		ID uuid.UUID
	}
	tests := []struct {
		name               string
		fixtureFiles       []string
		args               args
		wantImageMetadatas []models.ImageMetadata
		wantErr            error
	}{
		{
			name: "image metadata not found :NEG",
			fixtureFiles: []string{
				"topics.yml",
				"voxspheres.yml",
				"users.yml",
				"posts.yml",
				"post_medias.yml",
				"images.yml",
				"image_metadatas.yml",
			},
			args: args{
				ID: uuid.MustParse("00000000-0000-0000-0000-000000000009"),
			},
			wantImageMetadatas: []models.ImageMetadata{
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					ImageID:       uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Height:        1080,
					Width:         1920,
					Url:           "https://example.com/image1.jpg",
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix: 1725091100,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					ImageID:       uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Height:        720,
					Width:         1280,
					Url:           "https://example.com/image2.jpg",
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
					CreatedAtUnix: 1725091101,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
				},
			},
			wantErr: mediarepo.ErrNotFound,
		},
		{
			name: "image metadata deleted :POS",
			fixtureFiles: []string{
				"topics.yml",
				"voxspheres.yml",
				"users.yml",
				"posts.yml",
				"post_medias.yml",
				"images.yml",
				"image_metadatas.yml",
			},
			args: args{
				ID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			},
			wantImageMetadatas: []models.ImageMetadata{
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					ImageID:       uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Height:        720,
					Width:         1280,
					Url:           "https://example.com/image2.jpg",
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
					CreatedAtUnix: 1725091101,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupPostgres(t, tt.fixtureFiles...)
			pgrepo := mediarepo.NewRepo(db)

			gotErr := pgrepo.DeleteImageMetadata(context.Background(), tt.args.ID)

			assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")

			gotImageMetadatas, err := pgrepo.ImageMetadatas(context.Background())

			assert.NoError(t, err, "expect no error while getting image metadatas")
			assert.Equal(t, tt.wantImageMetadatas, gotImageMetadatas, "expect image metadatas to match")
		})
	}
}

func TestRepo_ImageMetadataForeignKeyCascade(t *testing.T) {
	t.Run("on deleting image from parent table , no child references should exist in image_metadatas table", func(t *testing.T) {
		db := setupPostgres(
			t,
			"topics.yml",
			"voxspheres.yml",
			"users.yml",
			"posts.yml",
			"post_medias.yml",
			"images.yml",
			"image_metadatas.yml",
		)
		pgrepo := mediarepo.NewRepo(db)

		err := pgrepo.DeleteImage(context.Background(), uuid.MustParse("00000000-0000-0000-0000-000000000001"))

		assert.NoError(t, err, "expect no error while deleting images")

		gotImageMetadatas, err := pgrepo.ImageMetadatas(context.Background())

		assert.NoError(t, err, "expect no error while getting image metadatas")
		assertImageMetadatasWithTimestamp(t, nil, gotImageMetadatas)
	})
}

func TestRepo_Links(t *testing.T) {
	tests := []struct {
		name         string
		fixtureFiles []string
		wantLinks    []models.Link
		wantErr      error
	}{
		{
			name: "links :POS",
			fixtureFiles: []string{
				"topics.yml",
				"voxspheres.yml",
				"users.yml",
				"posts.yml",
				"post_medias.yml",
				"links.yml",
			},
			wantLinks: []models.Link{
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					MediaID:       uuid.MustParse("00000000-0000-0000-0000-000000000005"),
					Link:          "https://example.com/video.mp4",
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix: 1725091100,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
			},
			wantErr: nil,
		},
		{
			name:         "no links :POS",
			fixtureFiles: []string{},
			wantLinks:    []models.Link{},
			wantErr:      nil,
		},
	}
	for _, tt := range tests {
		db := setupPostgres(t, tt.fixtureFiles...)
		pgrepo := mediarepo.NewRepo(db)

		gotLinks, gotErr := pgrepo.Links(context.Background())

		assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")
		assertLinksWithoutTimestamp(t, tt.wantLinks, gotLinks)
	}
}

func TestRepo_LinkByID(t *testing.T) {
	type args struct {
		ID uuid.UUID
	}
	tests := []struct {
		name         string
		fixtureFiles []string
		args         args
		wantLink     models.Link
		wantErr      error
	}{
		{
			name: "link not found :NEG",
			fixtureFiles: []string{
				"topics.yml",
				"voxspheres.yml",
				"users.yml",
				"posts.yml",
				"post_medias.yml",
				"links.yml",
			},
			args: args{
				ID: uuid.MustParse("00000000-0000-0000-0000-000000000009"),
			},
			wantLink: models.Link{},
			wantErr:  mediarepo.ErrNotFound,
		},
		{
			name: "link found :POS",
			fixtureFiles: []string{
				"topics.yml",
				"voxspheres.yml",
				"users.yml",
				"posts.yml",
				"post_medias.yml",
				"links.yml",
			},
			args: args{
				ID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			},
			wantLink: models.Link{
				ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				MediaID:       uuid.MustParse("00000000-0000-0000-0000-000000000005"),
				Link:          "https://example.com/video.mp4",
				CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				CreatedAtUnix: 1725091100,
				UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupPostgres(t, tt.fixtureFiles...)
			pgrepo := mediarepo.NewRepo(db)

			gotLink, gotErr := pgrepo.LinkByID(context.Background(), tt.args.ID)

			assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")
			assert.Equal(t, tt.wantLink, gotLink, "expect link to match")
		})
	}
}

func TestRepo_AddLink(t *testing.T) {
	type args struct {
		link models.Link
	}
	tests := []struct {
		name         string
		fixtureFiles []string
		args         args
		wantLink     models.Link
		wantLinks    []models.Link
		wantErr      error
	}{
		{
			name: "duplicte link :NEG",
			fixtureFiles: []string{
				"topics.yml",
				"voxspheres.yml",
				"users.yml",
				"posts.yml",
				"post_medias.yml",
				"links.yml",
			},
			args: args{
				link: models.Link{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					MediaID:       uuid.MustParse("00000000-0000-0000-0000-000000000005"),
					Link:          "https://example.com/video.mp4",
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix: 1725091100,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
			},
			wantLink: models.Link{},
			wantLinks: []models.Link{
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					MediaID:       uuid.MustParse("00000000-0000-0000-0000-000000000005"),
					Link:          "https://example.com/video.mp4",
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix: 1725091100,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
			},
			wantErr: mediarepo.ErrDuplicateID,
		},
		{
			name: "post media is not present in parent table :NEG",
			fixtureFiles: []string{
				"topics.yml",
				"voxspheres.yml",
				"users.yml",
				"posts.yml",
				"post_medias.yml",
				"links.yml",
			},
			args: args{
				link: models.Link{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					MediaID:       uuid.MustParse("00000000-0000-0000-0000-000000000009"),
					Link:          "https://example.com/video.mp4",
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix: 1725091100,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
			},
			wantLink: models.Link{},
			wantLinks: []models.Link{
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					MediaID:       uuid.MustParse("00000000-0000-0000-0000-000000000005"),
					Link:          "https://example.com/video.mp4",
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix: 1725091100,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
			},
			wantErr: mediarepo.ErrParentTableRecordNotFound,
		},
		{
			name: "add link :POS",
			fixtureFiles: []string{
				"topics.yml",
				"voxspheres.yml",
				"users.yml",
				"posts.yml",
				"post_medias.yml",
				"links.yml",
			},
			args: args{
				link: models.Link{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					MediaID:       uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Link:          "https://example.com/videonew.mp4",
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix: 1725091100,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
			},
			wantLink: models.Link{
				ID:            uuid.MustParse("00000000-0000-0000-0000-000000000003"),
				MediaID:       uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				Link:          "https://example.com/videonew.mp4",
				CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				CreatedAtUnix: 1725091100,
				UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
			},
			wantLinks: []models.Link{
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					MediaID:       uuid.MustParse("00000000-0000-0000-0000-000000000005"),
					Link:          "https://example.com/video.mp4",
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix: 1725091100,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					MediaID:       uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Link:          "https://example.com/videonew.mp4",
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix: 1725091100,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupPostgres(t, tt.fixtureFiles...)
			pgrepo := mediarepo.NewRepo(db)

			startTime := time.Now()
			gotLink, gotErr := pgrepo.AddLink(context.Background(), tt.args.link)
			endTime := time.Now()

			assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")
			assert.Equal(
				t,
				gotLink.UpdatedAt,
				gotLink.CreatedAt,
				"expect CreatedAt and UpdatedAt to be same",
			)
			if tt.wantErr == nil {
				assertTimeWithinRange(t, gotLink.CreatedAt, startTime, endTime)
				assertTimeWithinRange(t, gotLink.UpdatedAt, startTime, endTime)
			}

			gotLinks, err := pgrepo.Links(context.Background())

			assert.NoError(t, err, "expect no error while getting links")
			assertLinksWithoutTimestamp(t, tt.wantLinks, gotLinks)
		})
	}
}

func TestRepo_UpdateLink(t *testing.T) {
	type args struct {
		link models.Link
	}
	tests := []struct {
		name         string
		fixtureFiles []string
		args         args
		wantLink     models.Link
		wantLinks    []models.Link
		wantErr      error
	}{
		{
			name: "link not found :NEG",
			fixtureFiles: []string{
				"topics.yml",
				"voxspheres.yml",
				"users.yml",
				"posts.yml",
				"post_medias.yml",
				"links.yml",
			},
			args: args{
				link: models.Link{
					ID:      uuid.MustParse("00000000-0000-0000-0000-000000000009"),
					MediaID: uuid.MustParse("00000000-0000-0000-0000-000000000005"),
					Link:    "https://xyz.com",
				},
			},
			wantLink: models.Link{},
			wantLinks: []models.Link{
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					MediaID:       uuid.MustParse("00000000-0000-0000-0000-000000000005"),
					Link:          "https://example.com/video.mp4",
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix: 1725091100,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
			},
			wantErr: mediarepo.ErrNotFound,
		},
		{
			name: "media is not present in parent table :NEG",
			fixtureFiles: []string{
				"topics.yml",
				"voxspheres.yml",
				"users.yml",
				"posts.yml",
				"post_medias.yml",
				"links.yml",
			},
			args: args{
				link: models.Link{
					ID:      uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					MediaID: uuid.MustParse("00000000-0000-0000-0000-000000000009"),
					Link:    "https://xyz.com",
				},
			},
			wantLink: models.Link{},
			wantLinks: []models.Link{
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					MediaID:       uuid.MustParse("00000000-0000-0000-0000-000000000005"),
					Link:          "https://example.com/video.mp4",
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix: 1725091100,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
			},
			wantErr: mediarepo.ErrParentTableRecordNotFound,
		},
		{
			name: "update media :POS",
			fixtureFiles: []string{
				"topics.yml",
				"voxspheres.yml",
				"users.yml",
				"posts.yml",
				"post_medias.yml",
				"links.yml",
			},
			args: args{
				link: models.Link{
					ID:      uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					MediaID: uuid.MustParse("00000000-0000-0000-0000-000000000005"),
					Link:    "https://xyz.com",
				},
			},
			wantLink: models.Link{
				ID:      uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				MediaID: uuid.MustParse("00000000-0000-0000-0000-000000000005"),
				Link:    "https://xyz.com",
			},
			wantLinks: []models.Link{
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					MediaID:       uuid.MustParse("00000000-0000-0000-0000-000000000005"),
					Link:          "https://xyz.com",
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix: 1725091100,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupPostgres(t, tt.fixtureFiles...)
			pgrepo := mediarepo.NewRepo(db)

			startTime := time.Now()
			gotLink, gotErr := pgrepo.UpdateLink(context.Background(), tt.args.link)
			endTime := time.Now()

			assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")
			if tt.wantErr == nil {
				assertTimeWithinRange(t, gotLink.UpdatedAt, startTime, endTime)
			}

			gotLinks, err := pgrepo.Links(context.Background())

			assert.NoError(t, err, "expect no error while getting links")
			assertLinksWithoutTimestamp(t, tt.wantLinks, gotLinks)
		})
	}
}

func TestRepo_DeleteLink(t *testing.T) {
	type args struct {
		ID uuid.UUID
	}
	tests := []struct {
		name         string
		fixtureFiles []string
		args         args
		wantLinks    []models.Link
		wantErr      error
	}{
		{
			name:         "link not found :NEG",
			fixtureFiles: []string{"topics.yml", "voxspheres.yml", "users.yml", "posts.yml", "post_medias.yml", "links.yml"},
			args: args{
				ID: uuid.MustParse("00000000-0000-0000-0000-000000000009"),
			},
			wantLinks: []models.Link{
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					MediaID:       uuid.MustParse("00000000-0000-0000-0000-000000000005"),
					Link:          "https://example.com/video.mp4",
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix: 1725091100,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
			},
			wantErr: mediarepo.ErrNotFound,
		},
		{
			name:         "link deleted :POS",
			fixtureFiles: []string{"topics.yml", "voxspheres.yml", "users.yml", "posts.yml", "post_medias.yml", "links.yml"},
			args: args{
				ID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			},
			wantLinks: nil,
			wantErr:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupPostgres(t, tt.fixtureFiles...)
			pgrepo := mediarepo.NewRepo(db)

			gotErr := pgrepo.DeleteLink(context.Background(), tt.args.ID)

			assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")

			gotLinks, err := pgrepo.Links(context.Background())

			assert.NoError(t, err, "expect no error while getting links")
			assert.Equal(t, tt.wantLinks, gotLinks, "expect links to match")
		})
	}
}

func TestRepo_LinkForeignKeyCascade(t *testing.T) {
	t.Run("on deleting media from parent table , no child references should exist in links table", func(t *testing.T) {
		db := setupPostgres(
			t,
			"topics.yml",
			"voxspheres.yml",
			"users.yml",
			"posts.yml",
			"post_medias.yml",
			"links.yml",
		)
		pgrepo := mediarepo.NewRepo(db)

		err := pgrepo.DeletePostMedia(context.Background(), uuid.MustParse("00000000-0000-0000-0000-000000000005"))

		assert.NoError(t, err, "expect no error while deleting post media")

		gotLinks, err := pgrepo.Links(context.Background())

		assert.NoError(t, err, "expect no error while getting links")
		assert.Equal(t, []models.Link(nil), gotLinks, "expect links to match")
	})
}
