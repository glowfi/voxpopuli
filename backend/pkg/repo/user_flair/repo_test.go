package userflair_test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"slices"
	"testing"

	"github.com/glowfi/voxpopuli/backend/pkg/models"
	userflaireRepo "github.com/glowfi/voxpopuli/backend/pkg/repo/user_flair"
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

	db.RegisterModel((*models.Topic)(nil))
	db.RegisterModel((*models.Voxsphere)(nil))
	db.RegisterModel((*models.User)(nil))
	db.RegisterModel((*models.UserFlair)(nil))

	// drop all rows of the topic,voxsphere,user,user_flair table
	if _, err := db.NewTruncateTable().Cascade().Model((*models.Topic)(nil)).Exec(context.Background()); err != nil {
		t.Fatal("truncate table failed:", err)
	}
	if _, err := db.NewTruncateTable().Cascade().Model((*models.Voxsphere)(nil)).Exec(context.Background()); err != nil {
		t.Fatal("truncate table failed:", err)
	}
	if _, err := db.NewTruncateTable().Cascade().Model((*models.User)(nil)).Exec(context.Background()); err != nil {
		t.Fatal("truncate table failed:", err)
	}
	if _, err := db.NewTruncateTable().Cascade().Model((*models.UserFlair)(nil)).Exec(context.Background()); err != nil {
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

func assertUserFlairs(t *testing.T, wantUserFlairs, gotUserFlairs []models.UserFlair) {
	t.Helper()

	for _, userFlair := range wantUserFlairs {
		idx := slices.IndexFunc(gotUserFlairs, func(v models.UserFlair) bool {
			return v.ID == userFlair.ID
		})

		if idx == -1 {
			t.Fatal(fmt.Sprintf("user flair %v of ID %v is not present in gotUserFlairs", userFlair.FullText, userFlair.ID))
			return
		}
		assert.Equal(t, userFlair, gotUserFlairs[idx], "expect user flair to match")
	}
}

func TestRepo_UserFlairs(t *testing.T) {
	tests := []struct {
		name           string
		fixtureFiles   []string
		wantUserFlairs []models.UserFlair
		wantErr        error
	}{
		{
			name:         "user flairs :POS",
			fixtureFiles: []string{"topics.yml", "voxspheres.yml", "users.yml", "user_flairs.yml"},
			wantUserFlairs: []models.UserFlair{
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					UserID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					VoxsphereID:     uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					FullText:        "Flair 1",
					BackgroundColor: "#FFFFFF",
				},
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					UserID:          uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					VoxsphereID:     uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					FullText:        "Flair 2",
					BackgroundColor: "#000000",
				},
			},
			wantErr: nil,
		},
		{
			name:           "no user flairs :POS",
			fixtureFiles:   []string{},
			wantUserFlairs: nil,
			wantErr:        nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupPostgres(t, tt.fixtureFiles...)
			pgrepo := userflaireRepo.NewRepo(db)

			gotUserFlairs, gotErr := pgrepo.UserFlairs(context.Background())

			assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")
			assert.Equal(t, tt.wantUserFlairs, gotUserFlairs, "expect user flairs to match")
		})
	}
}

func TestRepo_UserFlairByID(t *testing.T) {
	type args struct {
		ID uuid.UUID
	}
	tests := []struct {
		name          string
		fixtureFiles  []string
		args          args
		wantUserFlair models.UserFlair
		wantErr       error
	}{
		{
			name:         "user flair not found :NEG",
			fixtureFiles: []string{"topics.yml", "voxspheres.yml", "users.yml", "user_flairs.yml"},
			args: args{
				ID: uuid.MustParse("00000000-0000-0000-0000-000000000006"),
			},
			wantUserFlair: models.UserFlair{},
			wantErr:       userflaireRepo.ErrUserFlairNotFound,
		},
		{
			name:         "user flair found :POS",
			fixtureFiles: []string{"topics.yml", "voxspheres.yml", "users.yml", "user_flairs.yml"},
			args: args{
				ID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			},
			wantUserFlair: models.UserFlair{
				ID:              uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				UserID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
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
			pgrepo := userflaireRepo.NewRepo(db)

			gotUserFlair, gotErr := pgrepo.UserFlairByID(context.Background(), tt.args.ID)

			assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")
			assert.Equal(t, tt.wantUserFlair, gotUserFlair, "expect user flair to match")
		})
	}
}

