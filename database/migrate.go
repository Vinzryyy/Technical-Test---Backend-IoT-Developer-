package database

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed schema.sql
var schemaSQL string

func Migrate(pool *pgxpool.Pool) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*1e9)
	defer cancel()
	if _, err := pool.Exec(ctx, schemaSQL); err != nil {
		return fmt.Errorf("run schema.sql: %w", err)
	}
	return nil
}
