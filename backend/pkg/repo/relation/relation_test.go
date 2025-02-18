package relation_test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/glowfi/voxpopuli/backend/pkg/models"
	relationrepo "github.com/glowfi/voxpopuli/backend/pkg/repo/relation"
	trophyrepo "github.com/glowfi/voxpopuli/backend/pkg/repo/trophy"
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

	db.RegisterModel((*models.Topic)(nil))
	db.RegisterModel((*models.Voxsphere)(nil))
	db.RegisterModel((*models.User)(nil))
	db.RegisterModel((*models.Trophy)(nil))
	db.RegisterModel((*models.UserTrophy)(nil))
	db.RegisterModel((*models.VoxsphereMember)(nil))

	// drop all rows of the user,trophies table
	if _, err := db.NewTruncateTable().Cascade().Model((*models.Topic)(nil)).Exec(context.Background()); err != nil {
		t.Fatal("truncate table failed:", err)
	}
	if _, err := db.NewTruncateTable().Cascade().Model((*models.Voxsphere)(nil)).Exec(context.Background()); err != nil {
		t.Fatal("truncate table failed:", err)
	}
	if _, err := db.NewTruncateTable().Cascade().Model((*models.User)(nil)).Exec(context.Background()); err != nil {
		t.Fatal("truncate table failed:", err)
	}
	if _, err := db.NewTruncateTable().Cascade().Model((*models.Trophy)(nil)).Exec(context.Background()); err != nil {
		t.Fatal("truncate table failed:", err)
	}
	if _, err := db.NewTruncateTable().Cascade().Model((*models.UserTrophy)(nil)).Exec(context.Background()); err != nil {
		t.Fatal("truncate table failed:", err)
	}
	if _, err := db.NewTruncateTable().Cascade().Model((*models.VoxsphereMember)(nil)).Exec(context.Background()); err != nil {
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

func TestRepo_UserTrophies(t *testing.T) {
	tests := []struct {
		name             string
		fixtureFiles     []string
		wantUserTrophies []models.UserTrophy
		wantErr          error
	}{
		{
			name:         "user trophies :POS",
			fixtureFiles: []string{"users.yml", "trophies.yml", "user_trophies.yml"},
			wantUserTrophies: []models.UserTrophy{
				{
					UserID:   uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					TrophyID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				},
			},
			wantErr: nil,
		},
		{
			name:             "no user trophies :POS",
			fixtureFiles:     []string{},
			wantUserTrophies: nil,
			wantErr:          nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupPostgres(t, tt.fixtureFiles...)
			pgrepo := relationrepo.NewRepo(db)

			gotUserTrophies, gotErr := pgrepo.UserTrophies(context.Background())

			assert.ErrorIs(t, tt.wantErr, gotErr, "expect error to match")
			assert.Equal(t, tt.wantUserTrophies, gotUserTrophies, "expect user trophies to match")
		})
	}
}

