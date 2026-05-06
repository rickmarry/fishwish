package service

import (
	"context"
	"errors"
	"testing"

	"fishwish/services/user-service/internal/model"
)

type mockUserRepo struct{}

func (m *mockUserRepo) GetByEmail(ctx context.Context, email string) (string, string, string, error) {
	if email == "exists@test.com" {
		return "1", "user", "hashed", nil
	}
	return "", "", "", errors.New("not found")
}

func (m *mockUserRepo) Create(ctx context.Context, id, email, username, hashedPassword string) error {
	return nil
}

func (m *mockUserRepo) GetByID(ctx context.Context, id string) (string, string, string, error) {
	return "test@test.com", "testuser", "", nil
}

func TestLogin(t *testing.T) {
	svc := &UserService{repo: &mockUserRepo{}}

	_, err := svc.Login(context.Background(), model.LoginRequest{Email: "exists@test.com", Password: "wrong"})
	if err == nil {
		t.Error("expected error for wrong password")
	}

	_, err = svc.Login(context.Background(), model.LoginRequest{Email: "new@test.com", Password: "pass"})
	if err == nil {
		t.Error("expected error for non-existent user")
	}
}
