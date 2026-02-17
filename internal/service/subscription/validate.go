package subscription

import (
	"strings"
	"time"

	"github.com/google/uuid"

	"subscription_service/internal/domain"
)

const monthYearLayout = "01-2006"

func validateCreateOrUpdateInput(sub domain.Subscription) (domain.Subscription, error) {
	if strings.TrimSpace(sub.ServiceName) == "" || strings.TrimSpace(sub.UserID) == "" || strings.TrimSpace(sub.StartDate) == "" {
		return domain.Subscription{}, &domain.ValidationError{Err: domain.ErrMissingRequiredFields}
	}

	if strings.TrimSpace(sub.ServiceName) != sub.ServiceName {
		return domain.Subscription{}, &domain.ValidationError{Err: domain.ErrInvalidServiceName}
	}

	if sub.Price <= 0 {
		return domain.Subscription{}, &domain.ValidationError{Err: domain.ErrInvalidPrice}
	}

	if _, err := uuid.Parse(sub.UserID); err != nil {
		return domain.Subscription{}, &domain.ValidationError{Err: domain.ErrInvalidUserID}
	}

	startDate, err := parseMonthYear(sub.StartDate)
	if err != nil {
		return domain.Subscription{}, &domain.ValidationError{Err: domain.ErrInvalidStartDate}
	}
	sub.StartDate = startDate.Format(monthYearLayout)

	if sub.EndDate != nil {
		if strings.TrimSpace(*sub.EndDate) == "" {
			return domain.Subscription{}, &domain.ValidationError{Err: domain.ErrInvalidEndDate}
		}
		if strings.TrimSpace(*sub.EndDate) != *sub.EndDate {
			return domain.Subscription{}, &domain.ValidationError{Err: domain.ErrInvalidEndDate}
		}

		endDate, err := parseMonthYear(*sub.EndDate)
		if err != nil {
			return domain.Subscription{}, &domain.ValidationError{Err: domain.ErrInvalidEndDate}
		}
		if endDate.Before(startDate) {
			return domain.Subscription{}, &domain.ValidationError{Err: domain.ErrInvalidPeriod}
		}
		formatted := endDate.Format(monthYearLayout)
		sub.EndDate = &formatted
	}

	return sub, nil
}

func validateID(id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return &domain.ValidationError{Err: domain.ErrInvalidID}
	}
	return nil
}

func validateListFilter(userID string, serviceName string) (string, string, error) {
	if userID != "" {
		if strings.TrimSpace(userID) != userID {
			return "", "", &domain.ValidationError{Err: domain.ErrInvalidUserID}
		}
		if _, err := uuid.Parse(userID); err != nil {
			return "", "", &domain.ValidationError{Err: domain.ErrInvalidUserID}
		}
	}

	if serviceName != "" {
		if strings.TrimSpace(serviceName) == "" || strings.TrimSpace(serviceName) != serviceName {
			return "", "", &domain.ValidationError{Err: domain.ErrInvalidServiceName}
		}
	}

	return userID, serviceName, nil
}

func validateTotalFilter(filter domain.Subscription) (domain.Subscription, error) {
	if strings.TrimSpace(filter.StartDate) == "" || filter.EndDate == nil || strings.TrimSpace(*filter.EndDate) == "" {
		return domain.Subscription{}, &domain.ValidationError{Err: domain.ErrMissingRequiredFields}
	}

	if strings.TrimSpace(filter.StartDate) != filter.StartDate {
		return domain.Subscription{}, &domain.ValidationError{Err: domain.ErrInvalidFromDate}
	}

	fromDate, err := parseMonthYear(filter.StartDate)
	if err != nil {
		return domain.Subscription{}, &domain.ValidationError{Err: domain.ErrInvalidFromDate}
	}

	if strings.TrimSpace(*filter.EndDate) != *filter.EndDate {
		return domain.Subscription{}, &domain.ValidationError{Err: domain.ErrInvalidToDate}
	}

	toDate, err := parseMonthYear(*filter.EndDate)
	if err != nil {
		return domain.Subscription{}, &domain.ValidationError{Err: domain.ErrInvalidToDate}
	}

	if toDate.Before(fromDate) {
		return domain.Subscription{}, &domain.ValidationError{Err: domain.ErrInvalidPeriod}
	}

	formattedFrom := fromDate.Format(monthYearLayout)
	formattedTo := toDate.Format(monthYearLayout)
	filter.StartDate = formattedFrom
	filter.EndDate = &formattedTo

	if filter.UserID != "" {
		if strings.TrimSpace(filter.UserID) != filter.UserID {
			return domain.Subscription{}, &domain.ValidationError{Err: domain.ErrInvalidUserID}
		}
		if _, err := uuid.Parse(filter.UserID); err != nil {
			return domain.Subscription{}, &domain.ValidationError{Err: domain.ErrInvalidUserID}
		}
	}

	if filter.ServiceName != "" {
		if strings.TrimSpace(filter.ServiceName) == "" || strings.TrimSpace(filter.ServiceName) != filter.ServiceName {
			return domain.Subscription{}, &domain.ValidationError{Err: domain.ErrInvalidServiceName}
		}
	}

	return filter, nil
}

func parseMonthYear(value string) (time.Time, error) {
	parsed, err := time.Parse(monthYearLayout, value)
	if err != nil {
		return time.Time{}, err
	}
	return parsed, nil
}
