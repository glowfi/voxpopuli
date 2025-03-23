package voxsphere_test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"slices"
	"testing"
	"time"

	"github.com/glowfi/voxpopuli/backend/pkg/models"
	"github.com/glowfi/voxpopuli/backend/pkg/repo/topic"
	voxrepo "github.com/glowfi/voxpopuli/backend/pkg/repo/voxsphere"
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

	// drop all rows of the topics,voxspheres table
	_, err := db.NewTruncateTable().Cascade().Model((*models.Topic)(nil)).Exec(context.Background())
	if err != nil {
		t.Fatal("truncate table failed:", err)
	}
	if _, err := db.NewTruncateTable().Cascade().Model((*models.Voxsphere)(nil)).Exec(context.Background()); err != nil {
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

func assertVoxsphereWithoutTimestamp(t *testing.T, wantVoxsphere, gotVoxsphere models.Voxsphere) {
	assert.Equal(t, wantVoxsphere.ID, gotVoxsphere.ID, "expected id to match")
	assert.Equal(t, wantVoxsphere.TopicID, gotVoxsphere.TopicID, "expected topic id to match")
	assert.Equal(t, wantVoxsphere.Topic, gotVoxsphere.Topic, "expected topic to match")
	assert.Equal(t, wantVoxsphere.Title, gotVoxsphere.Title, "expected title to match")
	assert.Equal(t, wantVoxsphere.PublicDescription, gotVoxsphere.PublicDescription, "expected public description to match")
	assert.Equal(t, wantVoxsphere.CommunityIcon, gotVoxsphere.CommunityIcon, "expected community icon to match")
	assert.Equal(t, wantVoxsphere.BannerBackgroundImage, gotVoxsphere.BannerBackgroundImage, "expected banner background image to match")
	assert.Equal(t, wantVoxsphere.BannerBackgroundColor, gotVoxsphere.BannerBackgroundColor, "expected banner background color to match")
	assert.Equal(t, wantVoxsphere.KeyColor, gotVoxsphere.KeyColor, "expected key color to match")
	assert.Equal(t, wantVoxsphere.PrimaryColor, gotVoxsphere.PrimaryColor, "expected primary color to match")
	assert.Equal(t, wantVoxsphere.Over18, gotVoxsphere.Over18, "expected over 18 to match")
	assert.Equal(t, wantVoxsphere.SpoilersEnabled, gotVoxsphere.SpoilersEnabled, "expected spoilers enabled to match")
}

func assertVoxspheresWithoutTimestamp(t *testing.T, wantVoxspheres, gotVoxspheres []models.Voxsphere) {
	t.Helper()

	if len(wantVoxspheres) != len(gotVoxspheres) {
		t.Fatal("length of wantVoxspheres and gotVoxspheres do not match")
	}

	for _, voxsphere := range wantVoxspheres {
		idx := slices.IndexFunc(gotVoxspheres, func(v models.Voxsphere) bool {
			return v.ID == voxsphere.ID
		})

		if idx == -1 {
			t.Fatalf("voxsphere %v of ID %v is not present in gotVoxspheres", voxsphere.Title, voxsphere.ID)
			return
		}
		assertVoxsphereWithoutTimestamp(t, voxsphere, gotVoxspheres[idx])
	}
}

func assertVoxspheresWithTimestamp(t *testing.T, wantVoxspheres, gotVoxspheres []models.Voxsphere) {
	t.Helper()

	if len(wantVoxspheres) != len(gotVoxspheres) {
		t.Fatal("length of wantVoxspheres and gotVoxspheres do not match")
	}

	for _, voxsphere := range wantVoxspheres {
		idx := slices.IndexFunc(gotVoxspheres, func(v models.Voxsphere) bool {
			return v.ID == voxsphere.ID
		})

		if idx == -1 {
			t.Fatalf("voxsphere %v of ID %v is not present in gotVoxspheres", voxsphere.Title, voxsphere.ID)
			return
		}
		assert.Equal(t, voxsphere, gotVoxspheres[idx])
	}
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
			assertVoxspheresWithTimestamp(t, tt.wantVoxspheres, gotVoxspheres)
		})
	}
}

