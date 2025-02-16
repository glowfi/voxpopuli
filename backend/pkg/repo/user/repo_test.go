package user_test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"slices"
	"testing"
	"time"

	"github.com/glowfi/voxpopuli/backend/pkg/models"
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

	db.RegisterModel((*models.User)(nil))

	// drop all rows of the user table
	_, err := db.NewTruncateTable().Cascade().Model((*models.User)(nil)).Exec(context.Background())
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

func assertUsers(t *testing.T, wantUsers, gotUsers []models.User) {
	t.Helper()

	for _, user := range wantUsers {
		idx := slices.IndexFunc(gotUsers, func(v models.User) bool {
			return v.ID == user.ID
		})

		if idx == -1 {
			t.Fatal(fmt.Sprintf("user %v of ID %v is not present in gotUsers", user.Name, user.ID))
			return
		}
		assert.Equal(t, user, gotUsers[idx], "expect user to match")
	}
}

func ptrof[T comparable](v T) *T {
	return &v
}

func TestRepo_Users(t *testing.T) {
	tests := []struct {
		name         string
		fixtureFiles []string
		wantUsers    []models.User
		wantErr      error
	}{
		{
			name:         "users :POS",
			fixtureFiles: []string{"users.yml"},
			wantUsers: []models.User{
				{
					ID:                uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Name:              "John Doe",
					PublicDescription: ptrof("This is a public description"),
					AvatarImg:         ptrof("https://example.com/avatar1.jpg"),
					BannerImg:         ptrof("https://example.com/banner1.jpg"),
					Iconcolor:         ptrof("#FF0000"),
					Keycolor:          ptrof("#00FF00"),
					Primarycolor:      ptrof("#0000FF"),
					Over18:            true,
					Suspended:         false,
					CreatedAt:         time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix:     1725091100,
					UpdatedAt:         time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:                uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					Name:              "Jane Doe",
					PublicDescription: ptrof("This is another public description"),
					AvatarImg:         ptrof("https://example.com/avatar2.jpg"),
					BannerImg:         ptrof("https://example.com/banner2.jpg"),
					Iconcolor:         ptrof("#FFFF00"),
					Keycolor:          ptrof("#FF00FF"),
					Primarycolor:      ptrof("#00FFFF"),
					Over18:            true,
					Suspended:         false,
					CreatedAt:         time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
					CreatedAtUnix:     1725091101,
					UpdatedAt:         time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
				},
			},
			wantErr: nil,
		},
		{
			name:         "no users :POS",
			fixtureFiles: []string{},
			wantUsers:    nil,
			wantErr:      nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupPostgres(t, tt.fixtureFiles...)
			pgrepo := userrepo.NewRepo(db)

			gotUsers, gotErr := pgrepo.Users(context.Background())

			assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")
			assert.Equal(t, tt.wantUsers, gotUsers, "expect users to match")
		})
	}
}

func TestRepo_UserByID(t *testing.T) {
	type args struct {
		ID uuid.UUID
	}
	tests := []struct {
		name         string
		fixtureFiles []string
		args         args
		wantUser     models.User
		wantErr      error
	}{
		{
			name:         "user not found :NEG",
			fixtureFiles: []string{"users.yml"},
			args: args{
				ID: uuid.MustParse("00000000-0000-0000-0000-000000000006"),
			},
			wantUser: models.User{},
			wantErr:  userrepo.ErrUserNotFound,
		},
		{
			name:         "user found :NEG",
			fixtureFiles: []string{"users.yml"},
			args: args{
				ID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			},
			wantUser: models.User{
				ID:                uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				Name:              "John Doe",
				PublicDescription: ptrof("This is a public description"),
				AvatarImg:         ptrof("https://example.com/avatar1.jpg"),
				BannerImg:         ptrof("https://example.com/banner1.jpg"),
				Iconcolor:         ptrof("#FF0000"),
				Keycolor:          ptrof("#00FF00"),
				Primarycolor:      ptrof("#0000FF"),
				Over18:            true,
				Suspended:         false,
				CreatedAt:         time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				CreatedAtUnix:     1725091100,
				UpdatedAt:         time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupPostgres(t, tt.fixtureFiles...)
			pgrepo := userrepo.NewRepo(db)

			gotUser, gotErr := pgrepo.UserByID(context.Background(), tt.args.ID)

			assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")
			assert.Equal(t, tt.wantUser, gotUser, "expect user to match")
		})
	}
}
