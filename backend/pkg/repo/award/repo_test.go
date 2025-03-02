package award_test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"slices"
	"testing"

	"github.com/glowfi/voxpopuli/backend/pkg/models"
	awardrepo "github.com/glowfi/voxpopuli/backend/pkg/repo/award"
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

	db.RegisterModel((*models.Award)(nil))

	// drop all rows of the award table
	_, err := db.NewTruncateTable().Cascade().Model((*models.Award)(nil)).Exec(context.Background())
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

func assertAwards(t *testing.T, wantAwards, gotAwards []models.Award) {
	t.Helper()

	if len(wantAwards) != len(gotAwards) {
		t.Fatal("length of wantAwards and gotAwards do not match")
	}

	for _, award := range wantAwards {
		idx := slices.IndexFunc(gotAwards, func(v models.Award) bool {
			return v.ID == award.ID
		})

		if idx == -1 {
			t.Fatal(fmt.Sprintf("award %v of ID %v is not present in gotAwards", award.Title, award.ID))
			return
		}
		assert.Equal(t, award, gotAwards[idx], "expect award to match")
	}
}

func TestRepo_Awards(t *testing.T) {
	tests := []struct {
		name         string
		fixtureFiles []string
		wantAwards   []models.Award
		wantErr      error
	}{
		{
			name:         "awards :POS",
			fixtureFiles: []string{"awards.yml"},
			wantAwards: []models.Award{
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Title:     "award_foo",
					ImageLink: "https:/fooimage.com",
				},
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					Title:     "award_bar",
					ImageLink: "https:/barimage.com",
				},
			},
			wantErr: nil,
		},
		{
			name:         "no awards :POS",
			fixtureFiles: []string{},
			wantAwards:   nil,
			wantErr:      nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupPostgres(t, tt.fixtureFiles...)
			pgrepo := awardrepo.NewRepo(db)

			gotAwards, gotErr := pgrepo.Awards(context.Background())

			assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")
			assert.Equal(t, tt.wantAwards, gotAwards, "expect awards to match")
		})
	}
}

func TestRepo_AwardByID(t *testing.T) {
	type args struct {
		ID uuid.UUID
	}
	tests := []struct {
		name         string
		fixtureFiles []string
		args         args
		wantAward    models.Award
		wantErr      error
	}{
		{
			name:         "award not found :NEG",
			fixtureFiles: []string{},
			args: args{
				ID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			},
			wantAward: models.Award{},
			wantErr:   awardrepo.ErrAwardNotFound,
		},
		{
			name:         "award found :POS",
			fixtureFiles: []string{"awards.yml"},
			args: args{
				ID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			},
			wantAward: models.Award{
				ID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				Title:     "award_foo",
				ImageLink: "https:/fooimage.com",
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupPostgres(t, tt.fixtureFiles...)
			pgrepo := awardrepo.NewRepo(db)

			gotAward, gotErr := pgrepo.AwardByID(context.Background(), tt.args.ID)

			assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")
			assert.Equal(t, tt.wantAward, gotAward, "expect award to match")
		})
	}
}

func TestRepo_AddAwards(t *testing.T) {
	type args struct {
		awards []models.Award
	}
	tests := []struct {
		name               string
		fixtureFiles       []string
		args               args
		wantInsertedAwards []models.Award
		wantAwards         []models.Award
		wantErr            error
	}{
		{
			name:         "duplicate award id :NEG",
			fixtureFiles: []string{"awards.yml"},
			args: args{
				awards: []models.Award{
					{
						ID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
						Title:     "new title",
						ImageLink: "new image link",
					},
				},
			},
			wantInsertedAwards: nil,
			wantAwards: []models.Award{
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Title:     "award_foo",
					ImageLink: "https:/fooimage.com",
				},
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					Title:     "award_bar",
					ImageLink: "https:/barimage.com",
				},
			},
			wantErr: awardrepo.ErrAwardDuplicateIDorTitle,
		},
		{
			name:         "add awards :POS",
			fixtureFiles: []string{"awards.yml"},
			args: args{
				awards: []models.Award{
					{
						ID:        uuid.MustParse("00000000-0000-0000-0000-000000000003"),
						Title:     "new title1",
						ImageLink: "new image link1",
					},
					{
						ID:        uuid.MustParse("00000000-0000-0000-0000-000000000004"),
						Title:     "new title2",
						ImageLink: "new image link2",
					},
				},
			},
			wantInsertedAwards: []models.Award{
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					Title:     "new title1",
					ImageLink: "new image link1",
				},
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000004"),
					Title:     "new title2",
					ImageLink: "new image link2",
				},
			},
			wantAwards: []models.Award{
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Title:     "award_foo",
					ImageLink: "https:/fooimage.com",
				},
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					Title:     "award_bar",
					ImageLink: "https:/barimage.com",
				},
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					Title:     "new title1",
					ImageLink: "new image link1",
				},
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000004"),
					Title:     "new title2",
					ImageLink: "new image link2",
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupPostgres(t, tt.fixtureFiles...)
			pgrepo := awardrepo.NewRepo(db)

			gotInsertedAward, gotErr := pgrepo.AddAwards(context.Background(), tt.args.awards...)

			assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")
			assert.Equal(t, tt.wantInsertedAwards, gotInsertedAward, "expect inserted awards to match")

			gotAwards, err := pgrepo.Awards(context.Background())

			assert.NoError(t, err, "expect no error while getting awards")
			assertAwards(t, tt.wantAwards, gotAwards)
		})
	}
}

