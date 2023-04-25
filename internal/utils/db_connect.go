package utils

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GetConnectToPg(ctx context.Context, cfg *StorageConfig) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", cfg.PgUsername, cfg.PgPassword, cfg.PgHost, cfg.PgPort, cfg.PgDatabase)
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return pool, err
	}
	return pool, nil
}