func TestRepo_AddUserFlair(t *testing.T) {
	type args struct {
		userFlair models.UserFlair
	}
	tests := []struct {
		name           string
		fixtureFiles   []string
		args           args
		wantUserFlair  models.UserFlair
		wantUserFlairs []models.UserFlair
		wantErr        error
	}{
		{
			name:         "duplicate user flair id :NEG",
			fixtureFiles: []string{"topics.yml", "voxspheres.yml", "users.yml", "user_flairs.yml"},
			args: args{
				userFlair: models.UserFlair{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					UserID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					VoxsphereID:     uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					FullText:        "new text",
					BackgroundColor: "#FFFFFF",
				},
			},
			wantUserFlair: models.UserFlair{},
			wantUserFlairs: []models.UserFlair{
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					UserID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					VoxsphereID:     uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					FullText:        "Flair 1",
					BackgroundColor: "#FFFFFF",
				},
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					UserID:          uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					VoxsphereID:     uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					FullText:        "Flair 2",
					BackgroundColor: "#000000",
				},
			},
			wantErr: userflaireRepo.ErrUserFlairDuplicateID,
		},
		{
			name:         "voxsphere not present in parent table :NEG",
			fixtureFiles: []string{"topics.yml", "voxspheres.yml", "users.yml", "user_flairs.yml"},
			args: args{
				userFlair: models.UserFlair{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					UserID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					VoxsphereID:     uuid.MustParse("00000000-0000-0000-0000-000000000006"),
					FullText:        "new text",
					BackgroundColor: "#FFFFFF",
				},
			},
			wantUserFlair: models.UserFlair{},
			wantUserFlairs: []models.UserFlair{
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					UserID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					VoxsphereID:     uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					FullText:        "Flair 1",
					BackgroundColor: "#FFFFFF",
				},
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					UserID:          uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					VoxsphereID:     uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					FullText:        "Flair 2",
					BackgroundColor: "#000000",
				},
			},
			wantErr: userflaireRepo.ErrUserFlairParentTableRecordNotFound,
		},
		{
			name:         "add user flair :POS",
			fixtureFiles: []string{"topics.yml", "voxspheres.yml", "users.yml", "user_flairs.yml"},
			args: args{
				userFlair: models.UserFlair{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					UserID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					VoxsphereID:     uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					FullText:        "new text",
					BackgroundColor: "#FFFFFF",
				},
			},
			wantUserFlair: models.UserFlair{
				ID:              uuid.MustParse("00000000-0000-0000-0000-000000000003"),
				UserID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				VoxsphereID:     uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				FullText:        "new text",
				BackgroundColor: "#FFFFFF",
			},
			wantUserFlairs: []models.UserFlair{
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					UserID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					VoxsphereID:     uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					FullText:        "Flair 1",
					BackgroundColor: "#FFFFFF",
				},
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					UserID:          uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					VoxsphereID:     uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					FullText:        "Flair 2",
					BackgroundColor: "#000000",
				},
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					UserID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					VoxsphereID:     uuid.MustParse("00000000-0000-0000-0000-000000000002"),
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
			pgrepo := userflaireRepo.NewRepo(db)

			gotUserFlair, gotErr := pgrepo.AddUserFlair(context.Background(), tt.args.userFlair)
			gotUserFlairs, err := pgrepo.UserFlairs(context.Background())
			if err != nil {
				t.Fatal("expect no error while getting user flairs")
			}

			assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")
			assert.Equal(t, tt.wantUserFlair, gotUserFlair, "expect user flair to match")
			assertUserFlairs(t, tt.wantUserFlairs, gotUserFlairs)
		})
	}
}