func TestRepo_LinkUserTrophy(t *testing.T) {
	t.Run("duplicate user id :NEG", func(t *testing.T) {
		db := setupPostgres(t, "users.yml", "trophies.yml", "user_trophies.yml")
		pgrepo := relationrepo.NewRepo(db)

		gotUserTrophy, gotErr := pgrepo.LinkUserTrophy(
			context.Background(),
			uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			uuid.MustParse("00000000-0000-0000-0000-000000000002"),
		)

		assert.ErrorIs(t, relationrepo.ErrDuplicateID, gotErr, "expect error to match")
		assert.Equal(t, models.UserTrophy{}, gotUserTrophy)
	})

	t.Run("trophy not found :NEG", func(t *testing.T) {
		db := setupPostgres(t, "users.yml", "trophies.yml", "user_trophies.yml")
		pgrepo := relationrepo.NewRepo(db)

		gotUserTrophy, gotErr := pgrepo.LinkUserTrophy(
			context.Background(),
			uuid.MustParse("00000000-0000-0000-0000-000000000002"),
			uuid.MustParse("00000000-0000-0000-0000-000000000003"),
		)

		assert.ErrorIs(t, relationrepo.ErrParentTableRecordNotFound, gotErr, "expect error to match")
		assert.Equal(t, models.UserTrophy{}, gotUserTrophy)
	})

	t.Run("user not found :NEG", func(t *testing.T) {
		db := setupPostgres(t, "users.yml", "trophies.yml", "user_trophies.yml")
		pgrepo := relationrepo.NewRepo(db)

		gotUserTrophy, gotErr := pgrepo.LinkUserTrophy(
			context.Background(),
			uuid.MustParse("00000000-0000-0000-0000-000000000003"),
			uuid.MustParse("00000000-0000-0000-0000-000000000002"),
		)

		assert.ErrorIs(t, relationrepo.ErrParentTableRecordNotFound, gotErr, "expect error to match")
		assert.Equal(t, models.UserTrophy{}, gotUserTrophy)
	})

	t.Run("link user and trophy :POS", func(t *testing.T) {
		db := setupPostgres(t, "users.yml", "trophies.yml", "user_trophies.yml")
		pgrepo := relationrepo.NewRepo(db)

		gotUserTrophy, gotErr := pgrepo.LinkUserTrophy(
			context.Background(),
			uuid.MustParse("00000000-0000-0000-0000-000000000002"),
			uuid.MustParse("00000000-0000-0000-0000-000000000001"),
		)

		assert.ErrorIs(t, nil, gotErr, "expect error to match")
		assert.Equal(t, models.UserTrophy{
			UserID:   uuid.MustParse("00000000-0000-0000-0000-000000000002"),
			TrophyID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
		}, gotUserTrophy)
	})

	t.Run("on deleting user id child refrences gets deleted :POS", func(t *testing.T) {
		db := setupPostgres(t, "users.yml", "trophies.yml", "user_trophies.yml")
		relationPgrepo := relationrepo.NewRepo(db)
		userPgrepo := userrepo.NewRepo(db)

		gotErr := userPgrepo.DeleteUser(context.Background(), uuid.MustParse("00000000-0000-0000-0000-000000000001"))

		assert.NoError(t, gotErr)

		gotUserTrophies, gotErr := relationPgrepo.UserTrophies(context.Background())

		assert.NoError(t, gotErr)
		assert.Equal(t, []models.UserTrophy(nil), gotUserTrophies, "expect user trophies to match")
	})

	t.Run("on deleting trophy id child refrences gets deleted :POS", func(t *testing.T) {
		db := setupPostgres(t, "users.yml", "trophies.yml", "user_trophies.yml")
		relationPgrepo := relationrepo.NewRepo(db)
		userPgrepo := trophyrepo.NewRepo(db)

		gotErr := userPgrepo.DeleteTrophy(context.Background(), uuid.MustParse("00000000-0000-0000-0000-000000000002"))

		assert.NoError(t, gotErr)

		gotUserTrophies, gotErr := relationPgrepo.UserTrophies(context.Background())

		assert.NoError(t, gotErr)
		assert.Equal(t, []models.UserTrophy(nil), gotUserTrophies, "expect user trophies to match")
	})
}

func TestRepo_VoxsphereMemebers(t *testing.T) {
	tests := []struct {
		name                 string
		fixtureFiles         []string
		wantVoxsphereMembers []models.VoxsphereMember
		wantErr              error
	}{
		{
			name:         "voxsphere members :POS",
			fixtureFiles: []string{"topics.yml", "voxspheres.yml", "users.yml", "voxsphere_members.yml"},
			wantVoxsphereMembers: []models.VoxsphereMember{
				{
					VoxsphereID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					UserID:      uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				},
			},
			wantErr: nil,
		},
		{
			name:                 "no voxsphere members :POS",
			fixtureFiles:         []string{},
			wantVoxsphereMembers: nil,
			wantErr:              nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupPostgres(t, tt.fixtureFiles...)
			pgrepo := relationrepo.NewRepo(db)

			gotVoxsphereMembers, gotErr := pgrepo.VoxsphereMemebers(context.Background())

			assert.ErrorIs(t, tt.wantErr, gotErr, "expect error to match")
			assert.Equal(t, tt.wantVoxsphereMembers, gotVoxsphereMembers, "expect voxspheres to match")
		})
	}
}
