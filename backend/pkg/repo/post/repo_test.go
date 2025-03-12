package postrepo_test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"slices"
	"testing"
	"time"

	"github.com/glowfi/voxpopuli/backend/pkg/models"
	postrepo "github.com/glowfi/voxpopuli/backend/pkg/repo/post"
	userrepo "github.com/glowfi/voxpopuli/backend/pkg/repo/user"
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

	// drop all rows of the topics,voxspheres table
	_, err := db.NewTruncateTable().Cascade().Model((*models.Topic)(nil)).Exec(context.Background())
	if err != nil {
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

func assertPostWithoutTimestamp(t *testing.T, wantPost, gotPost models.Post) {
	assert.Equal(t, wantPost.ID, gotPost.ID, "expected id to match")
	assert.Equal(t, wantPost.AuthorID, gotPost.AuthorID, "expected author id to match")
	assert.Equal(t, wantPost.VoxsphereID, gotPost.VoxsphereID, "expected voxsphere id to match")
	assert.Equal(t, wantPost.Title, gotPost.Title, "expected title to match")
	assert.Equal(t, wantPost.Text, gotPost.Text, "expected text to match")
	assert.Equal(t, wantPost.TextHtml, gotPost.TextHtml, "expected text html to match")
	assert.Equal(t, wantPost.Ups, gotPost.Ups, "expected ups to match")
	assert.Equal(t, wantPost.Over18, gotPost.Over18, "expected over 18 to match")
	assert.Equal(t, wantPost.Spoiler, gotPost.Spoiler, "expected spoiler to match")
}

func assertPostsWithoutTimestamp(t *testing.T, wantPosts, gotPosts []models.Post) {
	t.Helper()

	if len(wantPosts) != len(gotPosts) {
		t.Fatal("length of wantPosts and gotPosts do not match")
	}

	for _, post := range wantPosts {
		idx := slices.IndexFunc(gotPosts, func(p models.Post) bool {
			return p.ID == post.ID
		})

		if idx == -1 {
			t.Fatal(fmt.Sprintf("post %v of ID %v is not present in gotPosts", post.Title, post.ID))
			return
		}
		assertPostWithoutTimestamp(t, post, gotPosts[idx])
	}
}

func TestRepo_Posts(t *testing.T) {
	tests := []struct {
		name         string
		fixtureFiles []string
		wantPosts    []models.Post
		wantErr      error
	}{
		{
			name:         "posts :POS",
			fixtureFiles: []string{"topics.yml", "voxspheres.yml", "users.yml", "posts.yml"},
			wantPosts: []models.Post{
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					AuthorID:      uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					VoxsphereID:   uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Title:         "Example Post Title 1",
					Text:          "This is an example post text 1.",
					TextHtml:      "<p>This is an example post text 1 in HTML.</p>",
					Ups:           10,
					Over18:        false,
					Spoiler:       false,
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix: 1725091100,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					AuthorID:      uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					VoxsphereID:   uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					Title:         "Example Post Title 2",
					Text:          "This is an example post text 2.",
					TextHtml:      "<p>This is an example post text 2 in HTML.</p>",
					Ups:           20,
					Over18:        true,
					Spoiler:       true,
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
					CreatedAtUnix: 1725091120,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
				},
			},
			wantErr: nil,
		},
		{
			name:         "no posts :POS",
			fixtureFiles: []string{},
			wantPosts:    nil,
			wantErr:      nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupPostgres(t, tt.fixtureFiles...)
			pgrepo := postrepo.NewRepo(db)

			gotPosts, gotErr := pgrepo.Posts(context.Background())

			assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")
			assert.Equal(t, tt.wantPosts, gotPosts, "expect posts to match")
		})
	}
}

