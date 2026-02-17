package httpapi_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"subscription_service/internal/domain"
	"subscription_service/internal/httpapi"
	"subscription_service/pkg/logger"
)

func TestCreateSubscription_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := NewMocksubscriptionService(ctrl)
	svc.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, sub domain.Subscription) (string, error) {
			require.Equal(t, "Netflix", sub.ServiceName)
			require.Equal(t, 400, sub.Price)
			return "id-123", nil
		})
	log := logger.NewNoop()
	apiHandler := httpapi.NewSubscriptionHandler(log, svc)
	h := httpapi.NewHandler(log, apiHandler)

	body := []byte(`{"service_name":"Netflix","price":400,"user_id":"` + uuid.NewString() + `","start_date":"07-2025"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/subscriptions/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)
	var resp map[string]string
	require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
	require.Equal(t, "id-123", resp["id"])
}

func TestCreateSubscription_InvalidJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := NewMocksubscriptionService(ctrl)
	log := logger.NewNoop()
	apiHandler := httpapi.NewSubscriptionHandler(log, svc)
	h := httpapi.NewHandler(log, apiHandler)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/subscriptions/", bytes.NewBufferString("{"))
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
	var resp map[string]string
	require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
	require.Equal(t, httpapi.ErrInvalidJSON.Error(), resp["error"])
}

func TestGetSubscription_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := NewMocksubscriptionService(ctrl)
	svc.EXPECT().GetByID(gomock.Any(), gomock.Any()).
		Return(domain.Subscription{}, domain.ErrSubscriptionNotFound)
	log := logger.NewNoop()
	apiHandler := httpapi.NewSubscriptionHandler(log, svc)
	h := httpapi.NewHandler(log, apiHandler)

	id := uuid.NewString()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/subscriptions/"+id, nil)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	require.Equal(t, http.StatusNotFound, w.Code)
}

func TestUpdateSubscription_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := NewMocksubscriptionService(ctrl)
	svc.EXPECT().Update(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, sub domain.Subscription) error {
			require.Equal(t, "Netflix", sub.ServiceName)
			return nil
		})
	log := logger.NewNoop()
	apiHandler := httpapi.NewSubscriptionHandler(log, svc)
	h := httpapi.NewHandler(log, apiHandler)

	id := uuid.NewString()
	body := []byte(`{"service_name":"Netflix","price":500,"user_id":"` + uuid.NewString() + `","start_date":"07-2025"}`)
	req := httptest.NewRequest(http.MethodPut, "/api/v1/subscriptions/"+id, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
}

func TestTotalSubscriptions_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := NewMocksubscriptionService(ctrl)
	svc.EXPECT().Total(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, filter domain.Subscription) (int64, error) {
			require.Equal(t, "07-2025", filter.StartDate)
			require.NotNil(t, filter.EndDate)
			require.Equal(t, "08-2025", *filter.EndDate)
			return 300, nil
		})
	log := logger.NewNoop()
	apiHandler := httpapi.NewSubscriptionHandler(log, svc)
	h := httpapi.NewHandler(log, apiHandler)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/subscriptions/total?from=07-2025&to=08-2025", nil)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	var resp httpapi.TotalResponse
	require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
	require.Equal(t, int64(300), resp.Total)
}
