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

type SpotRepository struct {
	db *DB
}

func NewSpotRepository(db *DB) *SpotRepository {
	return &SpotRepository{db: db}
}

func (r *SpotRepository) List(ctx context.Context, params struct {
	Type      string
	Species   string
	Difficulty string
	Page      int
	Limit     int
}) ([]struct {
	ID          string
	Name        string
	Lat         float64
	Lon         float64
	Type        string
	AccessNotes string
	Difficulty  string
	Rating      float64
	ReviewCount int
	Species     []string
	BestSeasons []string
}, error) {
	return nil, nil
}

func (r *SpotRepository) GetByID(ctx context.Context, id string) (spot struct {
	ID          string
	Name        string
	Lat         float64
	Lon         float64
	Type        string
	AccessNotes string
	Difficulty  string
	Rating      float64
	ReviewCount int
	Description string
}, err error) {
	err = r.db.Pool.QueryRow(ctx,
		`SELECT id, name, ST_Y(location::geometry), ST_X(location::geometry), type,
		 COALESCE(access_notes,''), COALESCE(difficulty,'easy'), COALESCE(rating,0),
		 COALESCE(review_count,0), COALESCE(description,'')
		 FROM spots WHERE id = $1`,
		id,
	).Scan(&spot.ID, &spot.Name, &spot.Lat, &spot.Lon, &spot.Type,
		&spot.AccessNotes, &spot.Difficulty, &spot.Rating, &spot.ReviewCount, &spot.Description)
	return
}

func (r *SpotRepository) Nearby(ctx context.Context, lat, lon, radiusMi float64, limit int) ([]struct {
	ID          string
	Name        string
	Lat         float64
	Lon         float64
	Type        string
	Distance    float64
	Rating      float64
	Species     []string
}, error) {
	rows, err := r.db.Pool.Query(ctx, `
		SELECT id, name, ST_X(location::geometry), ST_Y(location::geometry), type,
		       COALESCE(rating, 0),
		       ST_Distance(location::geography, ST_SetSRID(ST_MakePoint($1, $2), 4326)::geography) * 0.000621371
		FROM spots
		WHERE ST_DWithin(location::geography, ST_SetSRID(ST_MakePoint($1, $2), 4326)::geography, $3 * 1609.34)
		ORDER BY location::geography <-> ST_SetSRID(ST_MakePoint($1, $2), 4326)::geography
		LIMIT $4`,
		lon, lat, radiusMi, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []struct {
		ID       string
		Name     string
		Lat      float64
		Lon      float64
		Type     string
		Distance float64
		Rating   float64
		Species  []string
	}
	for rows.Next() {
		var s struct {
			ID       string
			Name     string
			Lat      float64
			Lon      float64
			Type     string
			Distance float64
			Rating   float64
			Species  []string
		}
		if err := rows.Scan(&s.ID, &s.Name, &s.Lon, &s.Lat, &s.Type, &s.Rating, &s.Distance); err != nil {
			return nil, err
		}
		results = append(results, s)
	}
	return results, nil
}

func (r *SpotRepository) Create(ctx context.Context, name string, lat, lon float64, spotType, difficulty string) (string, error) {
	var id string
	err := r.db.Pool.QueryRow(ctx,
		`INSERT INTO spots (name, location, type, difficulty)
		 VALUES ($1, ST_SetSRID(ST_MakePoint($2, $3), 4326), $4, $5)
		 RETURNING id`,
		name, lon, lat, spotType, difficulty,
	).Scan(&id)
	return id, err
}

func (r *SpotRepository) SearchBySpecies(ctx context.Context, species string, limit int) ([]struct {
	ID   string
	Name string
	Lat  float64
	Lon  float64
	Type string
}, error) {
	return nil, nil
}
