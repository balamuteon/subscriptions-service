package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"subscription_service/internal/domain"
	"subscription_service/pkg/logger"
)

type SubscriptionHandler struct {
	log     logger.Logger
	service subscriptionService
}

func NewSubscriptionHandler(log logger.Logger, service subscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{log: log, service: service}
}

// CreateSubscription godoc
// @Summary Create subscription
// @Description Create a new subscription record.
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param request body SubscriptionRequest true "Subscription data"
// @Success 201 {object} IDResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /subscriptions [post]
func (h *SubscriptionHandler) CreateSubscription(w http.ResponseWriter, r *http.Request) {
	var reqDTO SubscriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&reqDTO); err != nil {
		newErrorResponse(w, http.StatusBadRequest, ErrInvalidJSON)
		return
	}

	id, err := h.service.Create(r.Context(), reqDTO.toDomain())
	if err != nil {
		h.handleError(w, err, "create subscription")
		return
	}

	if err := writeJSON(w, http.StatusCreated, IDResponse{ID: id}); err != nil {
		h.log.Error("failed to write response", "error", err)
	}
}

// GetSubscription godoc
// @Summary Get subscription by ID
// @Tags subscriptions
// @Produce json
// @Param id path string true "Subscription ID"
// @Success 200 {object} SubscriptionResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /subscriptions/{id} [get]
func (h *SubscriptionHandler) GetSubscription(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	sub, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		h.handleError(w, err, "get subscription")
		return
	}

	if err := writeJSON(w, http.StatusOK, fromDomain(sub)); err != nil {
		h.log.Error("failed to write response", "error", err)
	}
}

// ListSubscriptions godoc
// @Summary List subscriptions
// @Description List subscriptions with optional filters.
// @Tags subscriptions
// @Produce json
// @Param user_id query string false "User ID (UUID)"
// @Param service_name query string false "Service name"
// @Success 200 {array} SubscriptionResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /subscriptions [get]
func (h *SubscriptionHandler) ListSubscriptions(w http.ResponseWriter, r *http.Request) {
	items, err := h.service.List(r.Context(), r.URL.Query().Get("user_id"), r.URL.Query().Get("service_name"))
	if err != nil {
		h.handleError(w, err, "list subscriptions")
		return
	}

	if err := writeJSON(w, http.StatusOK, fromDomainList(items)); err != nil {
		h.log.Error("failed to write response", "error", err)
	}
}

// UpdateSubscription godoc
// @Summary Update subscription
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "Subscription ID"
// @Param request body SubscriptionRequest true "Subscription data"
// @Success 200 {string} string "updated successfully"
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /subscriptions/{id} [put]
func (h *SubscriptionHandler) UpdateSubscription(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var reqDTO SubscriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&reqDTO); err != nil {
		newErrorResponse(w, http.StatusBadRequest, ErrInvalidJSON)
		return
	}

	sub := reqDTO.toDomain()
	sub.ID = id

	if err := h.service.Update(r.Context(), sub); err != nil {
		h.handleError(w, err, "update subscription")
		return
	}

	if err := writeJSON(w, http.StatusOK, StatusResponse{Status: "updated successfully"}); err != nil {
		h.log.Error("failed to write response", "error", err)
	}
}

// DeleteSubscription godoc
// @Summary Delete subscription
// @Tags subscriptions
// @Produce json
// @Param id path string true "Subscription ID"
// @Success 200 {object} StatusResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /subscriptions/{id} [delete]
func (h *SubscriptionHandler) DeleteSubscription(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.service.Delete(r.Context(), id); err != nil {
		h.handleError(w, err, "delete subscription")
		return
	}

	if err := writeJSON(w, http.StatusOK, StatusResponse{Status: "ok"}); err != nil {
		h.log.Error("failed to write response", "error", err)
	}
}

// TotalSubscriptions godoc
// @Summary Calculate total subscriptions cost
// @Description Sum of subscription costs for the specified period with optional filters.
// @Tags subscriptions
// @Produce json
// @Param from query string true "Start period (MM-YYYY)"
// @Param to query string true "End period (MM-YYYY)"
// @Param user_id query string false "User ID (UUID)"
// @Param service_name query string false "Service name"
// @Success 200 {object} TotalResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /subscriptions/total [get]
func (h *SubscriptionHandler) TotalSubscriptions(w http.ResponseWriter, r *http.Request) {
	to := r.URL.Query().Get("to")
	var endDate *string
	if to != "" {
		endDate = &to
	}

	filter := domain.Subscription{
		UserID:      r.URL.Query().Get("user_id"),
		ServiceName: r.URL.Query().Get("service_name"),
		StartDate:   r.URL.Query().Get("from"),
		EndDate:     endDate,
	}

	total, err := h.service.Total(r.Context(), filter)
	if err != nil {
		h.handleError(w, err, "calculate subscriptions total")
		return
	}

	if err := writeJSON(w, http.StatusOK, TotalResponse{Total: total}); err != nil {
		h.log.Error("failed to write response", "error", err)
	}
}
