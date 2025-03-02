package rule_test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"slices"
	"testing"

	"github.com/glowfi/voxpopuli/backend/pkg/models"
	"github.com/glowfi/voxpopuli/backend/pkg/repo/rule"
	rulerrepo "github.com/glowfi/voxpopuli/backend/pkg/repo/rule"
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

	db.RegisterModel((*models.Voxsphere)(nil))
	db.RegisterModel((*models.Rule)(nil))

	// drop all rows of the rules,voxsphere table
	if _, err := db.NewTruncateTable().Cascade().Model((*models.Voxsphere)(nil)).Exec(context.Background()); err != nil {
		t.Fatal("truncate table failed:", err)
	}
	if _, err := db.NewTruncateTable().Cascade().Model((*models.Rule)(nil)).Exec(context.Background()); err != nil {
		t.Fatal("truncate table failed:", err)
	}

	// load fixture
	fixture := dbfixture.New(db)
	if err := fixture.Load(context.Background(), os.DirFS("testdata"), fixtureFiles...); err != nil {
		t.Fatal("failed to load fixtures", err)
	}

	return db
}

func assertRules(t *testing.T, wantRules, gotRules []models.Rule) {
	t.Helper()

	if len(wantRules) != len(gotRules) {
		t.Fatal("length of wantRules and gotRules do not match")
	}

	for _, rule := range wantRules {
		idx := slices.IndexFunc(gotRules, func(v models.Rule) bool {
			return v.ID == rule.ID
		})

		if idx == -1 {
			t.Fatal(fmt.Sprintf("rule %v of ID %v is not present in gotRules", rule.ShortName, rule.ID))
			return
		}
		assert.Equal(t, rule, gotRules[idx], "expect rule to match")
	}
}

