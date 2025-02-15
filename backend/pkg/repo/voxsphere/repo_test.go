package voxsphere_test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/glowfi/voxpopuli/backend/pkg/models"
	voxrepo "github.com/glowfi/voxpopuli/backend/pkg/repo/voxsphere"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dbfixture"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
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

	// drop all rows of the topics,voxspheres table
	_, err := db.NewTruncateTable().Cascade().Model((*models.Topic)(nil)).Exec(context.Background())
	if err != nil {
		t.Fatal("truncate table failed:", err)
	}
	if _, err := db.NewTruncateTable().Cascade().Model((*models.Voxsphere)(nil)).Exec(context.Background()); err != nil {
		t.Fatal("truncate table failed:", err)
	}

	fixture := dbfixture.New(db)
	if err := fixture.Load(context.Background(), os.DirFS("testdata"), fixtureFiles...); err != nil {
		t.Fatal("failed to load fixtures", err)
	}

	return db
}

func ptrof[T comparable](v T) *T {
	return &v
}

func TestRepo_Voxspheres(t *testing.T) {
	tests := []struct {
		name           string
		fixtureFiles   []string
		wantVoxspheres []models.Voxsphere
		wantErr        error
	}{
		{
			name:         "voxspheres :POS",
			fixtureFiles: []string{"topics.yml", "voxspheres.yml"},
			wantVoxspheres: []models.Voxsphere{
				{
					ID:      uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					TopicID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Topic: models.Topic{
						ID:   uuid.MustParse("00000000-0000-0000-0000-000000000001"),
						Name: "xyz",
					},
					Title:                 "v/foo",
					PublicDescription:     ptrof("foo PublicDescription"),
					CommunityIcon:         ptrof("foo icon"),
					BannerBackgroundImage: ptrof("foo BannerBackgroundImage"),
					BannerBackgroundColor: ptrof("#000000"),
					KeyColor:              ptrof("#000000"),
					PrimaryColor:          ptrof("#000000"),
					Over18:                true,
					SpoilersEnabled:       false,
					CreatedAt:             time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix:         1725091100,
					UpdatedAt:             time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
				},
				{
					ID:      uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					TopicID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					Topic: models.Topic{
						ID:   uuid.MustParse("00000000-0000-0000-0000-000000000002"),
						Name: "pqr",
					},
					Title:                 "v/bar",
					PublicDescription:     ptrof("bar PublicDescription"),
					CommunityIcon:         ptrof("bar icon"),
					BannerBackgroundImage: ptrof("bar BannerBackgroundImage"),
					BannerBackgroundColor: ptrof("#ffffff"),
					KeyColor:              ptrof("#ffffff"),
					PrimaryColor:          ptrof("#ffffff"),
					Over18:                false,
					SpoilersEnabled:       false,
					CreatedAt:             time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
					CreatedAtUnix:         1725091101,
					UpdatedAt:             time.Date(2024, 10, 10, 10, 10, 20, 0, time.UTC),
				},
			},
			wantErr: nil,
		},
		{
			name:           "no voxspheres :POS",
			fixtureFiles:   []string{},
			wantVoxspheres: nil,
			wantErr:        nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupPostgres(t, tt.fixtureFiles...)
			pgrepo := voxrepo.NewRepo(db)

			gotVoxspheres, gotErr := pgrepo.Voxspheres(context.Background())

			assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")
			assert.Equal(t, tt.wantVoxspheres, gotVoxspheres, "expect voxspheres to match")
		})
	}
}
