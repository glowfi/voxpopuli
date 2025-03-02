package relation_test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/glowfi/voxpopuli/backend/pkg/models"
	awardrepo "github.com/glowfi/voxpopuli/backend/pkg/repo/award"
	customemojirepo "github.com/glowfi/voxpopuli/backend/pkg/repo/custom_emoji"
	emojirepo "github.com/glowfi/voxpopuli/backend/pkg/repo/emoji"
	postrepo "github.com/glowfi/voxpopuli/backend/pkg/repo/post"
	postFlairrepo "github.com/glowfi/voxpopuli/backend/pkg/repo/post_flair"
	relationrepo "github.com/glowfi/voxpopuli/backend/pkg/repo/relation"
	trophyrepo "github.com/glowfi/voxpopuli/backend/pkg/repo/trophy"
	userrepo "github.com/glowfi/voxpopuli/backend/pkg/repo/user"
	userFlairrepo "github.com/glowfi/voxpopuli/backend/pkg/repo/user_flair"
	voxsphererepo "github.com/glowfi/voxpopuli/backend/pkg/repo/voxsphere"
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
	db.RegisterModel((*models.Award)(nil))
	db.RegisterModel((*models.User)(nil))
	db.RegisterModel((*models.Trophy)(nil))
	db.RegisterModel((*models.Post)(nil))
	db.RegisterModel((*models.Emoji)(nil))
	db.RegisterModel((*models.CustomEmoji)(nil))
	db.RegisterModel((*models.UserTrophy)(nil))
	db.RegisterModel((*models.VoxsphereMember)(nil))
	db.RegisterModel((*models.VoxsphereModerator)(nil))
	db.RegisterModel((*models.UserFlair)(nil))
	db.RegisterModel((*models.UserFlairCustomEmoji)(nil))
	db.RegisterModel((*models.UserFlairEmoji)(nil))
	db.RegisterModel((*models.UserFlairDescription)(nil))
	db.RegisterModel((*models.PostFlair)(nil))
	db.RegisterModel((*models.PostFlairCustomEmoji)(nil))
	db.RegisterModel((*models.PostFlairEmoji)(nil))
	db.RegisterModel((*models.PostFlairDescription)(nil))
	db.RegisterModel((*models.PostAward)(nil))

	// drop all rows of the user,trophies table
	if _, err := db.NewTruncateTable().Cascade().Model((*models.Topic)(nil)).Exec(context.Background()); err != nil {
		t.Fatal("truncate table failed:", err)
	}
	if _, err := db.NewTruncateTable().Cascade().Model((*models.Voxsphere)(nil)).Exec(context.Background()); err != nil {
		t.Fatal("truncate table failed:", err)
	}
	if _, err := db.NewTruncateTable().Cascade().Model((*models.Award)(nil)).Exec(context.Background()); err != nil {
		t.Fatal("truncate table failed:", err)
	}
	if _, err := db.NewTruncateTable().Cascade().Model((*models.User)(nil)).Exec(context.Background()); err != nil {
		t.Fatal("truncate table failed:", err)
	}
	if _, err := db.NewTruncateTable().Cascade().Model((*models.Trophy)(nil)).Exec(context.Background()); err != nil {
		t.Fatal("truncate table failed:", err)
	}
	if _, err := db.NewTruncateTable().Cascade().Model((*models.Emoji)(nil)).Exec(context.Background()); err != nil {
		t.Fatal("truncate table failed:", err)
	}
	if _, err := db.NewTruncateTable().Cascade().Model((*models.CustomEmoji)(nil)).Exec(context.Background()); err != nil {
		t.Fatal("truncate table failed:", err)
	}
	if _, err := db.NewTruncateTable().Cascade().Model((*models.UserTrophy)(nil)).Exec(context.Background()); err != nil {
		t.Fatal("truncate table failed:", err)
	}
	if _, err := db.NewTruncateTable().Cascade().Model((*models.VoxsphereMember)(nil)).Exec(context.Background()); err != nil {
		t.Fatal("truncate table failed:", err)
	}
	if _, err := db.NewTruncateTable().Cascade().Model((*models.VoxsphereModerator)(nil)).Exec(context.Background()); err != nil {
		t.Fatal("truncate table failed:", err)
	}
	if _, err := db.NewTruncateTable().Cascade().Model((*models.UserFlair)(nil)).Exec(context.Background()); err != nil {
		t.Fatal("truncate table failed:", err)
	}
	if _, err := db.NewTruncateTable().Cascade().Model((*models.UserFlairCustomEmoji)(nil)).Exec(context.Background()); err != nil {
		t.Fatal("truncate table failed:", err)
	}
	if _, err := db.NewTruncateTable().Cascade().Model((*models.UserFlairEmoji)(nil)).Exec(context.Background()); err != nil {
		t.Fatal("truncate table failed:", err)
	}
	if _, err := db.NewTruncateTable().Cascade().Model((*models.UserFlairDescription)(nil)).Exec(context.Background()); err != nil {
		t.Fatal("truncate table failed:", err)
	}
	if _, err := db.NewTruncateTable().Cascade().Model((*models.PostFlair)(nil)).Exec(context.Background()); err != nil {
		t.Fatal("truncate table failed:", err)
	}
	if _, err := db.NewTruncateTable().Cascade().Model((*models.PostFlairCustomEmoji)(nil)).Exec(context.Background()); err != nil {
		t.Fatal("truncate table failed:", err)
	}
	if _, err := db.NewTruncateTable().Cascade().Model((*models.PostFlairEmoji)(nil)).Exec(context.Background()); err != nil {
		t.Fatal("truncate table failed:", err)
	}
	if _, err := db.NewTruncateTable().Cascade().Model((*models.PostFlairDescription)(nil)).Exec(context.Background()); err != nil {
		t.Fatal("truncate table failed:", err)
	}
	if _, err := db.NewTruncateTable().Cascade().Model((*models.PostAward)(nil)).Exec(context.Background()); err != nil {
		t.Fatal("truncate table failed:", err)
	}

	// load fixture
	fixture := dbfixture.New(db)
	if err := fixture.Load(context.Background(), os.DirFS("testdata"), fixtureFiles...); err != nil {
		t.Fatal("failed to load fixtures", err)
	}

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

func TestRepo_LinkUserTrophies(t *testing.T) {
	t.Run("duplicate user_id,trophy_id while linking user and trophy :NEG", func(t *testing.T) {
		db := setupPostgres(t, "users.yml", "trophies.yml", "user_trophies.yml")
		pgrepo := relationrepo.NewRepo(db)

		gotUserTrophies, gotErr := pgrepo.LinkUserTrophies(
			context.Background(),
			models.UserTrophy{
				UserID:   uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				TrophyID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
			},
		)

		assert.ErrorIs(t, relationrepo.ErrDuplicateID, gotErr, "expect error to match")
		assert.Equal(t, []models.UserTrophy(nil), gotUserTrophies, "expect user trophies to match")
	})

	t.Run("trophy not found while linking user and trophy :NEG", func(t *testing.T) {
		db := setupPostgres(t, "users.yml", "trophies.yml", "user_trophies.yml")
		pgrepo := relationrepo.NewRepo(db)

		gotUserTrophies, gotErr := pgrepo.LinkUserTrophies(
			context.Background(),
			models.UserTrophy{
				UserID:   uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				TrophyID: uuid.MustParse("00000000-0000-0000-0000-000000000003"),
			},
		)

		assert.ErrorIs(t, relationrepo.ErrParentTableRecordNotFound, gotErr, "expect error to match")
		assert.Equal(t, []models.UserTrophy(nil), gotUserTrophies, "expect user trophies to match")
	})

	t.Run("user not found while linking user and trophy :NEG", func(t *testing.T) {
		db := setupPostgres(t, "users.yml", "trophies.yml", "user_trophies.yml")
		pgrepo := relationrepo.NewRepo(db)

		gotUserTrophies, gotErr := pgrepo.LinkUserTrophies(
			context.Background(),
			models.UserTrophy{
				UserID:   uuid.MustParse("00000000-0000-0000-0000-000000000009"),
				TrophyID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
			},
		)

		assert.ErrorIs(t, relationrepo.ErrParentTableRecordNotFound, gotErr, "expect error to match")
		assert.Equal(t, []models.UserTrophy(nil), gotUserTrophies, "expect user trophies to match")
	})

	t.Run("link user and trophy :POS", func(t *testing.T) {
		db := setupPostgres(t, "users.yml", "trophies.yml", "user_trophies.yml")
		pgrepo := relationrepo.NewRepo(db)

		gotUserTrophies, gotErr := pgrepo.LinkUserTrophies(
			context.Background(),
			models.UserTrophy{
				UserID:   uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				TrophyID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			},
		)

		assert.ErrorIs(t, nil, gotErr, "expect error to match")
		assert.Equal(t, []models.UserTrophy{
			{
				UserID:   uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				TrophyID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			},
		}, gotUserTrophies, "expect user trophies to match")
	})

	t.Run("on deleting user id child refrences gets deleted in user_trophies table :POS", func(t *testing.T) {
		db := setupPostgres(t, "users.yml", "trophies.yml", "user_trophies.yml")
		relationPgrepo := relationrepo.NewRepo(db)
		userPgrepo := userrepo.NewRepo(db)

		gotErr := userPgrepo.DeleteUser(context.Background(), uuid.MustParse("00000000-0000-0000-0000-000000000001"))

		assert.NoError(t, gotErr)

		gotUserTrophies, gotErr := relationPgrepo.UserTrophies(context.Background())

		assert.NoError(t, gotErr)
		assert.Equal(t, []models.UserTrophy(nil), gotUserTrophies, "expect user trophies to match")
	})

	t.Run("on deleting trophy id child refrences gets deleted in user_trophies table :POS", func(t *testing.T) {
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

func TestRepo_VoxsphereMembers(t *testing.T) {
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

			gotVoxsphereMembers, gotErr := pgrepo.VoxsphereMembers(context.Background())

			assert.ErrorIs(t, tt.wantErr, gotErr, "expect error to match")
			assert.Equal(t, tt.wantVoxsphereMembers, gotVoxsphereMembers, "expect voxsphere members to match")
		})
	}
}

