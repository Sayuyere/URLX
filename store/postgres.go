package store

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresStore struct {
	pool *pgxpool.Pool
}

func NewPostgresStore() (*PostgresStore, error) {
	dsn := os.Getenv("DATABASE_URL")
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, err
	}
	return &PostgresStore{pool: pool}, nil
}

func (p *PostgresStore) Set(short, long string) {
	_, _ = p.pool.Exec(context.Background(), "INSERT INTO urls (short, long) VALUES ($1, $2) ON CONFLICT (short) DO UPDATE SET long = $2", short, long)
}

func (p *PostgresStore) Get(short string) (string, bool) {
	var long string
	err := p.pool.QueryRow(context.Background(), "SELECT long FROM urls WHERE short = $1", short).Scan(&long)
	if err != nil {
		return "", false
	}
	return long, true
}

func (p *PostgresStore) Delete(short string) {
	_, _ = p.pool.Exec(context.Background(), "DELETE FROM urls WHERE short = $1", short)
}
