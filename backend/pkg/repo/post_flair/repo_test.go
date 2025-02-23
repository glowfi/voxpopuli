package postflair_test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"slices"
	"testing"

	"github.com/glowfi/voxpopuli/backend/pkg/models"
	postRepo "github.com/glowfi/voxpopuli/backend/pkg/repo/post"
	postflaireRepo "github.com/glowfi/voxpopuli/backend/pkg/repo/post_flair"
	voxrepo "github.com/glowfi/voxpopuli/backend/pkg/repo/voxsphere"
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
	db.RegisterModel((*models.PostFlair)(nil))

	// drop all rows of the topic,voxsphere,post,post_flairs table
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
	if _, err := db.NewTruncateTable().Cascade().Model((*models.PostFlair)(nil)).Exec(context.Background()); err != nil {
		t.Fatal("truncate table failed:", err)
	}

	// load fixture
	fixture := dbfixture.New(db)
	if err := fixture.Load(context.Background(), os.DirFS("testdata"), fixtureFiles...); err != nil {
		t.Fatal("failed to load fixtures", err)
	}

	return db
}

func assertPostFlairs(t *testing.T, wantPostFlairs, gotPostFlairs []models.PostFlair) {
	t.Helper()

	if len(wantPostFlairs) != len(gotPostFlairs) {
		t.Fatal("length of wantPostFlairs and gotPostFlairs do not match")
	}

	for _, wantpostFlair := range wantPostFlairs {
		idx := slices.IndexFunc(gotPostFlairs, func(v models.PostFlair) bool {
			return v.ID == wantpostFlair.ID
		})

		if idx == -1 {
			t.Fatal(fmt.Sprintf("post flair %v of ID %v is not present in gotPostFlairs", wantpostFlair.FullText, wantpostFlair.ID))
			return
		}
		assert.Equal(t, wantpostFlair, gotPostFlairs[idx], "expect post flair to match")
	}
}

func TestRepo_PostFlairs(t *testing.T) {
	tests := []struct {
		name           string
		fixtureFiles   []string
		wantPostFlairs []models.PostFlair
		wantErr        error
	}{
		{
			name:         "post flairs :POS",
			fixtureFiles: []string{"topics.yml", "users.yml", "voxspheres.yml", "posts.yml", "post_flairs.yml"},
			wantPostFlairs: []models.PostFlair{
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					VoxsphereID:     uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					FullText:        "Flair 1",
					BackgroundColor: "#FFFFFF",
				},
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					VoxsphereID:     uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					FullText:        "Flair 2",
					BackgroundColor: "#000000",
				},
			},
			wantErr: nil,
		},
		{
			name:           "no post flairs :POS",
			fixtureFiles:   []string{},
			wantPostFlairs: nil,
			wantErr:        nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupPostgres(t, tt.fixtureFiles...)
			pgrepo := postflaireRepo.NewRepo(db)

			gotPostFlairs, gotErr := pgrepo.PostFlairs(context.Background())

			assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")
			assertPostFlairs(t, tt.wantPostFlairs, gotPostFlairs)
		})
	}
}

func TestRepo_PostFlairByID(t *testing.T) {
	type args struct {
		ID uuid.UUID
	}
	tests := []struct {
		name          string
		fixtureFiles  []string
		args          args
		wantPostFlair models.PostFlair
		wantErr       error
	}{
		{
			name:         "post flair not found :NEG",
			fixtureFiles: []string{"topics.yml", "users.yml", "voxspheres.yml", "posts.yml", "post_flairs.yml"},
			args: args{
				ID: uuid.MustParse("00000000-0000-0000-0000-000000000006"),
			},
			wantPostFlair: models.PostFlair{},
			wantErr:       postflaireRepo.ErrPostFlairNotFound,
		},
		{
			name:         "post flair found :POS",
			fixtureFiles: []string{"topics.yml", "users.yml", "voxspheres.yml", "posts.yml", "post_flairs.yml"},
			args: args{
				ID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			},
			wantPostFlair: models.PostFlair{
				ID:              uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				VoxsphereID:     uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				FullText:        "Flair 1",
				BackgroundColor: "#FFFFFF",
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupPostgres(t, tt.fixtureFiles...)
			pgrepo := postflaireRepo.NewRepo(db)

			gotPostFlair, gotErr := pgrepo.PostFlairByID(context.Background(), tt.args.ID)

			assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")
			assert.Equal(t, tt.wantPostFlair, gotPostFlair, "expect post flair to match")
		})
	}
}

