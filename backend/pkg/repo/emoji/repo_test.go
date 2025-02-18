package emoji_test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"slices"
	"testing"

	"github.com/glowfi/voxpopuli/backend/pkg/models"
	emojirepo "github.com/glowfi/voxpopuli/backend/pkg/repo/emoji"
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

	db.RegisterModel((*models.Emoji)(nil))

	// drop all rows of the emoji table
	_, err := db.NewTruncateTable().Cascade().Model((*models.Emoji)(nil)).Exec(context.Background())
	if err != nil {
		t.Fatal("truncate table failed:", err)
	}

	// load fixture
	fixture := dbfixture.New(db)
	if err := fixture.Load(context.Background(), os.DirFS("testdata"), fixtureFiles...); err != nil {
		t.Fatal("failed to load fixtures", err)
	}

	// add query logging hook
	db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))

	return db
}

func assertEmojis(t *testing.T, wantEmojis, gotEmojis []models.Emoji) {
	t.Helper()

	for _, emoji := range wantEmojis {
		idx := slices.IndexFunc(gotEmojis, func(v models.Emoji) bool {
			return v.ID == emoji.ID
		})

		if idx == -1 {
			t.Fatal(fmt.Sprintf("emoji %v of ID %v is not present in gotEmojis", emoji.Title, emoji.ID))
			return
		}
		assert.Equal(t, emoji, gotEmojis[idx], "expect emoji to match")
	}
}

func TestRepo_Emojis(t *testing.T) {
	tests := []struct {
		name         string
		fixtureFiles []string
		wantEmojis   []models.Emoji
		wantErr      error
	}{
		{
			name:         "emojis :POS",
			fixtureFiles: []string{"emojis.yml"},
			wantEmojis: []models.Emoji{
				{
					ID:    uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Title: "emoji_foo",
				},
				{
					ID:    uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					Title: "emoji_bar",
				},
			},
			wantErr: nil,
		},
		{
			name:         "no emojis :POS",
			fixtureFiles: []string{},
			wantEmojis:   nil,
			wantErr:      nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupPostgres(t, tt.fixtureFiles...)
			pgrepo := emojirepo.NewRepo(db)

			gotEmojis, gotErr := pgrepo.Emojis(context.Background())

			assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")
			assert.Equal(t, tt.wantEmojis, gotEmojis, "expect emojis to match")
		})
	}
}

func TestRepo_EmojiByID(t *testing.T) {
	type args struct {
		ID uuid.UUID
	}
	tests := []struct {
		name         string
		fixtureFiles []string
		args         args
		wantEmoji    models.Emoji
		wantErr      error
	}{
		{
			name:         "emoji not found :NEG",
			fixtureFiles: []string{},
			args: args{
				ID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			},
			wantEmoji: models.Emoji{},
			wantErr:   emojirepo.ErrEmojiNotFound,
		},
		{
			name:         "emoji found :POS",
			fixtureFiles: []string{"emojis.yml"},
			args: args{
				ID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			},
			wantEmoji: models.Emoji{
				ID:    uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				Title: "emoji_foo",
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupPostgres(t, tt.fixtureFiles...)
			pgrepo := emojirepo.NewRepo(db)

			gotEmoji, gotErr := pgrepo.EmojiByID(context.Background(), tt.args.ID)

			assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")
			assert.Equal(t, tt.wantEmoji, gotEmoji, "expect emoji to match")
		})
	}
}

