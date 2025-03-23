package comments_test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"slices"
	"testing"
	"time"

	"github.com/glowfi/voxpopuli/backend/pkg/models"
	commentrepo "github.com/glowfi/voxpopuli/backend/pkg/repo/comments"
	postrepo "github.com/glowfi/voxpopuli/backend/pkg/repo/post"
	userrepo "github.com/glowfi/voxpopuli/backend/pkg/repo/user"
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
	db.RegisterModel((*models.Comment)(nil))

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
	if _, err := db.NewTruncateTable().Cascade().Model((*models.Comment)(nil)).Exec(context.Background()); err != nil {
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

func assertCommentWithoutTimestamp(t *testing.T, wantComment, gotComment models.Comment) {
	assert.Equal(t, wantComment.ID, gotComment.ID, "expect comment ID to match")
	assert.Equal(t, wantComment.AuthorID, gotComment.AuthorID, "expect author ID to match")
	assert.Equal(t, wantComment.ParentCommentID, gotComment.ParentCommentID, "expect parent comment ID to match")
	assert.Equal(t, wantComment.PostID, gotComment.PostID, "expect post ID to match")
	assert.Equal(t, wantComment.Body, gotComment.Body, "expect comment body to match")
	assert.Equal(t, wantComment.BodyHtml, gotComment.BodyHtml, "expect comment html to match")
	assert.Equal(t, wantComment.Ups, gotComment.Ups, "expect comment ups to match")
	assert.Equal(t, wantComment.Score, gotComment.Score, "expect comment score to match")
}

func assertCommentsWithoutTimestamp(t *testing.T, wantComments, gotComments []models.Comment) {
	t.Helper()

	if len(wantComments) != len(gotComments) {
		t.Fatal("length of wantComments and gotComments do not match")
	}

	for _, wantComment := range wantComments {
		idx := slices.IndexFunc(gotComments, func(c models.Comment) bool {
			return c.ID == wantComment.ID
		})

		if idx == -1 {
			t.Fatalf("comment of ID %v is not present in gotComments", wantComment.ID)
			return
		}
		assertCommentWithoutTimestamp(t, wantComment, gotComments[idx])
	}
}

func assertCommentsWithTimestamp(t *testing.T, wantComments, gotComments []models.Comment) {
	t.Helper()

	if len(wantComments) != len(gotComments) {
		t.Fatal("length of wantComments and gotComments do not match")
	}

	for _, wantComment := range wantComments {
		idx := slices.IndexFunc(gotComments, func(c models.Comment) bool {
			return c.ID == wantComment.ID
		})

		if idx == -1 {
			t.Fatalf("comment of ID %v is not present in gotComments", wantComment.ID)
			return
		}
		assert.Equal(t, wantComment, gotComments[idx])
	}
}

func TestRepo_Comments(t *testing.T) {
	tests := []struct {
		name         string
		fixtureFiles []string
		wantComments []models.Comment
		wantErr      error
	}{
		{
			name:         "comments :POS",
			fixtureFiles: []string{"topics.yml", "voxspheres.yml", "users.yml", "posts.yml", "comments.yml"},
			wantComments: []models.Comment{
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					AuthorID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					ParentCommentID: uuid.MustParse("00000000-0000-0000-0000-000000000000"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Body:            "This is a parent comment 1",
					BodyHtml:        "<p>This is a parent comment 1</p>",
					Ups:             1,
					Score:           1,
					CreatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix:   1725091100,
					UpdatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					AuthorID:        uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					ParentCommentID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Body:            "This is reply 1 to parent comment 1",
					BodyHtml:        "<p>This is reply 1 to parent comment 1</p>",
					Ups:             1,
					Score:           1,
					CreatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix:   1725091100,
					UpdatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					AuthorID:        uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					ParentCommentID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Body:            "This is reply to reply 1",
					BodyHtml:        "<p>This is reply to reply 1</p>",
					Ups:             1,
					Score:           1,
					CreatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix:   1725091100,
					UpdatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000004"),
					AuthorID:        uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					ParentCommentID: uuid.MustParse("00000000-0000-0000-0000-000000000000"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Body:            "This is a parent comment 2",
					BodyHtml:        "<p>This is a parent comment 2</p>",
					Ups:             1,
					Score:           1,
					CreatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix:   1725091100,
					UpdatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000005"),
					AuthorID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					ParentCommentID: uuid.MustParse("00000000-0000-0000-0000-000000000000"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					Body:            "This is a parent comment 3",
					BodyHtml:        "<p>This is a parent comment 3</p>",
					Ups:             1,
					Score:           1,
					CreatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix:   1725091100,
					UpdatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000006"),
					AuthorID:        uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					ParentCommentID: uuid.MustParse("00000000-0000-0000-0000-000000000005"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					Body:            "This is reply to parent comment 3",
					BodyHtml:        "<p>This is reply to parent comment 3</p>",
					Ups:             1,
					Score:           1,
					CreatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix:   1725091100,
					UpdatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000007"),
					AuthorID:        uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					ParentCommentID: uuid.MustParse("00000000-0000-0000-0000-000000000000"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					Body:            "This is a parent comment 4",
					BodyHtml:        "<p>This is a parent comment 4</p>",
					Ups:             1,
					Score:           1,
					CreatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix:   1725091100,
					UpdatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
			},
			wantErr: nil,
		},
		{
			name:         "no comments :POS",
			fixtureFiles: []string{},
			wantComments: []models.Comment{},
			wantErr:      nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupPostgres(t, tt.fixtureFiles...)
			pgrepo := commentrepo.NewRepo(db)

			gotComments, gotErr := pgrepo.Comments(context.Background())

			assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")
			assertCommentsWithTimestamp(t, tt.wantComments, gotComments)
		})
	}
}

func TestRepo_CommentByID(t *testing.T) {
	type args struct {
		ID uuid.UUID
	}

	tests := []struct {
		name         string
		fixtureFiles []string
		args         args
		wantComment  models.Comment
		wantErr      error
	}{
		{
			name:         "comment not found :NEG",
			fixtureFiles: []string{"topics.yml", "voxspheres.yml", "users.yml", "posts.yml", "comments.yml"},
			args: args{
				ID: uuid.MustParse("00000000-0000-0000-0000-000000000009"),
			},
			wantComment: models.Comment{},
			wantErr:     commentrepo.ErrCommentNotFound,
		},
		{
			name:         "comment found :POS",
			fixtureFiles: []string{"topics.yml", "voxspheres.yml", "users.yml", "posts.yml", "comments.yml"},
			args: args{
				ID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			},
			wantComment: models.Comment{
				ID:              uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				AuthorID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				ParentCommentID: uuid.MustParse("00000000-0000-0000-0000-000000000000"),
				PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				Body:            "This is a parent comment 1",
				BodyHtml:        "<p>This is a parent comment 1</p>",
				Ups:             1,
				Score:           1,
				CreatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				CreatedAtUnix:   1725091100,
				UpdatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupPostgres(t, tt.fixtureFiles...)
			pgrepo := commentrepo.NewRepo(db)

			gotComment, gotErr := pgrepo.CommentByID(context.Background(), tt.args.ID)

			assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")
			assert.Equal(t, tt.wantComment, gotComment, "expect comment to match")
		})
	}
}

func TestRepo_AddComments(t *testing.T) {
	type args struct {
		comments []models.Comment
	}
	tests := []struct {
		name                 string
		fixtureFiles         []string
		args                 args
		wantInsertedComments []models.Comment
		wantComments         []models.Comment
		wantErr              error
	}{
		{
			name:         "duplicate comment id :NEG",
			fixtureFiles: []string{"topics.yml", "voxspheres.yml", "users.yml", "posts.yml", "comments.yml"},
			args: args{
				comments: []models.Comment{
					{
						ID:              uuid.MustParse("00000000-0000-0000-0000-000000000001"),
						AuthorID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
						ParentCommentID: uuid.MustParse("00000000-0000-0000-0000-000000000000"),
						PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
						Body:            "This is a parent comment 1",
						BodyHtml:        "<p>This is a parent comment 1</p>",
						Ups:             1,
						Score:           1,
						CreatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
						CreatedAtUnix:   1725091100,
						UpdatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					},
				},
			},
			wantInsertedComments: nil,
			wantComments: []models.Comment{
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					AuthorID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					ParentCommentID: uuid.MustParse("00000000-0000-0000-0000-000000000000"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Body:            "This is a parent comment 1",
					BodyHtml:        "<p>This is a parent comment 1</p>",
					Ups:             1,
					Score:           1,
					CreatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix:   1725091100,
					UpdatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					AuthorID:        uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					ParentCommentID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Body:            "This is reply 1 to parent comment 1",
					BodyHtml:        "<p>This is reply 1 to parent comment 1</p>",
					Ups:             1,
					Score:           1,
					CreatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix:   1725091100,
					UpdatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					AuthorID:        uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					ParentCommentID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Body:            "This is reply to reply 1",
					BodyHtml:        "<p>This is reply to reply 1</p>",
					Ups:             1,
					Score:           1,
					CreatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix:   1725091100,
					UpdatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000004"),
					AuthorID:        uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					ParentCommentID: uuid.MustParse("00000000-0000-0000-0000-000000000000"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Body:            "This is a parent comment 2",
					BodyHtml:        "<p>This is a parent comment 2</p>",
					Ups:             1,
					Score:           1,
					CreatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix:   1725091100,
					UpdatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000005"),
					AuthorID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					ParentCommentID: uuid.MustParse("00000000-0000-0000-0000-000000000000"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					Body:            "This is a parent comment 3",
					BodyHtml:        "<p>This is a parent comment 3</p>",
					Ups:             1,
					Score:           1,
					CreatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix:   1725091100,
					UpdatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000006"),
					AuthorID:        uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					ParentCommentID: uuid.MustParse("00000000-0000-0000-0000-000000000005"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					Body:            "This is reply to parent comment 3",
					BodyHtml:        "<p>This is reply to parent comment 3</p>",
					Ups:             1,
					Score:           1,
					CreatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix:   1725091100,
					UpdatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000007"),
					AuthorID:        uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					ParentCommentID: uuid.MustParse("00000000-0000-0000-0000-000000000000"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					Body:            "This is a parent comment 4",
					BodyHtml:        "<p>This is a parent comment 4</p>",
					Ups:             1,
					Score:           1,
					CreatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix:   1725091100,
					UpdatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
			},
			wantErr: commentrepo.ErrCommentDuplicateID,
		},
		{
			name:         "author does not exist in parent table :NEG",
			fixtureFiles: []string{"topics.yml", "voxspheres.yml", "users.yml", "posts.yml", "comments.yml"},
			args: args{
				comments: []models.Comment{
					{
						ID:              uuid.MustParse("00000000-0000-0000-0000-000000000008"),
						AuthorID:        uuid.MustParse("00000000-0000-0000-0000-000000000009"),
						ParentCommentID: uuid.MustParse("00000000-0000-0000-0000-000000000000"),
						PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
						Body:            "This is a parent comment 5",
						BodyHtml:        "<p>This is a parent comment 5</p>",
						Ups:             1,
						Score:           1,
						CreatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
						CreatedAtUnix:   1725091100,
						UpdatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					},
				},
			},
			wantInsertedComments: nil,
			wantComments: []models.Comment{
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					AuthorID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					ParentCommentID: uuid.MustParse("00000000-0000-0000-0000-000000000000"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Body:            "This is a parent comment 1",
					BodyHtml:        "<p>This is a parent comment 1</p>",
					Ups:             1,
					Score:           1,
					CreatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix:   1725091100,
					UpdatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					AuthorID:        uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					ParentCommentID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Body:            "This is reply 1 to parent comment 1",
					BodyHtml:        "<p>This is reply 1 to parent comment 1</p>",
					Ups:             1,
					Score:           1,
					CreatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix:   1725091100,
					UpdatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					AuthorID:        uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					ParentCommentID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Body:            "This is reply to reply 1",
					BodyHtml:        "<p>This is reply to reply 1</p>",
					Ups:             1,
					Score:           1,
					CreatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix:   1725091100,
					UpdatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000004"),
					AuthorID:        uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					ParentCommentID: uuid.MustParse("00000000-0000-0000-0000-000000000000"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Body:            "This is a parent comment 2",
					BodyHtml:        "<p>This is a parent comment 2</p>",
					Ups:             1,
					Score:           1,
					CreatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix:   1725091100,
					UpdatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000005"),
					AuthorID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					ParentCommentID: uuid.MustParse("00000000-0000-0000-0000-000000000000"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					Body:            "This is a parent comment 3",
					BodyHtml:        "<p>This is a parent comment 3</p>",
					Ups:             1,
					Score:           1,
					CreatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix:   1725091100,
					UpdatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000006"),
					AuthorID:        uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					ParentCommentID: uuid.MustParse("00000000-0000-0000-0000-000000000005"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					Body:            "This is reply to parent comment 3",
					BodyHtml:        "<p>This is reply to parent comment 3</p>",
					Ups:             1,
					Score:           1,
					CreatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix:   1725091100,
					UpdatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000007"),
					AuthorID:        uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					ParentCommentID: uuid.MustParse("00000000-0000-0000-0000-000000000000"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					Body:            "This is a parent comment 4",
					BodyHtml:        "<p>This is a parent comment 4</p>",
					Ups:             1,
					Score:           1,
					CreatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix:   1725091100,
					UpdatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
			},
			wantErr: commentrepo.ErrCommentParentTableRecordNotFound,
		},
		{
			name:         "post does not exist in parent table :NEG",
			fixtureFiles: []string{"topics.yml", "voxspheres.yml", "users.yml", "posts.yml", "comments.yml"},
			args: args{
				comments: []models.Comment{
					{
						ID:              uuid.MustParse("00000000-0000-0000-0000-000000000008"),
						AuthorID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
						ParentCommentID: uuid.MustParse("00000000-0000-0000-0000-000000000000"),
						PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000009"),
						Body:            "This is a parent comment 5",
						BodyHtml:        "<p>This is a parent comment 5</p>",
						Ups:             1,
						Score:           1,
						CreatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
						CreatedAtUnix:   1725091100,
						UpdatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					},
				},
			},
			wantInsertedComments: nil,
			wantComments: []models.Comment{
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					AuthorID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					ParentCommentID: uuid.MustParse("00000000-0000-0000-0000-000000000000"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Body:            "This is a parent comment 1",
					BodyHtml:        "<p>This is a parent comment 1</p>",
					Ups:             1,
					Score:           1,
					CreatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix:   1725091100,
					UpdatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					AuthorID:        uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					ParentCommentID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Body:            "This is reply 1 to parent comment 1",
					BodyHtml:        "<p>This is reply 1 to parent comment 1</p>",
					Ups:             1,
					Score:           1,
					CreatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix:   1725091100,
					UpdatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					AuthorID:        uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					ParentCommentID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Body:            "This is reply to reply 1",
					BodyHtml:        "<p>This is reply to reply 1</p>",
					Ups:             1,
					Score:           1,
					CreatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix:   1725091100,
					UpdatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000004"),
					AuthorID:        uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					ParentCommentID: uuid.MustParse("00000000-0000-0000-0000-000000000000"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Body:            "This is a parent comment 2",
					BodyHtml:        "<p>This is a parent comment 2</p>",
					Ups:             1,
					Score:           1,
					CreatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix:   1725091100,
					UpdatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000005"),
					AuthorID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					ParentCommentID: uuid.MustParse("00000000-0000-0000-0000-000000000000"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					Body:            "This is a parent comment 3",
					BodyHtml:        "<p>This is a parent comment 3</p>",
					Ups:             1,
					Score:           1,
					CreatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix:   1725091100,
					UpdatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000006"),
					AuthorID:        uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					ParentCommentID: uuid.MustParse("00000000-0000-0000-0000-000000000005"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					Body:            "This is reply to parent comment 3",
					BodyHtml:        "<p>This is reply to parent comment 3</p>",
					Ups:             1,
					Score:           1,
					CreatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix:   1725091100,
					UpdatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000007"),
					AuthorID:        uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					ParentCommentID: uuid.MustParse("00000000-0000-0000-0000-000000000000"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					Body:            "This is a parent comment 4",
					BodyHtml:        "<p>This is a parent comment 4</p>",
					Ups:             1,
					Score:           1,
					CreatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix:   1725091100,
					UpdatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
			},
			wantErr: commentrepo.ErrCommentParentTableRecordNotFound,
		},
		{
			name:         "add comments :NEG",
			fixtureFiles: []string{"topics.yml", "voxspheres.yml", "users.yml", "posts.yml", "comments.yml"},
			args: args{
				comments: []models.Comment{
					{
						ID:              uuid.MustParse("00000000-0000-0000-0000-000000000008"),
						AuthorID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
						ParentCommentID: uuid.MustParse("00000000-0000-0000-0000-000000000000"),
						PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000002"),
						Body:            "This is a parent comment 5",
						BodyHtml:        "<p>This is a parent comment 5</p>",
						Ups:             1,
						Score:           1,
						CreatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
						CreatedAtUnix:   1725091100,
						UpdatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					},
				},
			},
			wantInsertedComments: []models.Comment{
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000008"),
					AuthorID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					ParentCommentID: uuid.MustParse("00000000-0000-0000-0000-000000000000"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					Body:            "This is a parent comment 5",
					BodyHtml:        "<p>This is a parent comment 5</p>",
					Ups:             1,
					Score:           1,
					CreatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix:   1725091100,
					UpdatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
			},
			wantComments: []models.Comment{
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					AuthorID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					ParentCommentID: uuid.MustParse("00000000-0000-0000-0000-000000000000"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Body:            "This is a parent comment 1",
					BodyHtml:        "<p>This is a parent comment 1</p>",
					Ups:             1,
					Score:           1,
					CreatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix:   1725091100,
					UpdatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					AuthorID:        uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					ParentCommentID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Body:            "This is reply 1 to parent comment 1",
					BodyHtml:        "<p>This is reply 1 to parent comment 1</p>",
					Ups:             1,
					Score:           1,
					CreatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix:   1725091100,
					UpdatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					AuthorID:        uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					ParentCommentID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Body:            "This is reply to reply 1",
					BodyHtml:        "<p>This is reply to reply 1</p>",
					Ups:             1,
					Score:           1,
					CreatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix:   1725091100,
					UpdatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000004"),
					AuthorID:        uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					ParentCommentID: uuid.MustParse("00000000-0000-0000-0000-000000000000"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Body:            "This is a parent comment 2",
					BodyHtml:        "<p>This is a parent comment 2</p>",
					Ups:             1,
					Score:           1,
					CreatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix:   1725091100,
					UpdatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000005"),
					AuthorID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					ParentCommentID: uuid.MustParse("00000000-0000-0000-0000-000000000000"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					Body:            "This is a parent comment 3",
					BodyHtml:        "<p>This is a parent comment 3</p>",
					Ups:             1,
					Score:           1,
					CreatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix:   1725091100,
					UpdatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000006"),
					AuthorID:        uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					ParentCommentID: uuid.MustParse("00000000-0000-0000-0000-000000000005"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					Body:            "This is reply to parent comment 3",
					BodyHtml:        "<p>This is reply to parent comment 3</p>",
					Ups:             1,
					Score:           1,
					CreatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix:   1725091100,
					UpdatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000007"),
					AuthorID:        uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					ParentCommentID: uuid.MustParse("00000000-0000-0000-0000-000000000000"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					Body:            "This is a parent comment 4",
					BodyHtml:        "<p>This is a parent comment 4</p>",
					Ups:             1,
					Score:           1,
					CreatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix:   1725091100,
					UpdatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000008"),
					AuthorID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					ParentCommentID: uuid.MustParse("00000000-0000-0000-0000-000000000000"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					Body:            "This is a parent comment 5",
					BodyHtml:        "<p>This is a parent comment 5</p>",
					Ups:             1,
					Score:           1,
					CreatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix:   1725091100,
					UpdatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupPostgres(t, tt.fixtureFiles...)
			pgrepo := commentrepo.NewRepo(db)

			startTime := time.Now()
			gotInsertedComments, gotErr := pgrepo.AddComments(context.Background(), tt.args.comments...)
			endTime := time.Now()

			for _, gotInsertedComment := range gotInsertedComments {
				assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")
				assert.Equal(
					t,
					gotInsertedComment.UpdatedAt,
					gotInsertedComment.CreatedAt,
					"expect CreatedAt and UpdatedAt to be same",
				)
				if tt.wantErr == nil {
					assertTimeWithinRange(t, gotInsertedComment.CreatedAt, startTime, endTime)
					assertTimeWithinRange(t, gotInsertedComment.UpdatedAt, startTime, endTime)
				}
			}

			gotComments, err := pgrepo.Comments(context.Background())

			assert.NoError(t, err, "expect no error while getting comments")
			assertCommentsWithoutTimestamp(t, tt.wantComments, gotComments)
		})
	}
}

func TestRepo_UpdateComment(t *testing.T) {
	type args struct {
		comment models.Comment
	}
	tests := []struct {
		name         string
		fixtureFiles []string
		args         args
		wantComment  models.Comment
		wantComments []models.Comment
		wantErr      error
	}{
		{
			name:         "comment not id :NEG",
			fixtureFiles: []string{"topics.yml", "voxspheres.yml", "users.yml", "posts.yml", "comments.yml"},
			args: args{
				comment: models.Comment{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000009"),
					AuthorID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					ParentCommentID: uuid.MustParse("00000000-0000-0000-0000-000000000000"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Body:            "This is a parent comment 1",
					BodyHtml:        "<p>This is a parent comment 1</p>",
					Ups:             1,
					Score:           1,
					CreatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix:   1725091100,
					UpdatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
			},
			wantComment: models.Comment{},
			wantComments: []models.Comment{
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					AuthorID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					ParentCommentID: uuid.MustParse("00000000-0000-0000-0000-000000000000"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Body:            "This is a parent comment 1",
					BodyHtml:        "<p>This is a parent comment 1</p>",
					Ups:             1,
					Score:           1,
					CreatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix:   1725091100,
					UpdatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					AuthorID:        uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					ParentCommentID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Body:            "This is reply 1 to parent comment 1",
					BodyHtml:        "<p>This is reply 1 to parent comment 1</p>",
					Ups:             1,
					Score:           1,
					CreatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix:   1725091100,
					UpdatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					AuthorID:        uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					ParentCommentID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Body:            "This is reply to reply 1",
					BodyHtml:        "<p>This is reply to reply 1</p>",
					Ups:             1,
					Score:           1,
					CreatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix:   1725091100,
					UpdatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000004"),
					AuthorID:        uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					ParentCommentID: uuid.MustParse("00000000-0000-0000-0000-000000000000"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Body:            "This is a parent comment 2",
					BodyHtml:        "<p>This is a parent comment 2</p>",
					Ups:             1,
					Score:           1,
					CreatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix:   1725091100,
					UpdatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000005"),
					AuthorID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					ParentCommentID: uuid.MustParse("00000000-0000-0000-0000-000000000000"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					Body:            "This is a parent comment 3",
					BodyHtml:        "<p>This is a parent comment 3</p>",
					Ups:             1,
					Score:           1,
					CreatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix:   1725091100,
					UpdatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000006"),
					AuthorID:        uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					ParentCommentID: uuid.MustParse("00000000-0000-0000-0000-000000000005"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					Body:            "This is reply to parent comment 3",
					BodyHtml:        "<p>This is reply to parent comment 3</p>",
					Ups:             1,
					Score:           1,
					CreatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix:   1725091100,
					UpdatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000007"),
					AuthorID:        uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					ParentCommentID: uuid.MustParse("00000000-0000-0000-0000-000000000000"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					Body:            "This is a parent comment 4",
					BodyHtml:        "<p>This is a parent comment 4</p>",
					Ups:             1,
					Score:           1,
					CreatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix:   1725091100,
					UpdatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
			},
			wantErr: commentrepo.ErrCommentNotFound,
		},
		{
			name:         "author does not exist in parent table :NEG",
			fixtureFiles: []string{"topics.yml", "voxspheres.yml", "users.yml", "posts.yml", "comments.yml"},
			args: args{
				comment: models.Comment{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					AuthorID:        uuid.MustParse("00000000-0000-0000-0000-000000000009"),
					ParentCommentID: uuid.MustParse("00000000-0000-0000-0000-000000000000"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Body:            "This is a parent comment 5",
					BodyHtml:        "<p>This is a parent comment 5</p>",
					Ups:             1,
					Score:           1,
					CreatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix:   1725091100,
					UpdatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
			},
			wantComment: models.Comment{},
			wantComments: []models.Comment{
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					AuthorID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					ParentCommentID: uuid.MustParse("00000000-0000-0000-0000-000000000000"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Body:            "This is a parent comment 1",
					BodyHtml:        "<p>This is a parent comment 1</p>",
					Ups:             1,
					Score:           1,
					CreatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix:   1725091100,
					UpdatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					AuthorID:        uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					ParentCommentID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Body:            "This is reply 1 to parent comment 1",
					BodyHtml:        "<p>This is reply 1 to parent comment 1</p>",
					Ups:             1,
					Score:           1,
					CreatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix:   1725091100,
					UpdatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					AuthorID:        uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					ParentCommentID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Body:            "This is reply to reply 1",
					BodyHtml:        "<p>This is reply to reply 1</p>",
					Ups:             1,
					Score:           1,
					CreatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix:   1725091100,
					UpdatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000004"),
					AuthorID:        uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					ParentCommentID: uuid.MustParse("00000000-0000-0000-0000-000000000000"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Body:            "This is a parent comment 2",
					BodyHtml:        "<p>This is a parent comment 2</p>",
					Ups:             1,
					Score:           1,
					CreatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix:   1725091100,
					UpdatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000005"),
					AuthorID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					ParentCommentID: uuid.MustParse("00000000-0000-0000-0000-000000000000"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					Body:            "This is a parent comment 3",
					BodyHtml:        "<p>This is a parent comment 3</p>",
					Ups:             1,
					Score:           1,
					CreatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix:   1725091100,
					UpdatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000006"),
					AuthorID:        uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					ParentCommentID: uuid.MustParse("00000000-0000-0000-0000-000000000005"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					Body:            "This is reply to parent comment 3",
					BodyHtml:        "<p>This is reply to parent comment 3</p>",
					Ups:             1,
					Score:           1,
					CreatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix:   1725091100,
					UpdatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000007"),
					AuthorID:        uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					ParentCommentID: uuid.MustParse("00000000-0000-0000-0000-000000000000"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					Body:            "This is a parent comment 4",
					BodyHtml:        "<p>This is a parent comment 4</p>",
					Ups:             1,
					Score:           1,
					CreatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix:   1725091100,
					UpdatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
			},
			wantErr: commentrepo.ErrCommentParentTableRecordNotFound,
		},
		{
			name:         "post does not exist in parent table :NEG",
			fixtureFiles: []string{"topics.yml", "voxspheres.yml", "users.yml", "posts.yml", "comments.yml"},
			args: args{
				comment: models.Comment{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					AuthorID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					ParentCommentID: uuid.MustParse("00000000-0000-0000-0000-000000000000"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000009"),
					Body:            "This is a parent comment 5",
					BodyHtml:        "<p>This is a parent comment 5</p>",
					Ups:             1,
					Score:           1,
					CreatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix:   1725091100,
					UpdatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
			},
			wantComment: models.Comment{},
			wantComments: []models.Comment{
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					AuthorID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					ParentCommentID: uuid.MustParse("00000000-0000-0000-0000-000000000000"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Body:            "This is a parent comment 1",
					BodyHtml:        "<p>This is a parent comment 1</p>",
					Ups:             1,
					Score:           1,
					CreatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix:   1725091100,
					UpdatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					AuthorID:        uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					ParentCommentID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Body:            "This is reply 1 to parent comment 1",
					BodyHtml:        "<p>This is reply 1 to parent comment 1</p>",
					Ups:             1,
					Score:           1,
					CreatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix:   1725091100,
					UpdatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					AuthorID:        uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					ParentCommentID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Body:            "This is reply to reply 1",
					BodyHtml:        "<p>This is reply to reply 1</p>",
					Ups:             1,
					Score:           1,
					CreatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix:   1725091100,
					UpdatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000004"),
					AuthorID:        uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					ParentCommentID: uuid.MustParse("00000000-0000-0000-0000-000000000000"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Body:            "This is a parent comment 2",
					BodyHtml:        "<p>This is a parent comment 2</p>",
					Ups:             1,
					Score:           1,
					CreatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix:   1725091100,
					UpdatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000005"),
					AuthorID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					ParentCommentID: uuid.MustParse("00000000-0000-0000-0000-000000000000"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					Body:            "This is a parent comment 3",
					BodyHtml:        "<p>This is a parent comment 3</p>",
					Ups:             1,
					Score:           1,
					CreatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix:   1725091100,
					UpdatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000006"),
					AuthorID:        uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					ParentCommentID: uuid.MustParse("00000000-0000-0000-0000-000000000005"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					Body:            "This is reply to parent comment 3",
					BodyHtml:        "<p>This is reply to parent comment 3</p>",
					Ups:             1,
					Score:           1,
					CreatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix:   1725091100,
					UpdatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000007"),
					AuthorID:        uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					ParentCommentID: uuid.MustParse("00000000-0000-0000-0000-000000000000"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					Body:            "This is a parent comment 4",
					BodyHtml:        "<p>This is a parent comment 4</p>",
					Ups:             1,
					Score:           1,
					CreatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix:   1725091100,
					UpdatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
			},
			wantErr: commentrepo.ErrCommentParentTableRecordNotFound,
		},
		{
			name:         "update comment :NEG",
			fixtureFiles: []string{"topics.yml", "voxspheres.yml", "users.yml", "posts.yml", "comments.yml"},
			args: args{
				comment: models.Comment{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					AuthorID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					ParentCommentID: uuid.MustParse("00000000-0000-0000-0000-000000000000"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Body:            "This is a parent comment 1 updated",
					BodyHtml:        "<p>This is a parent comment 1 updated</p>",
					Ups:             1,
					Score:           1,
					CreatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix:   1725091100,
					UpdatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
			},
			wantComment: models.Comment{
				ID:              uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				AuthorID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				ParentCommentID: uuid.MustParse("00000000-0000-0000-0000-000000000000"),
				PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				Body:            "This is a parent comment 1 updated",
				BodyHtml:        "<p>This is a parent comment 1 updated</p>",
				Ups:             1,
				Score:           1,
				CreatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				CreatedAtUnix:   1725091100,
				UpdatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
			},
			wantComments: []models.Comment{
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					AuthorID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					ParentCommentID: uuid.MustParse("00000000-0000-0000-0000-000000000000"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Body:            "This is a parent comment 1 updated",
					BodyHtml:        "<p>This is a parent comment 1 updated</p>",
					Ups:             1,
					Score:           1,
					CreatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix:   1725091100,
					UpdatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					AuthorID:        uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					ParentCommentID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Body:            "This is reply 1 to parent comment 1",
					BodyHtml:        "<p>This is reply 1 to parent comment 1</p>",
					Ups:             1,
					Score:           1,
					CreatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix:   1725091100,
					UpdatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					AuthorID:        uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					ParentCommentID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Body:            "This is reply to reply 1",
					BodyHtml:        "<p>This is reply to reply 1</p>",
					Ups:             1,
					Score:           1,
					CreatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix:   1725091100,
					UpdatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000004"),
					AuthorID:        uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					ParentCommentID: uuid.MustParse("00000000-0000-0000-0000-000000000000"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Body:            "This is a parent comment 2",
					BodyHtml:        "<p>This is a parent comment 2</p>",
					Ups:             1,
					Score:           1,
					CreatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix:   1725091100,
					UpdatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000005"),
					AuthorID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					ParentCommentID: uuid.MustParse("00000000-0000-0000-0000-000000000000"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					Body:            "This is a parent comment 3",
					BodyHtml:        "<p>This is a parent comment 3</p>",
					Ups:             1,
					Score:           1,
					CreatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix:   1725091100,
					UpdatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000006"),
					AuthorID:        uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					ParentCommentID: uuid.MustParse("00000000-0000-0000-0000-000000000005"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					Body:            "This is reply to parent comment 3",
					BodyHtml:        "<p>This is reply to parent comment 3</p>",
					Ups:             1,
					Score:           1,
					CreatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix:   1725091100,
					UpdatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000007"),
					AuthorID:        uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					ParentCommentID: uuid.MustParse("00000000-0000-0000-0000-000000000000"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					Body:            "This is a parent comment 4",
					BodyHtml:        "<p>This is a parent comment 4</p>",
					Ups:             1,
					Score:           1,
					CreatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix:   1725091100,
					UpdatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupPostgres(t, tt.fixtureFiles...)
			pgrepo := commentrepo.NewRepo(db)

			startTime := time.Now()
			gotComment, gotErr := pgrepo.UpdateComment(context.Background(), tt.args.comment)
			endTime := time.Now()

			assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")
			if tt.wantErr == nil {
				assertTimeWithinRange(t, gotComment.UpdatedAt, startTime, endTime)
			}

			gotComments, err := pgrepo.Comments(context.Background())

			assert.NoError(t, err, "expect no error while getting comments")
			assertCommentsWithoutTimestamp(t, tt.wantComments, gotComments)
		})
	}
}

func TestRepo_DeleteComment(t *testing.T) {
	type args struct {
		ID uuid.UUID
	}
	tests := []struct {
		name         string
		fixtureFiles []string
		args         args
		wantComments []models.Comment
		wantErr      error
	}{
		{
			name:         "comment not found :NEG",
			fixtureFiles: []string{"topics.yml", "voxspheres.yml", "users.yml", "posts.yml", "comments.yml"},
			args: args{
				ID: uuid.MustParse("00000000-0000-0000-0000-000000000009"),
			},
			wantComments: []models.Comment{
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					AuthorID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					ParentCommentID: uuid.MustParse("00000000-0000-0000-0000-000000000000"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Body:            "This is a parent comment 1",
					BodyHtml:        "<p>This is a parent comment 1</p>",
					Ups:             1,
					Score:           1,
					CreatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix:   1725091100,
					UpdatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					AuthorID:        uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					ParentCommentID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Body:            "This is reply 1 to parent comment 1",
					BodyHtml:        "<p>This is reply 1 to parent comment 1</p>",
					Ups:             1,
					Score:           1,
					CreatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix:   1725091100,
					UpdatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					AuthorID:        uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					ParentCommentID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Body:            "This is reply to reply 1",
					BodyHtml:        "<p>This is reply to reply 1</p>",
					Ups:             1,
					Score:           1,
					CreatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix:   1725091100,
					UpdatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000004"),
					AuthorID:        uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					ParentCommentID: uuid.MustParse("00000000-0000-0000-0000-000000000000"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Body:            "This is a parent comment 2",
					BodyHtml:        "<p>This is a parent comment 2</p>",
					Ups:             1,
					Score:           1,
					CreatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix:   1725091100,
					UpdatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000005"),
					AuthorID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					ParentCommentID: uuid.MustParse("00000000-0000-0000-0000-000000000000"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					Body:            "This is a parent comment 3",
					BodyHtml:        "<p>This is a parent comment 3</p>",
					Ups:             1,
					Score:           1,
					CreatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix:   1725091100,
					UpdatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000006"),
					AuthorID:        uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					ParentCommentID: uuid.MustParse("00000000-0000-0000-0000-000000000005"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					Body:            "This is reply to parent comment 3",
					BodyHtml:        "<p>This is reply to parent comment 3</p>",
					Ups:             1,
					Score:           1,
					CreatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix:   1725091100,
					UpdatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000007"),
					AuthorID:        uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					ParentCommentID: uuid.MustParse("00000000-0000-0000-0000-000000000000"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					Body:            "This is a parent comment 4",
					BodyHtml:        "<p>This is a parent comment 4</p>",
					Ups:             1,
					Score:           1,
					CreatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix:   1725091100,
					UpdatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
			},
			wantErr: commentrepo.ErrCommentNotFound,
		},
		{
			name:         "delete comment :POS",
			fixtureFiles: []string{"topics.yml", "voxspheres.yml", "users.yml", "posts.yml", "comments.yml"},
			args: args{
				ID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			},
			wantComments: []models.Comment{
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					AuthorID:        uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					ParentCommentID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Body:            "This is reply 1 to parent comment 1",
					BodyHtml:        "<p>This is reply 1 to parent comment 1</p>",
					Ups:             1,
					Score:           1,
					CreatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix:   1725091100,
					UpdatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					AuthorID:        uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					ParentCommentID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Body:            "This is reply to reply 1",
					BodyHtml:        "<p>This is reply to reply 1</p>",
					Ups:             1,
					Score:           1,
					CreatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix:   1725091100,
					UpdatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000004"),
					AuthorID:        uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					ParentCommentID: uuid.MustParse("00000000-0000-0000-0000-000000000000"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Body:            "This is a parent comment 2",
					BodyHtml:        "<p>This is a parent comment 2</p>",
					Ups:             1,
					Score:           1,
					CreatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix:   1725091100,
					UpdatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000005"),
					AuthorID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					ParentCommentID: uuid.MustParse("00000000-0000-0000-0000-000000000000"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					Body:            "This is a parent comment 3",
					BodyHtml:        "<p>This is a parent comment 3</p>",
					Ups:             1,
					Score:           1,
					CreatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix:   1725091100,
					UpdatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000006"),
					AuthorID:        uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					ParentCommentID: uuid.MustParse("00000000-0000-0000-0000-000000000005"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					Body:            "This is reply to parent comment 3",
					BodyHtml:        "<p>This is reply to parent comment 3</p>",
					Ups:             1,
					Score:           1,
					CreatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix:   1725091100,
					UpdatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000007"),
					AuthorID:        uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					ParentCommentID: uuid.MustParse("00000000-0000-0000-0000-000000000000"),
					PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					Body:            "This is a parent comment 4",
					BodyHtml:        "<p>This is a parent comment 4</p>",
					Ups:             1,
					Score:           1,
					CreatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix:   1725091100,
					UpdatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupPostgres(t, tt.fixtureFiles...)
			pgrepo := commentrepo.NewRepo(db)

			gotErr := pgrepo.DeleteComment(context.Background(), tt.args.ID)

			assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")

			gotComments, err := pgrepo.Comments(context.Background())

			assert.NoError(t, err, "expect no error while getting comment")
			assert.Equal(t, tt.wantComments, gotComments, "expect comments to match")
		})
	}
}

func TestRepo_ForeignKeyCascade(t *testing.T) {
	t.Run("on deleting post from parent table , no child references should exist in comments table", func(t *testing.T) {
		db := setupPostgres(t, "topics.yml", "voxspheres.yml", "users.yml", "posts.yml", "comments.yml")
		postPgrepo := postrepo.NewRepo(db)
		comments := commentrepo.NewRepo(db)

		wantComments := []models.Comment{
			{
				ID:              uuid.MustParse("00000000-0000-0000-0000-000000000005"),
				AuthorID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				ParentCommentID: uuid.MustParse("00000000-0000-0000-0000-000000000000"),
				PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				Body:            "This is a parent comment 3",
				BodyHtml:        "<p>This is a parent comment 3</p>",
				Ups:             1,
				Score:           1,
				CreatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				CreatedAtUnix:   1725091100,
				UpdatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
			},
			{
				ID:              uuid.MustParse("00000000-0000-0000-0000-000000000006"),
				AuthorID:        uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				ParentCommentID: uuid.MustParse("00000000-0000-0000-0000-000000000005"),
				PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				Body:            "This is reply to parent comment 3",
				BodyHtml:        "<p>This is reply to parent comment 3</p>",
				Ups:             1,
				Score:           1,
				CreatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				CreatedAtUnix:   1725091100,
				UpdatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
			},
			{
				ID:              uuid.MustParse("00000000-0000-0000-0000-000000000007"),
				AuthorID:        uuid.MustParse("00000000-0000-0000-0000-000000000003"),
				ParentCommentID: uuid.MustParse("00000000-0000-0000-0000-000000000000"),
				PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				Body:            "This is a parent comment 4",
				BodyHtml:        "<p>This is a parent comment 4</p>",
				Ups:             1,
				Score:           1,
				CreatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				CreatedAtUnix:   1725091100,
				UpdatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
			},
		}

		err := postPgrepo.DeletePost(context.Background(), uuid.MustParse("00000000-0000-0000-0000-000000000001"))

		assert.NoError(t, err, "expect no error while deleting post")

		gotComments, err := comments.Comments(context.Background())

		assert.NoError(t, err, "expect no error while getting comments")
		assertCommentsWithoutTimestamp(t, wantComments, gotComments)
	})

	t.Run("on deleting author from parent table , no child references should exist in comments table", func(t *testing.T) {
		db := setupPostgres(t, "topics.yml", "voxspheres.yml", "users.yml", "posts.yml", "comments.yml")
		userPgrepo := userrepo.NewRepo(db)
		comments := commentrepo.NewRepo(db)

		wantComments := []models.Comment{
			{
				ID:              uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				AuthorID:        uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				ParentCommentID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				Body:            "This is reply 1 to parent comment 1",
				BodyHtml:        "<p>This is reply 1 to parent comment 1</p>",
				Ups:             1,
				Score:           1,
				CreatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				CreatedAtUnix:   1725091100,
				UpdatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
			},
			{
				ID:              uuid.MustParse("00000000-0000-0000-0000-000000000003"),
				AuthorID:        uuid.MustParse("00000000-0000-0000-0000-000000000003"),
				ParentCommentID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				Body:            "This is reply to reply 1",
				BodyHtml:        "<p>This is reply to reply 1</p>",
				Ups:             1,
				Score:           1,
				CreatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				CreatedAtUnix:   1725091100,
				UpdatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
			},
			{
				ID:              uuid.MustParse("00000000-0000-0000-0000-000000000004"),
				AuthorID:        uuid.MustParse("00000000-0000-0000-0000-000000000003"),
				ParentCommentID: uuid.MustParse("00000000-0000-0000-0000-000000000000"),
				PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				Body:            "This is a parent comment 2",
				BodyHtml:        "<p>This is a parent comment 2</p>",
				Ups:             1,
				Score:           1,
				CreatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				CreatedAtUnix:   1725091100,
				UpdatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
			},
			{
				ID:              uuid.MustParse("00000000-0000-0000-0000-000000000006"),
				AuthorID:        uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				ParentCommentID: uuid.MustParse("00000000-0000-0000-0000-000000000005"),
				PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				Body:            "This is reply to parent comment 3",
				BodyHtml:        "<p>This is reply to parent comment 3</p>",
				Ups:             1,
				Score:           1,
				CreatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				CreatedAtUnix:   1725091100,
				UpdatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
			},
			{
				ID:              uuid.MustParse("00000000-0000-0000-0000-000000000007"),
				AuthorID:        uuid.MustParse("00000000-0000-0000-0000-000000000003"),
				ParentCommentID: uuid.MustParse("00000000-0000-0000-0000-000000000000"),
				PostID:          uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				Body:            "This is a parent comment 4",
				BodyHtml:        "<p>This is a parent comment 4</p>",
				Ups:             1,
				Score:           1,
				CreatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				CreatedAtUnix:   1725091100,
				UpdatedAt:       time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
			},
		}

		err := userPgrepo.DeleteUser(context.Background(), uuid.MustParse("00000000-0000-0000-0000-000000000001"))

		assert.NoError(t, err, "expect no error while deleting user")

		gotComments, err := comments.Comments(context.Background())

		assert.NoError(t, err, "expect no error while getting comments")
		assertCommentsWithoutTimestamp(t, wantComments, gotComments)
	})
}