func TestRepo_AddPostFlair(t *testing.T) {
	type args struct {
		postFlair models.PostFlair
	}
	tests := []struct {
		name           string
		fixtureFiles   []string
		args           args
		wantPostFlair  models.PostFlair
		wantPostFlairs []models.PostFlair
		wantErr        error
	}{
		{
			name:         "duplicate post flair id :NEG",
			fixtureFiles: []string{"topics.yml", "users.yml", "voxspheres.yml", "posts.yml", "post_flairs.yml"},
			args: args{
				postFlair: models.PostFlair{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					VoxsphereID:     uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					FullText:        "new text",
					BackgroundColor: "#FFFFFF",
				},
			},
			wantPostFlair: models.PostFlair{},
			wantPostFlairs: []models.PostFlair{
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					VoxsphereID:     uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					FullText:        "Flair 1",
					BackgroundColor: "#FFFFFF",
				},
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					VoxsphereID:     uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					FullText:        "Flair 2",
					BackgroundColor: "#000000",
				},
			},
			wantErr: postflaireRepo.ErrPostFlairDuplicateID,
		},
		{
			name:         "duplicate post id :NEG",
			fixtureFiles: []string{"topics.yml", "users.yml", "voxspheres.yml", "posts.yml", "post_flairs.yml"},
			args: args{
				postFlair: models.PostFlair{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					VoxsphereID:     uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					FullText:        "new text",
					BackgroundColor: "#FFFFFF",
				},
			},
			wantPostFlair: models.PostFlair{},
			wantPostFlairs: []models.PostFlair{
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					VoxsphereID:     uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					FullText:        "Flair 1",
					BackgroundColor: "#FFFFFF",
				},
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					VoxsphereID:     uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					FullText:        "Flair 2",
					BackgroundColor: "#000000",
				},
			},
			wantErr: postflaireRepo.ErrPostFlairDuplicateID,
		},
		{
			name:         "post not present in parent table :NEG",
			fixtureFiles: []string{"topics.yml", "users.yml", "voxspheres.yml", "posts.yml", "post_flairs.yml"},
			args: args{
				postFlair: models.PostFlair{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000006"),
					VoxsphereID:     uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					FullText:        "new text",
					BackgroundColor: "#FFFFFF",
				},
			},
			wantPostFlair: models.PostFlair{},
			wantPostFlairs: []models.PostFlair{
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					VoxsphereID:     uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					FullText:        "Flair 1",
					BackgroundColor: "#FFFFFF",
				},
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					VoxsphereID:     uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					FullText:        "Flair 2",
					BackgroundColor: "#000000",
				},
			},
			wantErr: postflaireRepo.ErrPostFlairParentTableRecordNotFound,
		},
		{
			name:         "voxsphere not present in parent table :NEG",
			fixtureFiles: []string{"topics.yml", "users.yml", "voxspheres.yml", "posts.yml", "post_flairs.yml"},
			args: args{
				postFlair: models.PostFlair{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					VoxsphereID:     uuid.MustParse("00000000-0000-0000-0000-000000000009"),
					FullText:        "new text",
					BackgroundColor: "#FFFFFF",
				},
			},
			wantPostFlair: models.PostFlair{},
			wantPostFlairs: []models.PostFlair{
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					VoxsphereID:     uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					FullText:        "Flair 1",
					BackgroundColor: "#FFFFFF",
				},
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					VoxsphereID:     uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					FullText:        "Flair 2",
					BackgroundColor: "#000000",
				},
			},
			wantErr: postflaireRepo.ErrPostFlairParentTableRecordNotFound,
		},
		{
			name:         "add post flair :POS",
			fixtureFiles: []string{"topics.yml", "users.yml", "voxspheres.yml", "posts.yml", "post_flairs.yml"},
			args: args{
				postFlair: models.PostFlair{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					VoxsphereID:     uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					FullText:        "new text",
					BackgroundColor: "#FFFFFF",
				},
			},
			wantPostFlair: models.PostFlair{
				ID:              uuid.MustParse("00000000-0000-0000-0000-000000000003"),
				PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000003"),
				VoxsphereID:     uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				FullText:        "new text",
				BackgroundColor: "#FFFFFF",
			},
			wantPostFlairs: []models.PostFlair{
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					VoxsphereID:     uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					FullText:        "Flair 1",
					BackgroundColor: "#FFFFFF",
				},
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					VoxsphereID:     uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					FullText:        "Flair 2",
					BackgroundColor: "#000000",
				},
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					VoxsphereID:     uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					FullText:        "new text",
					BackgroundColor: "#FFFFFF",
				},
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupPostgres(t, tt.fixtureFiles...)
			pgrepo := postflaireRepo.NewRepo(db)

			gotPostFlair, gotErr := pgrepo.AddPostFlair(context.Background(), tt.args.postFlair)

			assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")
			assert.Equal(t, tt.wantPostFlair, gotPostFlair, "expect post flair to match")

			gotPostFlairs, err := pgrepo.PostFlairs(context.Background())

			assert.NoError(t, err, "expect no error while getting post flairs")
			assertPostFlairs(t, tt.wantPostFlairs, gotPostFlairs)
		})
	}
}

