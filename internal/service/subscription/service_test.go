package subscription_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"subscription_service/internal/domain"
	subscriptionService "subscription_service/internal/service/subscription"
)

func TestServiceCreate_Invalid(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := NewMockrepository(ctrl)
	svc := subscriptionService.New(repo)

	_, err := svc.Create(context.Background(), domain.Subscription{})
	var vErr *domain.ValidationError
	require.ErrorAs(t, err, &vErr)
	require.ErrorIs(t, vErr, domain.ErrMissingRequiredFields)
}

func TestServiceCreate_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := NewMockrepository(ctrl)
	svc := subscriptionService.New(repo)

	userID := uuid.NewString()
	input := domain.Subscription{
		ServiceName: "Netflix",
		Price:       500,
		UserID:      userID,
		StartDate:   "07-2025",
	}

	repo.EXPECT().Create(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, sub domain.Subscription) (string, error) {
			require.Equal(t, input.ServiceName, sub.ServiceName)
			require.Equal(t, input.Price, sub.Price)
			require.Equal(t, input.UserID, sub.UserID)
			require.Equal(t, input.StartDate, sub.StartDate)
			return "id-1", nil
		})

	id, err := svc.Create(context.Background(), input)
	require.NoError(t, err)
	require.Equal(t, "id-1", id)
}

func TestServiceGetByID_Invalid(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := NewMockrepository(ctrl)
	svc := subscriptionService.New(repo)

	_, err := svc.GetByID(context.Background(), "bad-id")
	var vErr *domain.ValidationError
	require.ErrorAs(t, err, &vErr)
	require.ErrorIs(t, vErr, domain.ErrInvalidID)
}

func TestServiceList_InvalidFilter(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := NewMockrepository(ctrl)
	svc := subscriptionService.New(repo)

	_, err := svc.List(context.Background(), " bad ", "")
	var vErr *domain.ValidationError
	require.ErrorAs(t, err, &vErr)
	require.ErrorIs(t, vErr, domain.ErrInvalidUserID)
}

func TestServiceTotal_InvalidFilter(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := NewMockrepository(ctrl)
	svc := subscriptionService.New(repo)

	_, err := svc.Total(context.Background(), domain.Subscription{})
	var vErr *domain.ValidationError
	require.ErrorAs(t, err, &vErr)
	require.ErrorIs(t, vErr, domain.ErrMissingRequiredFields)
}

func TestServiceUpdate_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := NewMockrepository(ctrl)
	svc := subscriptionService.New(repo)

	err := svc.Update(context.Background(), domain.Subscription{ID: "bad"})
	var vErr *domain.ValidationError
	require.ErrorAs(t, err, &vErr)
	require.ErrorIs(t, vErr, domain.ErrInvalidID)
}

func TestServiceDelete_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := NewMockrepository(ctrl)
	svc := subscriptionService.New(repo)

	err := svc.Delete(context.Background(), "bad")
	var vErr *domain.ValidationError
	require.ErrorAs(t, err, &vErr)
	require.ErrorIs(t, vErr, domain.ErrInvalidID)
}