func TestRepo_PostByID(t *testing.T) {
	type args struct {
		ID uuid.UUID
	}
	tests := []struct {
		name         string
		fixtureFiles []string
		args         args
		wantPost     models.Post
		wantErr      error
	}{
		{
			name:         "post not found :NEG",
			fixtureFiles: []string{"topics.yml", "voxspheres.yml", "users.yml", "posts.yml"},
			args: args{
				ID: uuid.MustParse("00000000-0000-0000-0000-000000000009"),
			},
			wantPost: models.Post{},
			wantErr:  postrepo.ErrPostNotFound,
		},
		{
			name:         "post found :POS",
			fixtureFiles: []string{"topics.yml", "voxspheres.yml", "users.yml", "posts.yml"},
			args: args{
				ID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			},
			wantPost: models.Post{
				ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				AuthorID:      uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				VoxsphereID:   uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				Title:         "Example Post Title 1",
				Text:          "This is an example post text 1.",
				TextHtml:      "<p>This is an example post text 1 in HTML.</p>",
				Ups:           10,
				Over18:        false,
				Spoiler:       false,
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
			pgrepo := postrepo.NewRepo(db)

			gotPost, gotErr := pgrepo.PostByID(context.Background(), tt.args.ID)

			assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")
			assert.Equal(t, tt.wantPost, gotPost, "expect post to match")
		})
	}
}