func TestRepo_UpdatePostFlair(t *testing.T) {
	type args struct {
		postFlair models.PostFlair
	}
	tests := []struct {
		name           string
		fixtureFiles   []string
		args           args
		wantPostFlair  models.PostFlair
		wantPostFlairs []models.PostFlair
		wantErr        error
	}{
		{
			name:         "post flair not found :NEG",
			fixtureFiles: []string{"topics.yml", "users.yml", "voxspheres.yml", "posts.yml", "post_flairs.yml"},
			args: args{
				postFlair: models.PostFlair{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000009"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					VoxsphereID:     uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					FullText:        "updated text",
					BackgroundColor: "#FFFFFF",
				},
			},
			wantPostFlair: models.PostFlair{},
			wantPostFlairs: []models.PostFlair{
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					VoxsphereID:     uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					FullText:        "Flair 1",
					BackgroundColor: "#FFFFFF",
				},
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					VoxsphereID:     uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					FullText:        "Flair 2",
					BackgroundColor: "#000000",
				},
			},
			wantErr: postflaireRepo.ErrPostFlairNotFound,
		},
		{
			name:         "post not present in parent table :NEG",
			fixtureFiles: []string{"topics.yml", "users.yml", "voxspheres.yml", "posts.yml", "post_flairs.yml"},
			args: args{
				postFlair: models.PostFlair{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000009"),
					VoxsphereID:     uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					FullText:        "new text",
					BackgroundColor: "#FFFFFF",
				},
			},
			wantPostFlair: models.PostFlair{},
			wantPostFlairs: []models.PostFlair{
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					VoxsphereID:     uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					FullText:        "Flair 1",
					BackgroundColor: "#FFFFFF",
				},
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					VoxsphereID:     uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					FullText:        "Flair 2",
					BackgroundColor: "#000000",
				},
			},
			wantErr: postflaireRepo.ErrPostFlairParentTableRecordNotFound,
		},
		{
			name:         "voxsphere not present in parent table :NEG",
			fixtureFiles: []string{"topics.yml", "users.yml", "voxspheres.yml", "posts.yml", "post_flairs.yml"},
			args: args{
				postFlair: models.PostFlair{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					VoxsphereID:     uuid.MustParse("00000000-0000-0000-0000-000000000009"),
					FullText:        "new text",
					BackgroundColor: "#FFFFFF",
				},
			},
			wantPostFlair: models.PostFlair{},
			wantPostFlairs: []models.PostFlair{
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					VoxsphereID:     uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					FullText:        "Flair 1",
					BackgroundColor: "#FFFFFF",
				},
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					VoxsphereID:     uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					FullText:        "Flair 2",
					BackgroundColor: "#000000",
				},
			},
			wantErr: postflaireRepo.ErrPostFlairParentTableRecordNotFound,
		},
		{
			name:         "update post flair :POS",
			fixtureFiles: []string{"topics.yml", "users.yml", "voxspheres.yml", "posts.yml", "post_flairs.yml"},
			args: args{
				postFlair: models.PostFlair{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					VoxsphereID:     uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					FullText:        "updated text",
					BackgroundColor: "#000000",
				},
			},
			wantPostFlair: models.PostFlair{
				ID:              uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				VoxsphereID:     uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				FullText:        "updated text",
				BackgroundColor: "#000000",
			},
			wantPostFlairs: []models.PostFlair{
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					VoxsphereID:     uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					FullText:        "updated text",
					BackgroundColor: "#000000",
				},
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					VoxsphereID:     uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					FullText:        "Flair 2",
					BackgroundColor: "#000000",
				},
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupPostgres(t, tt.fixtureFiles...)
			pgrepo := postflaireRepo.NewRepo(db)

			gotPostFlair, gotErr := pgrepo.UpdatePostFlair(context.Background(), tt.args.postFlair)

			assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")
			assert.Equal(t, tt.wantPostFlair, gotPostFlair, "expect post flair to match")

			gotPostFlairs, err := pgrepo.PostFlairs(context.Background())

			assert.NoError(t, err, "expect no error while getting post flairs")
			assertPostFlairs(t, tt.wantPostFlairs, gotPostFlairs)
		})
	}
}

