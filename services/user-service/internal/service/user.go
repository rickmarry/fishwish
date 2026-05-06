package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"fishwish/pkg/auth"
	"fishwish/services/user-service/internal/model"
	"fishwish/services/user-service/internal/repository"
)

type UserRepositoryInterface interface {
	GetByEmail(ctx context.Context, email string) (string, string, string, error)
	Create(ctx context.Context, id, email, username, hashedPassword string) error
	GetByID(ctx context.Context, id string) (string, string, string, error)
}

type UserService struct {
	repo UserRepositoryInterface
}

func NewUserService(db *repository.DB) *UserService {
	return &UserService{
		repo: repository.NewUserRepository(db),
	}
}

func (s *UserService) Register(ctx context.Context, req model.RegisterRequest) (*auth.TokenPair, error) {
	existing, _, _, err := s.repo.GetByEmail(ctx, req.Email)
	if err == nil && existing != "" {
		return nil, errors.New("email already registered")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	id := uuid.New().String()
	if err := s.repo.Create(ctx, id, req.Email, req.Username, string(hashedPassword)); err != nil {
		return nil, err
	}

	return auth.GenerateTokens(id, req.Email, "user")
}

func (s *UserService) Login(ctx context.Context, req model.LoginRequest) (*auth.TokenPair, error) {
	id, _, hashedPassword, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	return auth.GenerateTokens(id, req.Email, "user")
}

func (s *UserService) GetProfile(ctx context.Context, userID string) (*model.User, error) {
	email, username, avatarURL, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &model.User{
		ID:        userID,
		Email:     email,
		Username:  username,
		AvatarURL: avatarURL,
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}, nil
}
