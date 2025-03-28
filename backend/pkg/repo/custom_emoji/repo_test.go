package customemoji_test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"slices"
	"testing"

	"github.com/glowfi/voxpopuli/backend/pkg/models"
	customemojirepo "github.com/glowfi/voxpopuli/backend/pkg/repo/custom_emoji"
	"github.com/glowfi/voxpopuli/backend/pkg/repo/voxsphere"
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
	db.RegisterModel((*models.CustomEmoji)(nil))

	// drop all rows of the custom emoji,topic and voxsphere table
	if _, err := db.NewTruncateTable().Cascade().Model((*models.Topic)(nil)).Exec(context.Background()); err != nil {
		t.Fatal("truncate table failed:", err)
	}
	if _, err := db.NewTruncateTable().Cascade().Model((*models.Voxsphere)(nil)).Exec(context.Background()); err != nil {
		t.Fatal("truncate table failed:", err)
	}
	if _, err := db.NewTruncateTable().Cascade().Model((*models.CustomEmoji)(nil)).Exec(context.Background()); err != nil {
		t.Fatal("truncate table failed:", err)
	}

	// load fixture
	fixture := dbfixture.New(db)
	if err := fixture.Load(context.Background(), os.DirFS("testdata"), fixtureFiles...); err != nil {
		t.Fatal("failed to load fixtures", err)
	}

	return db
}

func assertCustomEmojis(t *testing.T, wantCustomEmojis, gotCustomEmojis []models.CustomEmoji) {
	t.Helper()

	if len(wantCustomEmojis) != len(gotCustomEmojis) {
		t.Fatal("length of wantCustomEmojis and gotCustomEmojis do not match")
	}

	for _, customemoji := range wantCustomEmojis {
		idx := slices.IndexFunc(gotCustomEmojis, func(v models.CustomEmoji) bool {
			return v.ID == customemoji.ID
		})

		if idx == -1 {
			t.Fatalf("custom emoji %v of ID %v is not present in gotCustomEmojis", customemoji.Title, customemoji.ID)
			return
		}
		assert.Equal(t, customemoji, gotCustomEmojis[idx], "expect custom emoji to match")
	}
}

func TestRepo_CustomEmojis(t *testing.T) {
	tests := []struct {
		name             string
		fixtureFiles     []string
		wantCustomEmojis []models.CustomEmoji
		wantErr          error
	}{
		{
			name:         "custom emojis :POS",
			fixtureFiles: []string{"topics.yml", "voxspheres.yml", "custom_emojis.yml"},
			wantCustomEmojis: []models.CustomEmoji{
				{
					ID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					VoxsphereID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Url:         "https://example.com/emoji1.png",
					Title:       "emoji1",
				},
				{
					ID:          uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					VoxsphereID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					Url:         "https://example.com/emoji2.png",
					Title:       "emoji2",
				},
			},
			wantErr: nil,
		},
		{
			name:             "no custom emojis :POS",
			fixtureFiles:     []string{},
			wantCustomEmojis: nil,
			wantErr:          nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupPostgres(t, tt.fixtureFiles...)
			pgrepo := customemojirepo.NewRepo(db)

			gotCustomEmojis, gotErr := pgrepo.CustomEmojis(context.Background())

			assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")
			assertCustomEmojis(t, tt.wantCustomEmojis, gotCustomEmojis)
		})
	}
}

func TestRepo_CustomEmojiByID(t *testing.T) {
	type args struct {
		ID uuid.UUID
	}
	tests := []struct {
		name            string
		fixtureFiles    []string
		args            args
		wantCustomEmoji models.CustomEmoji
		wantErr         error
	}{
		{
			name:         "custom emoji not found :NEG",
			fixtureFiles: []string{"topics.yml", "voxspheres.yml", "custom_emojis.yml"},
			args: args{
				ID: uuid.MustParse("00000000-0000-0000-0000-000000000003"),
			},
			wantCustomEmoji: models.CustomEmoji{},
			wantErr:         customemojirepo.ErrCustomEmojiNotFound,
		},
		{
			name:         "custom emoji found :POS",
			fixtureFiles: []string{"topics.yml", "voxspheres.yml", "custom_emojis.yml"},
			args: args{
				ID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			},
			wantCustomEmoji: models.CustomEmoji{
				ID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				VoxsphereID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				Url:         "https://example.com/emoji1.png",
				Title:       "emoji1",
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupPostgres(t, tt.fixtureFiles...)
			pgrepo := customemojirepo.NewRepo(db)

			gotCustomEmoji, gotErr := pgrepo.CustomEmojiByID(context.Background(), tt.args.ID)

			assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")
			assert.Equal(t, tt.wantCustomEmoji, gotCustomEmoji, "expect custom emoji to match")
		})
	}
}

