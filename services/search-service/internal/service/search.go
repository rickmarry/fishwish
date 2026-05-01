package service

import (
	"context"

	"fishwish/services/search-service/internal/model"
	"fishwish/services/search-service/internal/repository"
)

type SearchService struct {
	repo *repository.SearchRepository
}

func NewSearchService(db *repository.DB) *SearchService {
	return &SearchService{
		repo: repository.NewSearchRepository(db),
	}
}

func (s *SearchService) Search(ctx context.Context, params model.SearchParams) ([]model.SearchResult, error) {
	return []model.SearchResult{}, nil
}

func (s *SearchService) GetSpecies(ctx context.Context) ([]model.SpeciesResult, error) {
	return []model.SpeciesResult{}, nil
}

func (s *SearchService) GetSuggestions(ctx context.Context, prefix string) ([]string, error) {
	return s.repo.GetSuggestions(ctx, prefix)
}