func TestRepo_LinkVoxsphereMembers(t *testing.T) {
	t.Run("duplicate voxsphere_id,user_id while linking voxsphere and member :NEG", func(t *testing.T) {
		db := setupPostgres(t, "topics.yml", "voxspheres.yml", "users.yml", "voxsphere_members.yml")
		pgrepo := relationrepo.NewRepo(db)

		gotVoxsphereMembers, gotErr := pgrepo.LinkVoxsphereMembers(
			context.Background(),
			models.VoxsphereMember{
				VoxsphereID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				UserID:      uuid.MustParse("00000000-0000-0000-0000-000000000002"),
			},
		)

		assert.ErrorIs(t, relationrepo.ErrDuplicateID, gotErr, "expect error to match")
		assert.Equal(t, []models.VoxsphereMember(nil), gotVoxsphereMembers, "expect voxsphere members to match")
	})

	t.Run("voxsphere not found while linking voxsphere and member :NEG", func(t *testing.T) {
		db := setupPostgres(t, "topics.yml", "voxspheres.yml", "users.yml", "voxsphere_members.yml")
		pgrepo := relationrepo.NewRepo(db)

		gotVoxsphereMembers, gotErr := pgrepo.LinkVoxsphereMembers(
			context.Background(),
			models.VoxsphereMember{
				VoxsphereID: uuid.MustParse("00000000-0000-0000-0000-000000000003"),
				UserID:      uuid.MustParse("00000000-0000-0000-0000-000000000002"),
			},
		)

		assert.ErrorIs(t, relationrepo.ErrParentTableRecordNotFound, gotErr, "expect error to match")
		assert.Equal(t, []models.VoxsphereMember(nil), gotVoxsphereMembers, "expect voxsphere members to match")
	})

	t.Run("member not found while linking voxsphere and member :NEG", func(t *testing.T) {
		db := setupPostgres(t, "topics.yml", "voxspheres.yml", "users.yml", "voxsphere_members.yml")
		pgrepo := relationrepo.NewRepo(db)

		gotVoxsphereMember, gotErr := pgrepo.LinkVoxsphereMembers(
			context.Background(),
			models.VoxsphereMember{
				VoxsphereID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				UserID:      uuid.MustParse("00000000-0000-0000-0000-000000000009"),
			},
		)

		assert.ErrorIs(t, relationrepo.ErrParentTableRecordNotFound, gotErr, "expect error to match")
		assert.Equal(t, []models.VoxsphereMember(nil), gotVoxsphereMember, "expect voxsphere members to match")
	})

	t.Run("link voxsphere and member :POS", func(t *testing.T) {
		db := setupPostgres(t, "topics.yml", "voxspheres.yml", "users.yml", "voxsphere_members.yml")
		pgrepo := relationrepo.NewRepo(db)

		gotVoxsphereMembers, gotErr := pgrepo.LinkVoxsphereMembers(
			context.Background(),
			models.VoxsphereMember{
				VoxsphereID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				UserID:      uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			},
		)

		assert.ErrorIs(t, nil, gotErr, "expect error to match")
		assert.Equal(t, []models.VoxsphereMember{
			{
				VoxsphereID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				UserID:      uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			},
		}, gotVoxsphereMembers, "expect voxsphere members to match")
	})

	t.Run("on deleting voxsphere id child refrences gets deleted in voxsphere_members table :POS", func(t *testing.T) {
		db := setupPostgres(t, "topics.yml", "voxspheres.yml", "users.yml", "voxsphere_members.yml")
		voxspherePgrepo := voxsphererepo.NewRepo(db)
		relationPgrepo := relationrepo.NewRepo(db)

		gotErr := voxspherePgrepo.DeleteVoxsphere(context.Background(), uuid.MustParse("00000000-0000-0000-0000-000000000001"))

		assert.NoError(t, gotErr)

		gotVoxsphereMembers, gotErr := relationPgrepo.VoxsphereMembers(context.Background())

		assert.NoError(t, gotErr)
		assert.Equal(t, []models.VoxsphereMember(nil), gotVoxsphereMembers, "expect voxsphere members to match")
	})

	t.Run("on deleting member id child refrences gets deleted in voxsphere_members table :POS", func(t *testing.T) {
		db := setupPostgres(t, "topics.yml", "voxspheres.yml", "users.yml", "voxsphere_members.yml")
		userPgrepo := userrepo.NewRepo(db)
		relationPgrepo := relationrepo.NewRepo(db)

		gotErr := userPgrepo.DeleteUser(context.Background(), uuid.MustParse("00000000-0000-0000-0000-000000000002"))

		assert.NoError(t, gotErr)

		gotVoxsphereMembers, gotErr := relationPgrepo.VoxsphereMembers(context.Background())

		assert.NoError(t, gotErr)
		assert.Equal(t, []models.VoxsphereMember(nil), gotVoxsphereMembers, "expect voxsphere members to match")
	})
}

func TestRepo_VoxsphereModerators(t *testing.T) {
	tests := []struct {
		name                    string
		fixtureFiles            []string
		wantVoxsphereModerators []models.VoxsphereModerator
		wantErr                 error
	}{
		{
			name:         "voxsphere moderators :POS",
			fixtureFiles: []string{"topics.yml", "voxspheres.yml", "users.yml", "voxsphere_moderators.yml"},
			wantVoxsphereModerators: []models.VoxsphereModerator{
				{
					VoxsphereID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					UserID:      uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				},
			},
			wantErr: nil,
		},
		{
			name:                    "no voxsphere moderators :POS",
			fixtureFiles:            []string{},
			wantVoxsphereModerators: nil,
			wantErr:                 nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupPostgres(t, tt.fixtureFiles...)
			pgrepo := relationrepo.NewRepo(db)

			gotVoxsphereModerators, gotErr := pgrepo.VoxsphereModerators(context.Background())

			assert.ErrorIs(t, tt.wantErr, gotErr, "expect error to match")
			assert.Equal(t, tt.wantVoxsphereModerators, gotVoxsphereModerators, "expect voxsphere moderators to match")
		})
	}
}

