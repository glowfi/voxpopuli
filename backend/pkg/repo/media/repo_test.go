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

func assertGifMetadataWithoutTimestamp(t *testing.T, wantGifMetadata, gotGifMetadata models.GifMetadata) {
	assert.Equal(t, wantGifMetadata.ID, gotGifMetadata.ID, "expected id to match")
	assert.Equal(t, wantGifMetadata.GifID, gotGifMetadata.GifID, "expected gif id to match")
	assert.Equal(t, wantGifMetadata.Height, gotGifMetadata.Height, "expected height to match")
	assert.Equal(t, wantGifMetadata.Width, gotGifMetadata.Width, "expected width to match")
	assert.Equal(t, wantGifMetadata.Url, gotGifMetadata.Url, "expected url to match")
}

func assertGifMetadatasWithoutTimestamp(t *testing.T, wantGifMetadatas, gotGifMetadatas []models.GifMetadata) {
	t.Helper()

	if len(wantGifMetadatas) != len(gotGifMetadatas) {
		t.Fatal("length of wantGifMetadatas and gotGifMetadatas do not match")
	}

	for _, gifMetadata := range wantGifMetadatas {
		idx := slices.IndexFunc(gotGifMetadatas, func(gm models.GifMetadata) bool {
			return gm.ID == gifMetadata.ID
		})

		if idx == -1 {
			t.Fatal(fmt.Sprintf("gif metadata %v of ID %v is not present in gotGifMetadatas", gifMetadata.Url, gifMetadata.ID))
			return
		}
		assertGifMetadataWithoutTimestamp(t, gifMetadata, gotGifMetadatas[idx])
	}
}

func assertGifMetadatasWithTimestamp(t *testing.T, wantGifMetadatas, gotGifMetadatas []models.GifMetadata) {
	t.Helper()

	for _, gifMetadata := range wantGifMetadatas {
		idx := slices.IndexFunc(gotGifMetadatas, func(gm models.GifMetadata) bool {
			return gm.ID == gifMetadata.ID
		})

		if idx == -1 {
			t.Fatal(fmt.Sprintf("gif metadata %v of ID %v is not present in gotGifMetadatas", gifMetadata.Url, gifMetadata.ID))
			return
		}
		assert.Equal(t, gifMetadata, gotGifMetadatas[idx], "expected gif metadata to match")
	}
}

func assertGif(t *testing.T, wantGif, gotGif models.Gif) {
	assert.Equal(t, wantGif.ID, gotGif.ID, "expected id to match")
	assert.Equal(t, wantGif.MediaID, gotGif.MediaID, "expected media id to match")
	assertGifMetadatasWithTimestamp(t, wantGif.GifMetadata, gotGif.GifMetadata)
}