func TestRepo_UpdateUserFlair(t *testing.T) {
	type args struct {
		userFlair models.UserFlair
	}
	tests := []struct {
		name           string
		fixtureFiles   []string
		args           args
		wantUserFlair  models.UserFlair
		wantUserFlairs []models.UserFlair
		wantErr        error
	}{
		{
			name:         "user flair not found :NEG",
			fixtureFiles: []string{"topics.yml", "voxspheres.yml", "users.yml", "user_flairs.yml"},
			args: args{
				userFlair: models.UserFlair{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					UserID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					VoxsphereID:     uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					FullText:        "updated text",
					BackgroundColor: "#000000",
				},
			},
			wantUserFlair: models.UserFlair{},
			wantUserFlairs: []models.UserFlair{
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					UserID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					VoxsphereID:     uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					FullText:        "Flair 1",
					BackgroundColor: "#FFFFFF",
				},
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					UserID:          uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					VoxsphereID:     uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					FullText:        "Flair 2",
					BackgroundColor: "#000000",
				},
			},
			wantErr: userflaireRepo.ErrUserFlairNotFound,
		},
		{
			name:         "voxsphere is not present in parent table :NEG",
			fixtureFiles: []string{"topics.yml", "voxspheres.yml", "users.yml", "user_flairs.yml"},
			args: args{
				userFlair: models.UserFlair{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					UserID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					VoxsphereID:     uuid.MustParse("00000000-0000-0000-0000-000000000006"),
					FullText:        "updated text",
					BackgroundColor: "#000000",
				},
			},
			wantUserFlair: models.UserFlair{},
			wantUserFlairs: []models.UserFlair{
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					UserID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					VoxsphereID:     uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					FullText:        "Flair 1",
					BackgroundColor: "#FFFFFF",
				},
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					UserID:          uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					VoxsphereID:     uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					FullText:        "Flair 2",
					BackgroundColor: "#000000",
				},
			},
			wantErr: userflaireRepo.ErrUserFlairParentTableRecordNotFound,
		},
		{
			name:         "user is not present in parent table :NEG",
			fixtureFiles: []string{"topics.yml", "voxspheres.yml", "users.yml", "user_flairs.yml"},
			args: args{
				userFlair: models.UserFlair{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					UserID:          uuid.MustParse("00000000-0000-0000-0000-000000000006"),
					VoxsphereID:     uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					FullText:        "updated text",
					BackgroundColor: "#000000",
				},
			},
			wantUserFlair: models.UserFlair{},
			wantUserFlairs: []models.UserFlair{
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					UserID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					VoxsphereID:     uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					FullText:        "Flair 1",
					BackgroundColor: "#FFFFFF",
				},
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					UserID:          uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					VoxsphereID:     uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					FullText:        "Flair 2",
					BackgroundColor: "#000000",
				},
			},
			wantErr: userflaireRepo.ErrUserFlairParentTableRecordNotFound,
		},
		{
			name:         "update user flair :POS",
			fixtureFiles: []string{"topics.yml", "voxspheres.yml", "users.yml", "user_flairs.yml"},
			args: args{
				userFlair: models.UserFlair{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					UserID:          uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					VoxsphereID:     uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					FullText:        "updated text",
					BackgroundColor: "#000000",
				},
			},
			wantUserFlair: models.UserFlair{
				ID:              uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				UserID:          uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				VoxsphereID:     uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				FullText:        "updated text",
				BackgroundColor: "#000000",
			},
			wantUserFlairs: []models.UserFlair{
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					UserID:          uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					VoxsphereID:     uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					FullText:        "updated text",
					BackgroundColor: "#000000",
				},
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					UserID:          uuid.MustParse("00000000-0000-0000-0000-000000000002"),
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
			pgrepo := userflaireRepo.NewRepo(db)

			gotUserFlair, gotErr := pgrepo.UpdateUserFlair(context.Background(), tt.args.userFlair)
			gotUserFlairs, err := pgrepo.UserFlairs(context.Background())
			if err != nil {
				t.Fatal("expect no error while getting user flairs")
			}

			assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")
			assert.Equal(t, tt.wantUserFlair, gotUserFlair, "expect user flair to match")
			assertUserFlairs(t, tt.wantUserFlairs, gotUserFlairs)
		})
	}
}

func TestRepo_DeleteUserFlair(t *testing.T) {
	type args struct {
		ID uuid.UUID
	}
	tests := []struct {
		name           string
		fixtureFiles   []string
		args           args
		wantUserFlairs []models.UserFlair
		wantErr        error
	}{
		{
			name:         "user flair not found :NEG",
			fixtureFiles: []string{"topics.yml", "voxspheres.yml", "users.yml", "user_flairs.yml"},
			args: args{
				ID: uuid.MustParse("00000000-0000-0000-0000-000000000006"),
			},
			wantUserFlairs: []models.UserFlair{
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					UserID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					VoxsphereID:     uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					FullText:        "Flair 1",
					BackgroundColor: "#FFFFFF",
				},
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					UserID:          uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					VoxsphereID:     uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					FullText:        "Flair 2",
					BackgroundColor: "#000000",
				},
			},
			wantErr: userflaireRepo.ErrUserFlairNotFound,
		},
		{
			name:         "user flair deleted :POS",
			fixtureFiles: []string{"topics.yml", "voxspheres.yml", "users.yml", "user_flairs.yml"},
			args: args{
				ID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			},
			wantUserFlairs: []models.UserFlair{
				{
					ID:              uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					UserID:          uuid.MustParse("00000000-0000-0000-0000-000000000002"),
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
			pgrepo := userflaireRepo.NewRepo(db)

			gotErr := pgrepo.DeleteUserFlair(context.Background(), tt.args.ID)
			gotUserFlairs, err := pgrepo.UserFlairs(context.Background())
			if err != nil {
				t.Fatal("expect no error while getting user flairs")
			}

			assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")
			assert.Equal(t, tt.wantUserFlairs, gotUserFlairs, "expect flairs to match")
		})
	}
}
