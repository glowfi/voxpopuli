package post

import (
	"context"

	"github.com/glowfi/voxpopuli/backend/pkg/models"
)

type PostService interface {
	PostsPaginated(ctx context.Context, skip int, limit int) ([]models.Post, error)
}

//counterfeiter:generate . PostRepository
type PostRepository interface {
	PostsPaginated(ctx context.Context, skip int, limit int) ([]models.Post, error)
}

type Service struct {
	repo PostRepository
}

func NewService(repo PostRepository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) PostsPaginated(ctx context.Context, skip, limit int) ([]models.Post, error) {
	return s.repo.PostsPaginated(ctx, skip, limit)
}