func TestRepo_LinkVoxsphereModerators(t *testing.T) {
	t.Run("duplicate voxsphere_id,user_id while linking voxsphere and moderators :NEG", func(t *testing.T) {
		db := setupPostgres(t, "topics.yml", "voxspheres.yml", "users.yml", "voxsphere_moderators.yml")
		pgrepo := relationrepo.NewRepo(db)

		gotVoxsphereModerators, gotErr := pgrepo.LinkVoxsphereModerators(
			context.Background(),
			models.VoxsphereModerator{
				VoxsphereID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				UserID:      uuid.MustParse("00000000-0000-0000-0000-000000000002"),
			},
		)

		assert.ErrorIs(t, relationrepo.ErrDuplicateID, gotErr, "expect error to match")
		assert.Equal(t, []models.VoxsphereModerator(nil), gotVoxsphereModerators, "expect voxsphere moderators to match")
	})

	t.Run("voxsphere not found while linking voxsphere and moderator :NEG", func(t *testing.T) {
		db := setupPostgres(t, "topics.yml", "voxspheres.yml", "users.yml", "voxsphere_moderators.yml")
		pgrepo := relationrepo.NewRepo(db)

		gotVoxsphereModerators, gotErr := pgrepo.LinkVoxsphereModerators(
			context.Background(),
			models.VoxsphereModerator{
				VoxsphereID: uuid.MustParse("00000000-0000-0000-0000-000000000003"),
				UserID:      uuid.MustParse("00000000-0000-0000-0000-000000000002"),
			},
		)

		assert.ErrorIs(t, relationrepo.ErrParentTableRecordNotFound, gotErr, "expect error to match")
		assert.Equal(t, []models.VoxsphereModerator(nil), gotVoxsphereModerators, "expect voxsphere moderators to match")
	})

	t.Run("moderator not found while linking voxsphere and moderator :NEG", func(t *testing.T) {
		db := setupPostgres(t, "topics.yml", "voxspheres.yml", "users.yml", "voxsphere_moderators.yml")
		pgrepo := relationrepo.NewRepo(db)

		gotVoxsphereModerators, gotErr := pgrepo.LinkVoxsphereModerators(
			context.Background(),
			models.VoxsphereModerator{
				VoxsphereID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				UserID:      uuid.MustParse("00000000-0000-0000-0000-000000000009"),
			},
		)

		assert.ErrorIs(t, relationrepo.ErrParentTableRecordNotFound, gotErr, "expect error to match")
		assert.Equal(t, []models.VoxsphereModerator(nil), gotVoxsphereModerators, "expect voxsphere moderators to match")
	})

	t.Run("link voxsphere and moderator :POS", func(t *testing.T) {
		db := setupPostgres(t, "topics.yml", "voxspheres.yml", "users.yml", "voxsphere_moderators.yml")
		pgrepo := relationrepo.NewRepo(db)

		gotVoxsphereModerators, gotErr := pgrepo.LinkVoxsphereModerators(
			context.Background(),
			models.VoxsphereModerator{
				VoxsphereID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				UserID:      uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			},
		)

		assert.ErrorIs(t, nil, gotErr, "expect error to match")
		assert.Equal(t, []models.VoxsphereModerator{
			{
				VoxsphereID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				UserID:      uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			},
		}, gotVoxsphereModerators, "expect voxsphere moderators to match")
	})

	t.Run("on deleting voxsphere id child refrences gets deleted in voxsphere_moderators table :POS", func(t *testing.T) {
		db := setupPostgres(t, "topics.yml", "voxspheres.yml", "users.yml", "voxsphere_moderators.yml")
		voxspherePgrepo := voxsphererepo.NewRepo(db)
		relationPgrepo := relationrepo.NewRepo(db)

		gotErr := voxspherePgrepo.DeleteVoxsphere(context.Background(), uuid.MustParse("00000000-0000-0000-0000-000000000001"))

		assert.NoError(t, gotErr)

		gotVoxsphereModerators, gotErr := relationPgrepo.VoxsphereModerators(context.Background())

		assert.NoError(t, gotErr)
		assert.Equal(t, []models.VoxsphereModerator(nil), gotVoxsphereModerators, "expect voxsphere moderators to match")
	})

	t.Run("on deleting moderator id child refrences gets deleted in voxsphere_moderators table :POS", func(t *testing.T) {
		db := setupPostgres(t, "topics.yml", "voxspheres.yml", "users.yml", "voxsphere_moderators.yml")
		userPgrepo := userrepo.NewRepo(db)
		relationPgrepo := relationrepo.NewRepo(db)

		gotErr := userPgrepo.DeleteUser(context.Background(), uuid.MustParse("00000000-0000-0000-0000-000000000002"))

		assert.NoError(t, gotErr)

		gotVoxsphereModerators, gotErr := relationPgrepo.VoxsphereModerators(context.Background())

		assert.NoError(t, gotErr)
		assert.Equal(t, []models.VoxsphereModerator(nil), gotVoxsphereModerators, "expect voxsphere moderators to match")
	})
}

func TestRepo_UserFlairEmojis(t *testing.T) {
	tests := []struct {
		name                string
		fixtureFiles        []string
		wantUserFlairEmojis []models.UserFlairEmoji
		wantErr             error
	}{
		{
			name:         "user flair emojis :POS",
			fixtureFiles: []string{"topics.yml", "voxspheres.yml", "emojis.yml", "users.yml", "user_flairs.yml", "user_flair_emojis.yml"},
			wantUserFlairEmojis: []models.UserFlairEmoji{
				{
					EmojiID:     uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					UserFlairID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					OrderIndex:  1,
				},
				{
					EmojiID:     uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					UserFlairID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					OrderIndex:  2,
				},
				{
					EmojiID:     uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					UserFlairID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					OrderIndex:  1,
				},
				{
					EmojiID:     uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					UserFlairID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					OrderIndex:  2,
				},
			},
			wantErr: nil,
		},
		{
			name:                "no user flair emojis :POS",
			fixtureFiles:        []string{},
			wantUserFlairEmojis: nil,
			wantErr:             nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupPostgres(t, tt.fixtureFiles...)
			pgrepo := relationrepo.NewRepo(db)

			gotUserFlairEmojis, gotErr := pgrepo.UserFlairEmojis(context.Background())

			assert.ErrorIs(t, tt.wantErr, gotErr, "expect error to match")
			assert.Equal(t, tt.wantUserFlairEmojis, gotUserFlairEmojis, "expect user flair emojis to match")
		})
	}
}

func TestRepo_UserFlairCustomEmojis(t *testing.T) {
	tests := []struct {
		name                      string
		fixtureFiles              []string
		wantUserFlairCustomEmojis []models.UserFlairCustomEmoji
		wantErr                   error
	}{
		{
			name: "user flair custom emojis :POS",
			fixtureFiles: []string{
				"topics.yml",
				"voxspheres.yml",
				"users.yml",
				"custom_emojis.yml",
				"user_flairs.yml",
				"user_flair_custom_emojis.yml",
			},
			wantUserFlairCustomEmojis: []models.UserFlairCustomEmoji{
				{
					CustomEmojiID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					UserFlairID:   uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					OrderIndex:    4,
				},
				{
					CustomEmojiID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					UserFlairID:   uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					OrderIndex:    5,
				},
				{
					CustomEmojiID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					UserFlairID:   uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					OrderIndex:    4,
				},
				{
					CustomEmojiID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					UserFlairID:   uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					OrderIndex:    5,
				},
			},
			wantErr: nil,
		},
		{
			name:                      "no user flair custom emojis :POS",
			fixtureFiles:              []string{},
			wantUserFlairCustomEmojis: nil,
			wantErr:                   nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupPostgres(t, tt.fixtureFiles...)
			pgrepo := relationrepo.NewRepo(db)

			gotUserFlairCustomEmojis, gotErr := pgrepo.UserFlairCustomEmojis(context.Background())

			assert.ErrorIs(t, tt.wantErr, gotErr, "expect error to match")
			assert.Equal(t, tt.wantUserFlairCustomEmojis, gotUserFlairCustomEmojis, "expect user flair custom emojis to match")
		})
	}
}

func TestRepo_UserFlairDescriptions(t *testing.T) {
	tests := []struct {
		name                      string
		fixtureFiles              []string
		wantUserFlairDescriptions []models.UserFlairDescription
		wantErr                   error
	}{
		{
			name: "user flair descriptions :POS",
			fixtureFiles: []string{
				"topics.yml",
				"voxspheres.yml",
				"users.yml",
				"custom_emojis.yml",
				"user_flairs.yml",
				"user_flair_descriptions.yml",
			},
			wantUserFlairDescriptions: []models.UserFlairDescription{
				{
					UserFlairID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					OrderIndex:  0,
					Description: "desc1 ",
				},
				{
					UserFlairID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					OrderIndex:  3,
					Description: " ",
				},
				{
					UserFlairID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					OrderIndex:  0,
					Description: "desc2 ",
				},
				{
					UserFlairID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					OrderIndex:  3,
					Description: " ",
				},
				{
					UserFlairID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					OrderIndex:  6,
					Description: " desc2",
				},
			},
			wantErr: nil,
		},
		{
			name:                      "no user flair descriptions :POS",
			fixtureFiles:              []string{},
			wantUserFlairDescriptions: nil,
			wantErr:                   nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupPostgres(t, tt.fixtureFiles...)
			pgrepo := relationrepo.NewRepo(db)

			gotUserFlairDescriptions, gotErr := pgrepo.UserFlairDescriptions(context.Background())

			assert.ErrorIs(t, tt.wantErr, gotErr, "expect error to match")
			assert.Equal(t, tt.wantUserFlairDescriptions, gotUserFlairDescriptions, "expect user flair descriptions to match")
		})
	}
}

