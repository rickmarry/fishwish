package service

import (
	"context"
	"testing"

	"fishwish/services/search-service/internal/model"
)

type mockSearchRepo struct{}

func (m *mockSearchRepo) Search(ctx context.Context, query string, limit int) ([]struct {
	ID      string
	Name    string
	Type    string
	Lat     float64
	Lon     float64
	Rating  float64
	Species []string
}, error) {
	return []struct {
		ID      string
		Name    string
		Type    string
		Lat     float64
		Lon     float64
		Rating  float64
		Species []string
	}{
		{ID: "1", Name: "Tahoe", Type: "lake", Lat: 39.1, Lon: -120.1, Rating: 4.5, Species: []string{"trout"}},
	}, nil
}

func (m *mockSearchRepo) GetSpeciesList(ctx context.Context) ([]struct {
	Name       string
	CommonName string
	SpotCount  int
}, error) {
	return []struct {
		Name       string
		CommonName string
		SpotCount  int
	}{
		{Name: "trout", CommonName: "Rainbow Trout", SpotCount: 5},
	}, nil
}

func (m *mockSearchRepo) GetSuggestions(ctx context.Context, prefix string) ([]string, error) {
	return []string{"Tahoe", "Tahoe River"}, nil
}

func TestSearch(t *testing.T) {
	svc := &SearchService{repo: &mockSearchRepo{}}

	results, err := svc.Search(context.Background(), model.SearchParams{Query: "Tah", Limit: 10})
	if err != nil {
		t.Fatalf("Search returned error: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("expected 1 result, got %d", len(results))
	}

	if results[0].Name != "Tahoe" {
		t.Errorf("expected Tahoe, got %s", results[0].Name)
	}
}

func TestGetSuggestions(t *testing.T) {
	svc := &SearchService{repo: &mockSearchRepo{}}

	suggestions, err := svc.GetSuggestions(context.Background(), "Ta")
	if err != nil {
		t.Fatalf("GetSuggestions returned error: %v", err)
	}

	if len(suggestions) != 2 {
		t.Errorf("expected 2 suggestions, got %d", len(suggestions))
	}
}
