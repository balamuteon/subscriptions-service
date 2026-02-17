package subscription

import (
	"context"

	"subscription_service/internal/domain"
)
//go:generate mockgen -source=contract.go -destination=mock_test.go -package=subscription_test
type repository interface {
	Create(ctx context.Context, sub domain.Subscription) (string, error)
	GetByID(ctx context.Context, id string) (domain.Subscription, error)
	List(ctx context.Context, userID string, serviceName string) ([]domain.Subscription, error)
	Update(ctx context.Context, sub domain.Subscription) error
	Delete(ctx context.Context, id string) error
	Total(ctx context.Context, filter domain.Subscription) (int64, error)
}