func TestRepo_DeletePostFlair(t *testing.T) {
	type args struct {
		ID uuid.UUID
	}
	tests := []struct {
		name           string
		fixtureFiles   []string
		args           args
		wantPostFlairs []models.PostFlair
		wantErr        error
	}{
		{
			name:         "post flair not found :NEG",
			fixtureFiles: []string{"topics.yml", "users.yml", "voxspheres.yml", "posts.yml", "post_flairs.yml"},
			args: args{
				ID: uuid.MustParse("00000000-0000-0000-0000-000000000006"),
			},
			wantPostFlairs: []models.PostFlair{
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					VoxsphereID:     uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					FullText:        "Flair 1",
					BackgroundColor: "#FFFFFF",
				},
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					VoxsphereID:     uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					FullText:        "Flair 2",
					BackgroundColor: "#000000",
				},
			},
			wantErr: postflaireRepo.ErrPostFlairNotFound,
		},
		{
			name:         "post flair deleted :POS",
			fixtureFiles: []string{"topics.yml", "users.yml", "voxspheres.yml", "posts.yml", "post_flairs.yml"},
			args: args{
				ID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			},
			wantPostFlairs: []models.PostFlair{
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					VoxsphereID:     uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					FullText:        "Flair 2",
					BackgroundColor: "#000000",
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupPostgres(t, tt.fixtureFiles...)
			pgrepo := postflaireRepo.NewRepo(db)

			gotErr := pgrepo.DeletePostFlair(context.Background(), tt.args.ID)

			assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")

			gotPostFlairs, err := pgrepo.PostFlairs(context.Background())

			assert.NoError(t, err, "expect no error while getting post flairs")
			assert.Equal(t, tt.wantPostFlairs, gotPostFlairs, "expect flairs to match")
		})
	}
}

func TestRepo_ForeignKeyCascade(t *testing.T) {
	t.Run("on deleting voxsphere from parent table , no child references should exist in post_flairs table", func(t *testing.T) {
		db := setupPostgres(t, "topics.yml", "voxspheres.yml", "users.yml", "posts.yml", "post_flairs.yml")
		postFlairPgrepo := postflaireRepo.NewRepo(db)
		voxspherePgrepo := voxrepo.NewRepo(db)

		wantPostFlairs := []models.PostFlair{
			{
				ID:              uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				VoxsphereID:     uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				FullText:        "Flair 2",
				BackgroundColor: "#000000",
			},
		}

		err := voxspherePgrepo.DeleteVoxsphere(context.Background(), uuid.MustParse("00000000-0000-0000-0000-000000000001"))

		assert.NoError(t, err, "expect no error while deleting voxsphere")

		gotPostFlairs, err := postFlairPgrepo.PostFlairs(context.Background())

		assert.NoError(t, err, "expect no error while getting post flairs")
		assertPostFlairs(t, wantPostFlairs, gotPostFlairs)
	})

	t.Run("on deleting post from parent table , no child references should exist in post_flairs table", func(t *testing.T) {
		db := setupPostgres(t, "topics.yml", "voxspheres.yml", "users.yml", "posts.yml", "post_flairs.yml")
		postFlairPgrepo := postflaireRepo.NewRepo(db)
		postPgrepo := postRepo.NewRepo(db)

		wantPostFlairs := []models.PostFlair{
			{
				ID:              uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				VoxsphereID:     uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				FullText:        "Flair 2",
				BackgroundColor: "#000000",
			},
		}

		err := postPgrepo.DeletePost(context.Background(), uuid.MustParse("00000000-0000-0000-0000-000000000001"))

		assert.NoError(t, err, "expect no error while deleting post")

		gotPostFlairs, err := postFlairPgrepo.PostFlairs(context.Background())

		assert.NoError(t, err, "expect no error while getting post flairs")
		assertPostFlairs(t, wantPostFlairs, gotPostFlairs)
	})
}