func TestRepo_LinkUserFlairEmojis(t *testing.T) {
	t.Run("duplicate user_flair_id,emoji_id,order_index while linking user flair and emoji :NEG", func(t *testing.T) {
		db := setupPostgres(t,
			"topics.yml",
			"voxspheres.yml",
			"users.yml",
			"emojis.yml",
			"user_flairs.yml",
			"user_flair_emojis.yml",
		)
		pgrepo := relationrepo.NewRepo(db)

		gotUserFlairEmojis, gotErr := pgrepo.LinkUserFlairEmojis(
			context.Background(),
			models.UserFlairEmoji{
				EmojiID:     uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				UserFlairID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				OrderIndex:  1,
			},
		)

		assert.ErrorIs(t, relationrepo.ErrDuplicateID, gotErr, "expect error to match")
		assert.Equal(t, []models.UserFlairEmoji(nil), gotUserFlairEmojis, "expect user flair emojis to match")
	})

	t.Run("emoji not found in parent table while linking user flair and custom emoji :NEG", func(t *testing.T) {
		db := setupPostgres(t,
			"topics.yml",
			"voxspheres.yml",
			"users.yml",
			"emojis.yml",
			"user_flairs.yml",
			"user_flair_emojis.yml",
		)
		pgrepo := relationrepo.NewRepo(db)

		gotUserFlairEmojis, gotErr := pgrepo.LinkUserFlairEmojis(
			context.Background(),
			models.UserFlairEmoji{
				EmojiID:     uuid.MustParse("00000000-0000-0000-0000-000000000003"),
				UserFlairID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				OrderIndex:  1,
			},
		)

		assert.ErrorIs(t, relationrepo.ErrParentTableRecordNotFound, gotErr, "expect error to match")
		assert.Equal(t, []models.UserFlairEmoji(nil), gotUserFlairEmojis, "expect user flair emojis to match")
	})

	t.Run("user flair not found in parent table while linking user flair and custom emoji :NEG", func(t *testing.T) {
		db := setupPostgres(t,
			"topics.yml",
			"voxspheres.yml",
			"users.yml",
			"emojis.yml",
			"user_flairs.yml",
			"user_flair_emojis.yml",
		)
		pgrepo := relationrepo.NewRepo(db)

		gotUserFlairEmojis, gotErr := pgrepo.LinkUserFlairEmojis(
			context.Background(),
			models.UserFlairEmoji{
				EmojiID:     uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				UserFlairID: uuid.MustParse("00000000-0000-0000-0000-000000000009"),
				OrderIndex:  1,
			},
		)

		assert.ErrorIs(t, relationrepo.ErrParentTableRecordNotFound, gotErr, "expect error to match")
		assert.Equal(t, []models.UserFlairEmoji(nil), gotUserFlairEmojis, "expect user flair emojis to match")
	})

	t.Run("link user flair and emoji :POS", func(t *testing.T) {
		db := setupPostgres(t,
			"topics.yml",
			"voxspheres.yml",
			"users.yml",
			"emojis.yml",
			"user_flairs.yml",
			"user_flair_emojis.yml",
		)
		pgrepo := relationrepo.NewRepo(db)

		gotUserFlairEmojis, gotErr := pgrepo.LinkUserFlairEmojis(
			context.Background(),
			models.UserFlairEmoji{
				EmojiID:     uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				UserFlairID: uuid.MustParse("00000000-0000-0000-0000-000000000003"),
				OrderIndex:  0,
			},
		)

		assert.ErrorIs(t, nil, gotErr, "expect error to match")
		assert.Equal(t, []models.UserFlairEmoji{
			{
				EmojiID:     uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				UserFlairID: uuid.MustParse("00000000-0000-0000-0000-000000000003"),
				OrderIndex:  0,
			},
		}, gotUserFlairEmojis, "expect user flair emojis to match")
	})

	t.Run("on deleting user flair child refrences gets deleted in user_flair_emojis table :POS", func(t *testing.T) {
		db := setupPostgres(t,
			"topics.yml",
			"voxspheres.yml",
			"users.yml",
			"emojis.yml",
			"user_flairs.yml",
			"user_flair_emojis.yml",
		)
		relationPgrepo := relationrepo.NewRepo(db)
		userFlairPgrepo := userFlairrepo.NewRepo(db)

		gotErr := userFlairPgrepo.DeleteUserFlair(context.Background(), uuid.MustParse("00000000-0000-0000-0000-000000000001"))
		assert.NoError(t, gotErr)

		gotUserFlairEmojis, gotErr := relationPgrepo.UserFlairEmojis(context.Background())

		assert.NoError(t, gotErr)
		assert.Equal(t, []models.UserFlairEmoji{
			{
				EmojiID:     uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				UserFlairID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				OrderIndex:  1,
			},
			{
				EmojiID:     uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				UserFlairID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				OrderIndex:  2,
			},
		}, gotUserFlairEmojis, "expect user flair emojis to match")
	})

	t.Run("on deleting emoji child refrences gets deleted in user_flair_emojis table :POS", func(t *testing.T) {
		db := setupPostgres(t,
			"topics.yml",
			"voxspheres.yml",
			"users.yml",
			"emojis.yml",
			"user_flairs.yml",
			"user_flair_emojis.yml",
		)
		relationPgrepo := relationrepo.NewRepo(db)
		emojiPgrepo := emojirepo.NewRepo(db)

		gotErr := emojiPgrepo.DeleteEmoji(context.Background(), uuid.MustParse("00000000-0000-0000-0000-000000000001"))

		assert.NoError(t, gotErr)

		gotUserFlairEmojis, gotErr := relationPgrepo.UserFlairEmojis(context.Background())

		assert.NoError(t, gotErr)
		assert.Equal(t, []models.UserFlairEmoji{
			{
				EmojiID:     uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				UserFlairID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				OrderIndex:  1,
			},
			{
				EmojiID:     uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				UserFlairID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				OrderIndex:  2,
			},
		}, gotUserFlairEmojis, "expect user flair emojis to match")
	})
}

