package post_test

import (
	"context"
	"testing"
	"time"

	"github.com/glowfi/voxpopuli/backend/pkg/models"
	postservice "github.com/glowfi/voxpopuli/backend/pkg/service/post"
	"github.com/glowfi/voxpopuli/backend/pkg/service/post/postfakes"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestService_PostsPaginated(t *testing.T) {
	type args struct {
		skip  int
		limit int
	}
	type mockReturns struct {
		posts     []models.Post
		postError error
	}

	tests := []struct {
		name        string
		args        args
		mockReturns mockReturns
		wantPosts   []models.Post
		wantErr     error
	}{
		{
			name: "paginated posts skip 3 limit 2 :POS",
			args: args{
				skip:  3,
				limit: 2,
			},
			mockReturns: mockReturns{
				posts: []models.Post{
					{
						ID:            uuid.MustParse("00000000-0000-0000-0000-000000000004"),
						AuthorID:      uuid.MustParse("00000000-0000-0000-0000-000000000001"),
						VoxsphereID:   uuid.MustParse("00000000-0000-0000-0000-000000000001"),
						Title:         "Example Post Title 4",
						Text:          "This is an example post text 4.",
						TextHtml:      "This is an example post text 4 in HTML.",
						Ups:           40,
						Over18:        true,
						Spoiler:       false,
						CreatedAt:     time.Date(2024, 10, 10, 10, 10, 40, 0, time.UTC),
						CreatedAtUnix: 1725091160,
						UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 40, 0, time.UTC),
					},
					{
						ID:            uuid.MustParse("00000000-0000-0000-0000-000000000005"),
						AuthorID:      uuid.MustParse("00000000-0000-0000-0000-000000000002"),
						VoxsphereID:   uuid.MustParse("00000000-0000-0000-0000-000000000002"),
						Title:         "Example Post Title 5",
						Text:          "This is an example post text 5.",
						TextHtml:      "This is an example post text 5 in HTML.",
						Ups:           50,
						Over18:        false,
						Spoiler:       true,
						CreatedAt:     time.Date(2024, 10, 10, 10, 10, 50, 0, time.UTC),
						CreatedAtUnix: 1725091180,
						UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 50, 0, time.UTC),
					},
				},
				postError: nil,
			},
			wantPosts: []models.Post{
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000004"),
					AuthorID:      uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					VoxsphereID:   uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Title:         "Example Post Title 4",
					Text:          "This is an example post text 4.",
					TextHtml:      "This is an example post text 4 in HTML.",
					Ups:           40,
					Over18:        true,
					Spoiler:       false,
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 40, 0, time.UTC),
					CreatedAtUnix: 1725091160,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 40, 0, time.UTC),
				},
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000005"),
					AuthorID:      uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					VoxsphereID:   uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					Title:         "Example Post Title 5",
					Text:          "This is an example post text 5.",
					TextHtml:      "This is an example post text 5 in HTML.",
					Ups:           50,
					Over18:        false,
					Spoiler:       true,
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 50, 0, time.UTC),
					CreatedAtUnix: 1725091180,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 50, 0, time.UTC),
				},
			},
			wantErr: nil,
		},
		{
			name: "paginated posts skip 3 limit 1 :POS",
			args: args{
				skip:  3,
				limit: 1,
			},
			mockReturns: mockReturns{
				posts: []models.Post{
					{
						ID:            uuid.MustParse("00000000-0000-0000-0000-000000000004"),
						AuthorID:      uuid.MustParse("00000000-0000-0000-0000-000000000001"),
						VoxsphereID:   uuid.MustParse("00000000-0000-0000-0000-000000000001"),
						Title:         "Example Post Title 4",
						Text:          "This is an example post text 4.",
						TextHtml:      "This is an example post text 4 in HTML.",
						Ups:           40,
						Over18:        true,
						Spoiler:       false,
						CreatedAt:     time.Date(2024, 10, 10, 10, 10, 40, 0, time.UTC),
						CreatedAtUnix: 1725091160,
						UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 40, 0, time.UTC),
					},
				},
				postError: nil,
			},
			wantPosts: []models.Post{
				{
					ID:            uuid.MustParse("00000000-0000-0000-0000-000000000004"),
					AuthorID:      uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					VoxsphereID:   uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Title:         "Example Post Title 4",
					Text:          "This is an example post text 4.",
					TextHtml:      "This is an example post text 4 in HTML.",
					Ups:           40,
					Over18:        true,
					Spoiler:       false,
					CreatedAt:     time.Date(2024, 10, 10, 10, 10, 40, 0, time.UTC),
					CreatedAtUnix: 1725091160,
					UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 40, 0, time.UTC),
				},
			},
			wantErr: nil,
		},
		{
			name: "paginated posts skip 100 limit 100 :POS",
			mockReturns: mockReturns{
				posts:     nil,
				postError: nil,
			},
			args: args{
				skip:  100,
				limit: 100,
			},
			wantPosts: nil,
			wantErr:   nil,
		},
		{
			name: "no posts :POS",
			mockReturns: mockReturns{
				posts:     nil,
				postError: nil,
			},
			args: args{
				skip:  100,
				limit: 100,
			},
			wantPosts: nil,
			wantErr:   nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakePostRepo := postfakes.FakePostRepository{}
			fakePostRepo.PostsPaginatedReturns(tt.mockReturns.posts, tt.mockReturns.postError)
			service := postservice.NewService(&fakePostRepo)

			gotPosts, gotErr := service.PostsPaginated(context.Background(), tt.args.skip, tt.args.limit)
			assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")
			assert.Equal(t, tt.wantPosts, gotPosts, "expect posts to match")
		})
	}
}
