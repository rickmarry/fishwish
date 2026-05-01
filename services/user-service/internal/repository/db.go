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

type UserRepository struct {
	db *DB
}

func NewUserRepository(db *DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, id, email, username, hashedPassword string) error {
	_, err := r.db.Pool.Exec(ctx,
		"INSERT INTO users (id, email, username, password_hash) VALUES ($1, $2, $3, $4)",
		id, email, username, hashedPassword,
	)
	return err
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (id, username, hashedPassword string, err error) {
	err = r.db.Pool.QueryRow(ctx,
		"SELECT id, username, password_hash FROM users WHERE email = $1",
		email,
	).Scan(&id, &username, &hashedPassword)
	return
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (email, username, avatarURL string, err error) {
	err = r.db.Pool.QueryRow(ctx,
		"SELECT email, username, COALESCE(avatar_url, '') FROM users WHERE id = $1",
		id,
	).Scan(&email, &username, &avatarURL)
	return
}
