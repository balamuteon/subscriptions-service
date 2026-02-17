package subscription

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"subscription_service/internal/domain"
)

type Repository struct {
	db dbExecutor
}

func New(db dbExecutor) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, sub domain.Subscription) (string, error) {
	query := `
		INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date)
		VALUES ($1, $2, $3, to_date($4, 'MM-YYYY'), to_date($5, 'MM-YYYY'))
		RETURNING id
	`

	var id uuid.UUID
	err := r.db.QueryRow(ctx, query, sub.ServiceName, sub.Price, sub.UserID, sub.StartDate, sub.EndDate).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("create subscription: %w", err)
	}

	return id.String(), nil
}

func (r *Repository) GetByID(ctx context.Context, id string) (domain.Subscription, error) {
	query := `
		SELECT
			id,
			service_name,
			price,
			user_id,
			to_char(start_date, 'MM-YYYY'),
			CASE WHEN end_date IS NULL THEN NULL ELSE to_char(end_date, 'MM-YYYY') END
		FROM subscriptions
		WHERE id = $1
	`

	var sub domain.Subscription
	var parsedID uuid.UUID
	var userID uuid.UUID
	var startDate string
	var endDate sql.NullString

	err := r.db.QueryRow(ctx, query, id).Scan(
		&parsedID,
		&sub.ServiceName,
		&sub.Price,
		&userID,
		&startDate,
		&endDate,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Subscription{}, domain.ErrSubscriptionNotFound
		}
		return domain.Subscription{}, fmt.Errorf("get subscription by id: %w", err)
	}

	sub.ID = parsedID.String()
	sub.UserID = userID.String()
	sub.StartDate = startDate
	if endDate.Valid {
		sub.EndDate = &endDate.String
	}

	return sub, nil
}

func (r *Repository) List(ctx context.Context, userID string, serviceName string) ([]domain.Subscription, error) {
	queryBuilder := strings.Builder{}
	queryBuilder.WriteString(`
		SELECT
			id,
			service_name,
			price,
			user_id,
			TO_CHAR(start_date, 'MM-YYYY'),
			CASE WHEN end_date IS NULL THEN NULL ELSE to_char(end_date, 'MM-YYYY') END
		FROM subscriptions
`)

	args := make([]any, 0, 2)
	conditions := make([]string, 0, 2)

	if userID != "" {
		args = append(args, userID)
		conditions = append(conditions, fmt.Sprintf("user_id = $%d", len(args)))
	}

	if serviceName != "" {
		args = append(args, serviceName)
		conditions = append(conditions, fmt.Sprintf("service_name = $%d", len(args)))
	}

	if len(conditions) > 0 {
		queryBuilder.WriteString(" WHERE ")
		queryBuilder.WriteString(strings.Join(conditions, " AND "))
	}

	queryBuilder.WriteString(" ORDER BY created_at DESC")

	rows, err := r.db.Query(ctx, queryBuilder.String(), args...)
	if err != nil {
		return nil, fmt.Errorf("list subscriptions: %w", err)
	}
	defer rows.Close()

	result := make([]domain.Subscription, 0)
	for rows.Next() {
		var sub domain.Subscription
		var parsedID uuid.UUID
		var parsedUserID uuid.UUID
		var startDate string
		var endDate sql.NullString

		if err := rows.Scan(
			&parsedID,
			&sub.ServiceName,
			&sub.Price,
			&parsedUserID,
			&startDate,
			&endDate,
		); err != nil {
			return nil, fmt.Errorf("scan listed subscription: %w", err)
		}

		sub.ID = parsedID.String()
		sub.UserID = parsedUserID.String()
		sub.StartDate = startDate
		if endDate.Valid {
			sub.EndDate = &endDate.String
		}

		result = append(result, sub)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate listed subscriptions: %w", err)
	}

	return result, nil
}

func (r *Repository) Update(ctx context.Context, sub domain.Subscription) error {
	query := `
		UPDATE subscriptions
		SET
			service_name = $2,
			price = $3,
			user_id = $4,
			start_date = to_date($5, 'MM-YYYY'),
			end_date = to_date($6, 'MM-YYYY')
		WHERE id = $1
	`

	result, err := r.db.Exec(ctx, query, sub.ID, sub.ServiceName, sub.Price, sub.UserID, sub.StartDate, sub.EndDate)
	if err != nil {
		return fmt.Errorf("update subscription: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrSubscriptionNotFound
	}

	return nil
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	result, err := r.db.Exec(ctx, `DELETE FROM subscriptions WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete subscription: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrSubscriptionNotFound
	}

	return nil
}

func (r *Repository) Total(ctx context.Context, filter domain.Subscription) (int64, error) {
	queryBuilder := strings.Builder{}
	queryBuilder.WriteString(`
		WITH bounds AS (
				SELECT
						to_date($1, 'MM-YYYY') AS from_date,
						to_date($2, 'MM-YYYY') AS to_date
		)
		SELECT COALESCE(SUM(
				s.price * (
				(
					(
						EXTRACT(YEAR FROM LEAST(COALESCE(s.end_date, b.to_date), b.to_date)) * 12 +
						EXTRACT(MONTH FROM LEAST(COALESCE(s.end_date, b.to_date), b.to_date))
					) - (
						EXTRACT(YEAR FROM GREATEST(s.start_date, b.from_date)) * 12 +
						EXTRACT(MONTH FROM GREATEST(s.start_date, b.from_date))
					) + 1
				)
			)
		), 0)::bigint
		FROM subscriptions s
		CROSS JOIN bounds b
		WHERE s.start_date <= b.to_date
			AND COALESCE(s.end_date, b.to_date) >= b.from_date
	`)

	args := []any{filter.StartDate, *filter.EndDate}

	if filter.UserID != "" {
		args = append(args, filter.UserID)
		fmt.Fprintf(&queryBuilder, " AND s.user_id = $%d", len(args))
	}

	if filter.ServiceName != "" {
		args = append(args, filter.ServiceName)
		fmt.Fprintf(&queryBuilder, " AND s.service_name = $%d", len(args))
	}

	var total int64
	if err := r.db.QueryRow(ctx, queryBuilder.String(), args...).Scan(&total); err != nil {
		return 0, fmt.Errorf("calculate subscriptions total: %w", err)
	}

	return total, nil
}