func TestRepo_AddEmoji(t *testing.T) {
	type args struct {
		emoji models.Emoji
	}
	tests := []struct {
		name         string
		fixtureFiles []string
		args         args
		wantEmoji    models.Emoji
		wantEmojis   []models.Emoji
		wantErr      error
	}{
		{
			name:         "duplicate emoji id :NEG",
			fixtureFiles: []string{"emojis.yml"},
			args: args{
				emoji: models.Emoji{
					ID:    uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Title: "new text",
				},
			},
			wantEmoji: models.Emoji{},
			wantEmojis: []models.Emoji{
				{
					ID:    uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Title: "emoji_foo",
				},
				{
					ID:    uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					Title: "emoji_bar",
				},
			},
			wantErr: emojirepo.ErrEmojiDuplicateIDorText,
		},
		{
			name:         "add emoji :POS",
			fixtureFiles: []string{"emojis.yml"},
			args: args{
				emoji: models.Emoji{
					ID:    uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					Title: "new text",
				},
			},
			wantEmoji: models.Emoji{
				ID:    uuid.MustParse("00000000-0000-0000-0000-000000000003"),
				Title: "new text",
			},
			wantEmojis: []models.Emoji{
				{
					ID:    uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Title: "emoji_foo",
				},
				{
					ID:    uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					Title: "emoji_bar",
				},
				{
					ID:    uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					Title: "new text",
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupPostgres(t, tt.fixtureFiles...)
			pgrepo := emojirepo.NewRepo(db)

			gotEmoji, gotErr := pgrepo.AddEmoji(context.Background(), tt.args.emoji)
			gotEmojis, err := pgrepo.Emojis(context.Background())
			if err != nil {
				t.Fatal("expect no error while getting emojis")
			}

			assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")
			assert.Equal(t, tt.wantEmoji, gotEmoji, "expect emoji to match")
			assertEmojis(t, tt.wantEmojis, gotEmojis)
		})
	}
}

func TestRepo_UpdateEmoji(t *testing.T) {
	type args struct {
		emoji models.Emoji
	}
	tests := []struct {
		name         string
		fixtureFiles []string
		args         args
		wantEmoji    models.Emoji
		wantEmojis   []models.Emoji
		wantErr      error
	}{
		{
			name:         "emoji id not found :NEG",
			fixtureFiles: []string{"emojis.yml"},
			args: args{
				emoji: models.Emoji{
					ID:    uuid.MustParse("00000000-0000-0000-0000-000000000006"),
					Title: "updated text",
				},
			},
			wantEmoji: models.Emoji{},
			wantEmojis: []models.Emoji{
				{
					ID:    uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Title: "emoji_foo",
				},
				{
					ID:    uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					Title: "emoji_bar",
				},
			},
			wantErr: emojirepo.ErrEmojiNotFound,
		},
		{
			name:         "emoji updated :POS",
			fixtureFiles: []string{"emojis.yml"},
			args: args{
				emoji: models.Emoji{
					ID:    uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Title: "updated text",
				},
			},
			wantEmoji: models.Emoji{
				ID:    uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				Title: "updated text",
			},
			wantEmojis: []models.Emoji{
				{
					ID:    uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Title: "updated text",
				},
				{
					ID:    uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					Title: "emoji_bar",
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupPostgres(t, tt.fixtureFiles...)
			pgrepo := emojirepo.NewRepo(db)

			gotEmoji, gotErr := pgrepo.UpdateEmoji(context.Background(), tt.args.emoji)
			gotEmojis, err := pgrepo.Emojis(context.Background())
			if err != nil {
				t.Fatal("expect no error while getting emojis")
			}

			assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")
			assert.Equal(t, tt.wantEmoji, gotEmoji, "expect emoji to match")
			assertEmojis(t, tt.wantEmojis, gotEmojis)
		})
	}
}

func TestRepo_DeleteEmoji(t *testing.T) {
	type args struct {
		ID uuid.UUID
	}
	tests := []struct {
		name         string
		fixtureFiles []string
		args         args
		wantEmojis   []models.Emoji
		wantErr      error
	}{
		{
			name:         "emoji id not found :NEG",
			fixtureFiles: []string{"emojis.yml"},
			args: args{
				ID: uuid.MustParse("00000000-0000-0000-0000-000000000006"),
			},
			wantEmojis: []models.Emoji{
				{
					ID:    uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Title: "emoji_foo",
				},
				{
					ID:    uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					Title: "emoji_bar",
				},
			},
			wantErr: emojirepo.ErrEmojiNotFound,
		},
		{
			name:         "emoji deleted :POS",
			fixtureFiles: []string{"emojis.yml"},
			args: args{
				ID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			},
			wantEmojis: []models.Emoji{
				{
					ID:    uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					Title: "emoji_bar",
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupPostgres(t, tt.fixtureFiles...)
			pgrepo := emojirepo.NewRepo(db)

			gotErr := pgrepo.DeleteEmoji(context.Background(), tt.args.ID)
			gotEmojis, err := pgrepo.Emojis(context.Background())
			if err != nil {
				t.Fatal("expect no error while getting emojis")
			}

			assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")
			assert.Equal(t, tt.wantEmojis, gotEmojis, "expect emojis to match")
		})
	}
}
