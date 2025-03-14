package post

import (
	"context"

	"github.com/glowfi/voxpopuli/backend/pkg/models"
)

type PostService interface {
	PostsPaginated(ctx context.Context, skip, limit int) ([]models.PostPaginated, error)
}

//counterfeiter:generate . PostRepository
type PostRepository interface {
	PostsPaginated(ctx context.Context, skip, limit int) ([]models.PostPaginated, error)
}

type Service struct {
	repo PostRepository
}

func NewService(repo PostRepository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) PostsPaginated(ctx context.Context, skip, limit int) ([]models.PostPaginated, error) {
	return s.repo.PostsPaginated(ctx, skip, limit)
}