func TestRepo_UpdateAward(t *testing.T) {
	type args struct {
		award models.Award
	}
	tests := []struct {
		name         string
		fixtureFiles []string
		args         args
		wantAward    models.Award
		wantAwards   []models.Award
		wantErr      error
	}{
		{
			name:         "award not found :NEG",
			fixtureFiles: []string{"awards.yml"},
			args: args{
				award: models.Award{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000006"),
					Title:     "updated title",
					ImageLink: "updated image link",
				},
			},
			wantAward: models.Award{},
			wantAwards: []models.Award{
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Title:     "award_foo",
					ImageLink: "https:/fooimage.com",
				},
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					Title:     "award_bar",
					ImageLink: "https:/barimage.com",
				},
			},
			wantErr: awardrepo.ErrAwardNotFound,
		},
		{
			name:         "update award :POS",
			fixtureFiles: []string{"awards.yml"},
			args: args{
				award: models.Award{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Title:     "updated title",
					ImageLink: "updated image link",
				},
			},
			wantAward: models.Award{
				ID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				Title:     "updated title",
				ImageLink: "updated image link",
			},
			wantAwards: []models.Award{
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Title:     "updated title",
					ImageLink: "updated image link",
				},
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					Title:     "award_bar",
					ImageLink: "https:/barimage.com",
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupPostgres(t, tt.fixtureFiles...)
			pgrepo := awardrepo.NewRepo(db)

			gotAward, gotErr := pgrepo.UpdateAward(context.Background(), tt.args.award)

			assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")
			assert.Equal(t, tt.wantAward, gotAward, "expect award to match")

			gotAwards, err := pgrepo.Awards(context.Background())

			assert.NoError(t, err, "expect no error while getting awards")
			assertAwards(t, tt.wantAwards, gotAwards)
		})
	}
}

func TestRepo_DeleteAward(t *testing.T) {
	type args struct {
		ID uuid.UUID
	}
	tests := []struct {
		name         string
		fixtureFiles []string
		args         args
		wantAwards   []models.Award
		wantErr      error
	}{
		{
			name:         "award not found :NEG",
			fixtureFiles: []string{"awards.yml"},
			args: args{
				ID: uuid.MustParse("00000000-0000-0000-0000-000000000006"),
			},
			wantAwards: []models.Award{
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Title:     "award_foo",
					ImageLink: "https:/fooimage.com",
				},
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					Title:     "award_bar",
					ImageLink: "https:/barimage.com",
				},
			},
			wantErr: awardrepo.ErrAwardNotFound,
		},
		{
			name:         "award deleted :POS",
			fixtureFiles: []string{"awards.yml"},
			args: args{
				ID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			},
			wantAwards: []models.Award{
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					Title:     "award_bar",
					ImageLink: "https:/barimage.com",
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupPostgres(t, tt.fixtureFiles...)
			pgrepo := awardrepo.NewRepo(db)

			gotErr := pgrepo.DeleteAward(context.Background(), tt.args.ID)

			assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")

			gotAwards, err := pgrepo.Awards(context.Background())

			assert.NoError(t, err, "expect no error while getting awards")
			assertAwards(t, tt.wantAwards, gotAwards)
		})
	}
}
