package service

import (
	"context"
	"testing"
)

func TestListReviews(t *testing.T) {
	svc := &SocialService{}
	reviews, err := svc.ListReviews(context.Background(), "spot1")
	if err != nil {
		t.Fatalf("ListReviews returned error: %v", err)
	}
	if reviews == nil {
		t.Error("expected non-nil slice")
	}
}

func TestGetUserCatches(t *testing.T) {
	svc := &SocialService{}
	catches, err := svc.GetUserCatches(context.Background(), "user1")
	if err != nil {
		t.Fatalf("GetUserCatches returned error: %v", err)
	}
	if catches == nil {
		t.Error("expected non-nil slice")
	}
}