func TestRepo_LinkUserFlairCustomEmojis(t *testing.T) {
	t.Run("duplicate user_flair_id,emoji_id,order_index while linking user flair and custom emoji :NEG", func(t *testing.T) {
		db := setupPostgres(t,
			"topics.yml",
			"voxspheres.yml",
			"users.yml",
			"custom_emojis.yml",
			"user_flairs.yml",
			"user_flair_custom_emojis.yml",
		)
		pgrepo := relationrepo.NewRepo(db)

		gotUserFlairCustomEmojis, gotErr := pgrepo.LinkUserFlairCustomEmojis(
			context.Background(),
			models.UserFlairCustomEmoji{
				CustomEmojiID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				UserFlairID:   uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				OrderIndex:    4,
			},
		)

		assert.ErrorIs(t, relationrepo.ErrDuplicateID, gotErr, "expect error to match")
		assert.Equal(t, []models.UserFlairCustomEmoji(nil), gotUserFlairCustomEmojis, "expect user flair custom emojis to match")
	})

	t.Run("custom emoji not found in parent table while linking user flair and custom emoji :NEG", func(t *testing.T) {
		db := setupPostgres(t,
			"topics.yml",
			"voxspheres.yml",
			"users.yml",
			"custom_emojis.yml",
			"user_flairs.yml",
			"user_flair_custom_emojis.yml",
		)
		pgrepo := relationrepo.NewRepo(db)

		gotUserFlairCustomEmojis, gotErr := pgrepo.LinkUserFlairCustomEmojis(
			context.Background(),
			models.UserFlairCustomEmoji{
				CustomEmojiID: uuid.MustParse("00000000-0000-0000-0000-000000000003"),
				UserFlairID:   uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				OrderIndex:    1,
			},
		)

		assert.ErrorIs(t, relationrepo.ErrParentTableRecordNotFound, gotErr, "expect error to match")
		assert.Equal(t, []models.UserFlairCustomEmoji(nil), gotUserFlairCustomEmojis, "expect user flair custom emojis to match")
	})

	t.Run("user flair not found in parent table while linking user flair and custom emoji :NEG", func(t *testing.T) {
		db := setupPostgres(t,
			"topics.yml",
			"voxspheres.yml",
			"users.yml",
			"custom_emojis.yml",
			"user_flairs.yml",
			"user_flair_custom_emojis.yml",
		)
		pgrepo := relationrepo.NewRepo(db)

		gotUserFlairCustomEmoji, gotErr := pgrepo.LinkUserFlairCustomEmojis(
			context.Background(),
			models.UserFlairCustomEmoji{
				CustomEmojiID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				UserFlairID:   uuid.MustParse("00000000-0000-0000-0000-000000000009"),
				OrderIndex:    1,
			},
		)

		assert.ErrorIs(t, relationrepo.ErrParentTableRecordNotFound, gotErr, "expect error to match")
		assert.Equal(t, []models.UserFlairCustomEmoji(nil), gotUserFlairCustomEmoji, "expect user flair custom emojis to match")
	})

	t.Run("link user flair and custom emoji :POS", func(t *testing.T) {
		db := setupPostgres(t,
			"topics.yml",
			"voxspheres.yml",
			"users.yml",
			"custom_emojis.yml",
			"user_flairs.yml",
			"user_flair_custom_emojis.yml",
		)
		pgrepo := relationrepo.NewRepo(db)

		gotUserFlairCustomEmoji, gotErr := pgrepo.LinkUserFlairCustomEmojis(
			context.Background(),
			models.UserFlairCustomEmoji{
				CustomEmojiID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				UserFlairID:   uuid.MustParse("00000000-0000-0000-0000-000000000003"),
				OrderIndex:    1,
			},
		)

		assert.ErrorIs(t, nil, gotErr, "expect error to match")
		assert.Equal(t, []models.UserFlairCustomEmoji{
			{
				CustomEmojiID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				UserFlairID:   uuid.MustParse("00000000-0000-0000-0000-000000000003"),
				OrderIndex:    1,
			},
		}, gotUserFlairCustomEmoji, "expect user flair custom emojis to match")
	})

	t.Run("on deleting user flair child refrences gets deleted in user_flair_custom_emojis table :POS", func(t *testing.T) {
		db := setupPostgres(t,
			"topics.yml",
			"voxspheres.yml",
			"users.yml",
			"custom_emojis.yml",
			"user_flairs.yml",
			"user_flair_custom_emojis.yml",
		)
		relationPgrepo := relationrepo.NewRepo(db)
		userFlairPgrepo := userFlairrepo.NewRepo(db)

		gotErr := userFlairPgrepo.DeleteUserFlair(context.Background(), uuid.MustParse("00000000-0000-0000-0000-000000000001"))
		assert.NoError(t, gotErr)

		gotUserFlairCustomEmojis, gotErr := relationPgrepo.UserFlairCustomEmojis(context.Background())

		assert.NoError(t, gotErr)
		assert.Equal(t, []models.UserFlairCustomEmoji{
			{
				CustomEmojiID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				UserFlairID:   uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				OrderIndex:    4,
			},
			{
				CustomEmojiID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				UserFlairID:   uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				OrderIndex:    5,
			},
		}, gotUserFlairCustomEmojis, "expect user flair custom emojis to match")
	})

	t.Run("on deleting custom emoji child refrences gets deleted in user_flair_custom_emojis table :POS", func(t *testing.T) {
		db := setupPostgres(t,
			"topics.yml",
			"voxspheres.yml",
			"users.yml",
			"custom_emojis.yml",
			"user_flairs.yml",
			"user_flair_custom_emojis.yml",
		)
		relationPgrepo := relationrepo.NewRepo(db)
		customemojiPgrepo := customemojirepo.NewRepo(db)

		gotErr := customemojiPgrepo.DeleteCustomEmoji(context.Background(), uuid.MustParse("00000000-0000-0000-0000-000000000001"))

		assert.NoError(t, gotErr)

		gotUserFlairCustomEmojis, gotErr := relationPgrepo.UserFlairCustomEmojis(context.Background())

		assert.NoError(t, gotErr)
		assert.Equal(t, []models.UserFlairCustomEmoji{
			{
				CustomEmojiID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				UserFlairID:   uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				OrderIndex:    4,
			},
			{
				CustomEmojiID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				UserFlairID:   uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				OrderIndex:    5,
			},
		}, gotUserFlairCustomEmojis, "expect user flair custom emojis to match")
	})
}

func TestRepo_LinkUserFlairDescriptions(t *testing.T) {
	t.Run("duplicate user_flair_id,order_index while linking user flair and description :NEG", func(t *testing.T) {
		db := setupPostgres(t,
			"topics.yml",
			"voxspheres.yml",
			"users.yml",
			"user_flairs.yml",
			"user_flair_descriptions.yml",
		)
		pgrepo := relationrepo.NewRepo(db)

		gotUserFlairDescriptions, gotErr := pgrepo.LinkUserFlairDescriptions(
			context.Background(),
			models.UserFlairDescription{
				UserFlairID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				OrderIndex:  0,
				Description: "text",
			},
		)

		assert.ErrorIs(t, relationrepo.ErrDuplicateID, gotErr, "expect error to match")
		assert.Equal(t, []models.UserFlairDescription(nil), gotUserFlairDescriptions, "expect user flair descriptions to match")
	})

	t.Run("user flair not found in parent table while linking user flair and description :NEG", func(t *testing.T) {
		db := setupPostgres(t,
			"topics.yml",
			"voxspheres.yml",
			"users.yml",
			"user_flairs.yml",
			"user_flair_descriptions.yml",
		)
		pgrepo := relationrepo.NewRepo(db)

		gotUserFlairDescriptions, gotErr := pgrepo.LinkUserFlairDescriptions(
			context.Background(),
			models.UserFlairDescription{
				UserFlairID: uuid.MustParse("00000000-0000-0000-0000-000000000009"),
				OrderIndex:  3,
				Description: "desc3",
			},
		)

		assert.ErrorIs(t, relationrepo.ErrParentTableRecordNotFound, gotErr, "expect error to match")
		assert.Equal(t, []models.UserFlairDescription(nil), gotUserFlairDescriptions, "expect user flair descriptions to match")
	})

	t.Run("link user flair and description :POS", func(t *testing.T) {
		db := setupPostgres(t,
			"topics.yml",
			"voxspheres.yml",
			"users.yml",
			"user_flairs.yml",
			"user_flair_descriptions.yml",
		)
		pgrepo := relationrepo.NewRepo(db)

		gotUserFlairDescriptions, gotErr := pgrepo.LinkUserFlairDescriptions(
			context.Background(),
			models.UserFlairDescription{
				UserFlairID: uuid.MustParse("00000000-0000-0000-0000-000000000003"),
				OrderIndex:  2,
				Description: "desc3",
			},
		)

		assert.ErrorIs(t, nil, gotErr, "expect error to match")
		assert.Equal(t, []models.UserFlairDescription{
			{
				UserFlairID: uuid.MustParse("00000000-0000-0000-0000-000000000003"),
				OrderIndex:  2,
				Description: "desc3",
			},
		}, gotUserFlairDescriptions, "expect user flair descriptions to match")
	})

	t.Run("on deleting user flair child refrences gets deleted in user_flair_descriptions table :POS", func(t *testing.T) {
		db := setupPostgres(t,
			"topics.yml",
			"voxspheres.yml",
			"users.yml",
			"user_flairs.yml",
			"user_flair_descriptions.yml",
		)
		relationPgrepo := relationrepo.NewRepo(db)
		userFlairPgrepo := userFlairrepo.NewRepo(db)

		gotErr := userFlairPgrepo.DeleteUserFlair(context.Background(), uuid.MustParse("00000000-0000-0000-0000-000000000001"))
		assert.NoError(t, gotErr)

		gotUserFlairDescriptions, gotErr := relationPgrepo.UserFlairDescriptions(context.Background())

		assert.NoError(t, gotErr)
		assert.Equal(t, []models.UserFlairDescription{
			{
				UserFlairID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				Description: "desc2 ",
				OrderIndex:  0,
			},
			{
				UserFlairID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				Description: " ",
				OrderIndex:  3,
			},
			{
				UserFlairID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				Description: " desc2",
				OrderIndex:  6,
			},
		}, gotUserFlairDescriptions, "expect user flair descriptions to match")
	})
}

func TestRepo_PostFlairEmojis(t *testing.T) {
	tests := []struct {
		name                string
		fixtureFiles        []string
		wantPostFlairEmojis []models.PostFlairEmoji
		wantErr             error
	}{
		{
			name: "post flair emojis :POS",
			fixtureFiles: []string{
				"topics.yml",
				"voxspheres.yml",
				"emojis.yml",
				"users.yml",
				"posts.yml",
				"post_flairs.yml",
				"post_flair_emojis.yml",
			},
			wantPostFlairEmojis: []models.PostFlairEmoji{
				{
					EmojiID:     uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					PostFlairID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					OrderIndex:  1,
				},
				{
					EmojiID:     uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					PostFlairID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					OrderIndex:  2,
				},
				{
					EmojiID:     uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					PostFlairID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					OrderIndex:  1,
				},
				{
					EmojiID:     uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					PostFlairID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					OrderIndex:  2,
				},
			},
			wantErr: nil,
		},
		{
			name:                "no post flair emojis :POS",
			fixtureFiles:        []string{},
			wantPostFlairEmojis: nil,
			wantErr:             nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupPostgres(t, tt.fixtureFiles...)
			pgrepo := relationrepo.NewRepo(db)

			gotPostFlairEmojis, gotErr := pgrepo.PostFlairEmojis(context.Background())

			assert.ErrorIs(t, tt.wantErr, gotErr, "expect error to match")
			assert.Equal(t, tt.wantPostFlairEmojis, gotPostFlairEmojis, "expect post flair emojis to match")
		})
	}
}