func TestRepo_AddPosts(t *testing.T) {
	type args struct {
		posts []models.Post
	}
	tests := []struct {
		name              string
		fixtureFiles      []string
		args              args
		wantInsertedPosts []models.Post
		wantPosts         []models.Post
		wantErr           error
	}{
		{
			name:         "duplicate post id :NEG",
			fixtureFiles: []string{"topics.yml", "voxspheres.yml", "users.yml", "posts.yml"},
			args: args{
				posts: []models.Post{
					{
						ID:            uuid.MustParse("00000000-0000-0000-0000-000000000002"),
						AuthorID:      uuid.MustParse("00000000-0000-0000-0000-000000000002"),
						VoxsphereID:   uuid.MustParse("00000000-0000-0000-0000-000000000002"),
						Title:         "Example Post Title 2",
						Text:          "This is an example post text 2.",
						TextHtml:      "<p>This is an example post text 2 in HTML.</p>",
						Ups:           20,
						Over18:        true,
						Spoiler:       true,
						CreatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
						CreatedAtUnix: 1725091120,
						UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
					},
				},
			},
			wantInsertedPosts: nil,
			wantPosts: []models.Post{
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					AuthorID:      uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					VoxsphereID:   uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Title:         "Example Post Title 1",
					Text:          "This is an example post text 1.",
					TextHtml:      "<p>This is an example post text 1 in HTML.</p>",
					Ups:           10,
					Over18:        false,
					Spoiler:       false,
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix: 1725091100,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					AuthorID:      uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					VoxsphereID:   uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					Title:         "Example Post Title 2",
					Text:          "This is an example post text 2.",
					TextHtml:      "<p>This is an example post text 2 in HTML.</p>",
					Ups:           20,
					Over18:        true,
					Spoiler:       true,
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
					CreatedAtUnix: 1725091120,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
				},
			},
			wantErr: postrepo.ErrPostDuplicateID,
		},
		{
			name:         "author does not exist in parent table :NEG",
			fixtureFiles: []string{"topics.yml", "voxspheres.yml", "users.yml", "posts.yml"},
			args: args{
				posts: []models.Post{
					{
						ID:            uuid.MustParse("00000000-0000-0000-0000-000000000003"),
						AuthorID:      uuid.MustParse("00000000-0000-0000-0000-000000000009"),
						VoxsphereID:   uuid.MustParse("00000000-0000-0000-0000-000000000002"),
						Title:         "Example Post Title 3",
						Text:          "This is an example post text 3.",
						TextHtml:      "<p>This is an example post text 3 in HTML.</p>",
						Ups:           30,
						Over18:        true,
						Spoiler:       true,
						CreatedAt:     time.Date(2024, 10, 10, 10, 10, 30, 0, time.UTC),
						CreatedAtUnix: 1725091120,
						UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 30, 0, time.UTC),
					},
				},
			},
			wantInsertedPosts: nil,
			wantPosts: []models.Post{
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					AuthorID:      uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					VoxsphereID:   uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Title:         "Example Post Title 1",
					Text:          "This is an example post text 1.",
					TextHtml:      "<p>This is an example post text 1 in HTML.</p>",
					Ups:           10,
					Over18:        false,
					Spoiler:       false,
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix: 1725091100,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					AuthorID:      uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					VoxsphereID:   uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					Title:         "Example Post Title 2",
					Text:          "This is an example post text 2.",
					TextHtml:      "<p>This is an example post text 2 in HTML.</p>",
					Ups:           20,
					Over18:        true,
					Spoiler:       true,
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
					CreatedAtUnix: 1725091120,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
				},
			},
			wantErr: postrepo.ErrPostParentTableRecordNotFound,
		},
		{
			name:         "voxsphere does not exist in parent table :NEG",
			fixtureFiles: []string{"topics.yml", "voxspheres.yml", "users.yml", "posts.yml"},
			args: args{
				posts: []models.Post{
					{
						ID:            uuid.MustParse("00000000-0000-0000-0000-000000000003"),
						AuthorID:      uuid.MustParse("00000000-0000-0000-0000-000000000001"),
						VoxsphereID:   uuid.MustParse("00000000-0000-0000-0000-000000000009"),
						Title:         "Example Post Title 3",
						Text:          "This is an example post text 3.",
						TextHtml:      "<p>This is an example post text 3 in HTML.</p>",
						Ups:           30,
						Over18:        true,
						Spoiler:       true,
						CreatedAt:     time.Date(2024, 10, 10, 10, 10, 30, 0, time.UTC),
						CreatedAtUnix: 1725091120,
						UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 30, 0, time.UTC),
					},
				},
			},
			wantInsertedPosts: nil,
			wantPosts: []models.Post{
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					AuthorID:      uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					VoxsphereID:   uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Title:         "Example Post Title 1",
					Text:          "This is an example post text 1.",
					TextHtml:      "<p>This is an example post text 1 in HTML.</p>",
					Ups:           10,
					Over18:        false,
					Spoiler:       false,
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix: 1725091100,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					AuthorID:      uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					VoxsphereID:   uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					Title:         "Example Post Title 2",
					Text:          "This is an example post text 2.",
					TextHtml:      "<p>This is an example post text 2 in HTML.</p>",
					Ups:           20,
					Over18:        true,
					Spoiler:       true,
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
					CreatedAtUnix: 1725091120,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
				},
			},
			wantErr: postrepo.ErrPostParentTableRecordNotFound,
		},
		{
			name:         "add posts :POS",
			fixtureFiles: []string{"topics.yml", "voxspheres.yml", "users.yml", "posts.yml"},
			args: args{
				posts: []models.Post{
					{
						ID:            uuid.MustParse("00000000-0000-0000-0000-000000000003"),
						AuthorID:      uuid.MustParse("00000000-0000-0000-0000-000000000001"),
						VoxsphereID:   uuid.MustParse("00000000-0000-0000-0000-000000000002"),
						Title:         "Example Post Title 3",
						Text:          "This is an example post text 3.",
						TextHtml:      "<p>This is an example post text 3 in HTML.</p>",
						Ups:           30,
						Over18:        true,
						Spoiler:       true,
						CreatedAt:     time.Date(2024, 10, 10, 10, 10, 30, 0, time.UTC),
						CreatedAtUnix: 1725091120,
						UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 30, 0, time.UTC),
					},
				},
			},
			wantInsertedPosts: []models.Post{
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					AuthorID:      uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					VoxsphereID:   uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					Title:         "Example Post Title 3",
					Text:          "This is an example post text 3.",
					TextHtml:      "<p>This is an example post text 3 in HTML.</p>",
					Ups:           30,
					Over18:        true,
					Spoiler:       true,
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 30, 0, time.UTC),
					CreatedAtUnix: 1725091120,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 30, 0, time.UTC),
				},
			},
			wantPosts: []models.Post{
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					AuthorID:      uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					VoxsphereID:   uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Title:         "Example Post Title 1",
					Text:          "This is an example post text 1.",
					TextHtml:      "<p>This is an example post text 1 in HTML.</p>",
					Ups:           10,
					Over18:        false,
					Spoiler:       false,
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix: 1725091100,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					AuthorID:      uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					VoxsphereID:   uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					Title:         "Example Post Title 2",
					Text:          "This is an example post text 2.",
					TextHtml:      "<p>This is an example post text 2 in HTML.</p>",
					Ups:           20,
					Over18:        true,
					Spoiler:       true,
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
					CreatedAtUnix: 1725091120,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
				},
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					AuthorID:      uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					VoxsphereID:   uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					Title:         "Example Post Title 3",
					Text:          "This is an example post text 3.",
					TextHtml:      "<p>This is an example post text 3 in HTML.</p>",
					Ups:           30,
					Over18:        true,
					Spoiler:       true,
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 30, 0, time.UTC),
					CreatedAtUnix: 1725091120,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 30, 0, time.UTC),
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupPostgres(t, tt.fixtureFiles...)
			pgrepo := postrepo.NewRepo(db)

			startTime := time.Now()
			gotInsertedPosts, gotErr := pgrepo.AddPosts(context.Background(), tt.args.posts...)
			endTime := time.Now()

			for _, gotInsertedPost := range gotInsertedPosts {
				assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")
				assert.Equal(
					t,
					gotInsertedPost.UpdatedAt,
					gotInsertedPost.CreatedAt,
					"expect CreatedAt and UpdatedAt to be same",
				)
				if tt.wantErr == nil {
					assertTimeWithinRange(t, gotInsertedPost.CreatedAt, startTime, endTime)
					assertTimeWithinRange(t, gotInsertedPost.UpdatedAt, startTime, endTime)
				}
			}

			gotPosts, err := pgrepo.Posts(context.Background())

			assert.NoError(t, err, "expect no error while getting posts")
			assertPostsWithoutTimestamp(t, tt.wantPosts, gotPosts)
		})
	}
}

