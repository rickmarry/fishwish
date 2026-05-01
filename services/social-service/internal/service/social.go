package service

import (
	"context"

	"fishwish/services/social-service/internal/model"
	"fishwish/services/social-service/internal/repository"

	"github.com/google/uuid"
)

type SocialService struct {
	repo *repository.SocialRepository
}

func NewSocialService(db *repository.DB) *SocialService {
	return &SocialService{
		repo: repository.NewSocialRepository(db),
	}
}

func (s *SocialService) CreateReview(ctx context.Context, userID, spotID string, req model.CreateReviewRequest) (*model.Review, error) {
	if req.Rating < 1 || req.Rating > 5 {
		return nil, nil
	}

	id := uuid.New().String()
	if err := s.repo.CreateReview(ctx, id, spotID, userID, req.Rating, req.Content); err != nil {
		return nil, err
	}

	return &model.Review{
		ID:      id,
		SpotID:  spotID,
		UserID:  userID,
		Rating:  req.Rating,
		Content: req.Content,
	}, nil
}

func (s *SocialService) ListReviews(ctx context.Context, spotID string) ([]model.Review, error) {
	return []model.Review{}, nil
}

func (s *SocialService) LogCatch(ctx context.Context, userID string, req model.CreateCatchRequest) (*model.CatchLog, error) {
	id := uuid.New().String()
	if err := s.repo.CreateCatch(ctx, id, userID, req.SpotID, req.Species, req.Weight, req.Length, req.BaitUsed); err != nil {
		return nil, err
	}

	return &model.CatchLog{
		ID:       id,
		UserID:   userID,
		SpotID:   req.SpotID,
		Species:  req.Species,
		Weight:   req.Weight,
		Length:   req.Length,
		BaitUsed: req.BaitUsed,
	}, nil
}

func (s *SocialService) GetUserCatches(ctx context.Context, userID string) ([]model.CatchLog, error) {
	return []model.CatchLog{}, nil
}
