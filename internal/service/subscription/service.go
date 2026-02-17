package subscription

import (
	"context"

	"subscription_service/internal/domain"
)

type Service struct {
	repo repository
}

func New(repo repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, sub domain.Subscription) (string, error) {
	normalized, err := validateCreateOrUpdateInput(sub)
	if err != nil {
		return "", err
	}

	return s.repo.Create(ctx, normalized)
}

func (s *Service) GetByID(ctx context.Context, id string) (domain.Subscription, error) {
	if err := validateID(id); err != nil {
		return domain.Subscription{}, err
	}

	return s.repo.GetByID(ctx, id)
}

func (s *Service) List(ctx context.Context, userID string, serviceName string) ([]domain.Subscription, error) {
	normalizedUserID, normalizedService, err := validateListFilter(userID, serviceName)
	if err != nil {
		return nil, err
	}

	return s.repo.List(ctx, normalizedUserID, normalizedService)
}

func (s *Service) Update(ctx context.Context, sub domain.Subscription) error {
	if err := validateID(sub.ID); err != nil {
		return err
	}

	normalized, err := validateCreateOrUpdateInput(sub)
	if err != nil {
		return err
	}

	return s.repo.Update(ctx, normalized)
}

func (s *Service) Delete(ctx context.Context, id string) error {
	if err := validateID(id); err != nil {
		return err
	}

	return s.repo.Delete(ctx, id)
}

func (s *Service) Total(ctx context.Context, filter domain.Subscription) (int64, error) {
	validated, err := validateTotalFilter(filter)
	if err != nil {
		return 0, err
	}

	return s.repo.Total(ctx, validated)
}