func TestRepo_PostFlairCustomEmojis(t *testing.T) {
	tests := []struct {
		name                      string
		fixtureFiles              []string
		wantPostFlairCustomEmojis []models.PostFlairCustomEmoji
		wantErr                   error
	}{
		{
			name: "post flair custom emojis :POS",
			fixtureFiles: []string{
				"topics.yml",
				"voxspheres.yml",
				"users.yml",
				"posts.yml",
				"custom_emojis.yml",
				"post_flairs.yml",
				"post_flair_custom_emojis.yml",
			},
			wantPostFlairCustomEmojis: []models.PostFlairCustomEmoji{
				{
					CustomEmojiID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					PostFlairID:   uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					OrderIndex:    4,
				},
				{
					CustomEmojiID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					PostFlairID:   uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					OrderIndex:    5,
				},
				{
					CustomEmojiID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					PostFlairID:   uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					OrderIndex:    4,
				},
				{
					CustomEmojiID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					PostFlairID:   uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					OrderIndex:    5,
				},
			},
			wantErr: nil,
		},
		{
			name:                      "no post flair custom emojis :POS",
			fixtureFiles:              []string{},
			wantPostFlairCustomEmojis: nil,
			wantErr:                   nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupPostgres(t, tt.fixtureFiles...)
			pgrepo := relationrepo.NewRepo(db)

			gotPostFlairCustomEmojis, gotErr := pgrepo.PostFlairCustomEmojis(context.Background())

			assert.ErrorIs(t, tt.wantErr, gotErr, "expect error to match")
			assert.Equal(t, tt.wantPostFlairCustomEmojis, gotPostFlairCustomEmojis, "expect post flair custom emojis to match")
		})
	}
}

func TestRepo_PostFlairDescriptions(t *testing.T) {
	tests := []struct {
		name                      string
		fixtureFiles              []string
		wantPostFlairDescriptions []models.PostFlairDescription
		wantErr                   error
	}{
		{
			name: "post flair descriptions :POS",
			fixtureFiles: []string{
				"topics.yml",
				"voxspheres.yml",
				"users.yml",
				"posts.yml",
				"post_flairs.yml",
				"post_flair_descriptions.yml",
			},
			wantPostFlairDescriptions: []models.PostFlairDescription{
				{
					PostFlairID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					OrderIndex:  0,
					Description: "desc1 ",
				},
				{
					PostFlairID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					OrderIndex:  3,
					Description: " ",
				},
				{
					PostFlairID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					OrderIndex:  0,
					Description: "desc2 ",
				},
				{
					PostFlairID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					OrderIndex:  3,
					Description: " ",
				},
				{
					PostFlairID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					OrderIndex:  6,
					Description: " desc2",
				},
			},
			wantErr: nil,
		},
		{
			name:                      "no post flair descriptions :POS",
			fixtureFiles:              []string{},
			wantPostFlairDescriptions: nil,
			wantErr:                   nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupPostgres(t, tt.fixtureFiles...)
			pgrepo := relationrepo.NewRepo(db)

			gotPostFlairDescriptions, gotErr := pgrepo.PostFlairDescriptions(context.Background())

			assert.ErrorIs(t, tt.wantErr, gotErr, "expect error to match")
			assert.Equal(t, tt.wantPostFlairDescriptions, gotPostFlairDescriptions, "expect post flair descriptions to match")
		})
	}
}

func TestRepo_LinkPostFlairEmojis(t *testing.T) {
	t.Run("duplicate post_flair_id,emoji_id,order_index while linking post flair and emoji :NEG", func(t *testing.T) {
		db := setupPostgres(t,
			"topics.yml",
			"voxspheres.yml",
			"users.yml",
			"posts.yml",
			"emojis.yml",
			"post_flairs.yml",
			"post_flair_emojis.yml",
		)
		pgrepo := relationrepo.NewRepo(db)

		gotPostFlairEmojis, gotErr := pgrepo.LinkPostFlairEmojis(
			context.Background(),
			models.PostFlairEmoji{
				EmojiID:     uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				PostFlairID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				OrderIndex:  1,
			},
		)

		assert.ErrorIs(t, relationrepo.ErrDuplicateID, gotErr, "expect error to match")
		assert.Equal(t, []models.PostFlairEmoji(nil), gotPostFlairEmojis, "expect post flair emojis to match")
	})

	t.Run("emoji not found in parent table while linking post flair and custom emoji :NEG", func(t *testing.T) {
		db := setupPostgres(t,
			"topics.yml",
			"voxspheres.yml",
			"users.yml",
			"posts.yml",
			"emojis.yml",
			"post_flairs.yml",
			"post_flair_emojis.yml",
		)
		pgrepo := relationrepo.NewRepo(db)

		gotPostFlairEmojis, gotErr := pgrepo.LinkPostFlairEmojis(
			context.Background(),
			models.PostFlairEmoji{
				EmojiID:     uuid.MustParse("00000000-0000-0000-0000-000000000003"),
				PostFlairID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				OrderIndex:  1,
			},
		)

		assert.ErrorIs(t, relationrepo.ErrParentTableRecordNotFound, gotErr, "expect error to match")
		assert.Equal(t, []models.PostFlairEmoji(nil), gotPostFlairEmojis, "expect post flair emojis to match")
	})

	t.Run("post flair not found in parent table while linking post flair and custom emoji :NEG", func(t *testing.T) {
		db := setupPostgres(t,
			"topics.yml",
			"voxspheres.yml",
			"users.yml",
			"posts.yml",
			"emojis.yml",
			"post_flairs.yml",
			"post_flair_emojis.yml",
		)
		pgrepo := relationrepo.NewRepo(db)

		gotPostFlairEmojis, gotErr := pgrepo.LinkPostFlairEmojis(
			context.Background(),
			models.PostFlairEmoji{
				EmojiID:     uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				PostFlairID: uuid.MustParse("00000000-0000-0000-0000-000000000009"),
				OrderIndex:  1,
			},
		)

		assert.ErrorIs(t, relationrepo.ErrParentTableRecordNotFound, gotErr, "expect error to match")
		assert.Equal(t, []models.PostFlairEmoji(nil), gotPostFlairEmojis, "expect post flair emojis to match")
	})

	t.Run("link post flair and emoji :POS", func(t *testing.T) {
		db := setupPostgres(t,
			"topics.yml",
			"voxspheres.yml",
			"users.yml",
			"posts.yml",
			"emojis.yml",
			"post_flairs.yml",
			"post_flair_emojis.yml",
		)
		pgrepo := relationrepo.NewRepo(db)

		gotPostFlairEmojis, gotErr := pgrepo.LinkPostFlairEmojis(
			context.Background(),
			models.PostFlairEmoji{
				EmojiID:     uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				PostFlairID: uuid.MustParse("00000000-0000-0000-0000-000000000003"),
				OrderIndex:  0,
			},
		)

		assert.ErrorIs(t, nil, gotErr, "expect error to match")
		assert.Equal(t, []models.PostFlairEmoji{
			{
				EmojiID:     uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				PostFlairID: uuid.MustParse("00000000-0000-0000-0000-000000000003"),
				OrderIndex:  0,
			},
		}, gotPostFlairEmojis, "expect post flair emojis to match")
	})

	t.Run("on deleting post flair child refrences gets deleted in post_flair_emojis table :POS", func(t *testing.T) {
		db := setupPostgres(t,
			"topics.yml",
			"voxspheres.yml",
			"users.yml",
			"posts.yml",
			"emojis.yml",
			"post_flairs.yml",
			"post_flair_emojis.yml",
		)
		relationPgrepo := relationrepo.NewRepo(db)
		postFlairPgrepo := postFlairrepo.NewRepo(db)

		gotErr := postFlairPgrepo.DeletePostFlair(context.Background(), uuid.MustParse("00000000-0000-0000-0000-000000000001"))
		assert.NoError(t, gotErr)

		gotPostFlairEmojis, gotErr := relationPgrepo.PostFlairEmojis(context.Background())

		assert.NoError(t, gotErr)
		assert.Equal(t, []models.PostFlairEmoji{
			{
				EmojiID:     uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				PostFlairID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				OrderIndex:  1,
			},
			{
				EmojiID:     uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				PostFlairID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				OrderIndex:  2,
			},
		}, gotPostFlairEmojis, "expect post flair emojis to match")
	})

	t.Run("on deleting emoji child refrences gets deleted in post_flair_emojis table :POS", func(t *testing.T) {
		db := setupPostgres(t,
			"topics.yml",
			"voxspheres.yml",
			"users.yml",
			"posts.yml",
			"emojis.yml",
			"post_flairs.yml",
			"post_flair_emojis.yml",
		)
		relationPgrepo := relationrepo.NewRepo(db)
		emojiPgrepo := emojirepo.NewRepo(db)

		gotErr := emojiPgrepo.DeleteEmoji(context.Background(), uuid.MustParse("00000000-0000-0000-0000-000000000001"))

		assert.NoError(t, gotErr)

		gotPostFlairEmojis, gotErr := relationPgrepo.PostFlairEmojis(context.Background())

		assert.NoError(t, gotErr)
		assert.Equal(t, []models.PostFlairEmoji{
			{
				EmojiID:     uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				PostFlairID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				OrderIndex:  1,
			},
			{
				EmojiID:     uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				PostFlairID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				OrderIndex:  2,
			},
		}, gotPostFlairEmojis, "expect post flair emojis to match")
	})
}

