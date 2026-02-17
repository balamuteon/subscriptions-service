package postgres

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	maxRetries = 5
	retryPause = 3 * time.Second
)

func NewConnection(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	const op = "db.NewConnection"

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Println("Failed setting connection config to postgres.")
		return nil, fmt.Errorf("%s %w", op, err)
	}

	var lastErr error
	for attempt := 1; attempt <= maxRetries; attempt++ {
		err := pool.Ping(ctx)
		if err == nil {
			log.Println("Success setting connection to postgres.")
			return pool, nil
		}

		lastErr = err
		if attempt == maxRetries {
			break
		}

		log.Printf(
			"Failed to connect to postgres (attempt %d/%d): %v. Retrying in %v...",
			attempt,
			maxRetries,
			err,
			retryPause,
		)

		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("%s: context cancelled during retry: %w", op, ctx.Err())
		case <-time.After(retryPause):
		}
	}

	log.Println("Failed setting connection to postgres after all retries.")
	return nil, fmt.Errorf("%s: all retries failed: %w", op, lastErr)
}
