package service

import (
	"context"

	"fishwish/services/spot-service/internal/model"
	"fishwish/services/spot-service/internal/repository"
)

type SpotService struct {
	repo *repository.SpotRepository
}

func NewSpotService(db *repository.DB) *SpotService {
	return &SpotService{
		repo: repository.NewSpotRepository(db),
	}
}

func (s *SpotService) ListSpots(ctx context.Context, params model.ListSpotsParams) ([]model.Spot, error) {
	return []model.Spot{}, nil
}

func (s *SpotService) GetSpot(ctx context.Context, id string) (*model.Spot, error) {
	spot, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &model.Spot{
		ID:          spot.ID,
		Name:        spot.Name,
		Lat:         spot.Lat,
		Lon:         spot.Lon,
		Type:        spot.Type,
		AccessNotes: spot.AccessNotes,
		Difficulty:  spot.Difficulty,
		Rating:      spot.Rating,
		ReviewCount: spot.ReviewCount,
	}, nil
}

func (s *SpotService) GetSpotDetails(ctx context.Context, id string) (*model.SpotDetail, error) {
	return nil, nil
}

func (s *SpotService) NearbySpots(ctx context.Context, params model.NearbyParams) ([]model.Spot, error) {
	return []model.Spot{}, nil
}
