package repository

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq" // postgres driver
)

const (
	tableBlacklist = "black_list"
	tableWhitelist = "white_list"
)

type PostgresIPRepository struct {
	db *sql.DB
}

func NewPostgresIPRepository(db *sql.DB) *PostgresIPRepository {
	return &PostgresIPRepository{db: db}
}

func (r *PostgresIPRepository) IsBlacklisted(ctx context.Context, ip string) (bool, error) {
	return r.exists(ctx, tableBlacklist, ip)
}

func (r *PostgresIPRepository) IsWhitelisted(ctx context.Context, ip string) (bool, error) {
	return r.exists(ctx, tableWhitelist, ip)
}

func (r *PostgresIPRepository) exists(ctx context.Context, table, ip string) (bool, error) {
	var exists bool
	var query string
	if table == tableBlacklist {
		query = "SELECT EXISTS(SELECT 1 FROM black_list WHERE subnet >> $1)"
	} else {
		query = "SELECT EXISTS(SELECT 1 FROM white_list WHERE subnet >> $1)"
	}
	err := r.db.QueryRowContext(ctx, query, ip).Scan(&exists)
	return exists, err
}

func (r *PostgresIPRepository) AddToBlacklist(ctx context.Context, subnet string) error {
	return r.add(ctx, tableBlacklist, subnet)
}

func (r *PostgresIPRepository) RemoveFromBlacklist(ctx context.Context, subnet string) error {
	return r.remove(ctx, tableBlacklist, subnet)
}

func (r *PostgresIPRepository) AddToWhitelist(ctx context.Context, subnet string) error {
	return r.add(ctx, tableWhitelist, subnet)
}

func (r *PostgresIPRepository) RemoveFromWhitelist(ctx context.Context, subnet string) error {
	return r.remove(ctx, tableWhitelist, subnet)
}

func (r *PostgresIPRepository) add(ctx context.Context, table, subnet string) error {
	var query string
	if table == tableBlacklist {
		query = "INSERT INTO black_list (subnet) VALUES ($1) ON CONFLICT DO NOTHING"
	} else {
		query = "INSERT INTO white_list (subnet) VALUES ($1) ON CONFLICT DO NOTHING"
	}
	_, err := r.db.ExecContext(ctx, query, subnet)
	return err
}

func (r *PostgresIPRepository) remove(ctx context.Context, table, subnet string) error {
	var query string
	if table == tableBlacklist {
		query = "DELETE FROM black_list WHERE subnet = $1"
	} else {
		query = "DELETE FROM white_list WHERE subnet = $1"
	}
	_, err := r.db.ExecContext(ctx, query, subnet)
	return err
}
