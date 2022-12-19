package adapters

import (
	"context"
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"log"
)

type PostgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(db *sql.DB) *PostgresRepo {
	return &PostgresRepo{
		db: db,
	}
}

func (r *PostgresRepo) GetByString(ctx context.Context, query string) (int, error) {
	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	default:
	}
	ctx1, cancel := context.WithCancel(ctx)
	count := sq.Select("letter_counts").From("messages").Where(sq.Eq{"text": query}).PlaceholderFormat(sq.Dollar)
	sql, args, err := count.ToSql()
	if err != nil {
		return 0, err
	}
	connection, err := r.db.Conn(ctx1)
	defer cancel()
	if err != nil {
		return 0, fmt.Errorf("Error with connection to db: %w", err)
	}
	rows := connection.QueryRowContext(ctx1, sql, args...)
	var res int
	err = rows.Scan(&res)
	if err != nil {
		return 0, err
	}
	return res, nil
}

func (r *PostgresRepo) PutStatistics(ctx context.Context, query string, count int) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	ctx1, cancel := context.WithCancel(ctx)
	sql, args, err := sq.Insert("messages").Columns("text", "letter_counts").Values(query, count).PlaceholderFormat(sq.Dollar).ToSql()
	log.Printf("huy %v", sql)
	if err != nil {
		return fmt.Errorf("Error in sql query: %w", err)
	}
	connection, err := r.db.Conn(ctx1)
	defer cancel()
	if err != nil {
		return fmt.Errorf("Error with connection to db: %w", err)
	}
	_, err = connection.ExecContext(ctx1, sql, args...)
	if err != nil {
		return fmt.Errorf("Error with executing: %w", err)
	}
	return err
}