func TestRepo_Rules(t *testing.T) {
	type args struct{}
	tests := []struct {
		name         string
		fixtureFiles []string
		args         args
		wantRules    []models.Rule
		wantErr      error
	}{
		{
			name:         "no rules :POS",
			fixtureFiles: []string{},
			args:         args{},
			wantRules:    []models.Rule{},
			wantErr:      nil,
		},
		{
			name:         "rules :POS",
			fixtureFiles: []string{"voxspheres.yml", "rules.yml"},
			args:         args{},
			wantRules: []models.Rule{
				{
					ID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					VoxsphereID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					ShortName:   "rule_foo",
					Description: "description_foo",
				},
				{
					ID:          uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					VoxsphereID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					ShortName:   "rule_bar",
					Description: "description_bar",
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupPostgres(t, tt.fixtureFiles...)
			pgrepo := rulerrepo.NewRepo(db)

			gotRules, gotErr := pgrepo.Rules(context.Background())
			assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")
			assertRules(t, tt.wantRules, gotRules)
		})
	}
}

func TestRepo_RuleByID(t *testing.T) {
	type args struct {
		ID uuid.UUID
	}
	tests := []struct {
		name         string
		fixtureFiles []string
		args         args
		wantRule     models.Rule
		wantErr      error
	}{
		{
			name:         "rule not found :NEG",
			fixtureFiles: []string{"voxspheres.yml", "rules.yml"},
			args: args{
				ID: uuid.MustParse("00000000-0000-0000-0000-000000000006"),
			},
			wantRule: models.Rule{},
			wantErr:  rulerrepo.ErrRuleNotFound,
		},
		{
			name:         "get rule by id :POS",
			fixtureFiles: []string{"voxspheres.yml", "rules.yml"},
			args: args{
				ID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			},
			wantRule: models.Rule{
				ID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				VoxsphereID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				ShortName:   "rule_foo",
				Description: "description_foo",
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupPostgres(t, tt.fixtureFiles...)
			pgrepo := rulerrepo.NewRepo(db)

			gotRule, gotErr := pgrepo.RuleByID(context.Background(), tt.args.ID)
			assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")
			assert.Equal(t, tt.wantRule, gotRule, "expect rule to match")
		})
	}
}

func TestRepo_UpdateRule(t *testing.T) {
	type args struct {
		rule models.Rule
	}
	tests := []struct {
		name         string
		fixtureFiles []string
		args         args
		wantRule     models.Rule
		wantRules    []models.Rule
		wantErr      error
	}{
		{
			name:         "rule not found :NEG",
			fixtureFiles: []string{"voxspheres.yml", "rules.yml"},
			args: args{
				rule: models.Rule{
					ID:          uuid.MustParse("00000000-0000-0000-0000-000000000006"),
					VoxsphereID: uuid.MustParse("00000000-0000-0000-0000-000000000006"),
					ShortName:   "updated rule",
					Description: "updated description",
				},
			},
			wantRule: models.Rule{},
			wantRules: []models.Rule{
				{
					ID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					VoxsphereID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					ShortName:   "rule_foo",
					Description: "description_foo",
				},
				{
					ID:          uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					VoxsphereID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					ShortName:   "rule_bar",
					Description: "description_bar",
				},
			},
			wantErr: rulerrepo.ErrRuleNotFound,
		},
		{
			name:         "update rule :POS",
			fixtureFiles: []string{"voxspheres.yml", "rules.yml"},
			args: args{
				rule: models.Rule{
					ID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					VoxsphereID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					ShortName:   "updated rule",
					Description: "updated description",
				},
			},
			wantRule: models.Rule{
				ID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				VoxsphereID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				ShortName:   "updated rule",
				Description: "updated description",
			},
			wantRules: []models.Rule{
				{
					ID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					VoxsphereID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					ShortName:   "updated rule",
					Description: "updated description",
				},
				{
					ID:          uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					VoxsphereID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					ShortName:   "rule_bar",
					Description: "description_bar",
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupPostgres(t, tt.fixtureFiles...)
			pgrepo := rulerrepo.NewRepo(db)

			gotRule, gotErr := pgrepo.UpdateRule(context.Background(), tt.args.rule)

			assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")
			assert.Equal(t, tt.wantRule, gotRule, "expect rule to match")

			gotRules, err := pgrepo.Rules(context.Background())

			assert.NoError(t, err, "expect no error while getting rules")
			assertRules(t, tt.wantRules, gotRules)
		})
	}
}

func TestRepo_AddRules(t *testing.T) {
	type args struct {
		rules []models.Rule
	}
	tests := []struct {
		name              string
		fixtureFiles      []string
		args              args
		wantInsertedRules []models.Rule
		wantRules         []models.Rule
		wantErr           error
	}{
		{
			name:         "duplicate rule id :POS",
			fixtureFiles: []string{"voxspheres.yml", "rules.yml"},
			args: args{
				rules: []models.Rule{
					{
						ID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
						VoxsphereID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
						ShortName:   "new rule",
						Description: "new description",
					},
				},
			},
			wantInsertedRules: nil,
			wantRules: []models.Rule{
				{
					ID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					VoxsphereID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					ShortName:   "rule_foo",
					Description: "description_foo",
				},
				{
					ID:          uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					VoxsphereID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					ShortName:   "rule_bar",
					Description: "description_bar",
				},
			},
			wantErr: rulerrepo.ErrRuleDuplicateIDorName,
		},
		{
			name:         "add rules :POS",
			fixtureFiles: []string{"voxspheres.yml", "rules.yml"},
			args: args{
				rules: []models.Rule{
					{
						ID:          uuid.MustParse("00000000-0000-0000-0000-000000000003"),
						VoxsphereID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
						ShortName:   "new rule",
						Description: "new description",
					},
				},
			},
			wantInsertedRules: []models.Rule{
				{
					ID:          uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					VoxsphereID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					ShortName:   "new rule",
					Description: "new description",
				},
			},
			wantRules: []models.Rule{
				{
					ID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					VoxsphereID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					ShortName:   "rule_foo",
					Description: "description_foo",
				},
				{
					ID:          uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					VoxsphereID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					ShortName:   "rule_bar",
					Description: "description_bar",
				},
				{
					ID:          uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					VoxsphereID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					ShortName:   "new rule",
					Description: "new description",
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupPostgres(t, tt.fixtureFiles...)
			pgrepo := rulerrepo.NewRepo(db)

			gotInsertedRules, gotErr := pgrepo.AddRules(context.Background(), tt.args.rules...)

			assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")
			assert.Equal(t, tt.wantInsertedRules, gotInsertedRules, "expect inserted rules to match")

			gotRules, err := pgrepo.Rules(context.Background())

			assert.NoError(t, err, "expect no error while getting rules")
			assertRules(t, tt.wantRules, gotRules)
		})
	}
}

func TestRepo_DeleteRule(t *testing.T) {
	type args struct {
		ID uuid.UUID
	}
	tests := []struct {
		name         string
		fixtureFiles []string
		args         args
		wantRules    []models.Rule
		wantErr      error
	}{
		{
			name:         "rule not found :NEG",
			fixtureFiles: []string{"voxspheres.yml", "rules.yml"},
			args: args{
				ID: uuid.MustParse("00000000-0000-0000-0000-000000000006"),
			},
			wantRules: []models.Rule{
				{
					ID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					VoxsphereID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					ShortName:   "rule_foo",
					Description: "description_foo",
				},
				{
					ID:          uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					VoxsphereID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					ShortName:   "rule_bar",
					Description: "description_bar",
				},
			},
			wantErr: rulerrepo.ErrRuleNotFound,
		},
		{
			name:         "delete rule :POS",
			fixtureFiles: []string{"voxspheres.yml", "rules.yml"},
			args: args{
				ID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			},
			wantRules: []models.Rule{
				{
					ID:          uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					VoxsphereID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					ShortName:   "rule_bar",
					Description: "description_bar",
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupPostgres(t, tt.fixtureFiles...)
			pgrepo := rulerrepo.NewRepo(db)

			gotErr := pgrepo.DeleteRule(context.Background(), tt.args.ID)

			assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")

			gotRules, err := pgrepo.Rules(context.Background())

			assert.NoError(t, err, "expect no error while getting rules")
			assertRules(t, tt.wantRules, gotRules)
		})
	}
}

func TestRepo_ForeignKeyCascade(t *testing.T) {
	t.Run("on deleting voxsphere from parent table , no child references should exist in rules table", func(t *testing.T) {
		db := setupPostgres(t, "voxspheres.yml", "rules.yml")
		rulePgrepo := rule.NewRepo(db)
		voxspherePgrepo := voxrepo.NewRepo(db)

		wantRules := []models.Rule{
			{
				ID:          uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				VoxsphereID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				ShortName:   "rule_bar",
				Description: "description_bar",
			},
		}

		err := voxspherePgrepo.DeleteVoxsphere(context.Background(), uuid.MustParse("00000000-0000-0000-0000-000000000001"))

		assert.NoError(t, err, "expect no error while deleting voxsphere")

		gotRules, err := rulePgrepo.Rules(context.Background())

		assert.NoError(t, err, "expect no error while getting rules")
		assertRules(t, wantRules, gotRules)
	})
}
