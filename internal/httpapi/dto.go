package httpapi

import "subscription_service/internal/domain"

type TotalResponse struct {
	Total int64 `json:"total"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type IDResponse struct {
	ID string `json:"id"`
}

type StatusResponse struct {
	Status string `json:"status"`
}

type SubscriptionRequest struct {
	ServiceName string  `json:"service_name"`
	Price       int     `json:"price"`
	UserID      string  `json:"user_id"`
	StartDate   string  `json:"start_date"`
	EndDate     *string `json:"end_date,omitempty"`
}

type SubscriptionResponse struct {
	ID          string  `json:"id"`
	ServiceName string  `json:"service_name"`
	Price       int     `json:"price"`
	UserID      string  `json:"user_id"`
	StartDate   string  `json:"start_date"`
	EndDate     *string `json:"end_date,omitempty"`
}

func (dto *SubscriptionRequest) toDomain() domain.Subscription {
	return domain.Subscription{
		ServiceName: dto.ServiceName,
		Price:       dto.Price,
		UserID:      dto.UserID,
		StartDate:   dto.StartDate,
		EndDate:     dto.EndDate,
	}
}

func fromDomain(sub domain.Subscription) SubscriptionResponse {
	return SubscriptionResponse{
		ID:          sub.ID,
		ServiceName: sub.ServiceName,
		Price:       sub.Price,
		UserID:      sub.UserID,
		StartDate:   sub.StartDate,
		EndDate:     sub.EndDate,
	}
}

func fromDomainList(items []domain.Subscription) []SubscriptionResponse {
	result := make([]SubscriptionResponse, len(items))
	for i, sub := range items {
		result[i] = fromDomain(sub)
	}
	return result
}
