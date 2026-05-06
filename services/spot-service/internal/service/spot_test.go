package service

import (
	"context"
	"testing"

	"fishwish/services/spot-service/internal/model"
)

func TestListSpots(t *testing.T) {
	svc := &SpotService{}
	spots, err := svc.ListSpots(context.Background(), model.ListSpotsParams{})
	if err != nil {
		t.Fatalf("ListSpots returned error: %v", err)
	}
	if spots == nil {
		t.Error("expected non-nil slice")
	}
}

func TestNearbySpots(t *testing.T) {
	svc := &SpotService{}
	spots, err := svc.NearbySpots(context.Background(), model.NearbyParams{})
	if err != nil {
		t.Fatalf("NearbySpots returned error: %v", err)
	}
	if spots == nil {
		t.Error("expected non-nil slice")
	}
}
