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

func (s *SpotService) CreateSpot(ctx context.Context, req model.CreateSpotRequest) (*model.Spot, error) {
	if req.Name == "" {
		req.Name = "New Spot"
	}
	if req.Type == "" {
		req.Type = "lake"
	}
	if req.Difficulty == "" {
		req.Difficulty = "easy"
	}

	id, err := s.repo.Create(ctx, req.Name, req.Lat, req.Lon, req.Type, req.Difficulty)
	if err != nil {
		return nil, err
	}

	return &model.Spot{
		ID:         id,
		Name:       req.Name,
		Lat:        req.Lat,
		Lon:        req.Lon,
		Type:       req.Type,
		Difficulty: req.Difficulty,
	}, nil
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
	rows, err := s.repo.Nearby(ctx, params.Lat, params.Lon, params.RadiusMi, params.Limit)
	if err != nil {
		return nil, err
	}
	spots := make([]model.Spot, 0, len(rows))
	for _, r := range rows {
		spots = append(spots, model.Spot{
			ID:     r.ID,
			Name:   r.Name,
			Lat:    r.Lat,
			Lon:    r.Lon,
			Type:   r.Type,
			Rating: r.Rating,
		})
	}
	return spots, nil
}
