package post_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/glowfi/voxpopuli/backend/pkg/models"
	tr "github.com/glowfi/voxpopuli/backend/pkg/transport"
	posttransport "github.com/glowfi/voxpopuli/backend/pkg/transport/post"
	"github.com/glowfi/voxpopuli/backend/pkg/transport/post/postfakes"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestTransport_PostsPaginated(t *testing.T) {
	type mockReturns struct {
		posts     []models.Post
		postError error
	}

	tests := []struct {
		name           string
		url            string
		mockReturns    mockReturns
		wantStatusCode int
		wantResponse   string
	}{
		{
			name:           "no query paramters :NEG",
			url:            fmt.Sprintf("/posts"),
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "paginated posts skip 3 limit 2 :POS",
			url:  fmt.Sprintf("/posts?skip=%v&limit=%v", 3, 2),
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
			wantStatusCode: http.StatusOK,
			wantResponse: `
                [
                  {
                    "ID": "00000000-0000-0000-0000-000000000004",
                    "AuthorID": "00000000-0000-0000-0000-000000000001",
                    "VoxsphereID": "00000000-0000-0000-0000-000000000001",
                    "Title": "Example Post Title 4",
                    "Text": "This is an example post text 4.",
                    "TextHtml": "This is an example post text 4 in HTML.",
                    "Ups": 40,
                    "Over18": true,
                    "Spoiler": false,
                    "CreatedAt": "2024-10-10T10:10:40Z",
                    "CreatedAtUnix": 1725091160,
                    "UpdatedAt": "2024-10-10T10:10:40Z"
                  },
                  {
                    "ID": "00000000-0000-0000-0000-000000000005",
                    "AuthorID": "00000000-0000-0000-0000-000000000002",
                    "VoxsphereID": "00000000-0000-0000-0000-000000000002",
                    "Title": "Example Post Title 5",
                    "Text": "This is an example post text 5.",
                    "TextHtml": "This is an example post text 5 in HTML.",
                    "Ups": 50,
                    "Over18": false,
                    "Spoiler": true,
                    "CreatedAt": "2024-10-10T10:10:50Z",
                    "CreatedAtUnix": 1725091180,
                    "UpdatedAt": "2024-10-10T10:10:50Z"
                  }
                ]
            `,
		},
		{
			name: "paginated posts skip 99 limit 99 :POS",
			url:  fmt.Sprintf("/posts?skip=%v&limit=%v", 3, 2),
			mockReturns: mockReturns{
				posts:     []models.Post{},
				postError: nil,
			},
			wantStatusCode: http.StatusOK,
			wantResponse: `
                [
                ]
            `,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakePostService := postfakes.FakePostService{}
			fakePostService.PostsPaginatedReturns(tt.mockReturns.posts, tt.mockReturns.postError)
			transport := posttransport.NewTransport(&fakePostService)

			server, err := tr.NewServer(tr.Services{
				Post: &fakePostService,
			}, tr.Transports{
				Post: transport,
			})
			if err != nil {
				t.Fatalf("error setting up server: %+v", err)
			}

			handler, err := server.HTTPHandler(context.Background())
			if err != nil {
				t.Fatalf("error setting up http handler: %+v", err)
			}

			request := httptest.NewRequest(
				"GET",
				tt.url,
				nil,
			)
			recorder := httptest.NewRecorder()
			handler.ServeHTTP(recorder, request)

			assert.Equal(
				t,
				tt.wantStatusCode,
				recorder.Result().StatusCode,
				"expect status code to match",
			)
		})
	}
}
