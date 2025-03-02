package topic_test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"slices"
	"testing"

	"github.com/glowfi/voxpopuli/backend/pkg/models"
	topicrepo "github.com/glowfi/voxpopuli/backend/pkg/repo/topic"
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

	// drop all rows of the topics table
	_, err := db.NewTruncateTable().Cascade().Model((*models.Topic)(nil)).Exec(context.Background())
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

func assertTopics(t *testing.T, wantTopics, gotTopics []models.Topic) {
	t.Helper()

	if len(wantTopics) != len(gotTopics) {
		t.Fatal("length of wantTopics and gotTopics do not match")
	}

	for _, topic := range wantTopics {
		idx := slices.IndexFunc(gotTopics, func(v models.Topic) bool {
			return v.ID == topic.ID
		})

		if idx == -1 {
			t.Fatal(fmt.Sprintf("topic %v of ID %v is not present in gotTopics", topic.Name, topic.ID))
			return
		}
		assert.Equal(t, topic, gotTopics[idx], "expect topic to match")
	}
}

func TestRepo_Topics(t *testing.T) {
	tests := []struct {
		name         string
		fixtureFiles []string
		wantTopics   []models.Topic
		wantErr      error
	}{
		{
			name:         "topics :POS",
			fixtureFiles: []string{"topics.yml"},
			wantTopics: []models.Topic{
				{
					ID:       uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Name:     "xyz",
					Category: "foo",
				},
				{
					ID:       uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					Name:     "pqr",
					Category: "bar",
				},
			},
			wantErr: nil,
		},
		{
			name:         "no topics :POS",
			fixtureFiles: []string{},
			wantTopics:   nil,
			wantErr:      nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupPostgres(t, tt.fixtureFiles...)
			pgrepo := topicrepo.NewRepo(db)

			gotTopics, gotErr := pgrepo.Topics(context.Background())

			assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")
			assert.Equal(t, tt.wantTopics, gotTopics, "expect topics to match")
		})
	}
}

func TestRepo_TopicByID(t *testing.T) {
	type args struct {
		ID uuid.UUID
	}
	tests := []struct {
		name         string
		fixtureFiles []string
		args         args
		wantTopic    models.Topic
		wantErr      error
	}{
		{
			name:         "topic not found :NEG",
			fixtureFiles: []string{},
			args: args{
				ID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			},
			wantTopic: models.Topic{},
			wantErr:   topicrepo.ErrTopicNotFound,
		},
		{
			name:         "topic found :POS",
			fixtureFiles: []string{"topics.yml"},
			args: args{
				ID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			},
			wantTopic: models.Topic{
				ID:       uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				Name:     "xyz",
				Category: "foo",
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupPostgres(t, tt.fixtureFiles...)
			pgrepo := topicrepo.NewRepo(db)

			gotTopic, gotErr := pgrepo.TopicByID(context.Background(), tt.args.ID)

			assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")
			assert.Equal(t, tt.wantTopic, gotTopic, "expect topic to match")
		})
	}
}

