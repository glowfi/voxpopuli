package trophy_test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"slices"
	"testing"

	"github.com/glowfi/voxpopuli/backend/pkg/models"
	trophyrepo "github.com/glowfi/voxpopuli/backend/pkg/repo/trophy"
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

	db.RegisterModel((*models.Trophy)(nil))

	// drop all rows of the trophy table
	_, err := db.NewTruncateTable().Cascade().Model((*models.Trophy)(nil)).Exec(context.Background())
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

func assertTrophies(t *testing.T, wantTrophies, gotTrophies []models.Trophy) {
	t.Helper()

	if len(wantTrophies) != len(gotTrophies) {
		t.Fatal("length of wantTrophies and gotTrophies do not match")
	}

	for _, trophy := range wantTrophies {
		idx := slices.IndexFunc(gotTrophies, func(v models.Trophy) bool {
			return v.ID == trophy.ID
		})

		if idx == -1 {
			t.Fatalf("trophy %v of ID %v is not present in gotTrophies", trophy.Title, trophy.ID)
			return
		}
		assert.Equal(t, trophy, gotTrophies[idx], "expect trophy to match")
	}
}

func TestRepo_Trophies(t *testing.T) {
	type args struct{}
	tests := []struct {
		name         string
		fixtureFiles []string
		args         args
		wantTrophies []models.Trophy
		wantErr      error
	}{
		{
			name:         "no trophies :NEG",
			fixtureFiles: []string{},
			args:         args{},
			wantTrophies: []models.Trophy{},
			wantErr:      nil,
		},
		{
			name:         "trophies :POS",
			fixtureFiles: []string{"trophies.yml"},
			args:         args{},
			wantTrophies: []models.Trophy{
				{
					ID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Title:       "trophy_foo",
					Description: "description_foo",
					ImageLink:   "image_link_foo",
				},
				{
					ID:          uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					Title:       "trophy_bar",
					Description: "description_bar",
					ImageLink:   "image_link_bar",
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupPostgres(t, tt.fixtureFiles...)
			pgrepo := trophyrepo.NewRepo(db)

			gotTrophies, gotErr := pgrepo.Trophies(context.Background())

			assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")
			assertTrophies(t, tt.wantTrophies, gotTrophies)
		})
	}
}

func TestRepo_TrophyByID(t *testing.T) {
	type args struct {
		ID uuid.UUID
	}
	tests := []struct {
		name         string
		fixtureFiles []string
		args         args
		wantTrophy   models.Trophy
		wantErr      error
	}{
		{
			name:         "trophy id not found :NEG",
			fixtureFiles: []string{"trophies.yml"},
			args: args{
				ID: uuid.MustParse("00000000-0000-0000-0000-000000000006"),
			},
			wantTrophy: models.Trophy{},
			wantErr:    trophyrepo.ErrTrophyNotFound,
		},
		{
			name:         "trophy by id :POS",
			fixtureFiles: []string{"trophies.yml"},
			args: args{
				ID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			},
			wantTrophy: models.Trophy{
				ID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				Title:       "trophy_foo",
				Description: "description_foo",
				ImageLink:   "image_link_foo",
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupPostgres(t, tt.fixtureFiles...)
			pgrepo := trophyrepo.NewRepo(db)

			gotTrophy, gotErr := pgrepo.TrophyByID(context.Background(), tt.args.ID)

			assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")
			assert.Equal(t, tt.wantTrophy, gotTrophy, "expect trophy to match")
		})
	}
}

func TestRepo_AddTrophies(t *testing.T) {
	type args struct {
		trophies []models.Trophy
	}
	tests := []struct {
		name                 string
		fixtureFiles         []string
		args                 args
		wantInsertedTrophies []models.Trophy
		wantTrophies         []models.Trophy
		wantErr              error
	}{
		{
			name:         "duplicate trophy title :NEG",
			fixtureFiles: []string{"trophies.yml"},
			args: args{
				trophies: []models.Trophy{
					{
						ID:          uuid.MustParse("00000000-0000-0000-0000-000000000003"),
						Title:       "trophy_foo",
						Description: "new description",
						ImageLink:   "new image link",
					},
				},
			},
			wantInsertedTrophies: nil,
			wantTrophies: []models.Trophy{
				{
					ID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Title:       "trophy_foo",
					Description: "description_foo",
					ImageLink:   "image_link_foo",
				},
				{
					ID:          uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					Title:       "trophy_bar",
					Description: "description_bar",
					ImageLink:   "image_link_bar",
				},
			},
			wantErr: trophyrepo.ErrTrophyDuplicateIDorTitle,
		},
		{
			name:         "add trophies :POS",
			fixtureFiles: []string{"trophies.yml"},
			args: args{
				trophies: []models.Trophy{
					{
						ID:          uuid.MustParse("00000000-0000-0000-0000-000000000003"),
						Title:       "new trophy",
						Description: "new description",
						ImageLink:   "new image link",
					},
				},
			},
			wantInsertedTrophies: []models.Trophy{
				{
					ID:          uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					Title:       "new trophy",
					Description: "new description",
					ImageLink:   "new image link",
				},
			},
			wantTrophies: []models.Trophy{
				{
					ID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Title:       "trophy_foo",
					Description: "description_foo",
					ImageLink:   "image_link_foo",
				},
				{
					ID:          uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					Title:       "trophy_bar",
					Description: "description_bar",
					ImageLink:   "image_link_bar",
				},
				{
					ID:          uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					Title:       "new trophy",
					Description: "new description",
					ImageLink:   "new image link",
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupPostgres(t, tt.fixtureFiles...)
			pgrepo := trophyrepo.NewRepo(db)

			gotInsertedTrophies, gotErr := pgrepo.AddTrophies(context.Background(), tt.args.trophies...)

			assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")
			assert.Equal(t, tt.wantInsertedTrophies, gotInsertedTrophies, "expect inserted trophies to match")

			gotTrophies, err := pgrepo.Trophies(context.Background())

			assert.NoError(t, err, "expect no error while getting trophies")
			assertTrophies(t, tt.wantTrophies, gotTrophies)
		})
	}
}

func TestRepo_UpdateTrophy(t *testing.T) {
	type args struct {
		trophy models.Trophy
	}
	tests := []struct {
		name         string
		fixtureFiles []string
		args         args
		wantTrophy   models.Trophy
		wantTrophies []models.Trophy
		wantErr      error
	}{
		{
			name:         "trophy not found :NEG",
			fixtureFiles: []string{"trophies.yml"},
			args: args{
				trophy: models.Trophy{
					ID:          uuid.MustParse("00000000-0000-0000-0000-000000000006"),
					Title:       "updated trophy",
					Description: "updated description",
					ImageLink:   "updated image link",
				},
			},
			wantTrophy: models.Trophy{},
			wantTrophies: []models.Trophy{
				{
					ID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Title:       "trophy_foo",
					Description: "description_foo",
					ImageLink:   "image_link_foo",
				},
				{
					ID:          uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					Title:       "trophy_bar",
					Description: "description_bar",
					ImageLink:   "image_link_bar",
				},
			},
			wantErr: trophyrepo.ErrTrophyNotFound,
		},
		{
			name:         "update trophy :POS",
			fixtureFiles: []string{"trophies.yml"},
			args: args{
				trophy: models.Trophy{
					ID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Title:       "updated trophy",
					Description: "updated description",
					ImageLink:   "updated image link",
				},
			},
			wantTrophy: models.Trophy{
				ID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				Title:       "updated trophy",
				Description: "updated description",
				ImageLink:   "updated image link",
			},
			wantTrophies: []models.Trophy{
				{
					ID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Title:       "updated trophy",
					Description: "updated description",
					ImageLink:   "updated image link",
				},
				{
					ID:          uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					Title:       "trophy_bar",
					Description: "description_bar",
					ImageLink:   "image_link_bar",
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupPostgres(t, tt.fixtureFiles...)
			pgrepo := trophyrepo.NewRepo(db)

			gotTrophy, gotErr := pgrepo.UpdateTrophy(context.Background(), tt.args.trophy)

			assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")
			assert.Equal(t, tt.wantTrophy, gotTrophy, "expect trophy to match")

			gotTrophies, err := pgrepo.Trophies(context.Background())

			assert.NoError(t, err, "expect no error while getting trophies")
			assertTrophies(t, tt.wantTrophies, gotTrophies)
		})
	}
}

func TestRepo_DeleteTrophy(t *testing.T) {
	type args struct {
		ID uuid.UUID
	}
	tests := []struct {
		name         string
		fixtureFiles []string
		args         args
		wantTrophies []models.Trophy
		wantErr      error
	}{
		{
			name:         "trophy not found :NEG",
			fixtureFiles: []string{"trophies.yml"},
			args: args{
				ID: uuid.MustParse("00000000-0000-0000-0000-000000000006"),
			},
			wantTrophies: []models.Trophy{
				{
					ID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Title:       "trophy_foo",
					Description: "description_foo",
					ImageLink:   "image_link_foo",
				},
				{
					ID:          uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					Title:       "trophy_bar",
					Description: "description_bar",
					ImageLink:   "image_link_bar",
				},
			},
			wantErr: trophyrepo.ErrTrophyNotFound,
		},
		{
			name:         "delete trophy :POS",
			fixtureFiles: []string{"trophies.yml"},
			args: args{
				ID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			},
			wantTrophies: []models.Trophy{
				{
					ID:          uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					Title:       "trophy_bar",
					Description: "description_bar",
					ImageLink:   "image_link_bar",
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupPostgres(t, tt.fixtureFiles...)
			pgrepo := trophyrepo.NewRepo(db)

			gotErr := pgrepo.DeleteTrophy(context.Background(), tt.args.ID)

			assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")

			gotTrophies, err := pgrepo.Trophies(context.Background())

			assert.NoError(t, err, "expect no error while getting trophies")
			assertTrophies(t, tt.wantTrophies, gotTrophies)
		})
	}
}
