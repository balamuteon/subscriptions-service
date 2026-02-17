package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"subscription_service/internal/domain"
)

func writeJSON(w http.ResponseWriter, status int, v any) error {
	data, err := json.Marshal(v)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
		return err
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_, err = w.Write(data)
	return err
}

func newErrorResponse(w http.ResponseWriter, status int, msg error) {
	resp := map[string]string{"error": msg.Error()}
	_ = writeJSON(w, status, resp)
}

func (h *SubscriptionHandler) handleError(w http.ResponseWriter, err error, operation string) {
	var vErr *domain.ValidationError
	if errors.As(err, &vErr) {
		newErrorResponse(w, http.StatusBadRequest, vErr)
		return
	}

	if errors.Is(err, domain.ErrSubscriptionNotFound) {
		newErrorResponse(w, http.StatusNotFound, err)
		return
	}

	if errors.Is(err, domain.ErrNotImplemented) {
		newErrorResponse(w, http.StatusNotImplemented, domain.ErrNotImplemented)
		return
	}

	h.log.Error("operation failed", "operation", operation, "error", err)
	newErrorResponse(w, http.StatusInternalServerError, ErrStatusInternalServerError)
}