func TestRepo_UpdatePost(t *testing.T) {
	type args struct {
		post models.Post
	}
	tests := []struct {
		name         string
		fixtureFiles []string
		args         args
		wantPost     models.Post
		wantPosts    []models.Post
		wantErr      error
	}{
		{
			name:         "post not found :NEG",
			fixtureFiles: []string{"topics.yml", "voxspheres.yml", "users.yml", "posts.yml"},
			args: args{
				post: models.Post{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					AuthorID:      uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					VoxsphereID:   uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					Title:         "Updated Title",
					Text:          "Update Text",
					TextHtml:      "<p>Updated HTML.</p>",
					Ups:           30,
					Over18:        true,
					Spoiler:       true,
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 30, 0, time.UTC),
					CreatedAtUnix: 1725091120,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 30, 0, time.UTC),
				},
			},
			wantPost: models.Post{},
			wantPosts: []models.Post{
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					AuthorID:      uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					VoxsphereID:   uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Title:         "Example Post Title 1",
					Text:          "This is an example post text 1.",
					TextHtml:      "<p>This is an example post text 1 in HTML.</p>",
					Ups:           10,
					Over18:        false,
					Spoiler:       false,
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix: 1725091100,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					AuthorID:      uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					VoxsphereID:   uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					Title:         "Example Post Title 2",
					Text:          "This is an example post text 2.",
					TextHtml:      "<p>This is an example post text 2 in HTML.</p>",
					Ups:           20,
					Over18:        true,
					Spoiler:       true,
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
					CreatedAtUnix: 1725091120,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
				},
			},
			wantErr: postrepo.ErrPostNotFound,
		},
		{
			name:         "author is not present in the parent table :NEG",
			fixtureFiles: []string{"topics.yml", "voxspheres.yml", "users.yml", "posts.yml"},
			args: args{
				post: models.Post{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					AuthorID:      uuid.MustParse("00000000-0000-0000-0000-000000000009"),
					VoxsphereID:   uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					Title:         "Updated Title",
					Text:          "Update Text",
					TextHtml:      "<p>Updated HTML.</p>",
					Ups:           30,
					Over18:        true,
					Spoiler:       true,
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 30, 0, time.UTC),
					CreatedAtUnix: 1725091120,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 30, 0, time.UTC),
				},
			},
			wantPost: models.Post{},
			wantPosts: []models.Post{
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					AuthorID:      uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					VoxsphereID:   uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Title:         "Example Post Title 1",
					Text:          "This is an example post text 1.",
					TextHtml:      "<p>This is an example post text 1 in HTML.</p>",
					Ups:           10,
					Over18:        false,
					Spoiler:       false,
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix: 1725091100,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					AuthorID:      uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					VoxsphereID:   uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					Title:         "Example Post Title 2",
					Text:          "This is an example post text 2.",
					TextHtml:      "<p>This is an example post text 2 in HTML.</p>",
					Ups:           20,
					Over18:        true,
					Spoiler:       true,
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
					CreatedAtUnix: 1725091120,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
				},
			},
			wantErr: postrepo.ErrPostParentTableRecordNotFound,
		},
		{
			name:         "voxsphere is not present in the parent table :NEG",
			fixtureFiles: []string{"topics.yml", "voxspheres.yml", "users.yml", "posts.yml"},
			args: args{
				post: models.Post{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					AuthorID:      uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					VoxsphereID:   uuid.MustParse("00000000-0000-0000-0000-000000000009"),
					Title:         "Updated Title",
					Text:          "Update Text",
					TextHtml:      "<p>Updated HTML.</p>",
					Ups:           30,
					Over18:        true,
					Spoiler:       true,
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 30, 0, time.UTC),
					CreatedAtUnix: 1725091120,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 30, 0, time.UTC),
				},
			},
			wantPost: models.Post{},
			wantPosts: []models.Post{
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					AuthorID:      uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					VoxsphereID:   uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Title:         "Example Post Title 1",
					Text:          "This is an example post text 1.",
					TextHtml:      "<p>This is an example post text 1 in HTML.</p>",
					Ups:           10,
					Over18:        false,
					Spoiler:       false,
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix: 1725091100,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					AuthorID:      uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					VoxsphereID:   uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					Title:         "Example Post Title 2",
					Text:          "This is an example post text 2.",
					TextHtml:      "<p>This is an example post text 2 in HTML.</p>",
					Ups:           20,
					Over18:        true,
					Spoiler:       true,
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
					CreatedAtUnix: 1725091120,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
				},
			},
			wantErr: postrepo.ErrPostParentTableRecordNotFound,
		},
		{
			name:         "update post :POS",
			fixtureFiles: []string{"topics.yml", "voxspheres.yml", "users.yml", "posts.yml"},
			args: args{
				post: models.Post{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					AuthorID:      uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					VoxsphereID:   uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Title:         "Updated Title",
					Text:          "Update Text",
					TextHtml:      "<p>Updated HTML.</p>",
					Ups:           30,
					Over18:        true,
					Spoiler:       true,
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 30, 0, time.UTC),
					CreatedAtUnix: 1725091120,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 30, 0, time.UTC),
				},
			},
			wantPost: models.Post{
				ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				AuthorID:      uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				VoxsphereID:   uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				Title:         "Updated Title",
				Text:          "Update Text",
				TextHtml:      "<p>Updated HTML.</p>",
				Ups:           30,
				Over18:        true,
				Spoiler:       true,
				CreatedAt:     time.Date(2024, 10, 10, 10, 10, 30, 0, time.UTC),
				CreatedAtUnix: 1725091120,
				UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 30, 0, time.UTC),
			},
			wantPosts: []models.Post{
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					AuthorID:      uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					VoxsphereID:   uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Title:         "Updated Title",
					Text:          "Update Text",
					TextHtml:      "<p>Updated HTML.</p>",
					Ups:           30,
					Over18:        true,
					Spoiler:       true,
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 30, 0, time.UTC),
					CreatedAtUnix: 1725091120,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 30, 0, time.UTC),
				},
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					AuthorID:      uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					VoxsphereID:   uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					Title:         "Example Post Title 2",
					Text:          "This is an example post text 2.",
					TextHtml:      "<p>This is an example post text 2 in HTML.</p>",
					Ups:           20,
					Over18:        true,
					Spoiler:       true,
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
					CreatedAtUnix: 1725091120,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupPostgres(t, tt.fixtureFiles...)
			pgrepo := postrepo.NewRepo(db)

			startTime := time.Now()
			gotPost, gotErr := pgrepo.UpdatePost(context.Background(), tt.args.post)
			endTime := time.Now()

			assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")
			if tt.wantErr == nil {
				assertTimeWithinRange(t, gotPost.UpdatedAt, startTime, endTime)
			}

			gotPosts, err := pgrepo.Posts(context.Background())

			assert.NoError(t, err, "expect no error while getting posts")
			assertPostsWithoutTimestamp(t, tt.wantPosts, gotPosts)
		})
	}
}

