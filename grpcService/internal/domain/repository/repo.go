package repository

import "context"

type Repository interface {
	GetByString(ctx context.Context, query string) (int, error)
	PutStatistics(ctx context.Context, query string, count int) error
}