func TestRepo_LinkPostFlairCustomEmojis(t *testing.T) {
	t.Run("duplicate post_flair_id,emoji_id,order_index while linking post flair and custom emoji :NEG", func(t *testing.T) {
		db := setupPostgres(t,
			"topics.yml",
			"voxspheres.yml",
			"users.yml",
			"posts.yml",
			"custom_emojis.yml",
			"post_flairs.yml",
			"post_flair_custom_emojis.yml",
		)
		pgrepo := relationrepo.NewRepo(db)

		gotPostFlairCustomEmojis, gotErr := pgrepo.LinkPostFlairCustomEmojis(
			context.Background(),
			models.PostFlairCustomEmoji{
				CustomEmojiID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				PostFlairID:   uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				OrderIndex:    4,
			},
		)

		assert.ErrorIs(t, relationrepo.ErrDuplicateID, gotErr, "expect error to match")
		assert.Equal(t, []models.PostFlairCustomEmoji(nil), gotPostFlairCustomEmojis, "expect post flair custom emojis to match")
	})

	t.Run("custom emoji not found in parent table while linking post flair and custom emoji :NEG", func(t *testing.T) {
		db := setupPostgres(t,
			"topics.yml",
			"voxspheres.yml",
			"users.yml",
			"posts.yml",
			"custom_emojis.yml",
			"post_flairs.yml",
			"post_flair_custom_emojis.yml",
		)
		pgrepo := relationrepo.NewRepo(db)

		gotPostFlairCustomEmojis, gotErr := pgrepo.LinkPostFlairCustomEmojis(
			context.Background(),
			models.PostFlairCustomEmoji{
				CustomEmojiID: uuid.MustParse("00000000-0000-0000-0000-000000000003"),
				PostFlairID:   uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				OrderIndex:    1,
			},
		)

		assert.ErrorIs(t, relationrepo.ErrParentTableRecordNotFound, gotErr, "expect error to match")
		assert.Equal(t, []models.PostFlairCustomEmoji(nil), gotPostFlairCustomEmojis, "expect post flair custom emojis to match")
	})

	t.Run("post flair not found in parent table while linking post flair and custom emoji :NEG", func(t *testing.T) {
		db := setupPostgres(t,
			"topics.yml",
			"voxspheres.yml",
			"users.yml",
			"posts.yml",
			"custom_emojis.yml",
			"post_flairs.yml",
			"post_flair_custom_emojis.yml",
		)
		pgrepo := relationrepo.NewRepo(db)

		gotPostFlairCustomEmojis, gotErr := pgrepo.LinkPostFlairCustomEmojis(
			context.Background(),
			models.PostFlairCustomEmoji{
				CustomEmojiID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				PostFlairID:   uuid.MustParse("00000000-0000-0000-0000-000000000009"),
				OrderIndex:    1,
			},
		)

		assert.ErrorIs(t, relationrepo.ErrParentTableRecordNotFound, gotErr, "expect error to match")
		assert.Equal(t, []models.PostFlairCustomEmoji(nil), gotPostFlairCustomEmojis, "expect post flair custom emojis to match")
	})

	t.Run("link post flair and custom emoji :POS", func(t *testing.T) {
		db := setupPostgres(t,
			"topics.yml",
			"voxspheres.yml",
			"users.yml",
			"posts.yml",
			"custom_emojis.yml",
			"post_flairs.yml",
			"post_flair_custom_emojis.yml",
		)
		pgrepo := relationrepo.NewRepo(db)

		gotPostFlairCustomEmojis, gotErr := pgrepo.LinkPostFlairCustomEmojis(
			context.Background(),
			models.PostFlairCustomEmoji{
				CustomEmojiID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				PostFlairID:   uuid.MustParse("00000000-0000-0000-0000-000000000003"),
				OrderIndex:    1,
			},
		)

		assert.ErrorIs(t, nil, gotErr, "expect error to match")
		assert.Equal(t, []models.PostFlairCustomEmoji{
			{
				CustomEmojiID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				PostFlairID:   uuid.MustParse("00000000-0000-0000-0000-000000000003"),
				OrderIndex:    1,
			},
		}, gotPostFlairCustomEmojis, "expect post flair custom emojis to match")
	})

	t.Run("on deleting post flair child refrences gets deleted in post_flair_custom_emojis table :POS", func(t *testing.T) {
		db := setupPostgres(t,
			"topics.yml",
			"voxspheres.yml",
			"users.yml",
			"posts.yml",
			"custom_emojis.yml",
			"post_flairs.yml",
			"post_flair_custom_emojis.yml",
		)
		relationPgrepo := relationrepo.NewRepo(db)
		postFlairPgrepo := postFlairrepo.NewRepo(db)

		gotErr := postFlairPgrepo.DeletePostFlair(context.Background(), uuid.MustParse("00000000-0000-0000-0000-000000000001"))
		assert.NoError(t, gotErr)

		gotPostFlairCustomEmojis, gotErr := relationPgrepo.PostFlairCustomEmojis(context.Background())

		assert.NoError(t, gotErr)
		assert.Equal(t, []models.PostFlairCustomEmoji{
			{
				CustomEmojiID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				PostFlairID:   uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				OrderIndex:    4,
			},
			{
				CustomEmojiID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				PostFlairID:   uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				OrderIndex:    5,
			},
		}, gotPostFlairCustomEmojis, "expect post flair custom emojis to match")
	})

	t.Run("on deleting custom emoji child refrences gets deleted in post_flair_custom_emojis table :POS", func(t *testing.T) {
		db := setupPostgres(t,
			"topics.yml",
			"voxspheres.yml",
			"users.yml",
			"posts.yml",
			"custom_emojis.yml",
			"post_flairs.yml",
			"post_flair_custom_emojis.yml",
		)
		relationPgrepo := relationrepo.NewRepo(db)
		customemojiPgrepo := customemojirepo.NewRepo(db)

		gotErr := customemojiPgrepo.DeleteCustomEmoji(context.Background(), uuid.MustParse("00000000-0000-0000-0000-000000000001"))

		assert.NoError(t, gotErr)

		gotPostFlairCustomEmojis, gotErr := relationPgrepo.PostFlairCustomEmojis(context.Background())

		assert.NoError(t, gotErr)
		assert.Equal(t, []models.PostFlairCustomEmoji{
			{
				CustomEmojiID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				PostFlairID:   uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				OrderIndex:    4,
			},
			{
				CustomEmojiID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				PostFlairID:   uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				OrderIndex:    5,
			},
		}, gotPostFlairCustomEmojis, "expect post flair custom emojis to match")
	})
}