func TestRepo_DeletePost(t *testing.T) {
	type args struct {
		ID uuid.UUID
	}
	tests := []struct {
		name         string
		fixtureFiles []string
		args         args
		wantPosts    []models.Post
		wantErr      error
	}{
		{
			name:         "post not found :NEG",
			fixtureFiles: []string{"topics.yml", "voxspheres.yml", "users.yml", "posts.yml"},
			args: args{
				ID: uuid.MustParse("00000000-0000-0000-0000-000000000009"),
			},
			wantPosts: []models.Post{
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					AuthorID:      uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					VoxsphereID:   uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Title:         "Example Post Title 1",
					Text:          "This is an example post text 1.",
					TextHtml:      "<p>This is an example post text 1 in HTML.</p>",
					Ups:           10,
					Over18:        false,
					Spoiler:       false,
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix: 1725091100,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					AuthorID:      uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					VoxsphereID:   uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					Title:         "Example Post Title 2",
					Text:          "This is an example post text 2.",
					TextHtml:      "<p>This is an example post text 2 in HTML.</p>",
					Ups:           20,
					Over18:        true,
					Spoiler:       true,
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
					CreatedAtUnix: 1725091120,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
				},
			},
			wantErr: postrepo.ErrPostNotFound,
		},
		{
			name:         "post deleted :POS",
			fixtureFiles: []string{"topics.yml", "voxspheres.yml", "users.yml", "posts.yml"},
			args: args{
				ID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			},
			wantPosts: []models.Post{
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					AuthorID:      uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					VoxsphereID:   uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					Title:         "Example Post Title 2",
					Text:          "This is an example post text 2.",
					TextHtml:      "<p>This is an example post text 2 in HTML.</p>",
					Ups:           20,
					Over18:        true,
					Spoiler:       true,
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
					CreatedAtUnix: 1725091120,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
				},
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupPostgres(t, tt.fixtureFiles...)
			pgrepo := postrepo.NewRepo(db)

			gotErr := pgrepo.DeletePost(context.Background(), tt.args.ID)

			assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")

			gotPosts, err := pgrepo.Posts(context.Background())

			assert.NoError(t, err, "expect no error while getting posts")
			assert.Equal(t, tt.wantPosts, gotPosts, "expect posts to match")
		})
	}
}

