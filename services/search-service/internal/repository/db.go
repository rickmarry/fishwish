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

type SearchRepositoryInterface interface {
	Search(ctx context.Context, query string, limit int) ([]struct {
		ID      string
		Name    string
		Type    string
		Lat     float64
		Lon     float64
		Rating  float64
		Species []string
	}, error)
	GetSpeciesList(ctx context.Context) ([]struct {
		Name       string
		CommonName string
		SpotCount  int
	}, error)
	GetSuggestions(ctx context.Context, prefix string) ([]string, error)
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
	rows, err := r.db.Pool.Query(ctx, `
		SELECT s.id, s.name, s.type,
		       ST_Y(s.location::geometry) as lat,
		       ST_X(s.location::geometry) as lon,
		       s.rating
		FROM spots s
		WHERE ($1 = '' OR s.name ILIKE '%' || $1 || '%')
		LIMIT $2`, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []struct {
		ID      string
		Name    string
		Type    string
		Lat     float64
		Lon     float64
		Rating  float64
		Species []string
	}

	for rows.Next() {
		var r struct {
			ID      string
			Name    string
			Type    string
			Lat     float64
			Lon     float64
			Rating  float64
			Species []string
		}
		if err := rows.Scan(&r.ID, &r.Name, &r.Type, &r.Lat, &r.Lon, &r.Rating); err != nil {
			return nil, err
		}
		r.Species = []string{} // TODO: join with spot_species when needed
		results = append(results, r)
	}

	return results, nil
}

func (r *SearchRepository) GetSpeciesList(ctx context.Context) ([]struct {
	Name       string
	CommonName string
	SpotCount  int
}, error) {
	rows, err := r.db.Pool.Query(ctx, `
		SELECT s.name, s.common_name, COUNT(ss.spot_id) as spot_count
		FROM species s
		LEFT JOIN spot_species ss ON s.id = ss.species_id
		GROUP BY s.id, s.name, s.common_name
		ORDER BY spot_count DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []struct {
		Name       string
		CommonName string
		SpotCount  int
	}

	for rows.Next() {
		var r struct {
			Name       string
			CommonName string
			SpotCount  int
		}
		if err := rows.Scan(&r.Name, &r.CommonName, &r.SpotCount); err != nil {
			return nil, err
		}
		results = append(results, r)
	}

	return results, nil
}

func (r *SearchRepository) GetSuggestions(ctx context.Context, prefix string) ([]string, error) {
	rows, err := r.db.Pool.Query(ctx, "SELECT DISTINCT name FROM spots WHERE name ILIKE $1 || '%' LIMIT 10", prefix)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var suggestions []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		suggestions = append(suggestions, name)
	}

	return suggestions, nil
}
