package service

import (
	"context"

	"fishwish/services/search-service/internal/model"
	"fishwish/services/search-service/internal/repository"
)

type SearchService struct {
	repo repository.SearchRepositoryInterface
}

func NewSearchService(db *repository.DB) *SearchService {
	return &SearchService{
		repo: repository.NewSearchRepository(db),
	}
}

func (s *SearchService) Search(ctx context.Context, params model.SearchParams) ([]model.SearchResult, error) {
	repoResults, err := s.repo.Search(ctx, params.Query, params.Limit)
	if err != nil {
		return nil, err
	}

	results := make([]model.SearchResult, len(repoResults))
	for i, r := range repoResults {
		results[i] = model.SearchResult{
			ID:      r.ID,
			Name:    r.Name,
			Type:    r.Type,
			Lat:     r.Lat,
			Lon:     r.Lon,
			Rating:  r.Rating,
			Species: r.Species,
		}
	}

	return results, nil
}

func (s *SearchService) GetSpecies(ctx context.Context) ([]model.SpeciesResult, error) {
	repoResults, err := s.repo.GetSpeciesList(ctx)
	if err != nil {
		return nil, err
	}

	results := make([]model.SpeciesResult, len(repoResults))
	for i, r := range repoResults {
		results[i] = model.SpeciesResult{
			Name:       r.Name,
			CommonName: r.CommonName,
			SpotCount:  r.SpotCount,
		}
	}

	return results, nil
}

func (s *SearchService) GetSuggestions(ctx context.Context, prefix string) ([]string, error) {
	return s.repo.GetSuggestions(ctx, prefix)
}