func TestRepo_ForeignKeyCascade(t *testing.T) {
	t.Run("on deleting voxsphere from parent table , no child references should exist in posts table", func(t *testing.T) {
		db := setupPostgres(t, "topics.yml", "voxspheres.yml", "users.yml", "posts.yml")
		postPgrepo := postrepo.NewRepo(db)
		voxspherePgrepo := voxrepo.NewRepo(db)

		wantPosts := []models.Post{
			{
				ID:            uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				AuthorID:      uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				VoxsphereID:   uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				Title:         "Example Post Title 2",
				Text:          "This is an example post text 2.",
				TextHtml:      "<p>This is an example post text 2 in HTML.</p>",
				Ups:           20,
				Over18:        true,
				Spoiler:       true,
				CreatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
				CreatedAtUnix: 1725091120,
				UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
			},
		}

		err := voxspherePgrepo.DeleteVoxsphere(context.Background(), uuid.MustParse("00000000-0000-0000-0000-000000000001"))

		assert.NoError(t, err, "expect no error while deleting voxsphere")

		gotPosts, err := postPgrepo.Posts(context.Background())

		assert.NoError(t, err, "expect no error while getting posts")
		assertPostsWithoutTimestamp(t, wantPosts, gotPosts)
	})

	t.Run("on deleting author from parent table , no child references should exist in posts table", func(t *testing.T) {
		db := setupPostgres(t, "topics.yml", "voxspheres.yml", "users.yml", "posts.yml")
		postPgrepo := postrepo.NewRepo(db)
		userPgrepo := userrepo.NewRepo(db)

		wantPosts := []models.Post{
			{
				ID:            uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				AuthorID:      uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				VoxsphereID:   uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				Title:         "Example Post Title 2",
				Text:          "This is an example post text 2.",
				TextHtml:      "<p>This is an example post text 2 in HTML.</p>",
				Ups:           20,
				Over18:        true,
				Spoiler:       true,
				CreatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
				CreatedAtUnix: 1725091120,
				UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
			},
		}

		err := userPgrepo.DeleteUser(context.Background(), uuid.MustParse("00000000-0000-0000-0000-000000000001"))

		assert.NoError(t, err, "expect no error while deleting user")

		gotPosts, err := postPgrepo.Posts(context.Background())

		assert.NoError(t, err, "expect no error while getting posts")
		assertPostsWithoutTimestamp(t, wantPosts, gotPosts)
	})
}

