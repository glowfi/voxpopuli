package award_test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
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

	db.RegisterModel((*models.Award)(nil))

	// drop all rows of the topics,award table
	_, err := db.NewTruncateTable().Cascade().Model((*models.Award)(nil)).Exec(context.Background())
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
