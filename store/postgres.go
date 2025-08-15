package store

import (
	"context"
	"os"

	"urlx/logging"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresStore struct {
	pool *pgxpool.Pool
}

func NewPostgresStore() (*PostgresStore, error) {
	dsn := os.Getenv("DATABASE_URL")
	logger := logging.NewLogger()
	logger.Debug("Connecting to Postgres", "dsn", dsn)
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		logger.Error("Failed to connect to Postgres", "error", err)
		return nil, err
	}
	logger.Debug("Connected to Postgres successfully")

	// Create schema if it doesn't exist
	schema := `CREATE TABLE IF NOT EXISTS urls (
		short VARCHAR(32) PRIMARY KEY,
		long TEXT NOT NULL
	);`
	_, err = pool.Exec(context.Background(), schema)
	if err != nil {
		logger.Error("Failed to create schema", "error", err)
		return nil, err
	}
	logger.Debug("Database schema ensured")

	return &PostgresStore{pool: pool}, nil
}

func (p *PostgresStore) Set(short, long string) {
	logger := logging.NewLogger()
	logger.Debug("Set short URL", "short", short, "long", long)
	_, err := p.pool.Exec(context.Background(), "INSERT INTO urls (short, long) VALUES ($1, $2) ON CONFLICT (short) DO UPDATE SET long = $2", short, long)
	if err != nil {
		logger.Error("Failed to set short URL", "short", short, "error", err)
	}
}

func (p *PostgresStore) Get(short string) (string, bool) {
	logger := logging.NewLogger()
	logger.Debug("Get short URL", "short", short)
	var long string
	err := p.pool.QueryRow(context.Background(), "SELECT long FROM urls WHERE short = $1", short).Scan(&long)
	if err != nil {
		logger.Error("Failed to get short URL", "short", short, "error", err)
		return "", false
	}
	return long, true
}

func (p *PostgresStore) Delete(short string) {
	logger := logging.NewLogger()
	logger.Debug("Delete short URL", "short", short)
	_, err := p.pool.Exec(context.Background(), "DELETE FROM urls WHERE short = $1", short)
	if err != nil {
		logger.Error("Failed to delete short URL", "short", short, "error", err)
	}
}