func TestRepo_PostsPaginated(t *testing.T) {
	type args struct {
		skip  int
		limit int
	}
	tests := []struct {
		name         string
		fixtureFiles []string
		args         args
		wantPosts    []models.Post
		wantErr      error
	}{
		{
			name:         "paginated posts skip 3 limit 2 :POS",
			fixtureFiles: []string{"topics.yml", "voxspheres.yml", "users.yml", "posts_paginated.yml"},
			args: args{
				skip:  3,
				limit: 2,
			},
			wantPosts: []models.Post{
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000004"),
					AuthorID:      uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					VoxsphereID:   uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Title:         "Example Post Title 4",
					Text:          "This is an example post text 4.",
					TextHtml:      "This is an example post text 4 in HTML.",
					Ups:           40,
					Over18:        true,
					Spoiler:       false,
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 40, 0, time.UTC),
					CreatedAtUnix: 1725091160,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 40, 0, time.UTC),
				},
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000005"),
					AuthorID:      uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					VoxsphereID:   uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					Title:         "Example Post Title 5",
					Text:          "This is an example post text 5.",
					TextHtml:      "This is an example post text 5 in HTML.",
					Ups:           50,
					Over18:        false,
					Spoiler:       true,
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 50, 0, time.UTC),
					CreatedAtUnix: 1725091180,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 50, 0, time.UTC),
				},
			},
			wantErr: nil,
		},
		{
			name:         "paginated posts skip 3 limit 1 :POS",
			fixtureFiles: []string{"topics.yml", "voxspheres.yml", "users.yml", "posts_paginated.yml"},
			args: args{
				skip:  3,
				limit: 1,
			},
			wantPosts: []models.Post{
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000004"),
					AuthorID:      uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					VoxsphereID:   uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Title:         "Example Post Title 4",
					Text:          "This is an example post text 4.",
					TextHtml:      "This is an example post text 4 in HTML.",
					Ups:           40,
					Over18:        true,
					Spoiler:       false,
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 40, 0, time.UTC),
					CreatedAtUnix: 1725091160,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 40, 0, time.UTC),
				},
			},
			wantErr: nil,
		},
		{
			name:         "paginated posts skip 100 limit 100 :POS",
			fixtureFiles: []string{"topics.yml", "voxspheres.yml", "users.yml", "posts_paginated.yml"},
			args: args{
				skip:  100,
				limit: 100,
			},
			wantPosts: nil,
			wantErr:   nil,
		},
		{
			name:         "no posts :POS",
			fixtureFiles: []string{"topics.yml", "voxspheres.yml", "users.yml"},
			args: args{
				skip:  100,
				limit: 100,
			},
			wantPosts: nil,
			wantErr:   nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupPostgres(t, tt.fixtureFiles...)
			pgrepo := postrepo.NewRepo(db)

			gotPosts, gotErr := pgrepo.PostsPaginated(context.Background(), tt.args.skip, tt.args.limit)

			assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")
			assert.Equal(t, tt.wantPosts, gotPosts, "expect posts to match")
		})
	}
}
