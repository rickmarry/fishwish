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

type SearchRepository struct {
	db *DB
}

func NewSearchRepository(db *DB) *SearchRepository {
	return &SearchRepository{db: db}
}

func (r *SearchRepository) Search(ctx context.Context, query string, limit int) ([]struct {
	ID      string
	Name    string
	Type    string
	Lat     float64
	Lon     float64
	Rating  float64
	Species []string
}, error) {
	return nil, nil
}

func (r *SearchRepository) GetSpeciesList(ctx context.Context) ([]struct {
	Name       string
	CommonName string
	SpotCount  int
}, error) {
	return nil, nil
}

func (r *SearchRepository) GetSuggestions(ctx context.Context, prefix string) ([]string, error) {
	return []string{}, nil
}
