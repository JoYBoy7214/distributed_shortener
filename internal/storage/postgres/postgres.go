package postgres

import (
	"context"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresStore struct {
	pool *pgxpool.Pool
}

func NewPostgresStore(dbUrl string) (*PostgresStore, error) {
	pool, err := pgxpool.New(context.Background(), dbUrl)
	if err != nil {
		return nil, err
	}
	return &PostgresStore{pool: pool}, nil
}
func (p *PostgresStore) CreateDb(ctx context.Context) error {
	_, err := p.pool.Exec(context.Background(), "CREATE TABLE IF NOT EXISTS urls(id SERIAL UNIQUE,short_code VARCHAR(10) PRIMARY KEY UNIQUE NOT NULL,original_url TEXT NOT NULL,click_count INTEGER DEFAULT 0)")
	return err
}
func (p *PostgresStore) SaveURL(ctx context.Context, shortCode string, originalUrl string) error {
	_, err := p.pool.Exec(ctx, "INSERT INTO urls (short_code, original_url, click_count) VALUES ($1, $2, 0)", shortCode, originalUrl)
	return err
}

func (p *PostgresStore) GetURL(ctx context.Context, shortCode string) (string, error) {
	var originalUrl string
	err := p.pool.QueryRow(ctx, "SELECT original_url FROM urls WHERE short_code = $1", shortCode).Scan(&originalUrl)
	return originalUrl, err
}

func (p *PostgresStore) GetClickCount(ctx context.Context, shortcode string) (int, error) {
	var count int
	Row := p.pool.QueryRow(ctx, "SELECT click_count from urls WHERE short_code= $1", shortcode)
	err := Row.Scan(&count)
	return count, err
}

func (p *PostgresStore) IncrementClick(ctx context.Context, shortCode string) error {
	_, err := p.pool.Exec(ctx, "UPDATE urls SET click_count = click_count + 1 WHERE short_code=$1", shortCode)
	return err
}

func (p *PostgresStore) Close() {
	p.pool.Close()
}

// Helper to detect duplicates
func IsUniqueViolation(err error) bool {
	return strings.Contains(err.Error(), "23505")
}
