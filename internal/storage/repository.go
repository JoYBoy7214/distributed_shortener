package storage

import "context"

type Repository interface {
	CreateDb(ctx context.Context) error
	SaveURL(ctx context.Context, shortCode string, originalUrl string) error
	GetURL(ctx context.Context, shortCode string) (string, error)
	GetClickCount(ctx context.Context, shortcode string) (int, error)
	IncrementClick(ctx context.Context, shortCode string) error
	Close()
}