func TestRepo_LinkPostFlairDescriptions(t *testing.T) {
	t.Run("duplicate post_flair_id,order_index while linking post flair and description :NEG", func(t *testing.T) {
		db := setupPostgres(t,
			"topics.yml",
			"voxspheres.yml",
			"users.yml",
			"posts.yml",
			"post_flairs.yml",
			"post_flair_descriptions.yml",
		)
		pgrepo := relationrepo.NewRepo(db)

		gotPostFlairDescriptions, gotErr := pgrepo.LinkPostFlairDescriptions(
			context.Background(),
			models.PostFlairDescription{
				PostFlairID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				OrderIndex:  0,
				Description: "text",
			},
		)

		assert.ErrorIs(t, relationrepo.ErrDuplicateID, gotErr, "expect error to match")
		assert.Equal(t, []models.PostFlairDescription(nil), gotPostFlairDescriptions, "expect post flair descriptions to match")
	})

	t.Run("post flair not found in parent table while linking post flair and description :NEG", func(t *testing.T) {
		db := setupPostgres(t,
			"topics.yml",
			"voxspheres.yml",
			"users.yml",
			"posts.yml",
			"post_flairs.yml",
			"post_flair_descriptions.yml",
		)
		pgrepo := relationrepo.NewRepo(db)

		gotPostFlairDescription, gotErr := pgrepo.LinkPostFlairDescriptions(
			context.Background(),
			models.PostFlairDescription{
				PostFlairID: uuid.MustParse("00000000-0000-0000-0000-000000000009"),
				OrderIndex:  3,
				Description: "desc3",
			},
		)

		assert.ErrorIs(t, relationrepo.ErrParentTableRecordNotFound, gotErr, "expect error to match")
		assert.Equal(t, []models.PostFlairDescription(nil), gotPostFlairDescription, "expect post flair descriptions to match")
	})

	t.Run("link user flair and description :POS", func(t *testing.T) {
		db := setupPostgres(t,
			"topics.yml",
			"voxspheres.yml",
			"users.yml",
			"posts.yml",
			"post_flairs.yml",
			"post_flair_descriptions.yml",
		)
		pgrepo := relationrepo.NewRepo(db)

		gotPostFlairDescriptions, gotErr := pgrepo.LinkPostFlairDescriptions(
			context.Background(),
			models.PostFlairDescription{
				PostFlairID: uuid.MustParse("00000000-0000-0000-0000-000000000003"),
				OrderIndex:  2,
				Description: "desc3",
			},
		)

		assert.ErrorIs(t, nil, gotErr, "expect error to match")
		assert.Equal(t, []models.PostFlairDescription{
			{
				PostFlairID: uuid.MustParse("00000000-0000-0000-0000-000000000003"),
				OrderIndex:  2,
				Description: "desc3",
			},
		}, gotPostFlairDescriptions, "expect post flair descriptions to match")
	})

	t.Run("on deleting post flair child refrences gets deleted in post_flair_descriptions table :POS", func(t *testing.T) {
		db := setupPostgres(t,
			"topics.yml",
			"voxspheres.yml",
			"users.yml",
			"posts.yml",
			"post_flairs.yml",
			"post_flair_descriptions.yml",
		)
		relationPgrepo := relationrepo.NewRepo(db)
		postFlairPgrepo := postFlairrepo.NewRepo(db)

		gotErr := postFlairPgrepo.DeletePostFlair(context.Background(), uuid.MustParse("00000000-0000-0000-0000-000000000001"))
		assert.NoError(t, gotErr)

		gotPostFlairDescriptions, gotErr := relationPgrepo.PostFlairDescriptions(context.Background())

		assert.NoError(t, gotErr)
		assert.Equal(t, []models.PostFlairDescription{
			{
				PostFlairID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				Description: "desc2 ",
				OrderIndex:  0,
			},
			{
				PostFlairID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				Description: " ",
				OrderIndex:  3,
			},
			{
				PostFlairID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				Description: " desc2",
				OrderIndex:  6,
			},
		}, gotPostFlairDescriptions, "expect post flair descriptions to match")
	})
}

func TestRepo_PostAwards(t *testing.T) {
	tests := []struct {
		name           string
		fixtureFiles   []string
		wantPostAwards []models.PostAward
		wantErr        error
	}{
		{
			name:         "post awards :POS",
			fixtureFiles: []string{"topics.yml", "voxspheres.yml", "users.yml", "awards.yml", "posts.yml", "post_awards.yml"},
			wantPostAwards: []models.PostAward{
				{
					PostID:  uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					AwardID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				},
			},
			wantErr: nil,
		},
		{
			name:           "no post awards :POS",
			fixtureFiles:   []string{},
			wantPostAwards: nil,
			wantErr:        nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupPostgres(t, tt.fixtureFiles...)
			pgrepo := relationrepo.NewRepo(db)

			gotPostAwards, gotErr := pgrepo.PostAwards(context.Background())

			assert.ErrorIs(t, tt.wantErr, gotErr, "expect error to match")
			assert.Equal(t, tt.wantPostAwards, gotPostAwards, "expect post awards to match")
		})
	}
}

func TestRepo_LinkPostAwards(t *testing.T) {
	t.Run("duplicate post_id,award_id while linking post and award :NEG", func(t *testing.T) {
		db := setupPostgres(
			t,
			"topics.yml",
			"voxspheres.yml",
			"users.yml",
			"awards.yml",
			"posts.yml",
			"post_awards.yml",
		)
		pgrepo := relationrepo.NewRepo(db)

		gotPostAwards, gotErr := pgrepo.LinkPostAwards(
			context.Background(),
			models.PostAward{
				PostID:  uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				AwardID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
			},
		)

		assert.ErrorIs(t, relationrepo.ErrDuplicateID, gotErr, "expect error to match")
		assert.Equal(t, []models.PostAward(nil), gotPostAwards, "expect post award to match")
	})

	t.Run("award not found while linking post and award :NEG", func(t *testing.T) {
		db := setupPostgres(t, "topics.yml", "voxspheres.yml", "users.yml", "awards.yml", "posts.yml", "post_awards.yml")
		pgrepo := relationrepo.NewRepo(db)

		gotPostAwards, gotErr := pgrepo.LinkPostAwards(
			context.Background(),
			models.PostAward{
				PostID:  uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				AwardID: uuid.MustParse("00000000-0000-0000-0000-000000000003"),
			},
		)

		assert.ErrorIs(t, relationrepo.ErrParentTableRecordNotFound, gotErr, "expect error to match")
		assert.Equal(t, []models.PostAward(nil), gotPostAwards, "expect post award to match")
	})

	t.Run("post not found while linking post and award :NEG", func(t *testing.T) {
		db := setupPostgres(t, "topics.yml", "voxspheres.yml", "users.yml", "awards.yml", "posts.yml", "post_awards.yml")
		pgrepo := relationrepo.NewRepo(db)

		gotPostAwards, gotErr := pgrepo.LinkPostAwards(
			context.Background(),
			models.PostAward{
				PostID:  uuid.MustParse("00000000-0000-0000-0000-000000000009"),
				AwardID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
			},
		)

		assert.ErrorIs(t, relationrepo.ErrParentTableRecordNotFound, gotErr, "expect error to match")
		assert.Equal(t, []models.PostAward(nil), gotPostAwards, "expect post award to match")
	})

	t.Run("link post and award :POS", func(t *testing.T) {
		db := setupPostgres(t, "topics.yml", "voxspheres.yml", "users.yml", "awards.yml", "posts.yml", "post_awards.yml")
		pgrepo := relationrepo.NewRepo(db)

		gotPostAwards, gotErr := pgrepo.LinkPostAwards(
			context.Background(),
			models.PostAward{
				PostID:  uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				AwardID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			},
		)

		assert.ErrorIs(t, nil, gotErr, "expect error to match")
		assert.Equal(t, []models.PostAward{
			{
				PostID:  uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				AwardID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			},
		}, gotPostAwards, "expect post award to match")
	})

	t.Run("on deleting post id child refrences gets deleted in post_awards table :POS", func(t *testing.T) {
		db := setupPostgres(t, "topics.yml", "voxspheres.yml", "users.yml", "awards.yml", "posts.yml", "post_awards.yml")
		relationPgrepo := relationrepo.NewRepo(db)
		postPgrepo := postrepo.NewRepo(db)

		gotErr := postPgrepo.DeletePost(context.Background(), uuid.MustParse("00000000-0000-0000-0000-000000000001"))

		assert.NoError(t, gotErr)

		gotPostAwards, gotErr := relationPgrepo.PostAwards(context.Background())

		assert.NoError(t, gotErr)
		assert.Equal(t, []models.PostAward(nil), gotPostAwards, "expect post awards to match")
	})

	t.Run("on deleting award id child refrences gets deleted in post_awards table :POS", func(t *testing.T) {
		db := setupPostgres(t, "topics.yml", "voxspheres.yml", "users.yml", "awards.yml", "posts.yml", "post_awards.yml")
		relationPgrepo := relationrepo.NewRepo(db)
		awardPgrepo := awardrepo.NewRepo(db)

		gotErr := awardPgrepo.DeleteAward(context.Background(), uuid.MustParse("00000000-0000-0000-0000-000000000002"))

		assert.NoError(t, gotErr)

		gotPostAwards, gotErr := relationPgrepo.PostAwards(context.Background())

		assert.NoError(t, gotErr)
		assert.Equal(t, []models.PostAward(nil), gotPostAwards, "expect post awards to match")
	})
}
