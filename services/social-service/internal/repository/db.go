package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	Pool *pgxpool.Pool
}

func NewDB(dsn string) (*DB, error) {
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(context.Background()); err != nil {
		return nil, err
	}

	return &DB{Pool: pool}, nil
}

func (db *DB) Close() {
	db.Pool.Close()
}

type SocialRepository struct {
	db *DB
}

func NewSocialRepository(db *DB) *SocialRepository {
	return &SocialRepository{db: db}
}

func (r *SocialRepository) CreateReview(ctx context.Context, id, spotID, userID string, rating int, content string) error {
	_, err := r.db.Pool.Exec(ctx,
		"INSERT INTO reviews (id, spot_id, user_id, rating, content) VALUES ($1, $2, $3, $4, $5)",
		id, spotID, userID, rating, content,
	)
	return err
}

func (r *SocialRepository) GetReviewsBySpotID(ctx context.Context, spotID string, limit int) ([]struct {
	ID       string
	UserID   string
	Rating   int
	Content  string
	CreatedAt string
}, error) {
	return nil, nil
}

func (r *SocialRepository) CreateCatch(ctx context.Context, id, userID, spotID, species string, weight, length float64, bait string) error {
	_, err := r.db.Pool.Exec(ctx,
		"INSERT INTO catch_logs (id, user_id, spot_id, species, weight_lbs, length_in, bait_used) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		id, userID, spotID, species, weight, length, bait,
	)
	return err
}

func (r *SocialRepository) GetCatchesByUserID(ctx context.Context, userID string, limit int) ([]struct {
	ID        string
	SpotID    string
	Species   string
	Weight    float64
	Length    float64
	BaitUsed  string
	CreatedAt string
}, error) {
	return nil, nil
}