func TestRepo_VoxsphereByID(t *testing.T) {
	type args struct {
		ID uuid.UUID
	}

	tests := []struct {
		name          string
		fixtureFiles  []string
		args          args
		wantVoxsphere models.Voxsphere
		wantErr       error
	}{
		{
			name:         "voxsphere not found :NEG",
			fixtureFiles: []string{},
			args: args{
				ID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			},
			wantVoxsphere: models.Voxsphere{},
			wantErr:       voxrepo.ErrVoxsphereNotFound,
		},
		{
			name:         "voxsphere found :POS",
			fixtureFiles: []string{"topics.yml", "voxspheres.yml"},
			args: args{
				ID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			},
			wantVoxsphere: models.Voxsphere{
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
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupPostgres(t, tt.fixtureFiles...)
			pgrepo := voxrepo.NewRepo(db)

			gotVoxsphere, gotErr := pgrepo.VoxsphereByID(context.Background(), tt.args.ID)

			assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")
			assert.Equal(t, tt.wantVoxsphere, gotVoxsphere, "expect voxsphere to match")
		})
	}
}

func TestRepo_AddVoxspheres(t *testing.T) {
	type args struct {
		voxspheres []models.Voxsphere
	}
	tests := []struct {
		name                   string
		fixtureFiles           []string
		args                   args
		wantInsertedVoxspheres []models.Voxsphere
		wantVoxspheres         []models.Voxsphere
		wantErr                error
	}{
		{
			name:         "duplicate voxsphere id :NEG",
			fixtureFiles: []string{"topics.yml", "voxspheres.yml"},
			args: args{
				voxspheres: []models.Voxsphere{
					{
						ID:      uuid.MustParse("00000000-0000-0000-0000-000000000001"),
						TopicID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
						Topic: models.Topic{
							ID:   uuid.MustParse("00000000-0000-0000-0000-000000000001"),
							Name: "xyz",
						},
						Title:                 "v/foobar",
						PublicDescription:     ptrof("foo PublicDescription"),
						CommunityIcon:         ptrof("foo icon"),
						BannerBackgroundImage: ptrof("foo BannerBackgroundImage"),
						BannerBackgroundColor: ptrof("#000000"),
						KeyColor:              ptrof("#000000"),
						PrimaryColor:          ptrof("#000000"),
						Over18:                true,
						SpoilersEnabled:       false,
					},
				},
			},
			wantInsertedVoxspheres: nil,
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
			wantErr: voxrepo.ErrVoxsphereDuplicateIDorTitle,
		},
		{
			name:         "topic does not exist in parent table :NEG",
			fixtureFiles: []string{"topics.yml", "voxspheres.yml"},
			args: args{
				voxspheres: []models.Voxsphere{
					{
						ID:      uuid.MustParse("00000000-0000-0000-0000-000000000003"),
						TopicID: uuid.MustParse("00000000-0000-0000-0000-000000000009"),
						Topic: models.Topic{
							ID:   uuid.MustParse("00000000-0000-0000-0000-000000000001"),
							Name: "xyz",
						},
						Title:                 "v/foobar",
						PublicDescription:     ptrof("foo PublicDescription"),
						CommunityIcon:         ptrof("foo icon"),
						BannerBackgroundImage: ptrof("foo BannerBackgroundImage"),
						BannerBackgroundColor: ptrof("#000000"),
						KeyColor:              ptrof("#000000"),
						PrimaryColor:          ptrof("#000000"),
						Over18:                true,
						SpoilersEnabled:       false,
					},
				},
			},
			wantInsertedVoxspheres: nil,
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
			wantErr: voxrepo.ErrVoxsphereParentTableRecordNotFound,
		},
		{
			name:         "add voxspheres :POS",
			fixtureFiles: []string{"topics.yml", "voxspheres.yml"},
			args: args{
				voxspheres: []models.Voxsphere{
					{
						ID:      uuid.MustParse("00000000-0000-0000-0000-000000000003"),
						TopicID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
						Topic: models.Topic{
							ID:   uuid.MustParse("00000000-0000-0000-0000-000000000001"),
							Name: "xyz",
						},
						Title:                 "v/foobar",
						PublicDescription:     ptrof("foo PublicDescription"),
						CommunityIcon:         ptrof("foo icon"),
						BannerBackgroundImage: ptrof("foo BannerBackgroundImage"),
						BannerBackgroundColor: ptrof("#000000"),
						KeyColor:              ptrof("#000000"),
						PrimaryColor:          ptrof("#000000"),
						Over18:                true,
						SpoilersEnabled:       false,
					},
				},
			},
			wantInsertedVoxspheres: []models.Voxsphere{
				{
					ID:      uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					TopicID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Topic: models.Topic{
						ID:   uuid.MustParse("00000000-0000-0000-0000-000000000001"),
						Name: "xyz",
					},
					Title:                 "v/foobar",
					PublicDescription:     ptrof("foo PublicDescription"),
					CommunityIcon:         ptrof("foo icon"),
					BannerBackgroundImage: ptrof("foo BannerBackgroundImage"),
					BannerBackgroundColor: ptrof("#000000"),
					KeyColor:              ptrof("#000000"),
					PrimaryColor:          ptrof("#000000"),
					Over18:                true,
					SpoilersEnabled:       false,
				},
			},
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
				{
					ID:      uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					TopicID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Topic: models.Topic{
						ID:   uuid.MustParse("00000000-0000-0000-0000-000000000001"),
						Name: "xyz",
					},
					Title:                 "v/foobar",
					PublicDescription:     ptrof("foo PublicDescription"),
					CommunityIcon:         ptrof("foo icon"),
					BannerBackgroundImage: ptrof("foo BannerBackgroundImage"),
					BannerBackgroundColor: ptrof("#000000"),
					KeyColor:              ptrof("#000000"),
					PrimaryColor:          ptrof("#000000"),
					Over18:                true,
					SpoilersEnabled:       false,
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupPostgres(t, tt.fixtureFiles...)
			pgrepo := voxrepo.NewRepo(db)

			startTime := time.Now()
			gotInsertedVoxspheres, gotErr := pgrepo.AddVoxspheres(context.Background(), tt.args.voxspheres...)
			endTime := time.Now()

			for _, gotInsertedVoxsphere := range gotInsertedVoxspheres {
				assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")
				assert.Equal(
					t,
					gotInsertedVoxsphere.UpdatedAt,
					gotInsertedVoxsphere.CreatedAt,
					"expect CreatedAt and UpdatedAt to be same",
				)
				if tt.wantErr == nil {
					assertTimeWithinRange(t, gotInsertedVoxsphere.CreatedAt, startTime, endTime)
					assertTimeWithinRange(t, gotInsertedVoxsphere.UpdatedAt, startTime, endTime)
				}
			}

			gotVoxspheres, err := pgrepo.Voxspheres(context.Background())

			assert.NoError(t, err, "expect no error while getting voxspheres")
			assertVoxspheresWithoutTimestamp(t, tt.wantVoxspheres, gotVoxspheres)
		})
	}
}

