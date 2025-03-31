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
	"github.com/glowfi/voxpopuli/backend/pkg/transport/post/postfakes"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestTransport_PostsPaginated(t *testing.T) {
	type mockReturns struct {
		posts     []models.PostPaginated
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
			name:           "no query parameters :NEG",
			url:            "/posts",
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "paginated posts skip 99 limit 99 :POS",
			url:  fmt.Sprintf("/posts?skip=%v&limit=%v", 99, 99),
			mockReturns: mockReturns{
				posts:     []models.PostPaginated{},
				postError: nil,
			},
			wantStatusCode: http.StatusOK,
			wantResponse: `
		                [
		                ]
		            `,
		},
		{
			name: "paginated posts skip 3 limit 2 :POS",
			url:  fmt.Sprintf("/posts?skip=%v&limit=%v", 3, 2),
			mockReturns: mockReturns{
				posts: []models.PostPaginated{
					{
						ID:          uuid.MustParse("00000000-0000-0000-0000-000000000004"),
						Author:      "John Doe",
						AuthorID:    uuid.MustParse("00000000-0000-0000-0000-000000000001"),
						Voxsphere:   "v/foo",
						VoxsphereID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
						Title:       "Example Post Title 4",
						Text:        "This is an example post text 4.",
						TextHtml:    "This is an example post text 4 in HTML.",
						MediaType:   models.MediaTypeVideo,
						Medias: []any{
							models.Video{
								ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
								MediaID:       uuid.MustParse("00000000-0000-0000-0000-000000000004"),
								Url:           "https://example.com/video.mp4",
								Height:        1080,
								Width:         1920,
								CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
								CreatedAtUnix: 1725091100,
								UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
							},
						},
						Ups:           40,
						NumComments:   10,
						NumAwards:     10,
						Over18:        true,
						Spoiler:       false,
						CreatedAt:     time.Date(2024, 10, 10, 10, 10, 40, 0, time.UTC),
						CreatedAtUnix: 1725091160,
						UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 40, 0, time.UTC),
					},
					{
						ID:          uuid.MustParse("00000000-0000-0000-0000-000000000005"),
						Author:      "Jane Doe",
						AuthorID:    uuid.MustParse("00000000-0000-0000-0000-000000000002"),
						Voxsphere:   "v/bar",
						VoxsphereID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
						Title:       "Example Post Title 5",
						Text:        "This is an example post text 5.",
						TextHtml:    "This is an example post text 5 in HTML.",
						MediaType:   models.MediaTypeLink,
						Medias: []any{
							models.Link{
								ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
								MediaID:       uuid.MustParse("00000000-0000-0000-0000-000000000005"),
								Link:          "https://example.com/video.mp4",
								CreatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
								CreatedAtUnix: 1725091100,
								UpdatedAt:     time.Date(2024, 10, 10, 10, 10, 10, 0, time.UTC),
							},
						},
						NumComments:   20,
						NumAwards:     10,
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
                    "id": "00000000-0000-0000-0000-000000000004",
                    "author":"John Doe",
                    "author_id": "00000000-0000-0000-0000-000000000001",
                    "voxsphere":"v/foo",
                    "voxsphere_id": "00000000-0000-0000-0000-000000000001",
                    "title": "Example Post Title 4",
                    "text": "This is an example post text 4.",
                    "text_html": "This is an example post text 4 in HTML.",
                    "media_type": "video",
                    "medias": [
                      {
                        "id": "00000000-0000-0000-0000-000000000001",
                        "media_id": "00000000-0000-0000-0000-000000000004",
                        "url": "https://example.com/video.mp4",
                        "height": 1080,
                        "width": 1920,
                        "created_at": "2024-10-10T10:10:10Z",
                        "created_at_unix": 1725091100,
                        "updated_at": "2024-10-10T10:10:10Z"
                      }
                    ],
                    "ups": 40,
                    "num_comments":10,
                    "num_awards":10,
                    "over18": true,
                    "spoiler": false,
                    "created_at": "2024-10-10T10:10:40Z",
                    "created_at_unix": 1725091160,
                    "updated_at": "2024-10-10T10:10:40Z"
                  },
                  {
                    "id": "00000000-0000-0000-0000-000000000005",
                    "author":"Jane Doe",
                    "author_id": "00000000-0000-0000-0000-000000000002",
                    "voxsphere":"v/bar",
                    "voxsphere_id": "00000000-0000-0000-0000-000000000002",
                    "title": "Example Post Title 5",
                    "text": "This is an example post text 5.",
                    "text_html": "This is an example post text 5 in HTML.",
                    "media_type": "link",
                    "medias": [
                      {
                        "id": "00000000-0000-0000-0000-000000000001",
                        "media_id": "00000000-0000-0000-0000-000000000005",
                        "link": "https://example.com/video.mp4",
                        "image": null,
                        "created_at": "2024-10-10T10:10:10Z",
                        "created_at_unix": 1725091100,
                        "updated_at": "2024-10-10T10:10:10Z"
                      }
                    ],
                    "ups": 50,
                    "num_comments":20,
                    "num_awards":10,
                    "over18": false,
                    "spoiler": true,
                    "created_at": "2024-10-10T10:10:50Z",
                    "created_at_unix": 1725091180,
                    "updated_at": "2024-10-10T10:10:50Z"
                  }
                ]
		            `,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakePostService := postfakes.FakePostService{}
			fakePostService.PostsPaginatedReturns(tt.mockReturns.posts, tt.mockReturns.postError)

			server, err := tr.NewServer(tr.Services{
				Post: &fakePostService,
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

			if tt.wantStatusCode == http.StatusOK {
				assert.JSONEq(t, tt.wantResponse, recorder.Body.String())
			}
		})
	}
}
