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

	// add query logging hook
	db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))

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

func assertUsersWithoutTimestamp(t *testing.T, wantUsers, gotUsers []models.User) {
	t.Helper()

	if len(wantUsers) != len(gotUsers) {
		t.Fatal("length of wantUsers and gotUsers do not match")
	}

	for _, user := range wantUsers {
		idx := slices.IndexFunc(gotUsers, func(v models.User) bool {
			return v.ID == user.ID
		})

		if idx == -1 {
			t.Fatalf("user %v of ID %v is not present in gotUsers", user.Name, user.ID)
			return
		}
		assertUserWithoutTimestamp(t, user, gotUsers[idx])
	}
}

func assertUserWithoutTimestamp(t *testing.T, wantUser, gotUser models.User) {
	assert.Equal(t, wantUser.ID, gotUser.ID, "expected id to match")
	assert.Equal(t, wantUser.Name, gotUser.Name, "expected name to match")
	assert.Equal(t, wantUser.PublicDescription, gotUser.PublicDescription, "expected public description to match")
	assert.Equal(t, wantUser.AvatarImg, gotUser.AvatarImg, "expected avatar image to match")
	assert.Equal(t, wantUser.BannerImg, gotUser.BannerImg, "expected banner image to match")
	assert.Equal(t, wantUser.Iconcolor, gotUser.Iconcolor, "expected icon color to match")
	assert.Equal(t, wantUser.Keycolor, gotUser.Keycolor, "expected key color to match")
	assert.Equal(t, wantUser.Primarycolor, gotUser.Primarycolor, "expected primary color to match")
	assert.Equal(t, wantUser.Over18, gotUser.Over18, "expected over 18 to match")
	assert.Equal(t, wantUser.Suspended, gotUser.Suspended, "expected suspended to match")
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

func TestRepo_AddUsers(t *testing.T) {
	type args struct {
		users []models.User
	}
	tests := []struct {
		name              string
		fixtureFiles      []string
		args              args
		wantInsertedUsers []models.User
		wantUsers         []models.User
		wantErr           error
	}{
		{
			name:         "duplicate user id :NEG",
			fixtureFiles: []string{"users.yml"},
			args: args{
				users: []models.User{
					{
						ID:                uuid.MustParse("00000000-0000-0000-0000-000000000001"),
						Name:              "new user",
						PublicDescription: ptrof("A public description"),
						AvatarImg:         ptrof("https://example.com/avatar.jpg"),
						BannerImg:         ptrof("https://example.com/banner.jpg"),
						Iconcolor:         ptrof("#ffffff"),
						Keycolor:          ptrof("#ffffff"),
						Primarycolor:      ptrof("#ffffff"),
						Over18:            false,
						Suspended:         false,
						CreatedAt:         time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
						CreatedAtUnix:     1725091101,
						UpdatedAt:         time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
					},
				},
			},
			wantInsertedUsers: nil,
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
			wantErr: userrepo.ErrUserDuplicateIDorName,
		},
		{
			name:         "add users :POS",
			fixtureFiles: []string{"users.yml"},
			args: args{
				users: []models.User{
					{
						ID:                uuid.MustParse("00000000-0000-0000-0000-000000000003"),
						Name:              "Jane Smith",
						PublicDescription: ptrof("A public description"),
						AvatarImg:         ptrof("https://example.com/avatar.jpg"),
						BannerImg:         ptrof("https://example.com/banner.jpg"),
						Iconcolor:         ptrof("#ffffff"),
						Keycolor:          ptrof("#ffffff"),
						Primarycolor:      ptrof("#ffffff"),
						Over18:            false,
						Suspended:         false,
						CreatedAt:         time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
						CreatedAtUnix:     1725091101,
						UpdatedAt:         time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
					},
				},
			},
			wantInsertedUsers: []models.User{
				{
					ID:                uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					Name:              "Jane Smith",
					PublicDescription: ptrof("A public description"),
					AvatarImg:         ptrof("https://example.com/avatar.jpg"),
					BannerImg:         ptrof("https://example.com/banner.jpg"),
					Iconcolor:         ptrof("#ffffff"),
					Keycolor:          ptrof("#ffffff"),
					Primarycolor:      ptrof("#ffffff"),
					Over18:            false,
					Suspended:         false,
					CreatedAt:         time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
					CreatedAtUnix:     1725091101,
					UpdatedAt:         time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
				},
			},
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
				{
					ID:                uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					Name:              "Jane Smith",
					PublicDescription: ptrof("A public description"),
					AvatarImg:         ptrof("https://example.com/avatar.jpg"),
					BannerImg:         ptrof("https://example.com/banner.jpg"),
					Iconcolor:         ptrof("#ffffff"),
					Keycolor:          ptrof("#ffffff"),
					Primarycolor:      ptrof("#ffffff"),
					Over18:            false,
					Suspended:         false,
					CreatedAt:         time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
					CreatedAtUnix:     1725091101,
					UpdatedAt:         time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupPostgres(t, tt.fixtureFiles...)
			pgrepo := userrepo.NewRepo(db)

			startTime := time.Now()
			gotInsertedUsers, gotErr := pgrepo.AddUsers(context.Background(), tt.args.users...)
			endTime := time.Now()

			for _, gotInsertedUser := range gotInsertedUsers {
				assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")
				assert.Equal(
					t,
					gotInsertedUser.UpdatedAt,
					gotInsertedUser.CreatedAt,
					"expect CreatedAt and UpdatedAt to be same",
				)
				if tt.wantErr == nil {
					assertTimeWithinRange(t, gotInsertedUser.CreatedAt, startTime, endTime)
					assertTimeWithinRange(t, gotInsertedUser.UpdatedAt, startTime, endTime)
				}
			}

			gotUsers, err := pgrepo.Users(context.Background())

			assert.NoError(t, err, "expect no error while getting users")
			assertUsersWithoutTimestamp(t, tt.wantUsers, gotUsers)
		})
	}
}

func TestRepo_UpdateUser(t *testing.T) {
	type args struct {
		user models.User
	}
	tests := []struct {
		name         string
		fixtureFiles []string
		args         args
		wantUser     models.User
		wantUsers    []models.User
		wantErr      error
	}{
		{
			name:         "user not found :NEG",
			fixtureFiles: []string{"users.yml"},
			args: args{
				user: models.User{
					ID:                uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					Name:              "Jane Smith",
					PublicDescription: ptrof("A public description"),
					AvatarImg:         ptrof("https://example.com/avatar.jpg"),
					BannerImg:         ptrof("https://example.com/banner.jpg"),
					Iconcolor:         ptrof("#ffffff"),
					Keycolor:          ptrof("#ffffff"),
					Primarycolor:      ptrof("#ffffff"),
					Over18:            false,
					Suspended:         false,
					CreatedAt:         time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
					CreatedAtUnix:     1725091101,
					UpdatedAt:         time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
				},
			},
			wantUser: models.User{},
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
			wantErr: userrepo.ErrUserNotFound,
		},
		{
			name:         "update user :POS",
			fixtureFiles: []string{"users.yml"},
			args: args{
				user: models.User{
					ID:                uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Name:              "Jane Smith",
					PublicDescription: ptrof("A public description updated"),
					AvatarImg:         ptrof("https://example.com/avatar.jpg"),
					BannerImg:         ptrof("https://example.com/banner.jpg"),
					Iconcolor:         ptrof("#ffffff"),
					Keycolor:          ptrof("#ffffff"),
					Primarycolor:      ptrof("#ffffff"),
					Over18:            false,
					Suspended:         false,
					CreatedAt:         time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
					CreatedAtUnix:     1725091101,
					UpdatedAt:         time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
				},
			},
			wantUser: models.User{
				ID:                uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				Name:              "Jane Smith",
				PublicDescription: ptrof("A public description updated"),
				AvatarImg:         ptrof("https://example.com/avatar.jpg"),
				BannerImg:         ptrof("https://example.com/banner.jpg"),
				Iconcolor:         ptrof("#ffffff"),
				Keycolor:          ptrof("#ffffff"),
				Primarycolor:      ptrof("#ffffff"),
				Over18:            false,
				Suspended:         false,
				CreatedAt:         time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
				CreatedAtUnix:     1725091101,
				UpdatedAt:         time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
			},
			wantUsers: []models.User{
				{
					ID:                uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Name:              "Jane Smith",
					PublicDescription: ptrof("A public description updated"),
					AvatarImg:         ptrof("https://example.com/avatar.jpg"),
					BannerImg:         ptrof("https://example.com/banner.jpg"),
					Iconcolor:         ptrof("#ffffff"),
					Keycolor:          ptrof("#ffffff"),
					Primarycolor:      ptrof("#ffffff"),
					Over18:            false,
					Suspended:         false,
					CreatedAt:         time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
					CreatedAtUnix:     1725091101,
					UpdatedAt:         time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupPostgres(t, tt.fixtureFiles...)
			pgrepo := userrepo.NewRepo(db)

			startTime := time.Now()
			gotUser, gotErr := pgrepo.UpdateUser(context.Background(), tt.args.user)
			endTime := time.Now()

			assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")
			if tt.wantErr == nil {
				assertTimeWithinRange(t, gotUser.UpdatedAt, startTime, endTime)
			}

			gotUsers, err := pgrepo.Users(context.Background())

			assert.NoError(t, err, "expect no error while getting users")
			assertUsersWithoutTimestamp(t, tt.wantUsers, gotUsers)
		})
	}
}

func TestRepo_DeleteUser(t *testing.T) {
	type args struct {
		ID uuid.UUID
	}
	tests := []struct {
		name         string
		fixtureFiles []string
		args         args
		wantUsers    []models.User
		wantErr      error
	}{
		{
			name:         "user not found :NEG",
			fixtureFiles: []string{"users.yml"},
			args: args{
				ID: uuid.MustParse("00000000-0000-0000-0000-000000000003"),
			},
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
			wantErr: userrepo.ErrUserNotFound,
		},
		{
			name:         "delete user :POS",
			fixtureFiles: []string{"users.yml"},
			args: args{
				ID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			},
			wantUsers: []models.User{
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupPostgres(t, tt.fixtureFiles...)
			pgrepo := userrepo.NewRepo(db)

			gotErr := pgrepo.DeleteUser(context.Background(), tt.args.ID)

			assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")

			gotUsers, err := pgrepo.Users(context.Background())

			assert.NoError(t, err, "expect no error while getting users")
			assert.Equal(t, tt.wantUsers, gotUsers, "expect users to match")
		})
	}
}
