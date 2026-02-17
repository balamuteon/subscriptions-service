//go:build integration
// +build integration

package subscription_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"

	"subscription_service/internal/domain"
	repository "subscription_service/internal/repository/subscription"
	"subscription_service/pkg/testdb"
)

var testPool *pgxpool.Pool
var teardown func()

func TestMain(m *testing.M) {
	ctx := context.Background()
	dsn, cleanup, err := testdb.SetupTestDatabase(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to setup test db: %v\n", err)
		os.Exit(1)
	}
	teardown = cleanup

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create pgx pool: %v\n", err)
		teardown()
		os.Exit(1)
	}
	testPool = pool

	code := m.Run()

	pool.Close()
	teardown()
	os.Exit(code)
}

func cleanupDB(t *testing.T) {
	t.Helper()
	_, err := testPool.Exec(context.Background(), "TRUNCATE TABLE subscriptions")
	require.NoError(t, err)
}

func TestRepositoryCreateGet(t *testing.T) {
	cleanupDB(t)

	repo := repository.New(testPool)
	userID := uuid.NewString()
	sub := domain.Subscription{
		ServiceName: "Netflix",
		Price:       500,
		UserID:      userID,
		StartDate:   "07-2025",
	}

	id, err := repo.Create(context.Background(), sub)
	require.NoError(t, err)
	require.NotEmpty(t, id)

	got, err := repo.GetByID(context.Background(), id)
	require.NoError(t, err)
	require.Equal(t, sub.ServiceName, got.ServiceName)
	require.Equal(t, sub.Price, got.Price)
	require.Equal(t, sub.UserID, got.UserID)
	require.Equal(t, sub.StartDate, got.StartDate)
	require.Nil(t, got.EndDate)
}

func TestRepositoryUpdate(t *testing.T) {
	cleanupDB(t)

	repo := repository.New(testPool)
	userID := uuid.NewString()
	sub := domain.Subscription{
		ServiceName: "Netflix",
		Price:       500,
		UserID:      userID,
		StartDate:   "07-2025",
	}

	id, err := repo.Create(context.Background(), sub)
	require.NoError(t, err)

	end := "12-2025"
	update := domain.Subscription{
		ID:          id,
		ServiceName: "HBO",
		Price:       700,
		UserID:      userID,
		StartDate:   "08-2025",
		EndDate:     &end,
	}
	require.NoError(t, repo.Update(context.Background(), update))

	got, err := repo.GetByID(context.Background(), id)
	require.NoError(t, err)
	require.Equal(t, "HBO", got.ServiceName)
	require.Equal(t, 700, got.Price)
	require.Equal(t, "08-2025", got.StartDate)
	require.NotNil(t, got.EndDate)
	require.Equal(t, "12-2025", *got.EndDate)
}

func TestRepositoryDelete(t *testing.T) {
	cleanupDB(t)

	repo := repository.New(testPool)
	userID := uuid.NewString()
	sub := domain.Subscription{
		ServiceName: "Netflix",
		Price:       500,
		UserID:      userID,
		StartDate:   "07-2025",
	}

	id, err := repo.Create(context.Background(), sub)
	require.NoError(t, err)

	require.NoError(t, repo.Delete(context.Background(), id))

	_, err = repo.GetByID(context.Background(), id)
	require.ErrorIs(t, err, domain.ErrSubscriptionNotFound)
}

func TestRepositoryListFilters(t *testing.T) {
	cleanupDB(t)

	repo := repository.New(testPool)
	userID := uuid.NewString()
	otherUser := uuid.NewString()

	_, err := repo.Create(context.Background(), domain.Subscription{
		ServiceName: "Netflix",
		Price:       500,
		UserID:      userID,
		StartDate:   "07-2025",
	})
	require.NoError(t, err)

	_, err = repo.Create(context.Background(), domain.Subscription{
		ServiceName: "Spotify",
		Price:       300,
		UserID:      otherUser,
		StartDate:   "07-2025",
	})
	require.NoError(t, err)

	items, err := repo.List(context.Background(), userID, "")
	require.NoError(t, err)
	require.Len(t, items, 1)
	require.Equal(t, "Netflix", items[0].ServiceName)
}

func TestRepositoryTotal(t *testing.T) {
	cleanupDB(t)

	repo := repository.New(testPool)
	userID := uuid.NewString()

	end := "09-2025"
	_, err := repo.Create(context.Background(), domain.Subscription{
		ServiceName: "Netflix",
		Price:       100,
		UserID:      userID,
		StartDate:   "07-2025",
		EndDate:     &end,
	})
	require.NoError(t, err)

	_, err = repo.Create(context.Background(), domain.Subscription{
		ServiceName: "Spotify",
		Price:       200,
		UserID:      userID,
		StartDate:   "08-2025",
	})
	require.NoError(t, err)

	to := "09-2025"
	total, err := repo.Total(context.Background(), domain.Subscription{
		UserID:    userID,
		StartDate: "07-2025",
		EndDate:   &to,
	})
	require.NoError(t, err)

	require.Equal(t, int64(700), total)
}