func TestRepo_UpdateVoxsphere(t *testing.T) {
	type args struct {
		voxsphere models.Voxsphere
	}
	tests := []struct {
		name           string
		fixtureFiles   []string
		args           args
		wantVoxsphere  models.Voxsphere
		wantVoxspheres []models.Voxsphere
		wantErr        error
	}{
		{
			name:         "voxsphere not found :NEG",
			fixtureFiles: []string{"topics.yml", "voxspheres.yml"},
			args: args{
				voxsphere: models.Voxsphere{
					ID:                    uuid.MustParse("00000000-0000-0000-0000-000000000006"),
					TopicID:               uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Title:                 "v/foo updated",
					PublicDescription:     ptrof("foo PublicDescription updated"),
					CommunityIcon:         ptrof("foo icon updated"),
					BannerBackgroundImage: ptrof("foo BannerBackgroundImage updated"),
					BannerBackgroundColor: ptrof("#ffffff"),
					KeyColor:              ptrof("#ffffff"),
					PrimaryColor:          ptrof("#ffffff"),
					Over18:                true,
					SpoilersEnabled:       false,
				},
			},
			wantVoxsphere: models.Voxsphere{},
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
			wantErr: voxrepo.ErrVoxsphereNotFound,
		},
		{
			name:         "topic is not present in parent table :NEG",
			fixtureFiles: []string{"topics.yml", "voxspheres.yml"},
			args: args{
				voxsphere: models.Voxsphere{
					ID:                    uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					TopicID:               uuid.MustParse("00000000-0000-0000-0000-000000000008"),
					Title:                 "v/foo updated",
					PublicDescription:     ptrof("foo PublicDescription updated"),
					CommunityIcon:         ptrof("foo icon updated"),
					BannerBackgroundImage: ptrof("foo BannerBackgroundImage updated"),
					BannerBackgroundColor: ptrof("#ffffff"),
					KeyColor:              ptrof("#ffffff"),
					PrimaryColor:          ptrof("#ffffff"),
					Over18:                true,
					SpoilersEnabled:       false,
				},
			},
			wantVoxsphere: models.Voxsphere{},
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
			wantErr: voxrepo.ErrVoxsphereParentTableRecordNotFound,
		},
		{
			name:         "update voxsphere :POS",
			fixtureFiles: []string{"topics.yml", "voxspheres.yml"},
			args: args{
				voxsphere: models.Voxsphere{
					ID:                    uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					TopicID:               uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Title:                 "v/foo updated",
					PublicDescription:     ptrof("foo PublicDescription updated"),
					CommunityIcon:         ptrof("foo icon updated"),
					BannerBackgroundImage: ptrof("foo BannerBackgroundImage updated"),
					BannerBackgroundColor: ptrof("#ffffff"),
					KeyColor:              ptrof("#ffffff"),
					PrimaryColor:          ptrof("#ffffff"),
					Over18:                true,
					SpoilersEnabled:       false,
				},
			},
			wantVoxsphere: models.Voxsphere{
				ID:                    uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				TopicID:               uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				Title:                 "v/foo updated",
				PublicDescription:     ptrof("foo PublicDescription updated"),
				CommunityIcon:         ptrof("foo icon updated"),
				BannerBackgroundImage: ptrof("foo BannerBackgroundImage updated"),
				BannerBackgroundColor: ptrof("#ffffff"),
				KeyColor:              ptrof("#ffffff"),
				PrimaryColor:          ptrof("#ffffff"),
				Over18:                true,
				SpoilersEnabled:       false,
			},
			wantVoxspheres: []models.Voxsphere{
				{
					ID:      uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					TopicID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Topic: models.Topic{
						ID:   uuid.MustParse("00000000-0000-0000-0000-000000000001"),
						Name: "xyz",
					},
					Title:                 "v/foo updated",
					PublicDescription:     ptrof("foo PublicDescription updated"),
					CommunityIcon:         ptrof("foo icon updated"),
					BannerBackgroundImage: ptrof("foo BannerBackgroundImage updated"),
					BannerBackgroundColor: ptrof("#ffffff"),
					KeyColor:              ptrof("#ffffff"),
					PrimaryColor:          ptrof("#ffffff"),
					Over18:                true,
					SpoilersEnabled:       false,
					CreatedAt:             time.Date(2024, time.October, 10, 10, 10, 10, 0, time.UTC),
					CreatedAtUnix:         1725091101,
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupPostgres(t, tt.fixtureFiles...)
			pgrepo := voxrepo.NewRepo(db)

			startTime := time.Now()
			gotVoxsphere, gotErr := pgrepo.UpdateVoxsphere(context.Background(), tt.args.voxsphere)
			endTime := time.Now()

			assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")
			if tt.wantErr == nil {
				assertTimeWithinRange(t, gotVoxsphere.UpdatedAt, startTime, endTime)
			}

			gotVoxspheres, err := pgrepo.Voxspheres(context.Background())

			assert.NoError(t, err, "expect no error while getting voxsphere")
			assertVoxspheresWithoutTimestamp(t, tt.wantVoxspheres, gotVoxspheres)
		})
	}
}

func TestRepo_DeleteVoxsphere(t *testing.T) {
	type args struct {
		ID uuid.UUID
	}
	tests := []struct {
		name           string
		fixtureFiles   []string
		args           args
		wantVoxspheres []models.Voxsphere
		wantErr        error
	}{
		{
			name:         "voxsphere not found :NEG",
			fixtureFiles: []string{"topics.yml", "voxspheres.yml"},
			args: args{
				ID: uuid.MustParse("00000000-0000-0000-0000-000000000006"),
			},
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
			wantErr: voxrepo.ErrVoxsphereNotFound,
		},
		{
			name:         "delete voxsphere :POS",
			fixtureFiles: []string{"topics.yml", "voxspheres.yml"},
			args: args{
				ID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			},
			wantVoxspheres: []models.Voxsphere{
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupPostgres(t, tt.fixtureFiles...)
			pgrepo := voxrepo.NewRepo(db)

			gotErr := pgrepo.DeleteVoxsphere(context.Background(), tt.args.ID)

			assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")

			gotVoxspheres, err := pgrepo.Voxspheres(context.Background())

			assert.NoError(t, err, "expect no error while getting voxsphere")
			assert.Equal(t, tt.wantVoxspheres, gotVoxspheres, "expect voxspheres to match")
		})
	}
}

func TestRepo_ForeignKeyCascade(t *testing.T) {
	t.Run("on deleting voxsphere from parent table , no child references should exist in voxspheres table", func(t *testing.T) {
		db := setupPostgres(t, "topics.yml", "voxspheres.yml")
		topicPgrepo := topic.NewRepo(db)
		voxspherePgrepo := voxrepo.NewRepo(db)

		wantVoxspheres := []models.Voxsphere{
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
		}

		err := topicPgrepo.DeleteTopic(context.Background(), uuid.MustParse("00000000-0000-0000-0000-000000000001"))

		assert.NoError(t, err, "expect no error while deleting topic")

		gotVoxspheres, err := voxspherePgrepo.Voxspheres(context.Background())

		assert.NoError(t, err, "expect no error while getting rules")
		assertVoxspheresWithoutTimestamp(t, wantVoxspheres, gotVoxspheres)
	})
}