func TestRepo_AddCustomEmojis(t *testing.T) {
	type args struct {
		customEmojis []models.CustomEmoji
	}
	tests := []struct {
		name                     string
		fixtureFiles             []string
		args                     args
		wantInsertedCustomEmojis []models.CustomEmoji
		wantCustomEmojis         []models.CustomEmoji
		wantErr                  error
	}{
		{
			name:         "add custom emoji with duplicate ID :NEG",
			fixtureFiles: []string{"topics.yml", "voxspheres.yml", "custom_emojis.yml"},
			args: args{
				customEmojis: []models.CustomEmoji{
					{
						ID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
						VoxsphereID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
						Url:         "https://example.com/emoji1.png",
						Title:       "emoji1",
					},
				},
			},
			wantInsertedCustomEmojis: nil,
			wantCustomEmojis: []models.CustomEmoji{
				{
					ID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					VoxsphereID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Url:         "https://example.com/emoji1.png",
					Title:       "emoji1",
				},
				{
					ID:          uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					VoxsphereID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					Url:         "https://example.com/emoji2.png",
					Title:       "emoji2",
				},
			},
			wantErr: customemojirepo.ErrCustomEmojiDuplicateID,
		},
		{
			name:         "voxsphere does not exist in the parent table :NEG",
			fixtureFiles: []string{"topics.yml", "voxspheres.yml", "custom_emojis.yml"},
			args: args{
				customEmojis: []models.CustomEmoji{
					{
						ID:          uuid.MustParse("00000000-0000-0000-0000-000000000003"),
						VoxsphereID: uuid.MustParse("00000000-0000-0000-0000-000000000009"),
						Url:         "https://example.com/emoji3.png",
						Title:       "emoji3",
					},
				},
			},
			wantInsertedCustomEmojis: nil,
			wantCustomEmojis: []models.CustomEmoji{
				{
					ID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					VoxsphereID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Url:         "https://example.com/emoji1.png",
					Title:       "emoji1",
				},
				{
					ID:          uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					VoxsphereID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					Url:         "https://example.com/emoji2.png",
					Title:       "emoji2",
				},
			},
			wantErr: customemojirepo.ErrCustomEmojiParentTableRecordNotFound,
		},
		{
			name:         "add custom emojis :POS",
			fixtureFiles: []string{"topics.yml", "voxspheres.yml", "custom_emojis.yml"},
			args: args{
				customEmojis: []models.CustomEmoji{
					{
						ID:          uuid.MustParse("00000000-0000-0000-0000-000000000003"),
						VoxsphereID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
						Url:         "https://example.com/emoji3.png",
						Title:       "emoji3",
					},
				},
			},
			wantInsertedCustomEmojis: []models.CustomEmoji{
				{
					ID:          uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					VoxsphereID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Url:         "https://example.com/emoji3.png",
					Title:       "emoji3",
				},
			},
			wantCustomEmojis: []models.CustomEmoji{
				{
					ID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					VoxsphereID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Url:         "https://example.com/emoji1.png",
					Title:       "emoji1",
				},
				{
					ID:          uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					VoxsphereID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					Url:         "https://example.com/emoji2.png",
					Title:       "emoji2",
				},
				{
					ID:          uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					VoxsphereID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Url:         "https://example.com/emoji3.png",
					Title:       "emoji3",
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupPostgres(t, tt.fixtureFiles...)
			pgrepo := customemojirepo.NewRepo(db)

			gotInsertedCustomEmojis, gotErr := pgrepo.AddCustomEmojis(context.Background(), tt.args.customEmojis...)

			assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")
			assert.Equal(t, tt.wantInsertedCustomEmojis, gotInsertedCustomEmojis, "expect inserted custom emojis to match")

			gotCustomEmojis, err := pgrepo.CustomEmojis(context.Background())

			assert.NoError(t, err, "expect no error while getting custom emojis")
			assertCustomEmojis(t, tt.wantCustomEmojis, gotCustomEmojis)
		})
	}
}

func TestRepo_UpdateCustomEmoji(t *testing.T) {
	type args struct {
		emoji models.CustomEmoji
	}
	tests := []struct {
		name             string
		fixtureFiles     []string
		args             args
		wantCustomEmoji  models.CustomEmoji
		wantCustomEmojis []models.CustomEmoji
		wantErr          error
	}{
		{
			name:         "custom emoji not found :NEG",
			fixtureFiles: []string{"topics.yml", "voxspheres.yml", "custom_emojis.yml"},
			args: args{
				emoji: models.CustomEmoji{
					ID:          uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					VoxsphereID: uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					Url:         "https://example.com/updated_emoji3.png",
					Title:       "updated_emoji3",
				},
			},
			wantCustomEmoji: models.CustomEmoji{},
			wantCustomEmojis: []models.CustomEmoji{
				{
					ID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					VoxsphereID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Url:         "https://example.com/emoji1.png",
					Title:       "emoji1",
				},
				{
					ID:          uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					VoxsphereID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					Url:         "https://example.com/emoji2.png",
					Title:       "emoji2",
				},
			},
			wantErr: customemojirepo.ErrCustomEmojiNotFound,
		},
		{
			name:         "voxsphere is not present in parent table :NEG",
			fixtureFiles: []string{"topics.yml", "voxspheres.yml", "custom_emojis.yml"},
			args: args{
				emoji: models.CustomEmoji{
					ID:          uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					VoxsphereID: uuid.MustParse("00000000-0000-0000-0000-000000000009"),
					Url:         "https://example.com/updated_emoji3.png",
					Title:       "updated_emoji3",
				},
			},
			wantCustomEmoji: models.CustomEmoji{},
			wantCustomEmojis: []models.CustomEmoji{
				{
					ID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					VoxsphereID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Url:         "https://example.com/emoji1.png",
					Title:       "emoji1",
				},
				{
					ID:          uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					VoxsphereID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					Url:         "https://example.com/emoji2.png",
					Title:       "emoji2",
				},
			},
			wantErr: customemojirepo.ErrCustomEmojiParentTableRecordNotFound,
		},
		{
			name:         "update custom emoji :POS",
			fixtureFiles: []string{"topics.yml", "voxspheres.yml", "custom_emojis.yml"},
			args: args{
				emoji: models.CustomEmoji{
					ID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					VoxsphereID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Url:         "https://example.com/updated_emoji1.png",
					Title:       "updated_emoji1",
				},
			},
			wantCustomEmoji: models.CustomEmoji{
				ID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				VoxsphereID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				Url:         "https://example.com/updated_emoji1.png",
				Title:       "updated_emoji1",
			},
			wantCustomEmojis: []models.CustomEmoji{
				{
					ID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					VoxsphereID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Url:         "https://example.com/updated_emoji1.png",
					Title:       "updated_emoji1",
				},
				{
					ID:          uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					VoxsphereID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					Url:         "https://example.com/emoji2.png",
					Title:       "emoji2",
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupPostgres(t, tt.fixtureFiles...)
			pgrepo := customemojirepo.NewRepo(db)

			gotCustomEmoji, gotErr := pgrepo.UpdateCustomEmoji(context.Background(), tt.args.emoji)

			assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")
			assert.Equal(t, tt.wantCustomEmoji, gotCustomEmoji, "expect custom emoji to match")

			gotCustomEmojis, err := pgrepo.CustomEmojis(context.Background())

			assert.NoError(t, err, "expect no error while getting custom emojis")
			assertCustomEmojis(t, tt.wantCustomEmojis, gotCustomEmojis)
		})
	}
}

func TestRepo_DeleteCustomEmoji(t *testing.T) {
	type args struct {
		ID uuid.UUID
	}
	tests := []struct {
		name             string
		fixtureFiles     []string
		args             args
		wantCustomEmojis []models.CustomEmoji
		wantErr          error
	}{
		{
			name:         "custom emoji not found :NEG",
			fixtureFiles: []string{"topics.yml", "voxspheres.yml", "custom_emojis.yml"},
			args: args{
				ID: uuid.MustParse("00000000-0000-0000-0000-000000000003"),
			},
			wantCustomEmojis: []models.CustomEmoji{
				{
					ID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					VoxsphereID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Url:         "https://example.com/emoji1.png",
					Title:       "emoji1",
				},
				{
					ID:          uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					VoxsphereID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					Url:         "https://example.com/emoji2.png",
					Title:       "emoji2",
				},
			},
			wantErr: customemojirepo.ErrCustomEmojiNotFound,
		},
		{
			name:         "delete custom emoji :POS",
			fixtureFiles: []string{"topics.yml", "voxspheres.yml", "custom_emojis.yml"},
			args: args{
				ID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			},
			wantCustomEmojis: []models.CustomEmoji{
				{
					ID:          uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					VoxsphereID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					Url:         "https://example.com/emoji2.png",
					Title:       "emoji2",
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupPostgres(t, tt.fixtureFiles...)
			pgrepo := customemojirepo.NewRepo(db)

			gotErr := pgrepo.DeleteCustomEmoji(context.Background(), tt.args.ID)

			assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")

			gotCustomEmojis, err := pgrepo.CustomEmojis(context.Background())

			assert.NoError(t, err, "expect no error while getting custom emojis")
			assertCustomEmojis(t, tt.wantCustomEmojis, gotCustomEmojis)
		})
	}
}

func TestRepo_ForeignKeyCascade(t *testing.T) {
	t.Run("on deleting voxsphere from parent table , no child references should exist in custom_emoji table", func(t *testing.T) {
		db := setupPostgres(t, "topics.yml", "voxspheres.yml", "custom_emojis.yml")
		customEmojiPgrepo := customemojirepo.NewRepo(db)
		voxspherePgrepo := voxsphere.NewRepo(db)

		wantCustomEmojis := []models.CustomEmoji{
			{
				ID:          uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				VoxsphereID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				Url:         "https://example.com/emoji2.png",
				Title:       "emoji2",
			},
		}

		err := voxspherePgrepo.DeleteVoxsphere(context.Background(), uuid.MustParse("00000000-0000-0000-0000-000000000001"))

		assert.NoError(t, err, "expect no erro while deleting voxsphere")

		gotCustomEmojis, err := customEmojiPgrepo.CustomEmojis(context.Background())

		assert.NoError(t, err, "expect no erro while deleting voxsphere")
		assertCustomEmojis(t, wantCustomEmojis, gotCustomEmojis)
	})
}
