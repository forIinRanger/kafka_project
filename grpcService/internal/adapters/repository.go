package adapters

import (
	"context"
	"database/sql"
	sq "github.com/Masterminds/squirrel"
)

type PostgresRepo struct {
	db *sql.DB
	sb sq.StatementBuilderType
}

func NewPostgresRepo(db *sql.DB) *PostgresRepo {
	return &PostgresRepo{
		db: db,
		sb: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *PostgresRepo) GetByString(ctx context.Context, query string) (int, error) {
	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	default:
	}
	count := sq.Select("count").From("messages").Where(sq.Eq{"value": query})
	sql, args, err := count.ToSql()
	if err != nil {
		return 0, err
	}
	rows, err := r.db.Query(sql, args)
	if err != nil {
		return 0, err
	}
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
	sql, args, err := sq.Insert("messages").Columns("value", "count").Values(query, count).ToSql()
	if err != nil {
		return err
	}
	_, err = r.db.Query(sql, args)
	return err
}