func TestRepo_AddTopics(t *testing.T) {
	type args struct {
		topics []models.Topic
	}
	tests := []struct {
		name               string
		fixtureFiles       []string
		args               args
		wantInsertedTopics []models.Topic
		wantTopics         []models.Topic
		wantErr            error
	}{
		{
			name:         "duplicate topic name :NEG",
			fixtureFiles: []string{"topics.yml"},
			args: args{
				topics: []models.Topic{
					{
						ID:       uuid.MustParse("00000000-0000-0000-0000-000000000003"),
						Name:     "xyz",
						Category: "new Category",
					},
				},
			},
			wantInsertedTopics: nil,
			wantTopics: []models.Topic{
				{
					ID:       uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Name:     "xyz",
					Category: "foo",
				},
				{
					ID:       uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					Name:     "pqr",
					Category: "bar",
				},
			},
			wantErr: topicrepo.ErrTopicDuplicateIDorName,
		},
		{
			name:         "add topics :POS",
			fixtureFiles: []string{"topics.yml"},
			args: args{
				topics: []models.Topic{
					{
						ID:       uuid.MustParse("00000000-0000-0000-0000-000000000003"),
						Name:     "new topic",
						Category: "new Category",
					},
				},
			},
			wantInsertedTopics: []models.Topic{
				{
					ID:       uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					Name:     "new topic",
					Category: "new Category",
				},
			},
			wantTopics: []models.Topic{
				{
					ID:       uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Name:     "xyz",
					Category: "foo",
				},
				{
					ID:       uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					Name:     "pqr",
					Category: "bar",
				},
				{
					ID:       uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					Name:     "new topic",
					Category: "new Category",
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupPostgres(t, tt.fixtureFiles...)
			pgrepo := topicrepo.NewRepo(db)

			gotInsertedTopic, gotErr := pgrepo.AddTopics(context.Background(), tt.args.topics...)
			assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")
			assert.Equal(t, tt.wantInsertedTopics, gotInsertedTopic, "expect inserted topics to match")

			gotTopics, err := pgrepo.Topics(context.Background())

			assert.NoError(t, err, "expect no error while getting topics")
			assertTopics(t, tt.wantTopics, gotTopics)
		})
	}
}

func TestRepo_UpdateTopic(t *testing.T) {
	type args struct {
		topic models.Topic
	}
	tests := []struct {
		name         string
		fixtureFiles []string
		args         args
		wantTopic    models.Topic
		wantTopics   []models.Topic
		wantErr      error
	}{
		{
			name:         "topic not found :NEG",
			fixtureFiles: []string{"topics.yml"},
			args: args{
				topic: models.Topic{
					ID:   uuid.MustParse("00000000-0000-0000-0000-000000000006"),
					Name: "updated topic",
				},
			},
			wantTopic: models.Topic{},
			wantTopics: []models.Topic{
				{
					ID:       uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Name:     "xyz",
					Category: "foo",
				},
				{
					ID:       uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					Name:     "pqr",
					Category: "bar",
				},
			},
			wantErr: topicrepo.ErrTopicNotFound,
		},
		{
			name:         "update topic :POS",
			fixtureFiles: []string{"topics.yml"},
			args: args{
				topic: models.Topic{
					ID:       uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Name:     "updated topic",
					Category: "updated category",
				},
			},
			wantTopic: models.Topic{
				ID:       uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				Name:     "updated topic",
				Category: "updated category",
			},
			wantTopics: []models.Topic{
				{
					ID:       uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Name:     "updated topic",
					Category: "updated category",
				},
				{
					ID:       uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					Name:     "pqr",
					Category: "bar",
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupPostgres(t, tt.fixtureFiles...)
			pgrepo := topicrepo.NewRepo(db)

			gotTopic, gotErr := pgrepo.UpdateTopic(context.Background(), tt.args.topic)

			assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")
			assert.Equal(t, tt.wantTopic, gotTopic, "expect topic to match")

			gotTopics, err := pgrepo.Topics(context.Background())

			assert.NoError(t, err, "expect no error while getting topics")
			assertTopics(t, tt.wantTopics, gotTopics)
		})
	}
}

func TestRepo_DeleteTopic(t *testing.T) {
	type args struct {
		ID uuid.UUID
	}
	tests := []struct {
		name         string
		fixtureFiles []string
		args         args
		wantTopics   []models.Topic
		wantErr      error
	}{
		{
			name:         "topic not found :NEG",
			fixtureFiles: []string{"topics.yml"},
			args: args{
				ID: uuid.MustParse("00000000-0000-0000-0000-000000000006"),
			},
			wantTopics: []models.Topic{
				{
					ID:       uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Name:     "xyz",
					Category: "foo",
				},
				{
					ID:       uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					Name:     "pqr",
					Category: "bar",
				},
			},
			wantErr: topicrepo.ErrTopicNotFound,
		},
		{
			name:         "delete topic :POS",
			fixtureFiles: []string{"topics.yml"},
			args: args{
				ID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			},
			wantTopics: []models.Topic{
				{
					ID:       uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					Name:     "pqr",
					Category: "bar",
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupPostgres(t, tt.fixtureFiles...)
			pgrepo := topicrepo.NewRepo(db)

			gotErr := pgrepo.DeleteTopic(context.Background(), tt.args.ID)

			assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")

			gotTopics, err := pgrepo.Topics(context.Background())

			assert.NoError(t, err, "expect no error while getting topics")

			assertTopics(t, tt.wantTopics, gotTopics)
		})
	}
}