func assertGifs(t *testing.T, wantGifs, gotGifs []models.Gif) {
	t.Helper()

	if len(wantGifs) != len(gotGifs) {
		t.Fatal("length of wantGifs and gotGifs do not match")
	}

	for _, gif := range wantGifs {
		idx := slices.IndexFunc(gotGifs, func(g models.Gif) bool {
			return g.ID == gif.ID
		})

		if idx == -1 {
			t.Fatal(fmt.Sprintf("gif %v of ID %v is not present in gotGifs", gif.MediaID, gif.ID))
			return
		}
		assertGif(t, gif, gotGifs[idx])
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

func assertLinksWitTimestamp(t *testing.T, wantLinks, gotLinks []models.Link) {
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
		assert.Equal(t, link, gotLinks[idx], "expect link to match")
	}
}

func assertVideoWithoutTimestamp(t *testing.T, wantVideo, gotVideo models.Video) {
	assert.Equal(t, wantVideo.ID, gotVideo.ID, "expect video id to match")
	assert.Equal(t, wantVideo.MediaID, gotVideo.MediaID, "expect video media id to match")
	assert.Equal(t, wantVideo.Url, gotVideo.Url, "expect video url to match")
	assert.Equal(t, wantVideo.Height, gotVideo.Height, "expect video height to match")
	assert.Equal(t, wantVideo.Width, gotVideo.Width, "expect video width to match")
}

func assertVideosWithoutTimestamp(t *testing.T, wantVideos, gotVideos []models.Video) {
	t.Helper()

	if len(wantVideos) != len(gotVideos) {
		t.Fatal("length of wantVideos and gotVideos do not match")
	}

	for _, video := range wantVideos {
		idx := slices.IndexFunc(gotVideos, func(im models.Video) bool {
			return im.ID == video.ID
		})

		if idx == -1 {
			t.Fatal(fmt.Sprintf("video %v of url %v is not present in gotVideos", video.ID, video.Url))
			return
		}
		assertVideoWithoutTimestamp(t, video, gotVideos[idx])
	}
}

func assertVideosWithTimestamp(t *testing.T, wantVideos, gotVideos []models.Video) {
	t.Helper()

	if len(wantVideos) != len(gotVideos) {
		t.Fatal("length of wantVideos and gotVideos do not match")
	}

	for _, video := range wantVideos {
		idx := slices.IndexFunc(gotVideos, func(im models.Video) bool {
			return im.ID == video.ID
		})

		if idx == -1 {
			t.Fatal(fmt.Sprintf("video %v of url %v is not present in gotVideos", video.ID, video.Url))
			return
		}
		assert.Equal(t, video, gotVideos[idx], "expect video to match")
	}
}

func assertGalleryMetadataWithoutTimestamp(t *testing.T, wantGalleryMetadata, gotGalleryMetadata models.GalleryMetadata) {
	assert.Equal(t, wantGalleryMetadata.ID, gotGalleryMetadata.ID, "expected id to match")
	assert.Equal(t, wantGalleryMetadata.GalleryID, gotGalleryMetadata.GalleryID, "expected gallery id to match")
	assert.Equal(t, wantGalleryMetadata.OrderIndex, gotGalleryMetadata.OrderIndex, "expected order index to match")
	assert.Equal(t, wantGalleryMetadata.Height, gotGalleryMetadata.Height, "expected height to match")
	assert.Equal(t, wantGalleryMetadata.Width, gotGalleryMetadata.Width, "expected width to match")
	assert.Equal(t, wantGalleryMetadata.Url, gotGalleryMetadata.Url, "expected url to match")
}

func assertGalleryMetadatasWithoutTimestamp(t *testing.T, wantGalleryMetadatas, gotGalleryMetadatas []models.GalleryMetadata) {
	t.Helper()

	if len(wantGalleryMetadatas) != len(gotGalleryMetadatas) {
		t.Fatal("length of wantGalleryMetadatas and gotGalleryMetadatas do not match")
	}

	for _, galleryMetadata := range wantGalleryMetadatas {
		idx := slices.IndexFunc(gotGalleryMetadatas, func(gm models.GalleryMetadata) bool {
			return gm.ID == galleryMetadata.ID
		})

		if idx == -1 {
			t.Fatal(fmt.Sprintf("gallery metadata %v of ID %v is not present in gotGalleryMetadatas", galleryMetadata.Url, galleryMetadata.ID))
			return
		}
		assertGalleryMetadataWithoutTimestamp(t, galleryMetadata, gotGalleryMetadatas[idx])
	}
}

func assertGalleryMetadatasWithTimestamp(t *testing.T, wantGalleryMetadatas, gotGalleryMetadatas []models.GalleryMetadata) {
	t.Helper()

	for _, galleryMetadata := range wantGalleryMetadatas {
		idx := slices.IndexFunc(gotGalleryMetadatas, func(gm models.GalleryMetadata) bool {
			return gm.ID == galleryMetadata.ID
		})

		if idx == -1 {
			t.Fatal(fmt.Sprintf("gallery metadata %v of ID %v is not present in gotGalleryMetadatas", galleryMetadata.Url, galleryMetadata.ID))
			return
		}
		assert.Equal(t, galleryMetadata, gotGalleryMetadatas[idx], "expected gallery metadata to match")
	}
}

func assertGallery(t *testing.T, wantGallery, gotGallery models.Gallery) {
	assert.Equal(t, wantGallery.ID, gotGallery.ID, "expected id to match")
	assert.Equal(t, wantGallery.MediaID, gotGallery.MediaID, "expected media id to match")
	assertGalleryMetadatasWithTimestamp(t, wantGallery.GalleryMetadata, gotGallery.GalleryMetadata)
}

func assertGalleries(t *testing.T, wantGalleries, gotGalleries []models.Gallery) {
	t.Helper()

	if len(wantGalleries) != len(gotGalleries) {
		t.Fatal("length of wantGalleries and gotGalleries do not match")
	}

	for _, gallery := range wantGalleries {
		idx := slices.IndexFunc(gotGalleries, func(g models.Gallery) bool {
			return g.ID == gallery.ID
		})

		if idx == -1 {
			t.Fatal(fmt.Sprintf("gallery %v of ID %v is not present in gotGalleries", gallery.MediaID, gallery.ID))
			return
		}
		assertGallery(t, gallery, gotGalleries[idx])
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

func TestRepo_AddPostMedias(t *testing.T) {
	type args struct {
		postMedias []models.PostMedia
	}
	tests := []struct {
		name                   string
		fixtureFiles           []string
		args                   args
		wantInsertedPostMedias []models.PostMedia
		wantPostMedias         []models.PostMedia
		wantErr                error
	}{
		{
			name:         "duplicate media id :NEG",
			fixtureFiles: []string{"topics.yml", "voxspheres.yml", "users.yml", "posts.yml", "post_medias.yml"},
			args: args{
				postMedias: []models.PostMedia{
					{
						ID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
						PostID:    uuid.MustParse("00000000-0000-0000-0000-000000000001"),
						MediaType: models.MediaTypeImage,
					},
				},
			},
			wantInsertedPostMedias: nil,
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
				postMedias: []models.PostMedia{
					{
						ID:        uuid.MustParse("00000000-0000-0000-0000-000000000006"),
						PostID:    uuid.MustParse("00000000-0000-0000-0000-000000000009"),
						MediaType: models.MediaTypeImage,
					},
				},
			},
			wantInsertedPostMedias: nil,
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
			name:         "add post medias :POS",
			fixtureFiles: []string{"topics.yml", "voxspheres.yml", "users.yml", "posts.yml", "post_medias.yml"},
			args: args{
				postMedias: []models.PostMedia{
					{
						ID:        uuid.MustParse("00000000-0000-0000-0000-000000000006"),
						PostID:    uuid.MustParse("00000000-0000-0000-0000-000000000001"),
						MediaType: models.MediaTypeImage,
					},
				},
			},
			wantInsertedPostMedias: []models.PostMedia{
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000006"),
					PostID:    uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					MediaType: models.MediaTypeImage,
				},
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

			gotInsertedPostMedias, gotErr := pgrepo.AddPostMedias(context.Background(), tt.args.postMedias...)

			assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")
			assert.Equal(t, tt.wantInsertedPostMedias, gotInsertedPostMedias, "expect inserted postMedias to match")

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

func TestRepo_AddImages(t *testing.T) {
	type args struct {
		images []models.Image
	}
	tests := []struct {
		name               string
		fixtureFiles       []string
		args               args
		wantInsertedImages []models.Image
		wantImages         []models.Image
		wantErr            error
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
				images: []models.Image{
					{
						ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
						MediaID:       uuid.MustParse("00000000-0000-0000-0000-000000000002"),
						ImageMetadata: []models.ImageMetadata{},
					},
				},
			},
			wantInsertedImages: nil,
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
				images: []models.Image{
					{
						ID:            uuid.MustParse("00000000-0000-0000-0000-000000000003"),
						MediaID:       uuid.MustParse("00000000-0000-0000-0000-000000000009"),
						ImageMetadata: []models.ImageMetadata{},
					},
				},
			},
			wantInsertedImages: nil,
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
				images: []models.Image{
					{
						ID:            uuid.MustParse("00000000-0000-0000-0000-000000000003"),
						MediaID:       uuid.MustParse("00000000-0000-0000-0000-000000000001"),
						ImageMetadata: []models.ImageMetadata{},
					},
				},
			},
			wantInsertedImages: []models.Image{
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					MediaID:       uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					ImageMetadata: []models.ImageMetadata{},
				},
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

		gotInsertedImages, gotErr := pgrepo.AddImages(context.Background(), tt.args.images...)

		assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")
		assert.Equal(t, tt.wantInsertedImages, gotInsertedImages, "expect inserted image to match")

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
			name: "image metadata found :POS",
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

func TestRepo_AddImageMetadatas(t *testing.T) {
	type args struct {
		imageMetadatas []models.ImageMetadata
	}
	tests := []struct {
		name                       string
		fixtureFiles               []string
		args                       args
		wantInsertedImageMetadatas []models.ImageMetadata
		wantImageMetadatas         []models.ImageMetadata
		wantErr                    error
	}{
		{
			name: "duplicate image metadata :NEG",
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
				imageMetadatas: []models.ImageMetadata{
					{
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
			},
			wantInsertedImageMetadatas: nil,
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
				imageMetadatas: []models.ImageMetadata{
					{
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
			},
			wantInsertedImageMetadatas: nil,
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
			name: "add image metadatas :POS",
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
				imageMetadatas: []models.ImageMetadata{
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
			},
			wantInsertedImageMetadatas: []models.ImageMetadata{
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
			gotInsertedImageMetadatas, gotErr := pgrepo.AddImageMetadatas(context.Background(), tt.args.imageMetadatas...)
			endTime := time.Now()

			for _, gotInsertedImageMetadata := range gotInsertedImageMetadatas {
				assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")
				assert.Equal(
					t,
					gotInsertedImageMetadata.UpdatedAt,
					gotInsertedImageMetadata.CreatedAt,
					"expect CreatedAt and UpdatedAt to be same",
				)
				if tt.wantErr == nil {
					assertTimeWithinRange(t, gotInsertedImageMetadata.CreatedAt, startTime, endTime)
					assertTimeWithinRange(t, gotInsertedImageMetadata.UpdatedAt, startTime, endTime)
				}
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
			assertImageMetadatasWithTimestamp(t, tt.wantImageMetadatas, gotImageMetadatas)
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

func TestRepo_Gifs(t *testing.T) {
	tests := []struct {
		name         string
		fixtureFiles []string
		wantGifs     []models.Gif
		wantErr      error
	}{
		{
			name: "gifs :POS",
			fixtureFiles: []string{
				"topics.yml",
				"voxspheres.yml",
				"users.yml",
				"posts.yml",
				"post_medias.yml",
				"gifs.yml",
				"gif_metadatas.yml",
			},
			wantGifs: []models.Gif{
				{
					ID:      uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					MediaID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					GifMetadata: []models.GifMetadata{
						{
							ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
							GifID:         uuid.MustParse("00000000-0000-0000-0000-000000000001"),
							Height:        1080,
							Width:         1920,
							Url:           "https://example.com/image1.jpg",
							CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
							CreatedAtUnix: 1725091100,
							UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
						},
						{
							ID:            uuid.MustParse("00000000-0000-0000-0000-000000000002"),
							GifID:         uuid.MustParse("00000000-0000-0000-0000-000000000001"),
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
			name:         "no gifs :POS",
			fixtureFiles: []string{},
			wantGifs:     nil,
			wantErr:      nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupPostgres(t, tt.fixtureFiles...)
			pgrepo := mediarepo.NewRepo(db)

			gotGifs, gotErr := pgrepo.Gifs(context.Background())

			assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")
			assertGifs(t, tt.wantGifs, gotGifs)
		})
	}
}

func TestRepo_GifByID(t *testing.T) {
	type args struct {
		ID uuid.UUID
	}
	tests := []struct {
		name         string
		fixtureFiles []string
		args         args
		wantGif      models.Gif
		wantErr      error
	}{
		{
			name: "gif not found :NEG",
			fixtureFiles: []string{
				"topics.yml",
				"voxspheres.yml",
				"users.yml",
				"posts.yml",
				"post_medias.yml",
				"gifs.yml",
				"gif_metadatas.yml",
			},
			args: args{
				ID: uuid.MustParse("00000000-0000-0000-0000-000000000009"),
			},
			wantGif: models.Gif{},
			wantErr: mediarepo.ErrNotFound,
		},
		{
			name: "gif found :POS",
			fixtureFiles: []string{
				"topics.yml",
				"voxspheres.yml",
				"users.yml",
				"posts.yml",
				"post_medias.yml",
				"gifs.yml",
				"gif_metadatas.yml",
			},
			args: args{
				ID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			},
			wantGif: models.Gif{
				ID:      uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				MediaID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				GifMetadata: []models.GifMetadata{
					{
						ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
						GifID:         uuid.MustParse("00000000-0000-0000-0000-000000000001"),
						Height:        1080,
						Width:         1920,
						Url:           "https://example.com/image1.jpg",
						CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
						CreatedAtUnix: 1725091100,
						UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					},
					{
						ID:            uuid.MustParse("00000000-0000-0000-0000-000000000002"),
						GifID:         uuid.MustParse("00000000-0000-0000-0000-000000000001"),
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

			gotGif, gotErr := pgrepo.GifByID(context.Background(), tt.args.ID)

			assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")
			assertGif(t, tt.wantGif, gotGif)
		})
	}
}

func TestRepo_AddGifs(t *testing.T) {
	type args struct {
		gifs []models.Gif
	}
	tests := []struct {
		name            string
		fixtureFiles    []string
		args            args
		wantInsertedGif []models.Gif
		wantGifs        []models.Gif
		wantErr         error
	}{
		{
			name: "duplicate gif :NEG",
			fixtureFiles: []string{
				"topics.yml",
				"voxspheres.yml",
				"users.yml",
				"posts.yml",
				"post_medias.yml",
				"gifs.yml",
				"gif_metadatas.yml",
			},
			args: args{
				gifs: []models.Gif{
					{
						ID:      uuid.MustParse("00000000-0000-0000-0000-000000000001"),
						MediaID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
						GifMetadata: []models.GifMetadata{
							{
								ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
								GifID:         uuid.MustParse("00000000-0000-0000-0000-000000000001"),
								Height:        2160,
								Width:         3840,
								Url:           "https://example.com/image3.jpg",
								CreatedAt:     time.Date(2024, 10, 10, 10, 10, 30, 0, time.UTC),
								CreatedAtUnix: 1725091300,
								UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 30, 0, time.UTC),
							},
						},
					},
				},
			},
			wantInsertedGif: nil,
			wantGifs: []models.Gif{
				{
					ID:      uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					MediaID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					GifMetadata: []models.GifMetadata{
						{
							ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
							GifID:         uuid.MustParse("00000000-0000-0000-0000-000000000001"),
							Height:        1080,
							Width:         1920,
							Url:           "https://example.com/image1.jpg",
							CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
							CreatedAtUnix: 1725091100,
							UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
						},
						{
							ID:            uuid.MustParse("00000000-0000-0000-0000-000000000002"),
							GifID:         uuid.MustParse("00000000-0000-0000-0000-000000000001"),
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
			name: "media not present in parent table :NEG",
			fixtureFiles: []string{
				"topics.yml",
				"voxspheres.yml",
				"users.yml",
				"posts.yml",
				"post_medias.yml",
				"gifs.yml",
				"gif_metadatas.yml",
			},
			args: args{
				gifs: []models.Gif{
					{
						ID:      uuid.MustParse("00000000-0000-0000-0000-000000000003"),
						MediaID: uuid.MustParse("00000000-0000-0000-0000-000000000009"),
						GifMetadata: []models.GifMetadata{
							{
								ID:            uuid.MustParse("00000000-0000-0000-0000-000000000003"),
								GifID:         uuid.MustParse("00000000-0000-0000-0000-000000000003"),
								Height:        2160,
								Width:         3840,
								Url:           "https://example.com/image3.jpg",
								CreatedAt:     time.Date(2024, 10, 10, 10, 10, 30, 0, time.UTC),
								CreatedAtUnix: 1725091300,
								UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 30, 0, time.UTC),
							},
						},
					},
				},
			},
			wantInsertedGif: nil,
			wantGifs: []models.Gif{
				{
					ID:      uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					MediaID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					GifMetadata: []models.GifMetadata{
						{
							ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
							GifID:         uuid.MustParse("00000000-0000-0000-0000-000000000001"),
							Height:        1080,
							Width:         1920,
							Url:           "https://example.com/image1.jpg",
							CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
							CreatedAtUnix: 1725091100,
							UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
						},
						{
							ID:            uuid.MustParse("00000000-0000-0000-0000-000000000002"),
							GifID:         uuid.MustParse("00000000-0000-0000-0000-000000000001"),
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
			name: "add gifs :POS",
			fixtureFiles: []string{
				"topics.yml",
				"voxspheres.yml",
				"users.yml",
				"posts.yml",
				"post_medias.yml",
				"gifs.yml",
				"gif_metadatas.yml",
			},
			args: args{
				gifs: []models.Gif{
					{
						ID:      uuid.MustParse("00000000-0000-0000-0000-000000000002"),
						MediaID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					},
				},
			},
			wantInsertedGif: []models.Gif{
				{
					ID:      uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					MediaID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				},
			},
			wantGifs: []models.Gif{
				{
					ID:      uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					MediaID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					GifMetadata: []models.GifMetadata{
						{
							ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
							GifID:         uuid.MustParse("00000000-0000-0000-0000-000000000001"),
							Height:        1080,
							Width:         1920,
							Url:           "https://example.com/image1.jpg",
							CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
							CreatedAtUnix: 1725091100,
							UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
						},
						{
							ID:            uuid.MustParse("00000000-0000-0000-0000-000000000002"),
							GifID:         uuid.MustParse("00000000-0000-0000-0000-000000000001"),
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
					ID:      uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					MediaID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupPostgres(t, tt.fixtureFiles...)
			pgrepo := mediarepo.NewRepo(db)

			gotInsertedGifs, gotErr := pgrepo.AddGifs(context.Background(), tt.args.gifs...)

			assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")
			assert.Equal(t, tt.wantInsertedGif, gotInsertedGifs, "expect inserted gif to match")

			gotGifs, err := pgrepo.Gifs(context.Background())

			assert.NoError(t, err, "expect no error while getting gifs")
			assertGifs(t, tt.wantGifs, gotGifs)
		})
	}
}

func TestRepo_UpdateGif(t *testing.T) {
	type args struct {
		gif models.Gif
	}
	tests := []struct {
		name         string
		fixtureFiles []string
		args         args
		wantGif      models.Gif
		wantGifs     []models.Gif
		wantErr      error
	}{
		{
			name: "not found :NEG",
			fixtureFiles: []string{
				"topics.yml",
				"voxspheres.yml",
				"users.yml",
				"posts.yml",
				"post_medias.yml",
				"gifs.yml",
				"gif_metadatas.yml",
			},
			args: args{
				gif: models.Gif{
					ID:      uuid.MustParse("00000000-0000-0000-0000-000000000009"),
					MediaID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					GifMetadata: []models.GifMetadata{
						{
							ID:            uuid.MustParse("00000000-0000-0000-0000-000000000003"),
							GifID:         uuid.MustParse("00000000-0000-0000-0000-000000000003"),
							Height:        2160,
							Width:         3840,
							Url:           "https://example.com/image3.jpg",
							CreatedAt:     time.Date(2024, 10, 10, 10, 10, 30, 0, time.UTC),
							CreatedAtUnix: 1725091300,
							UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 30, 0, time.UTC),
						},
					},
				},
			},
			wantGif: models.Gif{},
			wantGifs: []models.Gif{
				{
					ID:      uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					MediaID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					GifMetadata: []models.GifMetadata{
						{
							ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
							GifID:         uuid.MustParse("00000000-0000-0000-0000-000000000001"),
							Height:        1080,
							Width:         1920,
							Url:           "https://example.com/image1.jpg",
							CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
							CreatedAtUnix: 1725091100,
							UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
						},
						{
							ID:            uuid.MustParse("00000000-0000-0000-0000-000000000002"),
							GifID:         uuid.MustParse("00000000-0000-0000-0000-000000000001"),
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
			name: "media not present in parent table :NEG",
			fixtureFiles: []string{
				"topics.yml",
				"voxspheres.yml",
				"users.yml",
				"posts.yml",
				"post_medias.yml",
				"gifs.yml",
				"gif_metadatas.yml",
			},
			args: args{
				gif: models.Gif{
					ID:      uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					MediaID: uuid.MustParse("00000000-0000-0000-0000-000000000009"),
					GifMetadata: []models.GifMetadata{
						{
							ID:            uuid.MustParse("00000000-0000-0000-0000-000000000003"),
							GifID:         uuid.MustParse("00000000-0000-0000-0000-000000000003"),
							Height:        2160,
							Width:         3840,
							Url:           "https://example.com/image3.jpg",
							CreatedAt:     time.Date(2024, 10, 10, 10, 10, 30, 0, time.UTC),
							CreatedAtUnix: 1725091300,
							UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 30, 0, time.UTC),
						},
					},
				},
			},
			wantGif: models.Gif{},
			wantGifs: []models.Gif{
				{
					ID:      uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					MediaID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					GifMetadata: []models.GifMetadata{
						{
							ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
							GifID:         uuid.MustParse("00000000-0000-0000-0000-000000000001"),
							Height:        1080,
							Width:         1920,
							Url:           "https://example.com/image1.jpg",
							CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
							CreatedAtUnix: 1725091100,
							UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
						},
						{
							ID:            uuid.MustParse("00000000-0000-0000-0000-000000000002"),
							GifID:         uuid.MustParse("00000000-0000-0000-0000-000000000001"),
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
			name: "update gif :POS",
			fixtureFiles: []string{
				"topics.yml",
				"voxspheres.yml",
				"users.yml",
				"posts.yml",
				"post_medias.yml",
				"gifs.yml",
				"gif_metadatas.yml",
			},
			args: args{
				gif: models.Gif{
					ID:      uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					MediaID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				},
			},
			wantGif: models.Gif{
				ID:      uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				MediaID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			},
			wantGifs: []models.Gif{
				{
					ID:      uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					MediaID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					GifMetadata: []models.GifMetadata{
						{
							ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
							GifID:         uuid.MustParse("00000000-0000-0000-0000-000000000001"),
							Height:        1080,
							Width:         1920,
							Url:           "https://example.com/image1.jpg",
							CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
							CreatedAtUnix: 1725091100,
							UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
						},
						{
							ID:            uuid.MustParse("00000000-0000-0000-0000-000000000002"),
							GifID:         uuid.MustParse("00000000-0000-0000-0000-000000000001"),
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

			gotGif, gotErr := pgrepo.UpdateGif(context.Background(), tt.args.gif)

			assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")
			assert.Equal(t, tt.wantGif, gotGif, "expect gif to match")

			gotGifs, err := pgrepo.Gifs(context.Background())

			assert.NoError(t, err, "expect no error while getting gifs")
			assertGifs(t, tt.wantGifs, gotGifs)
		})
	}
}

func TestRepo_DeleteGif(t *testing.T) {
	type args struct {
		ID uuid.UUID
	}
	tests := []struct {
		name         string
		fixtureFiles []string
		args         args
		wantGifs     []models.Gif
		wantErr      error
	}{
		{
			name: "not found :NEG",
			fixtureFiles: []string{
				"topics.yml",
				"voxspheres.yml",
				"users.yml",
				"posts.yml",
				"post_medias.yml",
				"gifs.yml",
				"gif_metadatas.yml",
			},
			args: args{
				ID: uuid.MustParse("00000000-0000-0000-0000-000000000009"),
			},
			wantGifs: []models.Gif{
				{
					ID:      uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					MediaID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					GifMetadata: []models.GifMetadata{
						{
							ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
							GifID:         uuid.MustParse("00000000-0000-0000-0000-000000000001"),
							Height:        1080,
							Width:         1920,
							Url:           "https://example.com/image1.jpg",
							CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
							CreatedAtUnix: 1725091100,
							UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
						},
						{
							ID:            uuid.MustParse("00000000-0000-0000-0000-000000000002"),
							GifID:         uuid.MustParse("00000000-0000-0000-0000-000000000001"),
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
			name: "gif deleted :POS",
			fixtureFiles: []string{
				"topics.yml",
				"voxspheres.yml",
				"users.yml",
				"posts.yml",
				"post_medias.yml",
				"gifs.yml",
				"gif_metadatas.yml",
			},
			args: args{
				ID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			},
			wantGifs: nil,
			wantErr:  nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupPostgres(t, tt.fixtureFiles...)
			pgrepo := mediarepo.NewRepo(db)

			gotErr := pgrepo.DeleteGif(context.Background(), tt.args.ID)

			assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")

			gotGifs, err := pgrepo.Gifs(context.Background())

			assert.NoError(t, err, "expect no error while getting gifs")
			assertGifs(t, tt.wantGifs, gotGifs)
		})
	}
}

func TestRepo_GifForeignKeyCascade(t *testing.T) {
	t.Run("on deleting media from parent table, no child references should exist in gifs table", func(t *testing.T) {
		db := setupPostgres(
			t,
			"topics.yml",
			"voxspheres.yml",
			"users.yml",
			"posts.yml",
			"post_medias.yml",
			"gifs.yml",
			"gif_metadatas.yml",
		)
		pgrepo := mediarepo.NewRepo(db)

		err := pgrepo.DeletePostMedia(context.Background(), uuid.MustParse("00000000-0000-0000-0000-000000000002"))

		assert.NoError(t, err, "expect no error while deleting post media")

		gotGifs, err := pgrepo.Gifs(context.Background())

		assert.NoError(t, err, "expect no error while getting gifs")
		assertGifs(t, nil, gotGifs)
	})
}

func TestRepo_GifMetadatas(t *testing.T) {
	tests := []struct {
		name             string
		fixtureFiles     []string
		wantGifMetadatas []models.GifMetadata
		wantErr          error
	}{
		{
			name: "gif metadatas :POS",
			fixtureFiles: []string{
				"topics.yml",
				"voxspheres.yml",
				"users.yml",
				"posts.yml",
				"post_medias.yml",
				"gifs.yml",
				"gif_metadatas.yml",
			},
			wantGifMetadatas: []models.GifMetadata{
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					GifID:         uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Height:        1080,
					Width:         1920,
					Url:           "https://example.com/image1.jpg",
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix: 1725091100,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					GifID:         uuid.MustParse("00000000-0000-0000-0000-000000000001"),
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
			name:             "no gif metadatas :POS",
			fixtureFiles:     []string{},
			wantGifMetadatas: nil,
			wantErr:          nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupPostgres(t, tt.fixtureFiles...)
			pgrepo := mediarepo.NewRepo(db)

			gotGifMetadatas, gotErr := pgrepo.GifMetadatas(context.Background())

			assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")
			assertGifMetadatasWithTimestamp(t, tt.wantGifMetadatas, gotGifMetadatas)
		})
	}
}

func TestRepo_GifMetadataByID(t *testing.T) {
	type args struct {
		ID uuid.UUID
	}
	tests := []struct {
		name            string
		fixtureFiles    []string
		args            args
		wantGifMetadata models.GifMetadata
		wantErr         error
	}{
		{
			name: "gif metadata not found :NEG",
			fixtureFiles: []string{
				"topics.yml",
				"voxspheres.yml",
				"users.yml",
				"posts.yml",
				"post_medias.yml",
				"gifs.yml",
				"gif_metadatas.yml",
			},
			args: args{
				ID: uuid.MustParse("00000000-0000-0000-0000-000000000009"),
			},
			wantGifMetadata: models.GifMetadata{},
			wantErr:         mediarepo.ErrNotFound,
		},
		{
			name: "gif found :POS",
			fixtureFiles: []string{
				"topics.yml",
				"voxspheres.yml",
				"users.yml",
				"posts.yml",
				"post_medias.yml",
				"gifs.yml",
				"gif_metadatas.yml",
			},
			args: args{
				ID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			},
			wantGifMetadata: models.GifMetadata{
				ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				GifID:         uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				Height:        1080,
				Width:         1920,
				Url:           "https://example.com/image1.jpg",
				CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				CreatedAtUnix: 1725091100,
				UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
			},
			wantErr: nil,
		},
		{
			name: "gif metadata found :POS",
			fixtureFiles: []string{
				"topics.yml",
				"voxspheres.yml",
				"users.yml",
				"posts.yml",
				"post_medias.yml",
				"gifs.yml",
				"gif_metadatas.yml",
			},
			args: args{
				ID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
			},
			wantGifMetadata: models.GifMetadata{
				ID:            uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				GifID:         uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				Height:        720,
				Width:         1280,
				Url:           "https://example.com/image2.jpg",
				CreatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
				CreatedAtUnix: 1725091101,
				UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupPostgres(t, tt.fixtureFiles...)
			pgrepo := mediarepo.NewRepo(db)

			gotGifMetadata, gotErr := pgrepo.GifMetadataByID(context.Background(), tt.args.ID)

			assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")
			assert.Equal(t, tt.wantGifMetadata, gotGifMetadata, "expect gif metadata to match")
		})
	}
}

func TestRepo_AddGifMetadatas(t *testing.T) {
	type args struct {
		gifMetadatas []models.GifMetadata
	}
	tests := []struct {
		name                     string
		fixtureFiles             []string
		args                     args
		wantInsertedGifMetadatas []models.GifMetadata
		wantGifMetadatas         []models.GifMetadata
		wantErr                  error
	}{
		{
			name: "duplicate gif metadata :NEG",
			fixtureFiles: []string{
				"topics.yml",
				"voxspheres.yml",
				"users.yml",
				"posts.yml",
				"post_medias.yml",
				"gifs.yml",
				"gif_metadatas.yml",
			},
			args: args{
				gifMetadatas: []models.GifMetadata{
					{
						ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
						GifID:         uuid.MustParse("00000000-0000-0000-0000-000000000001"),
						Height:        2160,
						Width:         3840,
						Url:           "https://example.com/image3.jpg",
						CreatedAt:     time.Date(2024, 10, 10, 10, 10, 30, 0, time.UTC),
						CreatedAtUnix: 1725091300,
						UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 30, 0, time.UTC),
					},
				},
			},
			wantInsertedGifMetadatas: nil,
			wantGifMetadatas: []models.GifMetadata{
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					GifID:         uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Height:        1080,
					Width:         1920,
					Url:           "https://example.com/image1.jpg",
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix: 1725091100,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					GifID:         uuid.MustParse("00000000-0000-0000-0000-000000000001"),
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
			name: "gif not present in parent table :NEG",
			fixtureFiles: []string{
				"topics.yml",
				"voxspheres.yml",
				"users.yml",
				"posts.yml",
				"post_medias.yml",
				"gifs.yml",
				"gif_metadatas.yml",
			},
			args: args{
				gifMetadatas: []models.GifMetadata{
					{
						ID:            uuid.MustParse("00000000-0000-0000-0000-000000000003"),
						GifID:         uuid.MustParse("00000000-0000-0000-0000-000000000003"),
						Height:        2160,
						Width:         3840,
						Url:           "https://example.com/image3.jpg",
						CreatedAt:     time.Date(2024, 10, 10, 10, 10, 30, 0, time.UTC),
						CreatedAtUnix: 1725091300,
						UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 30, 0, time.UTC),
					},
				},
			},
			wantInsertedGifMetadatas: nil,
			wantGifMetadatas: []models.GifMetadata{
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					GifID:         uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Height:        1080,
					Width:         1920,
					Url:           "https://example.com/image1.jpg",
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix: 1725091100,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					GifID:         uuid.MustParse("00000000-0000-0000-0000-000000000001"),
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
			name: "add gif metadatas :NEG",
			fixtureFiles: []string{
				"topics.yml",
				"voxspheres.yml",
				"users.yml",
				"posts.yml",
				"post_medias.yml",
				"gifs.yml",
				"gif_metadatas.yml",
			},
			args: args{
				gifMetadatas: []models.GifMetadata{
					{
						ID:            uuid.MustParse("00000000-0000-0000-0000-000000000003"),
						GifID:         uuid.MustParse("00000000-0000-0000-0000-000000000001"),
						Height:        2160,
						Width:         3840,
						Url:           "https://example.com/image3.jpg",
						CreatedAt:     time.Date(2024, 10, 10, 10, 10, 30, 0, time.UTC),
						CreatedAtUnix: 1725091300,
						UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 30, 0, time.UTC),
					},
				},
			},
			wantInsertedGifMetadatas: nil,
			wantGifMetadatas: []models.GifMetadata{
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					GifID:         uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Height:        1080,
					Width:         1920,
					Url:           "https://example.com/image1.jpg",
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix: 1725091100,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					GifID:         uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Height:        720,
					Width:         1280,
					Url:           "https://example.com/image2.jpg",
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
					CreatedAtUnix: 1725091101,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
				},
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					GifID:         uuid.MustParse("00000000-0000-0000-0000-000000000001"),
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
			gotInsertedGifMetadatas, gotErr := pgrepo.AddGifMetadatas(context.Background(), tt.args.gifMetadatas...)
			endTime := time.Now()

			for _, gotInsertedGifMetadata := range gotInsertedGifMetadatas {
				assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")
				assert.Equal(
					t,
					gotInsertedGifMetadata.UpdatedAt,
					gotInsertedGifMetadata.CreatedAt,
					"expect CreatedAt and UpdatedAt to be same",
				)
				if tt.wantErr == nil {
					assertTimeWithinRange(t, gotInsertedGifMetadata.CreatedAt, startTime, endTime)
					assertTimeWithinRange(t, gotInsertedGifMetadata.UpdatedAt, startTime, endTime)
				}
			}

			gotGifMetadatas, err := pgrepo.GifMetadatas(context.Background())

			assert.NoError(t, err, "expect no error while getting gif metadatas")
			assertGifMetadatasWithoutTimestamp(t, tt.wantGifMetadatas, gotGifMetadatas)
		})
	}
}

func TestRepo_UpdateGifMetadata(t *testing.T) {
	type args struct {
		gifMetadata models.GifMetadata
	}
	tests := []struct {
		name             string
		fixtureFiles     []string
		args             args
		wantGifMetadata  models.GifMetadata
		wantGifMetadatas []models.GifMetadata
		wantErr          error
	}{
		{
			name: "gif metadata not found :NEG",
			fixtureFiles: []string{
				"topics.yml",
				"voxspheres.yml",
				"users.yml",
				"posts.yml",
				"post_medias.yml",
				"gifs.yml",
				"gif_metadatas.yml",
			},
			args: args{
				gifMetadata: models.GifMetadata{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000009"),
					GifID:         uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Height:        2160,
					Width:         3840,
					Url:           "https://example.com/image3.jpg",
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 30, 0, time.UTC),
					CreatedAtUnix: 1725091300,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 30, 0, time.UTC),
				},
			},
			wantGifMetadata: models.GifMetadata{},
			wantGifMetadatas: []models.GifMetadata{
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					GifID:         uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Height:        1080,
					Width:         1920,
					Url:           "https://example.com/image1.jpg",
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix: 1725091100,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					GifID:         uuid.MustParse("00000000-0000-0000-0000-000000000001"),
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
			name: "gif not present in parent table :NEG",
			fixtureFiles: []string{
				"topics.yml",
				"voxspheres.yml",
				"users.yml",
				"posts.yml",
				"post_medias.yml",
				"gifs.yml",
				"gif_metadatas.yml",
			},
			args: args{
				gifMetadata: models.GifMetadata{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					GifID:         uuid.MustParse("00000000-0000-0000-0000-000000000009"),
					Height:        2160,
					Width:         3840,
					Url:           "https://example.com/image3.jpg",
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 30, 0, time.UTC),
					CreatedAtUnix: 1725091300,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 30, 0, time.UTC),
				},
			},
			wantGifMetadata: models.GifMetadata{},
			wantGifMetadatas: []models.GifMetadata{
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					GifID:         uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Height:        1080,
					Width:         1920,
					Url:           "https://example.com/image1.jpg",
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix: 1725091100,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					GifID:         uuid.MustParse("00000000-0000-0000-0000-000000000001"),
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
			name: "update gif metadata :POS",
			fixtureFiles: []string{
				"topics.yml",
				"voxspheres.yml",
				"users.yml",
				"posts.yml",
				"post_medias.yml",
				"gifs.yml",
				"gif_metadatas.yml",
			},
			args: args{
				gifMetadata: models.GifMetadata{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					GifID:         uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Height:        640,
					Width:         360,
					Url:           "https://example.com/image1-updated.jpg",
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 30, 0, time.UTC),
					CreatedAtUnix: 1725091300,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 30, 0, time.UTC),
				},
			},
			wantGifMetadata: models.GifMetadata{
				ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				GifID:         uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				Height:        640,
				Width:         360,
				Url:           "https://example.com/image1-updated.jpg",
				CreatedAt:     time.Date(2024, 10, 10, 10, 10, 30, 0, time.UTC),
				CreatedAtUnix: 1725091300,
				UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 30, 0, time.UTC),
			},
			wantGifMetadatas: []models.GifMetadata{
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					GifID:         uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Height:        640,
					Width:         360,
					Url:           "https://example.com/image1-updated.jpg",
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 30, 0, time.UTC),
					CreatedAtUnix: 1725091300,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 30, 0, time.UTC),
				},
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					GifID:         uuid.MustParse("00000000-0000-0000-0000-000000000001"),
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

			startTime := time.Now()
			gotGifMetadata, gotErr := pgrepo.UpdateGifMetadata(context.Background(), tt.args.gifMetadata)
			endTime := time.Now()

			assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")
			if tt.wantErr == nil {
				assertTimeWithinRange(t, gotGifMetadata.UpdatedAt, startTime, endTime)
			}

			gotGifMetadatas, err := pgrepo.GifMetadatas(context.Background())

			assert.NoError(t, err, "expect no error while getting gif metadatas")
			assertGifMetadatasWithoutTimestamp(t, tt.wantGifMetadatas, gotGifMetadatas)
		})
	}
}

func TestRepo_DeleteGifMetadata(t *testing.T) {
	type args struct {
		ID uuid.UUID
	}
	tests := []struct {
		name             string
		fixtureFiles     []string
		args             args
		wantGifMetadatas []models.GifMetadata
		wantErr          error
	}{
		{
			name: "gif metadata not found :NEG",
			fixtureFiles: []string{
				"topics.yml",
				"voxspheres.yml",
				"users.yml",
				"posts.yml",
				"post_medias.yml",
				"gifs.yml",
				"gif_metadatas.yml",
			},
			args: args{
				ID: uuid.MustParse("00000000-0000-0000-0000-000000000009"),
			},
			wantGifMetadatas: []models.GifMetadata{
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					GifID:         uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Height:        1080,
					Width:         1920,
					Url:           "https://example.com/image1.jpg",
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix: 1725091100,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					GifID:         uuid.MustParse("00000000-0000-0000-0000-000000000001"),
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
			name: "gif metadata deleted :POS",
			fixtureFiles: []string{
				"topics.yml",
				"voxspheres.yml",
				"users.yml",
				"posts.yml",
				"post_medias.yml",
				"gifs.yml",
				"gif_metadatas.yml",
			},
			args: args{
				ID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			},
			wantGifMetadatas: nil,
			wantErr:          nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupPostgres(t, tt.fixtureFiles...)
			pgrepo := mediarepo.NewRepo(db)

			gotErr := pgrepo.DeleteGifMetadata(context.Background(), tt.args.ID)

			assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")

			gotGifMetadatas, err := pgrepo.GifMetadatas(context.Background())

			assert.NoError(t, err, "expect no error while getting gif metadatas")
			assertGifMetadatasWithTimestamp(t, tt.wantGifMetadatas, gotGifMetadatas)
		})
	}
}

func TestRepo_GifMetadataForeignKeyCascade(t *testing.T) {
	t.Run("on deleting gif from parent table, no child references should exist in gif_metadatas table", func(t *testing.T) {
		db := setupPostgres(
			t,
			"topics.yml",
			"voxspheres.yml",
			"users.yml",
			"posts.yml",
			"post_medias.yml",
			"gifs.yml",
			"gif_metadatas.yml",
		)
		pgrepo := mediarepo.NewRepo(db)

		err := pgrepo.DeleteGif(context.Background(), uuid.MustParse("00000000-0000-0000-0000-000000000001"))

		assert.NoError(t, err, "expect no error while deleting gifs")

		gotGifMetadatas, err := pgrepo.GifMetadatas(context.Background())

		assert.NoError(t, err, "expect no error while getting gif metadatas")
		assertGifMetadatasWithTimestamp(t, nil, gotGifMetadatas)
	})
}

func TestRepo_Videos(t *testing.T) {
	tests := []struct {
		name         string
		fixtureFiles []string
		wantVideos   []models.Video
		wantErr      error
	}{
		{
			name: "videos :POS",
			fixtureFiles: []string{
				"topics.yml",
				"voxspheres.yml",
				"users.yml",
				"posts.yml",
				"post_medias.yml",
				"videos.yml",
			},
			wantVideos: []models.Video{
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					MediaID:       uuid.MustParse("00000000-0000-0000-0000-000000000004"),
					Url:           "https://example.com/video.mp4",
					Height:        1080,
					Width:         1920,
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix: 1725091100,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
			},
			wantErr: nil,
		},
		{
			name:         "no videos :POS",
			fixtureFiles: []string{},
			wantVideos:   []models.Video{},
			wantErr:      nil,
		},
	}
	for _, tt := range tests {
		db := setupPostgres(t, tt.fixtureFiles...)
		pgrepo := mediarepo.NewRepo(db)

		gotVideos, gotErr := pgrepo.Videos(context.Background())

		assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")
		assertVideosWithoutTimestamp(t, tt.wantVideos, gotVideos)
	}
}

func TestRepo_VideoByID(t *testing.T) {
	type args struct {
		ID uuid.UUID
	}
	tests := []struct {
		name         string
		fixtureFiles []string
		args         args
		wantVideo    models.Video
		wantErr      error
	}{
		{
			name: "video not found :NEG",
			fixtureFiles: []string{
				"topics.yml",
				"voxspheres.yml",
				"users.yml",
				"posts.yml",
				"post_medias.yml",
				"videos.yml",
			},
			args: args{
				ID: uuid.MustParse("00000000-0000-0000-0000-000000000009"),
			},
			wantVideo: models.Video{},
			wantErr:   mediarepo.ErrNotFound,
		},
		{
			name: "video found :POS",
			fixtureFiles: []string{
				"topics.yml",
				"voxspheres.yml",
				"users.yml",
				"posts.yml",
				"post_medias.yml",
				"videos.yml",
			},
			args: args{
				ID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			},
			wantVideo: models.Video{
				ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				MediaID:       uuid.MustParse("00000000-0000-0000-0000-000000000004"),
				Url:           "https://example.com/video.mp4",
				Height:        1080,
				Width:         1920,
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

			gotVideo, gotErr := pgrepo.VideoByID(context.Background(), tt.args.ID)

			assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")
			assert.Equal(t, tt.wantVideo, gotVideo, "expect video to match")
		})
	}
}

func TestRepo_AddVideos(t *testing.T) {
	type args struct {
		videos []models.Video
	}
	tests := []struct {
		name               string
		fixtureFiles       []string
		args               args
		wantInsertedVideos []models.Video
		wantVideos         []models.Video
		wantErr            error
	}{
		{
			name: "duplicate video :NEG",
			fixtureFiles: []string{
				"topics.yml",
				"voxspheres.yml",
				"users.yml",
				"posts.yml",
				"post_medias.yml",
				"videos.yml",
			},
			args: args{
				videos: []models.Video{
					{
						ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
						MediaID:       uuid.MustParse("00000000-0000-0000-0000-000000000004"),
						Url:           "https://example.com/video.mp4",
						Height:        1080,
						Width:         1920,
						CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
						CreatedAtUnix: 1725091100,
						UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					},
				},
			},
			wantInsertedVideos: nil,
			wantVideos: []models.Video{
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					MediaID:       uuid.MustParse("00000000-0000-0000-0000-000000000004"),
					Url:           "https://example.com/video.mp4",
					Height:        1080,
					Width:         1920,
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
				"videos.yml",
			},
			args: args{
				videos: []models.Video{
					{
						ID:            uuid.MustParse("00000000-0000-0000-0000-000000000002"),
						MediaID:       uuid.MustParse("00000000-0000-0000-0000-000000000009"),
						Url:           "https://example.com/video2.mp4",
						Height:        1080,
						Width:         1920,
						CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
						CreatedAtUnix: 1725091100,
						UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					},
				},
			},
			wantInsertedVideos: nil,
			wantVideos: []models.Video{
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					MediaID:       uuid.MustParse("00000000-0000-0000-0000-000000000004"),
					Url:           "https://example.com/video.mp4",
					Height:        1080,
					Width:         1920,
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix: 1725091100,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
			},
			wantErr: mediarepo.ErrParentTableRecordNotFound,
		},
		{
			name: "add videos :POS",
			fixtureFiles: []string{
				"topics.yml",
				"voxspheres.yml",
				"users.yml",
				"posts.yml",
				"post_medias.yml",
				"videos.yml",
			},
			args: args{
				videos: []models.Video{
					{
						ID:            uuid.MustParse("00000000-0000-0000-0000-000000000003"),
						MediaID:       uuid.MustParse("00000000-0000-0000-0000-000000000001"),
						Url:           "https://example.com/videonew.mp4",
						Height:        1080,
						Width:         1920,
						CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
						CreatedAtUnix: 1725091100,
						UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					},
				},
			},
			wantInsertedVideos: []models.Video{
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					MediaID:       uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Url:           "https://example.com/videonew.mp4",
					Height:        1080,
					Width:         1920,
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix: 1725091100,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
			},
			wantVideos: []models.Video{
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					MediaID:       uuid.MustParse("00000000-0000-0000-0000-000000000004"),
					Url:           "https://example.com/video.mp4",
					Height:        1080,
					Width:         1920,
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix: 1725091100,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					MediaID:       uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Url:           "https://example.com/videonew.mp4",
					Height:        1080,
					Width:         1920,
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
			gotInsertedVideos, gotErr := pgrepo.AddVideos(context.Background(), tt.args.videos...)
			endTime := time.Now()

			for _, gotInsertedVideo := range gotInsertedVideos {
				assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")
				assert.Equal(
					t,
					gotInsertedVideo.UpdatedAt,
					gotInsertedVideo.CreatedAt,
					"expect CreatedAt and UpdatedAt to be same",
				)
				if tt.wantErr == nil {
					assertTimeWithinRange(t, gotInsertedVideo.CreatedAt, startTime, endTime)
					assertTimeWithinRange(t, gotInsertedVideo.UpdatedAt, startTime, endTime)
				}
			}

			gotVideos, err := pgrepo.Videos(context.Background())

			assert.NoError(t, err, "expect no error while getting videos")
			assertVideosWithoutTimestamp(t, tt.wantVideos, gotVideos)
		})
	}
}

func TestRepo_UpdateVideo(t *testing.T) {
	type args struct {
		video models.Video
	}
	tests := []struct {
		name         string
		fixtureFiles []string
		args         args
		wantVideo    models.Video
		wantVideos   []models.Video
		wantErr      error
	}{
		{
			name: "video not found :NEG",
			fixtureFiles: []string{
				"topics.yml",
				"voxspheres.yml",
				"users.yml",
				"posts.yml",
				"post_medias.yml",
				"videos.yml",
			},
			args: args{
				video: models.Video{
					ID:      uuid.MustParse("00000000-0000-0000-0000-000000000009"),
					MediaID: uuid.MustParse("00000000-0000-0000-0000-000000000004"),
					Url:     "https://xyz.com",
				},
			},
			wantVideo: models.Video{},
			wantVideos: []models.Video{
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					MediaID:       uuid.MustParse("00000000-0000-0000-0000-000000000004"),
					Url:           "https://example.com/video.mp4",
					Height:        1080,
					Width:         1920,
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
				"videos.yml",
			},
			args: args{
				video: models.Video{
					ID:      uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					MediaID: uuid.MustParse("00000000-0000-0000-0000-000000000009"),
					Url:     "https://xyz.com",
				},
			},
			wantVideo: models.Video{},
			wantVideos: []models.Video{
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					MediaID:       uuid.MustParse("00000000-0000-0000-0000-000000000004"),
					Url:           "https://example.com/video.mp4",
					Height:        1080,
					Width:         1920,
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix: 1725091100,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
			},
			wantErr: mediarepo.ErrParentTableRecordNotFound,
		},
		{
			name: "update video :POS",
			fixtureFiles: []string{
				"topics.yml",
				"voxspheres.yml",
				"users.yml",
				"posts.yml",
				"post_medias.yml",
				"videos.yml",
			},
			args: args{
				video: models.Video{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					MediaID:       uuid.MustParse("00000000-0000-0000-0000-000000000004"),
					Url:           "https://updatedxyz.com",
					Height:        1080,
					Width:         1920,
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix: 1725091100,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
			},
			wantVideo: models.Video{
				ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				MediaID:       uuid.MustParse("00000000-0000-0000-0000-000000000004"),
				Url:           "https://updatedxyz.com",
				Height:        1080,
				Width:         1920,
				CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				CreatedAtUnix: 1725091100,
				UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
			},
			wantVideos: []models.Video{
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					MediaID:       uuid.MustParse("00000000-0000-0000-0000-000000000004"),
					Url:           "https://updatedxyz.com",
					Height:        1080,
					Width:         1920,
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
			gotVideo, gotErr := pgrepo.UpdateVideo(context.Background(), tt.args.video)
			endTime := time.Now()

			assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")
			if tt.wantErr == nil {
				assertTimeWithinRange(t, gotVideo.UpdatedAt, startTime, endTime)
			}

			gotVideos, err := pgrepo.Videos(context.Background())

			assert.NoError(t, err, "expect no error while getting videos")
			assertVideosWithoutTimestamp(t, tt.wantVideos, gotVideos)
		})
	}
}

func TestRepo_DeleteVideo(t *testing.T) {
	type args struct {
		ID uuid.UUID
	}
	tests := []struct {
		name         string
		fixtureFiles []string
		args         args
		wantVideos   []models.Video
		wantErr      error
	}{
		{
			name: "video not found :NEG",
			fixtureFiles: []string{
				"topics.yml",
				"voxspheres.yml",
				"users.yml",
				"posts.yml",
				"post_medias.yml",
				"videos.yml",
			},
			args: args{
				ID: uuid.MustParse("00000000-0000-0000-0000-000000000009"),
			},
			wantVideos: []models.Video{
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					MediaID:       uuid.MustParse("00000000-0000-0000-0000-000000000004"),
					Url:           "https://example.com/video.mp4",
					Height:        1080,
					Width:         1920,
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix: 1725091100,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
			},
			wantErr: mediarepo.ErrNotFound,
		},
		{
			name: "video deleted :POS",
			fixtureFiles: []string{
				"topics.yml",
				"voxspheres.yml",
				"users.yml",
				"posts.yml",
				"post_medias.yml",
				"videos.yml",
			},
			args: args{
				ID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			},
			wantVideos: []models.Video{},
			wantErr:    nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupPostgres(t, tt.fixtureFiles...)
			pgrepo := mediarepo.NewRepo(db)

			gotErr := pgrepo.DeleteVideo(context.Background(), tt.args.ID)

			assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")

			gotVideos, err := pgrepo.Videos(context.Background())

			assert.NoError(t, err, "expect no error while getting videos")
			assertVideosWithTimestamp(t, tt.wantVideos, gotVideos)
		})
	}
}

func TestRepo_VideoForeignKeyCascade(t *testing.T) {
	t.Run("on deleting media from parent table, no child references should exist in videos table", func(t *testing.T) {
		db := setupPostgres(
			t,
			"topics.yml",
			"voxspheres.yml",
			"users.yml",
			"posts.yml",
			"post_medias.yml",
			"videos.yml",
		)
		pgrepo := mediarepo.NewRepo(db)

		err := pgrepo.DeletePostMedia(context.Background(), uuid.MustParse("00000000-0000-0000-0000-000000000004"))

		assert.NoError(t, err, "expect no error while deleting post media")

		gotVideos, err := pgrepo.Videos(context.Background())

		assert.NoError(t, err, "expect no error while getting videos")
		assert.Equal(t, []models.Video(nil), gotVideos, "expect videos to match")
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

func TestRepo_AddLinks(t *testing.T) {
	type args struct {
		links []models.Link
	}
	tests := []struct {
		name             string
		fixtureFiles     []string
		args             args
		wantInsertedLink []models.Link
		wantLinks        []models.Link
		wantErr          error
	}{
		{
			name: "duplicate link :NEG",
			fixtureFiles: []string{
				"topics.yml",
				"voxspheres.yml",
				"users.yml",
				"posts.yml",
				"post_medias.yml",
				"links.yml",
			},
			args: args{
				links: []models.Link{
					{
						ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
						MediaID:       uuid.MustParse("00000000-0000-0000-0000-000000000005"),
						Link:          "https://example.com/video.mp4",
						CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
						CreatedAtUnix: 1725091100,
						UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					},
				},
			},
			wantInsertedLink: nil,
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
				links: []models.Link{
					{
						ID:            uuid.MustParse("00000000-0000-0000-0000-000000000002"),
						MediaID:       uuid.MustParse("00000000-0000-0000-0000-000000000009"),
						Link:          "https://example.com/video.mp4",
						CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
						CreatedAtUnix: 1725091100,
						UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					},
				},
			},
			wantInsertedLink: nil,
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
			name: "add links :POS",
			fixtureFiles: []string{
				"topics.yml",
				"voxspheres.yml",
				"users.yml",
				"posts.yml",
				"post_medias.yml",
				"links.yml",
			},
			args: args{
				links: []models.Link{
					{
						ID:            uuid.MustParse("00000000-0000-0000-0000-000000000003"),
						MediaID:       uuid.MustParse("00000000-0000-0000-0000-000000000001"),
						Link:          "https://example.com/videonew.mp4",
						CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
						CreatedAtUnix: 1725091100,
						UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					},
				},
			},
			wantInsertedLink: []models.Link{
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					MediaID:       uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Link:          "https://example.com/videonew.mp4",
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix: 1725091100,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
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
			gotInsertedLinks, gotErr := pgrepo.AddLinks(context.Background(), tt.args.links...)
			endTime := time.Now()

			for _, gotInsertedLink := range gotInsertedLinks {
				assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")
				assert.Equal(
					t,
					gotInsertedLink.UpdatedAt,
					gotInsertedLink.CreatedAt,
					"expect CreatedAt and UpdatedAt to be same",
				)
				if tt.wantErr == nil {
					assertTimeWithinRange(t, gotInsertedLink.CreatedAt, startTime, endTime)
					assertTimeWithinRange(t, gotInsertedLink.UpdatedAt, startTime, endTime)
				}
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
			assertLinksWitTimestamp(t, tt.wantLinks, gotLinks)
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

func TestRepo_Galleries(t *testing.T) {
	tests := []struct {
		name          string
		fixtureFiles  []string
		wantGalleries []models.Gallery
		wantErr       error
	}{
		{
			name: "galleries :POS",
			fixtureFiles: []string{
				"topics.yml",
				"voxspheres.yml",
				"users.yml",
				"posts.yml",
				"post_medias.yml",
				"galleries.yml",
				"gallery_metadatas.yml",
			},
			wantGalleries: []models.Gallery{
				{
					ID:      uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					MediaID: uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					GalleryMetadata: []models.GalleryMetadata{
						{
							ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
							GalleryID:     uuid.MustParse("00000000-0000-0000-0000-000000000001"),
							OrderIndex:    0,
							Height:        1080,
							Width:         1920,
							Url:           "https://example.com/gallery11.jpg",
							CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
							CreatedAtUnix: 1725091100,
							UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
						},
						{
							ID:            uuid.MustParse("00000000-0000-0000-0000-000000000002"),
							GalleryID:     uuid.MustParse("00000000-0000-0000-0000-000000000001"),
							OrderIndex:    0,
							Height:        720,
							Width:         1280,
							Url:           "https://example.com/gallery12.jpg",
							CreatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
							CreatedAtUnix: 1725091101,
							UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
						},
					},
				},
				{
					ID:      uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					MediaID: uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					GalleryMetadata: []models.GalleryMetadata{
						{
							ID:            uuid.MustParse("00000000-0000-0000-0000-000000000003"),
							GalleryID:     uuid.MustParse("00000000-0000-0000-0000-000000000002"),
							OrderIndex:    1,
							Height:        720,
							Width:         1280,
							Url:           "https://example.com/gallery21.jpg",
							CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
							CreatedAtUnix: 1725091100,
							UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
						},
						{
							ID:            uuid.MustParse("00000000-0000-0000-0000-000000000004"),
							GalleryID:     uuid.MustParse("00000000-0000-0000-0000-000000000002"),
							OrderIndex:    1,
							Height:        1080,
							Width:         1920,
							Url:           "https://example.com/gallery22.jpg",
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
			name:          "no images :POS",
			fixtureFiles:  []string{},
			wantGalleries: nil,
			wantErr:       nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupPostgres(t, tt.fixtureFiles...)
			pgrepo := mediarepo.NewRepo(db)

			gotGalleries, gotErr := pgrepo.Galleries(context.Background())

			assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")
			assertGalleries(t, tt.wantGalleries, gotGalleries)
		})
	}
}

func TestRepo_GalleryByID(t *testing.T) {
	type args struct {
		ID uuid.UUID
	}
	tests := []struct {
		name         string
		fixtureFiles []string
		args         args
		wantGallery  models.Gallery
		wantErr      error
	}{
		{
			name: "gallery not found :NEG",
			fixtureFiles: []string{
				"topics.yml",
				"voxspheres.yml",
				"users.yml",
				"posts.yml",
				"post_medias.yml",
				"galleries.yml",
				"gallery_metadatas.yml",
			},
			args: args{
				ID: uuid.MustParse("00000000-0000-0000-0000-000000000009"),
			},
			wantGallery: models.Gallery{},
			wantErr:     mediarepo.ErrNotFound,
		},
		{
			name: "gallery found :POS",
			fixtureFiles: []string{
				"topics.yml",
				"voxspheres.yml",
				"users.yml",
				"posts.yml",
				"post_medias.yml",
				"galleries.yml",
				"gallery_metadatas.yml",
			},
			args: args{
				ID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			},
			wantGallery: models.Gallery{
				ID:      uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				MediaID: uuid.MustParse("00000000-0000-0000-0000-000000000003"),
				GalleryMetadata: []models.GalleryMetadata{
					{
						ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
						GalleryID:     uuid.MustParse("00000000-0000-0000-0000-000000000001"),
						OrderIndex:    0,
						Height:        1080,
						Width:         1920,
						Url:           "https://example.com/gallery11.jpg",
						CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
						CreatedAtUnix: 1725091100,
						UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					},
					{
						ID:            uuid.MustParse("00000000-0000-0000-0000-000000000002"),
						GalleryID:     uuid.MustParse("00000000-0000-0000-0000-000000000001"),
						OrderIndex:    0,
						Height:        720,
						Width:         1280,
						Url:           "https://example.com/gallery12.jpg",
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

			gotGallery, gotErr := pgrepo.GalleryByID(context.Background(), tt.args.ID)

			assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")
			assertGallery(t, tt.wantGallery, gotGallery)
		})
	}
}

func TestRepo_AddGalleries(t *testing.T) {
	type args struct {
		galleries []models.Gallery
	}
	tests := []struct {
		name                  string
		fixtureFiles          []string
		args                  args
		wantInsertedGalleries []models.Gallery
		wantGalleries         []models.Gallery
		wantErr               error
	}{
		{
			name: "duplicate gallery id :NEG",
			fixtureFiles: []string{
				"topics.yml",
				"voxspheres.yml",
				"users.yml",
				"posts.yml",
				"post_medias.yml",
				"galleries.yml",
				"gallery_metadatas.yml",
			},
			args: args{
				galleries: []models.Gallery{
					{
						ID:      uuid.MustParse("00000000-0000-0000-0000-000000000001"),
						MediaID: uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					},
				},
			},
			wantInsertedGalleries: nil,
			wantGalleries: []models.Gallery{
				{
					ID:      uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					MediaID: uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					GalleryMetadata: []models.GalleryMetadata{
						{
							ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
							GalleryID:     uuid.MustParse("00000000-0000-0000-0000-000000000001"),
							OrderIndex:    0,
							Height:        1080,
							Width:         1920,
							Url:           "https://example.com/gallery11.jpg",
							CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
							CreatedAtUnix: 1725091100,
							UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
						},
						{
							ID:            uuid.MustParse("00000000-0000-0000-0000-000000000002"),
							GalleryID:     uuid.MustParse("00000000-0000-0000-0000-000000000001"),
							OrderIndex:    0,
							Height:        720,
							Width:         1280,
							Url:           "https://example.com/gallery12.jpg",
							CreatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
							CreatedAtUnix: 1725091101,
							UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
						},
					},
				},
				{
					ID:      uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					MediaID: uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					GalleryMetadata: []models.GalleryMetadata{
						{
							ID:            uuid.MustParse("00000000-0000-0000-0000-000000000003"),
							GalleryID:     uuid.MustParse("00000000-0000-0000-0000-000000000002"),
							OrderIndex:    1,
							Height:        720,
							Width:         1280,
							Url:           "https://example.com/gallery21.jpg",
							CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
							CreatedAtUnix: 1725091100,
							UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
						},
						{
							ID:            uuid.MustParse("00000000-0000-0000-0000-000000000004"),
							GalleryID:     uuid.MustParse("00000000-0000-0000-0000-000000000002"),
							OrderIndex:    1,
							Height:        1080,
							Width:         1920,
							Url:           "https://example.com/gallery22.jpg",
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
				"galleries.yml",
				"gallery_metadatas.yml",
			},
			args: args{
				galleries: []models.Gallery{
					{
						ID:      uuid.MustParse("00000000-0000-0000-0000-000000000003"),
						MediaID: uuid.MustParse("00000000-0000-0000-0000-000000000009"),
					},
				},
			},
			wantInsertedGalleries: nil,
			wantGalleries: []models.Gallery{
				{
					ID:      uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					MediaID: uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					GalleryMetadata: []models.GalleryMetadata{
						{
							ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
							GalleryID:     uuid.MustParse("00000000-0000-0000-0000-000000000001"),
							OrderIndex:    0,
							Height:        1080,
							Width:         1920,
							Url:           "https://example.com/gallery11.jpg",
							CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
							CreatedAtUnix: 1725091100,
							UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
						},
						{
							ID:            uuid.MustParse("00000000-0000-0000-0000-000000000002"),
							GalleryID:     uuid.MustParse("00000000-0000-0000-0000-000000000001"),
							OrderIndex:    0,
							Height:        720,
							Width:         1280,
							Url:           "https://example.com/gallery12.jpg",
							CreatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
							CreatedAtUnix: 1725091101,
							UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
						},
					},
				},
				{
					ID:      uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					MediaID: uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					GalleryMetadata: []models.GalleryMetadata{
						{
							ID:            uuid.MustParse("00000000-0000-0000-0000-000000000003"),
							GalleryID:     uuid.MustParse("00000000-0000-0000-0000-000000000002"),
							OrderIndex:    1,
							Height:        720,
							Width:         1280,
							Url:           "https://example.com/gallery21.jpg",
							CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
							CreatedAtUnix: 1725091100,
							UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
						},
						{
							ID:            uuid.MustParse("00000000-0000-0000-0000-000000000004"),
							GalleryID:     uuid.MustParse("00000000-0000-0000-0000-000000000002"),
							OrderIndex:    1,
							Height:        1080,
							Width:         1920,
							Url:           "https://example.com/gallery22.jpg",
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
			name: "add galleries :POS",
			fixtureFiles: []string{
				"topics.yml",
				"voxspheres.yml",
				"users.yml",
				"posts.yml",
				"post_medias.yml",
				"galleries.yml",
				"gallery_metadatas.yml",
			},
			args: args{
				galleries: []models.Gallery{
					{
						ID:      uuid.MustParse("00000000-0000-0000-0000-000000000003"),
						MediaID: uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					},
				},
			},
			wantInsertedGalleries: []models.Gallery{
				{
					ID:      uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					MediaID: uuid.MustParse("00000000-0000-0000-0000-000000000003"),
				},
			},
			wantGalleries: []models.Gallery{
				{
					ID:      uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					MediaID: uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					GalleryMetadata: []models.GalleryMetadata{
						{
							ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
							GalleryID:     uuid.MustParse("00000000-0000-0000-0000-000000000001"),
							OrderIndex:    0,
							Height:        1080,
							Width:         1920,
							Url:           "https://example.com/gallery11.jpg",
							CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
							CreatedAtUnix: 1725091100,
							UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
						},
						{
							ID:            uuid.MustParse("00000000-0000-0000-0000-000000000002"),
							GalleryID:     uuid.MustParse("00000000-0000-0000-0000-000000000001"),
							OrderIndex:    0,
							Height:        720,
							Width:         1280,
							Url:           "https://example.com/gallery12.jpg",
							CreatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
							CreatedAtUnix: 1725091101,
							UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
						},
					},
				},
				{
					ID:      uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					MediaID: uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					GalleryMetadata: []models.GalleryMetadata{
						{
							ID:            uuid.MustParse("00000000-0000-0000-0000-000000000003"),
							GalleryID:     uuid.MustParse("00000000-0000-0000-0000-000000000002"),
							OrderIndex:    1,
							Height:        720,
							Width:         1280,
							Url:           "https://example.com/gallery21.jpg",
							CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
							CreatedAtUnix: 1725091100,
							UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
						},
						{
							ID:            uuid.MustParse("00000000-0000-0000-0000-000000000004"),
							GalleryID:     uuid.MustParse("00000000-0000-0000-0000-000000000002"),
							OrderIndex:    1,
							Height:        1080,
							Width:         1920,
							Url:           "https://example.com/gallery22.jpg",
							CreatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
							CreatedAtUnix: 1725091101,
							UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
						},
					},
				},
				{
					ID:      uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					MediaID: uuid.MustParse("00000000-0000-0000-0000-000000000003"),
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		db := setupPostgres(t, tt.fixtureFiles...)
		pgrepo := mediarepo.NewRepo(db)

		gotInsertedGalleries, gotErr := pgrepo.AddGalleries(context.Background(), tt.args.galleries...)

		assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")
		assert.Equal(t, tt.wantInsertedGalleries, gotInsertedGalleries, "expect inserted galleries to match")

		gotGalleries, err := pgrepo.Galleries(context.Background())

		assert.NoError(t, err, "expect no error while getting galleries")
		assertGalleries(t, tt.wantGalleries, gotGalleries)
	}
}

func TestRepo_UpdateGallery(t *testing.T) {
	type args struct {
		gallery models.Gallery
	}
	tests := []struct {
		name          string
		fixtureFiles  []string
		args          args
		wantGallery   models.Gallery
		wantGalleries []models.Gallery
		wantErr       error
	}{
		{
			name: "gallery not found :NEG",
			fixtureFiles: []string{
				"topics.yml",
				"voxspheres.yml",
				"users.yml",
				"posts.yml",
				"post_medias.yml",
				"galleries.yml",
				"gallery_metadatas.yml",
			},
			args: args{
				gallery: models.Gallery{
					ID:      uuid.MustParse("00000000-0000-0000-0000-000000000009"),
					MediaID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				},
			},
			wantGallery: models.Gallery{},
			wantGalleries: []models.Gallery{
				{
					ID:      uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					MediaID: uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					GalleryMetadata: []models.GalleryMetadata{
						{
							ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
							GalleryID:     uuid.MustParse("00000000-0000-0000-0000-000000000001"),
							OrderIndex:    0,
							Height:        1080,
							Width:         1920,
							Url:           "https://example.com/gallery11.jpg",
							CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
							CreatedAtUnix: 1725091100,
							UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
						},
						{
							ID:            uuid.MustParse("00000000-0000-0000-0000-000000000002"),
							GalleryID:     uuid.MustParse("00000000-0000-0000-0000-000000000001"),
							OrderIndex:    0,
							Height:        720,
							Width:         1280,
							Url:           "https://example.com/gallery12.jpg",
							CreatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
							CreatedAtUnix: 1725091101,
							UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
						},
					},
				},
				{
					ID:      uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					MediaID: uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					GalleryMetadata: []models.GalleryMetadata{
						{
							ID:            uuid.MustParse("00000000-0000-0000-0000-000000000003"),
							GalleryID:     uuid.MustParse("00000000-0000-0000-0000-000000000002"),
							OrderIndex:    1,
							Height:        720,
							Width:         1280,
							Url:           "https://example.com/gallery21.jpg",
							CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
							CreatedAtUnix: 1725091100,
							UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
						},
						{
							ID:            uuid.MustParse("00000000-0000-0000-0000-000000000004"),
							GalleryID:     uuid.MustParse("00000000-0000-0000-0000-000000000002"),
							OrderIndex:    1,
							Height:        1080,
							Width:         1920,
							Url:           "https://example.com/gallery22.jpg",
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
			name: "media not found in parent table :NEG",
			fixtureFiles: []string{
				"topics.yml",
				"voxspheres.yml",
				"users.yml",
				"posts.yml",
				"post_medias.yml",
				"galleries.yml",
				"gallery_metadatas.yml",
			},
			args: args{
				gallery: models.Gallery{
					ID:      uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					MediaID: uuid.MustParse("00000000-0000-0000-0000-000000000009"),
				},
			},
			wantGallery: models.Gallery{},
			wantGalleries: []models.Gallery{
				{
					ID:      uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					MediaID: uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					GalleryMetadata: []models.GalleryMetadata{
						{
							ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
							GalleryID:     uuid.MustParse("00000000-0000-0000-0000-000000000001"),
							OrderIndex:    0,
							Height:        1080,
							Width:         1920,
							Url:           "https://example.com/gallery11.jpg",
							CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
							CreatedAtUnix: 1725091100,
							UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
						},
						{
							ID:            uuid.MustParse("00000000-0000-0000-0000-000000000002"),
							GalleryID:     uuid.MustParse("00000000-0000-0000-0000-000000000001"),
							OrderIndex:    0,
							Height:        720,
							Width:         1280,
							Url:           "https://example.com/gallery12.jpg",
							CreatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
							CreatedAtUnix: 1725091101,
							UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
						},
					},
				},
				{
					ID:      uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					MediaID: uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					GalleryMetadata: []models.GalleryMetadata{
						{
							ID:            uuid.MustParse("00000000-0000-0000-0000-000000000003"),
							GalleryID:     uuid.MustParse("00000000-0000-0000-0000-000000000002"),
							OrderIndex:    1,
							Height:        720,
							Width:         1280,
							Url:           "https://example.com/gallery21.jpg",
							CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
							CreatedAtUnix: 1725091100,
							UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
						},
						{
							ID:            uuid.MustParse("00000000-0000-0000-0000-000000000004"),
							GalleryID:     uuid.MustParse("00000000-0000-0000-0000-000000000002"),
							OrderIndex:    1,
							Height:        1080,
							Width:         1920,
							Url:           "https://example.com/gallery22.jpg",
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
			name: "update gallery :NEG",
			fixtureFiles: []string{
				"topics.yml",
				"voxspheres.yml",
				"users.yml",
				"posts.yml",
				"post_medias.yml",
				"galleries.yml",
				"gallery_metadatas.yml",
			},
			args: args{
				gallery: models.Gallery{
					ID:      uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					MediaID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				},
			},
			wantGallery: models.Gallery{
				ID:      uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				MediaID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			},
			wantGalleries: []models.Gallery{
				{
					ID:      uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					MediaID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					GalleryMetadata: []models.GalleryMetadata{
						{
							ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
							GalleryID:     uuid.MustParse("00000000-0000-0000-0000-000000000001"),
							OrderIndex:    0,
							Height:        1080,
							Width:         1920,
							Url:           "https://example.com/gallery11.jpg",
							CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
							CreatedAtUnix: 1725091100,
							UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
						},
						{
							ID:            uuid.MustParse("00000000-0000-0000-0000-000000000002"),
							GalleryID:     uuid.MustParse("00000000-0000-0000-0000-000000000001"),
							OrderIndex:    0,
							Height:        720,
							Width:         1280,
							Url:           "https://example.com/gallery12.jpg",
							CreatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
							CreatedAtUnix: 1725091101,
							UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
						},
					},
				},
				{
					ID:      uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					MediaID: uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					GalleryMetadata: []models.GalleryMetadata{
						{
							ID:            uuid.MustParse("00000000-0000-0000-0000-000000000003"),
							GalleryID:     uuid.MustParse("00000000-0000-0000-0000-000000000002"),
							OrderIndex:    1,
							Height:        720,
							Width:         1280,
							Url:           "https://example.com/gallery21.jpg",
							CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
							CreatedAtUnix: 1725091100,
							UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
						},
						{
							ID:            uuid.MustParse("00000000-0000-0000-0000-000000000004"),
							GalleryID:     uuid.MustParse("00000000-0000-0000-0000-000000000002"),
							OrderIndex:    1,
							Height:        1080,
							Width:         1920,
							Url:           "https://example.com/gallery22.jpg",
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
		db := setupPostgres(t, tt.fixtureFiles...)
		pgrepo := mediarepo.NewRepo(db)

		gotGallery, gotErr := pgrepo.UpdateGallery(context.Background(), tt.args.gallery)

		assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")
		assert.Equal(t, tt.wantGallery, gotGallery, "expect gallery to match")

		gotGalleries, err := pgrepo.Galleries(context.Background())

		assert.NoError(t, err, "expect no error while getting galleries")
		assertGalleries(t, tt.wantGalleries, gotGalleries)
	}
}

func TestRepo_DeleteGallery(t *testing.T) {
	type args struct {
		ID uuid.UUID
	}
	tests := []struct {
		name          string
		fixtureFiles  []string
		args          args
		wantGalleries []models.Gallery
		wantErr       error
	}{
		{
			name: "gallery not found :NEG",
			fixtureFiles: []string{
				"topics.yml",
				"voxspheres.yml",
				"users.yml",
				"posts.yml",
				"post_medias.yml",
				"galleries.yml",
				"gallery_metadatas.yml",
			},
			args: args{
				ID: uuid.MustParse("00000000-0000-0000-0000-000000000009"),
			},
			wantGalleries: []models.Gallery{
				{
					ID:      uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					MediaID: uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					GalleryMetadata: []models.GalleryMetadata{
						{
							ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
							GalleryID:     uuid.MustParse("00000000-0000-0000-0000-000000000001"),
							OrderIndex:    0,
							Height:        1080,
							Width:         1920,
							Url:           "https://example.com/gallery11.jpg",
							CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
							CreatedAtUnix: 1725091100,
							UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
						},
						{
							ID:            uuid.MustParse("00000000-0000-0000-0000-000000000002"),
							GalleryID:     uuid.MustParse("00000000-0000-0000-0000-000000000001"),
							OrderIndex:    0,
							Height:        720,
							Width:         1280,
							Url:           "https://example.com/gallery12.jpg",
							CreatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
							CreatedAtUnix: 1725091101,
							UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
						},
					},
				},
				{
					ID:      uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					MediaID: uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					GalleryMetadata: []models.GalleryMetadata{
						{
							ID:            uuid.MustParse("00000000-0000-0000-0000-000000000003"),
							GalleryID:     uuid.MustParse("00000000-0000-0000-0000-000000000002"),
							OrderIndex:    1,
							Height:        720,
							Width:         1280,
							Url:           "https://example.com/gallery21.jpg",
							CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
							CreatedAtUnix: 1725091100,
							UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
						},
						{
							ID:            uuid.MustParse("00000000-0000-0000-0000-000000000004"),
							GalleryID:     uuid.MustParse("00000000-0000-0000-0000-000000000002"),
							OrderIndex:    1,
							Height:        1080,
							Width:         1920,
							Url:           "https://example.com/gallery22.jpg",
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
			name: "gallery deleted :POS",
			fixtureFiles: []string{
				"topics.yml",
				"voxspheres.yml",
				"users.yml",
				"posts.yml",
				"post_medias.yml",
				"galleries.yml",
				"gallery_metadatas.yml",
			},
			args: args{
				ID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			},
			wantGalleries: []models.Gallery{
				{
					ID:      uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					MediaID: uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					GalleryMetadata: []models.GalleryMetadata{
						{
							ID:            uuid.MustParse("00000000-0000-0000-0000-000000000003"),
							GalleryID:     uuid.MustParse("00000000-0000-0000-0000-000000000002"),
							OrderIndex:    1,
							Height:        720,
							Width:         1280,
							Url:           "https://example.com/gallery21.jpg",
							CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
							CreatedAtUnix: 1725091100,
							UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
						},
						{
							ID:            uuid.MustParse("00000000-0000-0000-0000-000000000004"),
							GalleryID:     uuid.MustParse("00000000-0000-0000-0000-000000000002"),
							OrderIndex:    1,
							Height:        1080,
							Width:         1920,
							Url:           "https://example.com/gallery22.jpg",
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

			gotErr := pgrepo.DeleteGallery(context.Background(), tt.args.ID)

			assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")

			gotGalleries, err := pgrepo.Galleries(context.Background())

			assert.NoError(t, err, "expect no error while getting galleries")
			assertGalleries(t, tt.wantGalleries, gotGalleries)
		})
	}
}

func TestRepo_GalleryForeignKeyCascade(t *testing.T) {
	t.Run("on deleting media from parent table , no child references should exist in galleries table", func(t *testing.T) {
		db := setupPostgres(
			t,
			"topics.yml",
			"voxspheres.yml",
			"users.yml",
			"posts.yml",
			"post_medias.yml",
			"galleries.yml",
			"gallery_metadatas.yml",
		)
		pgrepo := mediarepo.NewRepo(db)

		err := pgrepo.DeletePostMedia(context.Background(), uuid.MustParse("00000000-0000-0000-0000-000000000003"))

		assert.NoError(t, err, "expect no error while deleting post media")

		gotGalleries, err := pgrepo.Galleries(context.Background())

		assert.NoError(t, err, "expect no error while getting galleries")
		assertGalleries(t, nil, gotGalleries)
	})
}

func TestRepo_GalleryMetadatas(t *testing.T) {
	tests := []struct {
		name                 string
		fixtureFiles         []string
		wantGalleryMetadatas []models.GalleryMetadata
		wantErr              error
	}{
		{
			name: "gallery metadatas :POS",
			fixtureFiles: []string{
				"topics.yml",
				"voxspheres.yml",
				"users.yml",
				"posts.yml",
				"post_medias.yml",
				"galleries.yml",
				"gallery_metadatas.yml",
			},
			wantGalleryMetadatas: []models.GalleryMetadata{
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					GalleryID:     uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					OrderIndex:    0,
					Height:        1080,
					Width:         1920,
					Url:           "https://example.com/gallery11.jpg",
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix: 1725091100,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					GalleryID:     uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					OrderIndex:    0,
					Height:        720,
					Width:         1280,
					Url:           "https://example.com/gallery12.jpg",
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
					CreatedAtUnix: 1725091101,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
				},
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					GalleryID:     uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					OrderIndex:    1,
					Height:        720,
					Width:         1280,
					Url:           "https://example.com/gallery21.jpg",
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix: 1725091100,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000004"),
					GalleryID:     uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					OrderIndex:    1,
					Height:        1080,
					Width:         1920,
					Url:           "https://example.com/gallery22.jpg",
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
					CreatedAtUnix: 1725091101,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
				},
			},
			wantErr: nil,
		},
		{
			name:                 "no gallery metadatas :POS",
			fixtureFiles:         []string{},
			wantGalleryMetadatas: nil,
			wantErr:              nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupPostgres(t, tt.fixtureFiles...)
			pgrepo := mediarepo.NewRepo(db)

			gotGalleryMetadatas, gotErr := pgrepo.GalleryMetadatas(context.Background())

			assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")
			assertGalleryMetadatasWithTimestamp(t, tt.wantGalleryMetadatas, gotGalleryMetadatas)
		})
	}
}

func TestRepo_GalleryMetadataByID(t *testing.T) {
	type args struct {
		ID uuid.UUID
	}
	tests := []struct {
		name                string
		fixtureFiles        []string
		args                args
		wantGalleryMetadata models.GalleryMetadata
		wantErr             error
	}{
		{
			name: "gallery metadata not found :NEG",
			fixtureFiles: []string{
				"topics.yml",
				"voxspheres.yml",
				"users.yml",
				"posts.yml",
				"post_medias.yml",
				"galleries.yml",
				"gallery_metadatas.yml",
			},
			args: args{
				ID: uuid.MustParse("00000000-0000-0000-0000-000000000009"),
			},
			wantGalleryMetadata: models.GalleryMetadata{},
			wantErr:             mediarepo.ErrNotFound,
		},
		{
			name: "gallery metadata found :POS",
			fixtureFiles: []string{
				"topics.yml",
				"voxspheres.yml",
				"users.yml",
				"posts.yml",
				"post_medias.yml",
				"galleries.yml",
				"gallery_metadatas.yml",
			},
			args: args{
				ID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			},
			wantGalleryMetadata: models.GalleryMetadata{
				ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				GalleryID:     uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				OrderIndex:    0,
				Height:        1080,
				Width:         1920,
				Url:           "https://example.com/gallery11.jpg",
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

			gotGalleryMetadata, gotErr := pgrepo.GalleryMetadataByID(context.Background(), tt.args.ID)

			assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")
			assert.Equal(t, tt.wantGalleryMetadata, gotGalleryMetadata, "expect gallery metadata to match")
		})
	}
}

func TestRepo_AddGalleryMetadatas(t *testing.T) {
	type args struct {
		galleryMetadatas []models.GalleryMetadata
	}
	tests := []struct {
		name                        string
		fixtureFiles                []string
		args                        args
		wantInsertedGalleryMetadata []models.GalleryMetadata
		wantGalleryMetadatas        []models.GalleryMetadata
		wantErr                     error
	}{
		{
			name: "duplicate image metadata :NEG",
			fixtureFiles: []string{
				"topics.yml",
				"voxspheres.yml",
				"users.yml",
				"posts.yml",
				"post_medias.yml",
				"galleries.yml",
				"gallery_metadatas.yml",
			},
			args: args{
				galleryMetadatas: []models.GalleryMetadata{
					{
						ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
						GalleryID:     uuid.MustParse("00000000-0000-0000-0000-000000000001"),
						OrderIndex:    0,
						Height:        1080,
						Width:         1920,
						Url:           "https://example.com/gallery11.jpg",
						CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
						CreatedAtUnix: 1725091100,
						UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					},
				},
			},
			wantInsertedGalleryMetadata: nil,
			wantGalleryMetadatas: []models.GalleryMetadata{
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					GalleryID:     uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					OrderIndex:    0,
					Height:        1080,
					Width:         1920,
					Url:           "https://example.com/gallery11.jpg",
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix: 1725091100,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					GalleryID:     uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					OrderIndex:    0,
					Height:        720,
					Width:         1280,
					Url:           "https://example.com/gallery12.jpg",
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
					CreatedAtUnix: 1725091101,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
				},
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					GalleryID:     uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					OrderIndex:    1,
					Height:        720,
					Width:         1280,
					Url:           "https://example.com/gallery21.jpg",
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix: 1725091100,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000004"),
					GalleryID:     uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					OrderIndex:    1,
					Height:        1080,
					Width:         1920,
					Url:           "https://example.com/gallery22.jpg",
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
					CreatedAtUnix: 1725091101,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
				},
			},
			wantErr: mediarepo.ErrDuplicateID,
		},
		{
			name: "gallery not present in parent table :NEG",
			fixtureFiles: []string{
				"topics.yml",
				"voxspheres.yml",
				"users.yml",
				"posts.yml",
				"post_medias.yml",
				"galleries.yml",
				"gallery_metadatas.yml",
			},
			args: args{
				galleryMetadatas: []models.GalleryMetadata{
					{
						ID:            uuid.MustParse("00000000-0000-0000-0000-000000000005"),
						GalleryID:     uuid.MustParse("00000000-0000-0000-0000-000000000009"),
						Height:        2160,
						Width:         3840,
						Url:           "https://example.com/image3.jpg",
						CreatedAt:     time.Date(2024, 10, 10, 10, 10, 30, 0, time.UTC),
						CreatedAtUnix: 1725091300,
						UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 30, 0, time.UTC),
					},
				},
			},
			wantInsertedGalleryMetadata: nil,
			wantGalleryMetadatas: []models.GalleryMetadata{
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					GalleryID:     uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					OrderIndex:    0,
					Height:        1080,
					Width:         1920,
					Url:           "https://example.com/gallery11.jpg",
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix: 1725091100,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					GalleryID:     uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					OrderIndex:    0,
					Height:        720,
					Width:         1280,
					Url:           "https://example.com/gallery12.jpg",
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
					CreatedAtUnix: 1725091101,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
				},
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					GalleryID:     uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					OrderIndex:    1,
					Height:        720,
					Width:         1280,
					Url:           "https://example.com/gallery21.jpg",
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix: 1725091100,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000004"),
					GalleryID:     uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					OrderIndex:    1,
					Height:        1080,
					Width:         1920,
					Url:           "https://example.com/gallery22.jpg",
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
					CreatedAtUnix: 1725091101,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
				},
			},
			wantErr: mediarepo.ErrParentTableRecordNotFound,
		},
		{
			name: "add gallery metadata :POS",
			fixtureFiles: []string{
				"topics.yml",
				"voxspheres.yml",
				"users.yml",
				"posts.yml",
				"post_medias.yml",
				"galleries.yml",
				"gallery_metadatas.yml",
			},
			args: args{
				galleryMetadatas: []models.GalleryMetadata{
					{
						ID:            uuid.MustParse("00000000-0000-0000-0000-000000000005"),
						GalleryID:     uuid.MustParse("00000000-0000-0000-0000-000000000001"),
						Height:        2160,
						Width:         3840,
						Url:           "https://example.com/imagenew.jpg",
						CreatedAt:     time.Date(2024, 10, 10, 10, 10, 30, 0, time.UTC),
						CreatedAtUnix: 1725091300,
						UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 30, 0, time.UTC),
					},
				},
			},
			wantInsertedGalleryMetadata: []models.GalleryMetadata{
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000005"),
					GalleryID:     uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Height:        2160,
					Width:         3840,
					Url:           "https://example.com/imagenew.jpg",
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 30, 0, time.UTC),
					CreatedAtUnix: 1725091300,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 30, 0, time.UTC),
				},
			},
			wantGalleryMetadatas: []models.GalleryMetadata{
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					GalleryID:     uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					OrderIndex:    0,
					Height:        1080,
					Width:         1920,
					Url:           "https://example.com/gallery11.jpg",
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix: 1725091100,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					GalleryID:     uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					OrderIndex:    0,
					Height:        720,
					Width:         1280,
					Url:           "https://example.com/gallery12.jpg",
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
					CreatedAtUnix: 1725091101,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
				},
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					GalleryID:     uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					OrderIndex:    1,
					Height:        720,
					Width:         1280,
					Url:           "https://example.com/gallery21.jpg",
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix: 1725091100,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000004"),
					GalleryID:     uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					OrderIndex:    1,
					Height:        1080,
					Width:         1920,
					Url:           "https://example.com/gallery22.jpg",
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
					CreatedAtUnix: 1725091101,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
				},
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000005"),
					GalleryID:     uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Height:        2160,
					Width:         3840,
					Url:           "https://example.com/imagenew.jpg",
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
			gotInsertedGalleryMetadatas, gotErr := pgrepo.AddGalleryMetadatas(context.Background(), tt.args.galleryMetadatas...)
			endTime := time.Now()

			for _, gotInsertedGalleryMetadata := range gotInsertedGalleryMetadatas {
				assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")
				assert.Equal(
					t,
					gotInsertedGalleryMetadata.UpdatedAt,
					gotInsertedGalleryMetadata.CreatedAt,
					"expect CreatedAt and UpdatedAt to be same",
				)
				if tt.wantErr == nil {
					assertTimeWithinRange(t, gotInsertedGalleryMetadata.CreatedAt, startTime, endTime)
					assertTimeWithinRange(t, gotInsertedGalleryMetadata.UpdatedAt, startTime, endTime)
				}
			}

			gotGalleryMetadatas, err := pgrepo.GalleryMetadatas(context.Background())

			assert.NoError(t, err, "expect no error while getting gallery metadatas")
			assertGalleryMetadatasWithoutTimestamp(t, tt.wantGalleryMetadatas, gotGalleryMetadatas)
		})
	}
}

func TestRepo_UpdateGalleryMetadata(t *testing.T) {
	type args struct {
		galleryMetadata models.GalleryMetadata
	}
	tests := []struct {
		name                 string
		fixtureFiles         []string
		args                 args
		wantGalleryMetadata  models.GalleryMetadata
		wantGalleryMetadatas []models.GalleryMetadata
		wantErr              error
	}{
		{
			name: "gallery metadata not found :NEG",
			fixtureFiles: []string{
				"topics.yml",
				"voxspheres.yml",
				"users.yml",
				"posts.yml",
				"post_medias.yml",
				"galleries.yml",
				"gallery_metadatas.yml",
			},
			args: args{
				galleryMetadata: models.GalleryMetadata{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000009"),
					GalleryID:     uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					OrderIndex:    0,
					Height:        1080,
					Width:         1920,
					Url:           "https://example.com/gallery11-updated.jpg",
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix: 1725091100,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
			},
			wantGalleryMetadata: models.GalleryMetadata{},
			wantGalleryMetadatas: []models.GalleryMetadata{
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					GalleryID:     uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					OrderIndex:    0,
					Height:        1080,
					Width:         1920,
					Url:           "https://example.com/gallery11.jpg",
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix: 1725091100,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					GalleryID:     uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					OrderIndex:    0,
					Height:        720,
					Width:         1280,
					Url:           "https://example.com/gallery12.jpg",
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
					CreatedAtUnix: 1725091101,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
				},
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					GalleryID:     uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					OrderIndex:    1,
					Height:        720,
					Width:         1280,
					Url:           "https://example.com/gallery21.jpg",
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix: 1725091100,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000004"),
					GalleryID:     uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					OrderIndex:    1,
					Height:        1080,
					Width:         1920,
					Url:           "https://example.com/gallery22.jpg",
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
					CreatedAtUnix: 1725091101,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
				},
			},
			wantErr: mediarepo.ErrNotFound,
		},
		{
			name: "gallery not present in parent table :NEG",
			fixtureFiles: []string{
				"topics.yml",
				"voxspheres.yml",
				"users.yml",
				"posts.yml",
				"post_medias.yml",
				"galleries.yml",
				"gallery_metadatas.yml",
			},
			args: args{
				galleryMetadata: models.GalleryMetadata{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					GalleryID:     uuid.MustParse("00000000-0000-0000-0000-000000000009"),
					Height:        2160,
					Width:         3840,
					Url:           "https://example.com/gallery11-updated.jpg",
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 30, 0, time.UTC),
					CreatedAtUnix: 1725091300,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 30, 0, time.UTC),
				},
			},
			wantGalleryMetadata: models.GalleryMetadata{},
			wantGalleryMetadatas: []models.GalleryMetadata{
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					GalleryID:     uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					OrderIndex:    0,
					Height:        1080,
					Width:         1920,
					Url:           "https://example.com/gallery11.jpg",
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix: 1725091100,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					GalleryID:     uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					OrderIndex:    0,
					Height:        720,
					Width:         1280,
					Url:           "https://example.com/gallery12.jpg",
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
					CreatedAtUnix: 1725091101,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
				},
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					GalleryID:     uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					OrderIndex:    1,
					Height:        720,
					Width:         1280,
					Url:           "https://example.com/gallery21.jpg",
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix: 1725091100,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000004"),
					GalleryID:     uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					OrderIndex:    1,
					Height:        1080,
					Width:         1920,
					Url:           "https://example.com/gallery22.jpg",
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
					CreatedAtUnix: 1725091101,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
				},
			},
			wantErr: mediarepo.ErrParentTableRecordNotFound,
		},
		{
			name: "update gallery metadata :POS",
			fixtureFiles: []string{
				"topics.yml",
				"voxspheres.yml",
				"users.yml",
				"posts.yml",
				"post_medias.yml",
				"galleries.yml",
				"gallery_metadatas.yml",
			},
			args: args{
				galleryMetadata: models.GalleryMetadata{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					GalleryID:     uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Height:        1080,
					Width:         1920,
					Url:           "https://example.com/image-updated.jpg",
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 30, 0, time.UTC),
					CreatedAtUnix: 1725091300,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 30, 0, time.UTC),
				},
			},
			wantGalleryMetadata: models.GalleryMetadata{
				ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				GalleryID:     uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				Height:        1080,
				Width:         1920,
				Url:           "https://example.com/image-updated.jpg",
				CreatedAt:     time.Date(2024, 10, 10, 10, 10, 30, 0, time.UTC),
				CreatedAtUnix: 1725091300,
				UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 30, 0, time.UTC),
			},
			wantGalleryMetadatas: []models.GalleryMetadata{
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					GalleryID:     uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Height:        1080,
					Width:         1920,
					Url:           "https://example.com/image-updated.jpg",
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 30, 0, time.UTC),
					CreatedAtUnix: 1725091300,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 30, 0, time.UTC),
				},
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					GalleryID:     uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					OrderIndex:    0,
					Height:        720,
					Width:         1280,
					Url:           "https://example.com/gallery12.jpg",
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
					CreatedAtUnix: 1725091101,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
				},
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					GalleryID:     uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					OrderIndex:    1,
					Height:        720,
					Width:         1280,
					Url:           "https://example.com/gallery21.jpg",
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix: 1725091100,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000004"),
					GalleryID:     uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					OrderIndex:    1,
					Height:        1080,
					Width:         1920,
					Url:           "https://example.com/gallery22.jpg",
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

			startTime := time.Now()
			gotGalleryMetadata, gotErr := pgrepo.UpdateGalleryMetadata(context.Background(), tt.args.galleryMetadata)
			endTime := time.Now()

			assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")
			if tt.wantErr == nil {
				assertTimeWithinRange(t, gotGalleryMetadata.UpdatedAt, startTime, endTime)
			}

			gotGalleryMetadatas, err := pgrepo.GalleryMetadatas(context.Background())

			assert.NoError(t, err, "expect no error while getting gallery metadatas")
			assertGalleryMetadatasWithoutTimestamp(t, tt.wantGalleryMetadatas, gotGalleryMetadatas)
		})
	}
}

func TestRepo_DeleteGalleryMetadata(t *testing.T) {
	type args struct {
		ID uuid.UUID
	}
	tests := []struct {
		name                 string
		fixtureFiles         []string
		args                 args
		wantGalleryMetadatas []models.GalleryMetadata
		wantErr              error
	}{
		{
			name: "gallery metadata not found :NEG",
			fixtureFiles: []string{
				"topics.yml",
				"voxspheres.yml",
				"users.yml",
				"posts.yml",
				"post_medias.yml",
				"galleries.yml",
				"gallery_metadatas.yml",
			},
			args: args{
				ID: uuid.MustParse("00000000-0000-0000-0000-000000000009"),
			},
			wantGalleryMetadatas: []models.GalleryMetadata{
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					GalleryID:     uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					OrderIndex:    0,
					Height:        1080,
					Width:         1920,
					Url:           "https://example.com/gallery11.jpg",
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix: 1725091100,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					GalleryID:     uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					OrderIndex:    0,
					Height:        720,
					Width:         1280,
					Url:           "https://example.com/gallery12.jpg",
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
					CreatedAtUnix: 1725091101,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
				},
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					GalleryID:     uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					OrderIndex:    1,
					Height:        720,
					Width:         1280,
					Url:           "https://example.com/gallery21.jpg",
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix: 1725091100,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000004"),
					GalleryID:     uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					OrderIndex:    1,
					Height:        1080,
					Width:         1920,
					Url:           "https://example.com/gallery22.jpg",
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
					CreatedAtUnix: 1725091101,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
				},
			},
			wantErr: mediarepo.ErrNotFound,
		},
		{
			name: "gallery metadata deleted :POS",
			fixtureFiles: []string{
				"topics.yml",
				"voxspheres.yml",
				"users.yml",
				"posts.yml",
				"post_medias.yml",
				"galleries.yml",
				"gallery_metadatas.yml",
			},
			args: args{
				ID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			},
			wantGalleryMetadatas: []models.GalleryMetadata{
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					GalleryID:     uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					OrderIndex:    0,
					Height:        720,
					Width:         1280,
					Url:           "https://example.com/gallery12.jpg",
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
					CreatedAtUnix: 1725091101,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
				},
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					GalleryID:     uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					OrderIndex:    1,
					Height:        720,
					Width:         1280,
					Url:           "https://example.com/gallery21.jpg",
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix: 1725091100,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000004"),
					GalleryID:     uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					OrderIndex:    1,
					Height:        1080,
					Width:         1920,
					Url:           "https://example.com/gallery22.jpg",
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

			gotErr := pgrepo.DeleteGalleryMetadata(context.Background(), tt.args.ID)

			assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")

			gotGalleryMetadatas, err := pgrepo.GalleryMetadatas(context.Background())

			assert.NoError(t, err, "expect no error while getting gallery metadatas")
			assertGalleryMetadatasWithTimestamp(t, tt.wantGalleryMetadatas, gotGalleryMetadatas)
		})
	}
}

func TestRepo_GalleryMetadataForeignKeyCascade(t *testing.T) {
	t.Run("on deleting gallery from parent table , no child references should exist in gallery_metadatas table", func(t *testing.T) {
		db := setupPostgres(
			t,
			"topics.yml",
			"voxspheres.yml",
			"users.yml",
			"posts.yml",
			"post_medias.yml",
			"galleries.yml",
			"gallery_metadatas.yml",
		)
		pgrepo := mediarepo.NewRepo(db)

		err := pgrepo.DeleteGallery(context.Background(), uuid.MustParse("00000000-0000-0000-0000-000000000001"))

		assert.NoError(t, err, "expect no error while deleting gallery")

		gotGalleryMetadatas, err := pgrepo.GalleryMetadatas(context.Background())

		assert.NoError(t, err, "expect no error while getting gallery metadatas")
		wantGalleryMetadatas := []models.GalleryMetadata{
			{
				ID:            uuid.MustParse("00000000-0000-0000-0000-000000000003"),
				GalleryID:     uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				OrderIndex:    1,
				Height:        720,
				Width:         1280,
				Url:           "https://example.com/gallery21.jpg",
				CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				CreatedAtUnix: 1725091100,
				UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
			},
			{
				ID:            uuid.MustParse("00000000-0000-0000-0000-000000000004"),
				GalleryID:     uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				OrderIndex:    1,
				Height:        1080,
				Width:         1920,
				Url:           "https://example.com/gallery22.jpg",
				CreatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
				CreatedAtUnix: 1725091101,
				UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
			},
		}
		assertGalleryMetadatasWithTimestamp(t, wantGalleryMetadatas, gotGalleryMetadatas)
	})
}
