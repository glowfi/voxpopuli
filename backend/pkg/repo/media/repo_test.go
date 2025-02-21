package media_test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"slices"
	"testing"

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

	// load fixture
	fixture := dbfixture.New(db)
	if err := fixture.Load(context.Background(), os.DirFS("testdata"), fixtureFiles...); err != nil {
		t.Fatal("failed to load fixtures", err)
	}

	return db
}

func assertPostMedias(t *testing.T, wantPostMedias, gotPostMedias []models.PostMedia) {
	t.Helper()

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

func TestRepo_ForeignKeyCascade(t *testing.T) {
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
