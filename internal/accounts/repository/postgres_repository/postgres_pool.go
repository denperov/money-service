package postgres_repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

type postgresPool struct {
	postgresRepository
	url  string
	pool *pgxpool.Pool
}

func New(
	address string,
	name string,
	user string,
	password string,
) *postgresPool {
	return &postgresPool{
		url: fmt.Sprintf("postgresql://%s:%s@%s/%s?pool_max_conns=3", user, password, address, name),
	}
}

func (r *postgresPool) Start(ctx context.Context) error {
	if r.pool != nil {
		return nil
	}

	pool, err := connectWithRetries(ctx, r.url)
	if err != nil {
		return err
	}
	r.pool, r.postgresRepository.qi = pool, pool

	return nil
}

func (r *postgresPool) Stop() {
	if r.pool == nil {
		return
	}
	r.pool.Close()
}

func connectWithRetries(ctx context.Context, url string) (*pgxpool.Pool, error) {
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	cfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, err
	}

	for {
		pool, err := pgxpool.ConnectConfig(ctx, cfg)
		if err != nil {
			// wait
			select {
			case <-ctx.Done(): // cancellation
				return nil, ctx.Err()
			case <-ticker.C:
				continue
			}
		}
		return pool, nil
	}
}
